package main

import (
	"context"
	"fmt"
	"io"
	"os"
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

func (a *App) Messages(conversationID int) ([]*database.Message, error) {
	return a.queries.ListMessages(a.ctx, conversationID)
}

func (a *App) Conversations() ([]*database.Conversation, error) {
	return a.queries.ListConversations(a.ctx)
}

// func (a *App) CreateConversation(params database.CreateConversationParams) (*database.Conversation, error) {
// 	return a.queries.CreateConversation(a.ctx, params)
// }

func (a *App) ResetConversation(conversationID int) error {
	if err := a.queries.DeleteMessages(a.ctx); err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", 42))
	return nil
}

func (a *App) SendMessage(conversationID *int, content string) (err error) {
	// TODO: setConversation callback in both the chat view and the conversations sidebar.
	//		 In one it will "open" the newly created conversation. In the other it will open the selected conversation.
	//		 Also, we need to be listening for "conversations-updated".
	if conversationID == nil {
		title := content
		if len(title) > 20 {
			title = title[:13] + "..."
		}
		conversation, err := a.queries.CreateConversation(a.ctx, database.CreateConversationParams{
			Title:           title,
			LastMessageTime: time.Now(),
		})
		if err != nil {
			return err
		}
		conversationID = &conversation.ID
	}

	if _, err := a.queries.CreateMessage(a.ctx, database.CreateMessageParams{
		ConversationID: *conversationID,
		Content:        content,
		SentBySelf:     true,
	}); err != nil {
		return err
	}
	allMessages, err := a.queries.ListMessages(a.ctx, *conversationID)
	if err != nil {
		return err
	}
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", 42))
	stream, err := a.openAICli.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   500,
		Temperature: 0.7,
		TopP:        1,
		Messages:    MessagesToGPTMessages(allMessages),
		// TODO: For tools you'll need to pass "Observation" as a stop phrase.
	})
	if err != nil {
		return fmt.Errorf("couldn't create chat completion stream: %w", err)
	}
	gptMessage, err := a.queries.CreateMessage(a.ctx, database.CreateMessageParams{
		ConversationID: *conversationID,
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
			runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", 42))
		}
	}
	return nil
}

func MessagesToGPTMessages(messages []*database.Message) []openai.ChatCompletionMessage {
	var gptMessages []openai.ChatCompletionMessage
	gptMessages = append(gptMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: "You are a helpful assistant. Please respond to the user's messages as best as you can.",
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
