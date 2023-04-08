package terminal

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"gptui/tools"
)

type Tool struct {
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"command": "<bash command to run>",
	}
}

func (t *Tool) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
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
