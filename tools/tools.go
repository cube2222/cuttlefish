package tools

import (
	"context"

	"gptui/database"
)

// TODO: Some tools might like to keep state for the whole duration of the "Chain" of commands that GPT does. i.e. a Python REPL.

type Tool interface {
	Name() string
	Description() string
	ArgumentDescriptions() map[string]string
	Instantiate(ctx context.Context, settings database.Settings) (ToolInstance, error)
}

type ToolInstance interface {
	Run(ctx context.Context, args map[string]interface{}) (*RunResult, error)
	Shutdown() error
}

// TODO: In the future we might want Tools to be able to dynamically append to the message.
//       So basically kind of a result output stream.

type RunResult struct {
	Result          string
	CustomResultTag string
	Output          string
}
