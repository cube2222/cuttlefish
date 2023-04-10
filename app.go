package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/sashabaranov/go-openai"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slices"

	"gptui/database"
	"gptui/tools"
	"gptui/tools/chart"
	"gptui/tools/dalle2"
	"gptui/tools/geturl"
	"gptui/tools/python"
	"gptui/tools/search"
	"gptui/tools/terminal"
)

// App struct
type App struct {
	ctx      context.Context
	queries  *database.Queries
	tools    map[string]tools.Tool
	settings database.Settings

	m                       sync.Mutex
	generationContextCancel map[int]context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp(ctx context.Context, queries *database.Queries) *App {
	out := &App{
		ctx:     ctx,
		queries: queries,
		tools: map[string]tools.Tool{
			"terminal":       &terminal.Tool{},
			"generate_image": &dalle2.Tool{},
			"search":         &search.Tool{},
			"get_url":        &geturl.Tool{},
			"chart":          &chart.Tool{},
			"python":         &python.Tool{},
		},
		generationContextCancel: map[int]context.CancelFunc{},
	}

	settings, err := out.getSettingsRaw()
	if err != nil {
		log.Printf("couldn't load settings: %w", err)
	}
	out.settings = settings

	return out
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

func (a *App) openAICli() *openai.Client {
	a.m.Lock()
	defer a.m.Unlock()
	return openai.NewClient(a.settings.OpenAIAPIKey)
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
		return fmt.Errorf("coudn't delete conversation: %w", err)
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
		defaultConversationSettings, err := a.GetDefaultConversationSettings()
		if err != nil {
			return database.Message{}, fmt.Errorf("couldn't get default conversation settings: %w", err)
		}
		settings, err := a.queries.CreateConversationSettings(a.ctx, database.CreateConversationSettingsParams{
			SystemPromptTemplate: defaultConversationSettings.SystemPromptTemplate,
			ToolsEnabled:         defaultConversationSettings.ToolsEnabled,
		})
		if err != nil {
			return database.Message{}, fmt.Errorf("couldn't create conversation settings: %w", err)
		}
		conversation, err := a.queries.CreateConversation(a.ctx, database.CreateConversationParams{
			ConversationSettingsID: settings.ID,
			Title:                  title,
			LastMessageTime:        time.Now(),
		})
		if err != nil {
			return database.Message{}, fmt.Errorf("couldn't create conversation: %w", err)
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
		return database.Message{}, fmt.Errorf("couldn't create message: %w", err)
	}
	runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))

	go func() {
		if err := a.runChainOfMessages(conversationID); err != nil && !errors.Is(err, context.Canceled) {
			runtime.EventsEmit(a.ctx, "async-error", err.Error())
		}
	}()

	return msg, nil
}

func (a *App) runChainOfMessages(conversationID int) error {
	genCtx, cancelGeneration := context.WithCancel(a.ctx)
	defer cancelGeneration()

	a.m.Lock()
	a.generationContextCancel[conversationID] = cancelGeneration
	a.m.Unlock()

	if err := a.queries.MarkGenerationStarted(genCtx, conversationID); err != nil {
		return err
	}
	runtime.EventsEmit(genCtx, fmt.Sprintf("conversation-%d-updated", conversationID))

	curConversation, err := a.queries.GetConversation(genCtx, conversationID)
	if err != nil {
		return fmt.Errorf("couldn't get conversation: %w", err)
	}
	curConversationSettings, err := a.queries.GetConversationSettings(genCtx, curConversation.ConversationSettingsID)
	if err != nil {
		return fmt.Errorf("couldn't get conversation settings: %w", err)
	}

	defer func() {
		if err := a.queries.MarkGenerationDone(a.ctx, conversationID); err != nil {
			runtime.EventsEmit(a.ctx, "async-error", fmt.Errorf("couldn't mark conversation as done generating: %w", err).Error())
		}
		runtime.EventsEmit(a.ctx, fmt.Sprintf("conversation-%d-updated", conversationID))
	}()

	cachedToolInstances := map[string]tools.ToolInstance{}
	defer func() {
		for name, instance := range cachedToolInstances {
			if err := instance.Shutdown(); err != nil {
				runtime.EventsEmit(a.ctx, "async-error", fmt.Errorf("couldn't shut down tool `%s`: %w", name, err).Error())
			}
		}
	}()

	stop := []string{"Observation", "Response"}
	retries := 0
	for {
		allMessages, err := a.queries.ListMessages(genCtx, conversationID)
		if err != nil {
			return fmt.Errorf("couldn't list conversation messages: %w", err)
		}
		gptMessages, err := a.messagesToGPTMessages(curConversationSettings, allMessages)
		if err != nil {
			return fmt.Errorf("couldn't convert messages to GPT messages: %w", err)
		}
		stream, err := a.openAICli().CreateChatCompletionStream(genCtx, openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			MaxTokens:   500,
			Temperature: 0.7,
			TopP:        1,
			Messages:    gptMessages,
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
			return fmt.Errorf("couldn't create response message: %w", err)
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
					return fmt.Errorf("couldn't append to message: %w", err)
				}
				runtime.EventsEmit(genCtx, fmt.Sprintf("conversation-%d-updated", conversationID))
			}
		}
		gptMessage, err = a.queries.GetMessage(genCtx, gptMessage.ID)
		if err != nil {
			return fmt.Errorf("couldn't get response message: %w", err)
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
				// TODO: respond as observation
				return fmt.Errorf("couldn't decode action: %w", err)
			}

			toolInstance, ok := cachedToolInstances[action.Tool]
			if !ok {
				tool, ok := a.tools[action.Tool]
				if !ok {
					// TODO: respond as observation
					return fmt.Errorf("tool `%s` not found", action.Tool)
				}
				toolInstance, err = tool.Instantiate(genCtx, a.settings)
				if err != nil {
					// TODO: respond as observation
					return fmt.Errorf("couldn't instantiate tool `%s`: %w", action.Tool, err)
				}
			}

			result, err := toolInstance.Run(genCtx, action.Args)
			if err != nil {
				// TODO: respond as observation
				return fmt.Errorf("couldn't run tool `%s`: %w", action.Tool, err)
			}
			observationString := "Observation: "
			observationString += result.Result
			observationString += "\n"
			observationString += "```"
			if result.CustomResultTag != "" {
				observationString += result.CustomResultTag
			}
			observationString += "\n"
			observationString += result.Output + "\n```"

			if _, err := a.queries.CreateMessage(genCtx, database.CreateMessageParams{
				ConversationID: conversationID,
				Content:        observationString,
				Author:         a.tools[action.Tool].Name(),
			}); err != nil {
				return fmt.Errorf("couldn't create observation message: %w", err)
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

func (a *App) messagesToGPTMessages(conversationSettings database.ConversationSetting, messages []database.Message) ([]openai.ChatCompletionMessage, error) {
	generatedSystemPrompt, err := a.generateSystemPrompt(conversationSettings)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate system prompt: %w", err)
	}

	var gptMessages []openai.ChatCompletionMessage
	gptMessages = append(gptMessages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: generatedSystemPrompt,
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
			gptMessage.Content = fmt.Sprintf("`%s` response:", message.Author) + gptMessage.Content
			// gptMessage.Content += "\nMake sure not to start your next response with `Observation:`, nor by thanking for this reminder."
		}
		gptMessages = append(gptMessages, gptMessage)
	}
	return gptMessages, nil
}

func (a *App) RerunFromMessage(conversationID int, messageID int) error {
	if err := a.queries.ResetConversationFrom(a.ctx, database.ResetConversationFromParams{
		ConversationID: conversationID,
		ID:             messageID,
	}); err != nil {
		return fmt.Errorf("couldn't reset conversation: %w", err)
	}

	go func() {
		if err := a.runChainOfMessages(conversationID); err != nil && !errors.Is(err, context.Canceled) {
			runtime.EventsEmit(a.ctx, "async-error", err.Error())
		}
	}()

	return nil
}

//go:embed default_system_prompt.gotmpl
var defaultSystemPromptTemplate string

func (a *App) generateSystemPrompt(conversationSettings database.ConversationSetting) (string, error) {
	var params struct {
		ToolsDescription string
		AnyToolsEnabled  bool
	}

	type toolDescription struct {
		Tool        string            `json:"tool"`
		Description string            `json:"description"`
		Args        map[string]string `json:"args"`
	}

	toolsDescription := []toolDescription{}
	for toolName, tool := range a.tools {
		if !slices.Contains(conversationSettings.ToolsEnabled, toolName) {
			continue
		}
		toolsDescription = append(toolsDescription, toolDescription{
			Tool:        toolName,
			Description: tool.Description(),
			Args:        tool.ArgumentDescriptions(),
		})
	}
	data, err := json.MarshalIndent(toolsDescription, "", "  ")
	if err != nil {
		return "", fmt.Errorf("couldn't encode tools description: %w", err)
	}
	params.ToolsDescription = string(data)
	params.AnyToolsEnabled = len(toolsDescription) > 0

	var buf bytes.Buffer
	// TODO: Use custom template.
	tmpl, err := template.New("system_prompt").Parse(conversationSettings.SystemPromptTemplate)
	if err != nil {
		return "", fmt.Errorf("couldn't parse system prompt template: %w", err)
	}
	if err := tmpl.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("couldn't execute system prompt template: %w", err)
	}

	return buf.String(), nil
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

func (a *App) UpdateConversationSettings(params database.UpdateConversationSettingsParams) (database.ConversationSetting, error) {
	return a.queries.UpdateConversationSettings(a.ctx, params)
}

func (a *App) GetDefaultConversationSettings() (database.ConversationSetting, error) {
	conversationSettings, err := a.queries.GetDefaultConversationSettings(a.ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return database.ConversationSetting{
			ID:                   -1,
			SystemPromptTemplate: defaultSystemPromptTemplate,
			ToolsEnabled:         []string{"terminal", "python", "get_url", "chart"},
		}, nil
	} else if err != nil {
		return database.ConversationSetting{}, fmt.Errorf("couldn't get default conversation settings: %w", err)
	}
	return conversationSettings, nil
}

func (a *App) SetDefaultConversationSettings(params database.CreateDefaultConversationSettingsParams) (database.ConversationSetting, error) {
	defaultConversationSettings, err := a.queries.GetDefaultConversationSettings(a.ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return a.queries.CreateDefaultConversationSettings(a.ctx, params)
	} else if err != nil {
		return database.ConversationSetting{}, fmt.Errorf("couldn't get default conversation settings: %w", err)
	}

	return a.queries.UpdateConversationSettings(a.ctx, database.UpdateConversationSettingsParams{
		ID:                   defaultConversationSettings.ID,
		SystemPromptTemplate: params.SystemPromptTemplate,
		ToolsEnabled:         params.ToolsEnabled,
	})
}

func (a *App) getSettingsRaw() (database.Settings, error) {
	keyValue, err := a.queries.GetKeyValue(a.ctx, "settings")
	if errors.Is(err, sql.ErrNoRows) {
		return database.Settings{
			Model: "gpt-3.5-turbo",
			Python: database.PythonSettings{
				InterpreterPath: "python3",
			},
		}, nil
	} else if err != nil {
		return database.Settings{}, err
	}
	var settings database.Settings
	if err := json.Unmarshal([]byte(keyValue.Value), &settings); err != nil {
		return database.Settings{}, err
	}
	a.m.Lock()
	defer a.m.Unlock()
	a.settings = settings // Just to be safe.

	return settings, nil
}

func (a *App) GetSettings() (database.Settings, error) {
	settings, err := a.getSettingsRaw()
	if err != nil {
		return database.Settings{}, err
	}
	if settings.OpenAIAPIKey != "" {
		settings.OpenAIAPIKey = "*****"
	}
	return settings, nil
}

func (a *App) SaveSettings(settings database.Settings) error {
	_, err := a.queries.GetKeyValue(a.ctx, "settings")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("couldn't check if settings keyvalue exists: %w", err)
	}
	settingsExist := !errors.Is(err, sql.ErrNoRows)

	oldSettings, err := a.getSettingsRaw()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if settings.OpenAIAPIKey == "*****" {
		settings.OpenAIAPIKey = oldSettings.OpenAIAPIKey
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	if settingsExist {
		if err := a.queries.UpdateKeyValue(a.ctx, database.UpdateKeyValueParams{
			Key:   "settings",
			Value: string(settingsJSON),
		}); err != nil {
			return fmt.Errorf("couldn't save settings: %w", err)
		}
	} else {
		if err := a.queries.CreateKeyValue(a.ctx, database.CreateKeyValueParams{
			Key:   "settings",
			Value: string(settingsJSON),
		}); err != nil {
			return fmt.Errorf("couldn't save settings: %w", err)
		}
	}

	a.m.Lock()
	defer a.m.Unlock()
	a.settings = settings

	return nil
}

type AvailableTool struct {
	Name string `json:"name"`
	ID   string `json:"ID"`
}

func (a *App) GetAvailableTools() []AvailableTool {
	var out []AvailableTool
	for id, tool := range a.tools {
		out = append(out, AvailableTool{
			Name: tool.Name(),
			ID:   id,
		})
	}
	slices.SortFunc(out, func(a, b AvailableTool) bool {
		return a.Name < b.Name
	})
	return out
}
