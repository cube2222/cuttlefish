// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package database

import (
	"database/sql"
	"time"
)

type Conversation struct {
	ID                     int       `json:"id"`
	ConversationSettingsID int       `json:"conversationSettingsID"`
	Title                  string    `json:"title"`
	LastMessageTime        time.Time `json:"lastMessageTime"`
	Generating             bool      `json:"generating"`
}

type ConversationSetting struct {
	ID                   int          `json:"id"`
	IsDefault            sql.NullBool `json:"isDefault"`
	SystemPromptTemplate string       `json:"systemPromptTemplate"`
	ToolsEnabled         StringArray  `json:"toolsEnabled"`
}

type ConversationTemplate struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	ConversationSettingsID int    `json:"conversationSettingsID"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Message struct {
	ID             int    `json:"id"`
	ConversationID int    `json:"conversationID"`
	Content        string `json:"content"`
	Author         string `json:"author"`
}
