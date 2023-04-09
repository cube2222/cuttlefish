package chart

import (
	"context"
	"encoding/json"
	"fmt"

	"gptui/database"
	"gptui/tools"
)

type Tool struct {
}

func (t *Tool) Description() string {
	return "plot data on charts using Apache ECharts"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"options": "Options to pass to the Apache ECharts chart.",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings) (tools.ToolInstance, error) {
	return &ToolInstance{}, nil
}

type ToolInstance struct {
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	options := args["options"]
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not encode options: %w", err)
	}
	return &tools.RunResult{
		Result:          "here's the chart; this will already be visible to the user",
		CustomResultTag: "chart",
		Output:          string(data) + "\n",
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}
