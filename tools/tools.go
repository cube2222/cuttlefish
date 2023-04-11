package tools

import (
	"context"

	"cuttlefish/database"
)

type AppRuntime interface {
	WaitForApproval(ctx context.Context, message string) error
}

type Tool interface {
	Name() string
	Description() string
	ArgumentDescriptions() map[string]string
	Instantiate(ctx context.Context, settings database.Settings, runtime AppRuntime) (ToolInstance, error)
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
