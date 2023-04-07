package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"gptui/database"
)

// App struct
type App struct {
	ctx       context.Context
	openAICli *openai.Client
	queries   *database.Queries

	m                       sync.Mutex
	generationContextCancel map[int]context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp(ctx context.Context, queries *database.Queries) *App {
	return &App{
		ctx:                     ctx,
		openAICli:               openai.NewClient(os.Getenv("OPENAI_API_KEY")),
		queries:                 queries,
		generationContextCancel: map[int]context.CancelFunc{},
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

func (a *App) GetConversation(conversationID int) (database.Conversation, error) {
	return a.queries.GetConversation(a.ctx, conversationID)
}

func (a *App) Conversations() ([]database.Conversation, error) {
	return a.queries.ListConversations(a.ctx)
}

func (a *App) DeleteConversation(conversationID int) error {
	if err := a.queries.DeleteConversation(a.ctx, conversationID); err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, "conversations-updated")
	return nil
}

func (a *App) SendMessage(conversationID int, content string) (database.Message, error) {
	if conversationID == -1 {
		title := content
		if len(title) > 20 {
			title = title[:13] + "..."
		}
		settings, err := a.queries.CreateConversationSettings(a.ctx, database.CreateConversationSettingsParams{
			SystemPromptTemplate: "You are a helpful assistant. Respond to the queries as best as you can.",
			ToolsEnabled:         []string{},
		})
		if err != nil {
			return database.Message{}, err
		}
		conversation, err := a.queries.CreateConversation(a.ctx, database.CreateConversationParams{
			ConversationSettingsID: settings.ID,
			Title:                  title,
			LastMessageTime:        time.Now(),
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
		Author:         "user",
	})
	if err != nil {
		return database.Message{}, err
	}
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))

	go func() {
		if err := a.runChainOfMessages(conversationID); err != nil && !errors.Is(err, context.Canceled) {
			runtime.EventsEmit(a.ctx, "async-error", err.Error())
			log.Println("error generating streaming chatgpt response:", err)
		}
	}()

	return msg, nil
}

func (a *App) runChainOfMessages(conversationID int) error {
	// TODO: Add Dalle2
	genCtx, cancelGeneration := context.WithCancel(a.ctx)
	defer cancelGeneration()

	a.m.Lock()
	a.generationContextCancel[conversationID] = cancelGeneration
	a.m.Unlock()

	if err := a.queries.MarkGenerationStarted(genCtx, conversationID); err != nil {
		return err
	}
	runtime.EventsEmit(genCtx, fmt.Sprintf("conversation-%d-updated", conversationID))

	defer func() {
		if err := a.queries.MarkGenerationDone(a.ctx, conversationID); err != nil {
			runtime.EventsEmit(a.ctx, "async-error", fmt.Errorf("couldn't mark conversation as done generating: %w", err).Error())
		}
		runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))
	}()

	stop := []string{"Observation"}
	retries := 0
	for {
		allMessages, err := a.queries.ListMessages(genCtx, conversationID)
		if err != nil {
			return err
		}
		stream, err := a.openAICli.CreateChatCompletionStream(genCtx, openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   500,
			Temperature: 0.7,
			TopP:        1,
			Messages:    MessagesToGPTMessages(allMessages),
			Stop:        stop,
		})
		if err != nil {
			return fmt.Errorf("couldn't create chat completion stream: %w", err)
		}
		gptMessage, err := a.queries.CreateMessage(genCtx, database.CreateMessageParams{
			ConversationID: conversationID,
			Content:        "",
			Author:         "assistant",
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
				if _, err := a.queries.AppendMessage(genCtx, database.AppendMessageParams{
					ID:      gptMessage.ID,
					Content: res.Choices[0].Delta.Content,
				}); err != nil {
					return err
				}
				runtime.EventsEmit(genCtx, fmt.Sprintf("conversation-%d-updated", conversationID))
			}
		}
		gptMessage, err = a.queries.GetMessage(genCtx, gptMessage.ID)
		if err != nil {
			return err
		}
		if strings.TrimSpace(gptMessage.Content) == "" {
			stop = []string{}
			if retries > 2 {
				return fmt.Errorf("couldn't generate a response after %d retries", retries)
			}
			retries++
			continue
		}
		toolUseEnabled := true
		if toolUseEnabled && (strings.Contains(gptMessage.Content, "```action") || strings.Contains(gptMessage.Content, "Action:")) {
			// TODO: Make it so that each tool use can be approved by the user.

			// A tool has been called upon!
			// We match on either, cause ChatGPT doesn't always use the same format.
			content := gptMessage.Content
			if strings.Contains(gptMessage.Content, "```action") {
				content = content[strings.Index(content, "```action")+len("```action"):]
				content = content[:strings.Index(content, "```")]
				content = strings.TrimSpace(content)
			} else {
				content = content[strings.Index(content, "Action:"):]
				content = content[strings.Index(content, "```"):]
				content = content[strings.Index(content, "\n"):]
				content = content[:strings.Index(content, "```")]
				content = strings.TrimSpace(content)
			}

			var action Action
			if err := json.Unmarshal([]byte(content), &action); err != nil {
				return err // TODO: respond as observation
			}

			switch action.Tool {
			// TODO: Put these into a map. Implementing an interface.
			case "terminal":
				command, ok := action.Args["command"].(string)
				if !ok {
					return fmt.Errorf("command is not a string")
				}
				cmd := exec.CommandContext(genCtx, "bash", "-c", command)
				var buf bytes.Buffer // TODO: Stream output to a message.
				cmd.Stdout = &buf
				cmd.Stderr = &buf
				err := cmd.Run()
				observationString := "Observation: "
				if err == context.Canceled {
					return err
				} else if err != nil {
					observationString += err.Error()
				} else {
					observationString += "successfully executed `" + command + "`"
				}
				observationString += "\n"
				observationString += "```\n" + buf.String() + "\n```"

				if _, err := a.queries.CreateMessage(genCtx, database.CreateMessageParams{
					ConversationID: conversationID,
					Content:        observationString,
					Author:         action.Tool, // TODO: Fixme
				}); err != nil {
					return err
				}
			}
		} else {
			break
		}
	}
	return nil
}

type Action struct {
	Tool string                 `json:"tool"`
	Args map[string]interface{} `json:"args"`
}

func MessagesToGPTMessages(messages []database.Message) []openai.ChatCompletionMessage {
	// TODO: By default, have two modes, "tool use" and "casual".
	//       In tool use, don't show all of this if no tools are available.
	// TODO: In the upper right corner you should be able to select the list of tools.
	// TODO: Make this a Go template.
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
		if strings.TrimSpace(message.Content) == "" {
			continue
		}
		gptMessage := openai.ChatCompletionMessage{
			Content: message.Content,
		}
		if message.Author == "assistant" {
			gptMessage.Role = openai.ChatMessageRoleAssistant
		} else if message.Author == "user" {
			gptMessage.Role = openai.ChatMessageRoleUser
		} else {
			gptMessage.Role = openai.ChatMessageRoleUser
			gptMessage.Content += "\nMake sure not to start your next response with `Observation:`"
		}
		gptMessages = append(gptMessages, gptMessage)
	}
	return gptMessages
}

func (a *App) CancelGeneration(conversationID int) {
	a.m.Lock()
	defer a.m.Unlock()
	cancel, ok := a.generationContextCancel[conversationID]
	if !ok {
		return
	}
	cancel()
}

func (a *App) GetConversationSettings(conversationSettingsID int) (database.ConversationSetting, error) {
	return a.queries.GetConversationSettings(a.ctx, conversationSettingsID)
}

func (a *App) UpdateConversationSettings(params database.UpdateConversationSettingsParams) error {
	return a.queries.UpdateConversationSettings(a.ctx, params)
}

type Settings struct {
	OpenAIAPIKey string `json:"openAiApiKey"`
	Model        string `json:"model"`
	// Add nested struct per configurable plugin below.
}

func (a *App) GetSettings() (Settings, error) {
	keyValue, err := a.queries.GetKeyValue(a.ctx, "settings")
	if err != nil {
		return Settings{}, err
	}
	var settings Settings
	if err := json.Unmarshal([]byte(keyValue.Value), &settings); err != nil {
		return Settings{}, err
	}
	settings.OpenAIAPIKey = "*****"
	return settings, nil
}

func (a *App) SaveSettings(settings Settings) error {
	oldSettings, err := a.GetSettings()
	if err != nil {
		// TODO: Handle not found.
		return err
	}
	if settings.OpenAIAPIKey == "*****" {
		settings.OpenAIAPIKey = oldSettings.OpenAIAPIKey
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	return a.queries.SetKeyValue(a.ctx, database.SetKeyValueParams{
		Key:   "settings",
		Value: string(settingsJSON),
	})
}
