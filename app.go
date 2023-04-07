package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gptui/database"
)

// App struct
type App struct {
	ctx context.Context
	sync.Mutex
	openAICli *openai.Client
	queries   *database.Queries
}

// NewApp creates a new App application struct
func NewApp(ctx context.Context, queries *database.Queries) *App {
	return &App{
		ctx:       ctx,
		openAICli: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		queries:   queries,
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) Messages(conversationID int) ([]database.Message, error) {
	return a.queries.ListMessages(a.ctx, conversationID)
}

func (a *App) Conversations() ([]database.Conversation, error) {
	return a.queries.ListConversations(a.ctx)
}

// func (a *App) CreateConversation(params database.CreateConversationParams) (*database.Conversation, error) {
// 	return a.queries.CreateConversation(a.ctx, params)
// }

func (a *App) DeleteConversation(conversationID int) error {
	if err := a.queries.DeleteConversation(a.ctx, conversationID); err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "conversations-updated")
	return nil
}

func (a *App) SendMessage(conversationID int, content string) (database.Message, error) {
	// TODO: setConversation callback in both the chat view and the conversations sidebar.
	//		 In one it will "open" the newly created conversation. In the other it will open the selected conversation.
	//		 Also, we need to be listening for "conversations-updated".
	if conversationID == -1 {
		title := content
		if len(title) > 20 {
			title = title[:13] + "..."
		}
		conversation, err := a.queries.CreateConversation(a.ctx, database.CreateConversationParams{
			Title:           title,
			LastMessageTime: time.Now(),
		})
		if err != nil {
			return database.Message{}, err
		}
		conversationID = conversation.ID
		runtime.EventsEmit(a.ctx, "conversations-updated")
	}

	msg, err := a.queries.CreateMessage(a.ctx, database.CreateMessageParams{
		ConversationID: conversationID,
		Content:        content,
		SentBySelf:     true,
	})
	if err != nil {
		return database.Message{}, err
	}
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))

	go func() {
		if err := func() error {
			allMessages, err := a.queries.ListMessages(a.ctx, conversationID)
			if err != nil {
				return err
			}
			stream, err := a.openAICli.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
				Model:       openai.GPT3Dot5Turbo,
				MaxTokens:   500,
				Temperature: 0.7,
				TopP:        1,
				Messages:    MessagesToGPTMessages(allMessages),
				Stop:        []string{"Observation"},
				// TODO: For tools you'll need to pass "Observation" as a stop phrase.
			})
			if err != nil {
				return fmt.Errorf("couldn't create chat completion stream: %w", err)
			}
			gptMessage, err := a.queries.CreateMessage(a.ctx, database.CreateMessageParams{
				ConversationID: conversationID,
				Content:        "",
				SentBySelf:     false,
			})
			if err != nil {
				return err
			}
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					break
				} else if err != nil {
					return fmt.Errorf("couldn't receive from chat completion stream: %w", err)
				}

				if len(res.Choices) > 0 {
					gptMessage.Content += res.Choices[0].Delta.Content
					if _, err := a.queries.AppendMessage(a.ctx, database.AppendMessageParams{
						ID:      gptMessage.ID,
						Content: res.Choices[0].Delta.Content,
					}); err != nil {
						return err
					}
					runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))
				}
			}
			// TODO: If the message contains the string "```action" then we should interpret that and automatically respond with the action's result.
			//       Then, rinse and repeat.
			return nil
		}(); err != nil {
			runtime.EventsEmit(a.ctx, "async-error", err.Error())
			log.Println("error generating streaming chatgpt response:", err)
		}
	}()

	return msg, nil
}

func MessagesToGPTMessages(messages []database.Message) []openai.ChatCompletionMessage {
	// TODO: By default, have two modes, "tool use" and "casual".
	// TODO: In the upper right corner you should be able to select the list of tools.
	systemMessage := `List of available tools:
	
[
  {
    "tool": "terminal",
    "args": {
      "command": "<bash command to run>"
    }
  }
]
	
You are a helpful assistant on a MacOS system. You may additionally use tools repeatedly to aid your responses, but should always first describe your thought process, like this:
Thought: <always write out what you think>
Action:
<backticks>action
{
  "tool": "<tool name>",
  "args": {
	"<arg name>": <arg value>,
	...
  }
}
<backticks>
Then you'll receive a response as follows:
Observation:
<backticks>
<The tool's response>
<backticks>

For example (this tool doesn't necessarily exist):
Thought: I need to use the add tool to add 5 and 7.
Action:
<backticks>action
{
  "tool": "add",
  "args": {
    "num1": 5,
    "num2": 7
  }
}
<backticks>
Observation:
<backticks>
12
<backticks>

You can use tools repeatedly, or provide a final answer to the user.
Please respond to the user's messages as best as you can.`
	systemMessage = strings.ReplaceAll(systemMessage, "<backticks>", "```")

	var gptMessages []openai.ChatCompletionMessage
	gptMessages = append(gptMessages, openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleSystem,
		// Content: "You are a helpful assistant. Please respond to the user's messages as best as you can.",
		Content: systemMessage,
	})
	for _, message := range messages {
		gptMessage := openai.ChatCompletionMessage{
			Content: message.Content,
		}
		if message.SentBySelf {
			gptMessage.Role = openai.ChatMessageRoleUser
		} else {
			gptMessage.Role = openai.ChatMessageRoleSystem
		}
		gptMessages = append(gptMessages, gptMessage)
	}
	return gptMessages
}

// TODO: Add event "async error" for background processing.
