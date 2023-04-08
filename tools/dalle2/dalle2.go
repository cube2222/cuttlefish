package dalle2

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Andrew-peng/go-dalle2/dalle2"

	"gptui/tools"
)

type Tool struct {
	OpenAIToken string
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"prompt": "<prompt to use to generate the image; the prompt should be detailed, and include keywords regarding styling; it shouldn't be a proper sentence, rather, a bag of keywords>",
	}
}

func (t *Tool) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	cli, err := dalle2.MakeNewClientV1(t.OpenAIToken)
	if err != nil {
		return nil, fmt.Errorf("error creating dalle2 client: %s", err)
	}

	res, err := cli.Create(
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
