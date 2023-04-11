package python

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"

	"cuttlefish/database"
	"cuttlefish/tools"
)

type Tool struct {
}

func (t *Tool) Name() string {
	return "Python"
}

func (t *Tool) Description() string {
	return "run python3 scripts"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"code": "python code to run; keep in mind that this is a json string, so you'll need to escape newlines and relevant special characters",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings, runtime tools.AppRuntime) (tools.ToolInstance, error) {
	return &ToolInstance{
		pythonInterpreterPath: settings.Python.InterpreterPath,
	}, nil
}

type ToolInstance struct {
	pythonInterpreterPath string
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	code, ok := args["code"].(string)
	if !ok {
		return nil, fmt.Errorf("command is not a string")
	}
	cmd := exec.CommandContext(ctx, t.pythonInterpreterPath, "-c", code)
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
		result = "successfully executed `" + code + "`"
	}
	return &tools.RunResult{
		Result: result,
		Output: buf.String(),
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}
