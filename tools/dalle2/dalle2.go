package dalle2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Andrew-peng/go-dalle2/dalle2"

	"cuttlefish/database"
	"cuttlefish/tools"
)

type Tool struct {
}

func (t *Tool) Name() string {
	return "Image Generator"
}

func (t *Tool) Description() string {
	return "generate images using dalle2"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"prompt": "prompt to use to generate the image; the prompt should be detailed, and include keywords regarding styling; it shouldn't be a proper sentence, rather, a bag of keywords",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings, runtime tools.AppRuntime) (tools.ToolInstance, error) {
	dalleCli, err := dalle2.MakeNewClientV1(settings.OpenAIAPIKey)
	if err != nil {
		return nil, fmt.Errorf("couldn't create dalle2 client: %w", err)
	}
	return &ToolInstance{
		dalleCli: dalleCli,
	}, nil
}

type ToolInstance struct {
	dalleCli dalle2.Client
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	res, err := t.dalleCli.Create(
		ctx,
		args["prompt"].(string),
		dalle2.WithFormat(dalle2.URL),
		dalle2.WithSize(dalle2.LARGE),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating image: %s", err)
	}

	data, err := json.MarshalIndent(map[string]interface{}{
		"image_url": res.Data[0].Url,
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error encoding output json: %s", err)
	}

	return &tools.RunResult{
		Result: "successfully generated image",
		Output: string(data),
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}
