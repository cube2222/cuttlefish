package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

// TODO: Some tools might like to keep state for the whole duration of the "Chain" of commands that GPT does.

type Tool interface {
	ArgumentDescriptions() map[string]string
	Run(ctx context.Context, args map[string]interface{}) (*RunResult, error)
}

// TODO: In the future we might want Tools to be able to dynamically append to the message.
//       So basically kind of a result output stream.

type RunResult struct {
	Result string
	Output string
}

type Terminal struct {
}

func (t *Terminal) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"command": "<bash command to run>",
	}
}

func (t *Terminal) Run(ctx context.Context, args map[string]interface{}) (*RunResult, error) {
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
	return &RunResult{
		Result: result,
		Output: buf.String() + "\n",
	}, nil
}
