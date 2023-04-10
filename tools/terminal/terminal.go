package terminal

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"gptui/database"
	"gptui/tools"
)

type Tool struct {
}

func (t *Tool) Name() string {
	return "Terminal"
}

func (t *Tool) Description() string {
	return "run commands in a bash shell"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"command": "bash command to run",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings) (tools.ToolInstance, error) {
	return &ToolInstance{}, nil
}

type ToolInstance struct {
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	command, ok := args["command"].(string)
	if !ok {
		return nil, fmt.Errorf("command is not a string")
	}
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	var result string
	if err == context.Canceled {
		return nil, err
	} else if err != nil {
		result = err.Error()
	} else {
		result = "successfully executed `" + command + "`"
	}
	return &tools.RunResult{
		Result: result,
		Output: buf.String() + "\n",
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}
