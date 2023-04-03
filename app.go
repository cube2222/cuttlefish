package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
	sync.Mutex
	messages  []ChatMessage
	openAICli *openai.Client
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		openAICli: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
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

type ChatMessage struct {
	SentBySelf bool   `json:"sentBySelf"`
	Content    string `json:"message"`
}

// Greet returns a greeting for the given name
func (a *App) Messages(conversationID int) []ChatMessage {
	return a.messages
}

func (a *App) SendMessage(conversationID int, message ChatMessage) error {
	// TODO: concurrency, locking
	a.messages = append(a.messages, message)
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", 42))
	stream, err := a.openAICli.CreateChatCompletionStream(context.Background(), openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		MaxTokens:   500,
		Temperature: 0.7,
		TopP:        1,
		Messages:    MessagesToGPTMessages(a.messages),
	})
	if err != nil {
		return fmt.Errorf("couldn't create chat completion stream: %w", err)
	}
	a.messages = append(a.messages, ChatMessage{
		SentBySelf: false,
		Content:    " ",
	})
	gptMessage := &a.messages[len(a.messages)-1]
	i := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("couldn't receive from chat completion stream: %w", err)
		}

		if len(res.Choices) > 0 {
			gptMessage.Content += res.Choices[0].Delta.Content
			runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", 42))
			i++
		}
	}
	log.Println("events: ", i)
	return nil
}

func MessagesToGPTMessages(messages []ChatMessage) []openai.ChatCompletionMessage {
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
