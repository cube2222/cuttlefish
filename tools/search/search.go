package search

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/api/customsearch/v1"
	"google.golang.org/api/option"

	"cuttlefish/database"
	"cuttlefish/tools"
)

// TODO: Write how-to guide for getting the API key.
// https://stackoverflow.com/questions/37083058/programmatically-searching-google-in-python-using-custom-search
// https://support.google.com/googleapi/answer/6158862?hl=en

type Tool struct {
}

func (t *Tool) Name() string {
	return "Search"
}

func (t *Tool) Description() string {
	return "search using google; when the user asks about facts, or you need facts to answer, double-check them here instead of relying on your memory; you can also use it to find web pages relevant to a topic"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"query": "phrase to search for",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings, runtime tools.AppRuntime) (tools.ToolInstance, error) {
	return &ToolInstance{
		CustomSearchEngineID: settings.Search.GoogleCustomSearch.CustomSearchEngineID,
		GoogleCloudAPIKey:    settings.Search.GoogleCustomSearch.GoogleCloudAPIKey,
	}, nil
}

type ToolInstance struct {
	CustomSearchEngineID string
	GoogleCloudAPIKey    string
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	query, ok := args["query"].(string)
	if !ok {
		return nil, fmt.Errorf("search phrase is not a string")
	}
	searchService, err := customsearch.NewService(ctx, option.WithAPIKey(t.GoogleCloudAPIKey))
	if err != nil {
		return nil, fmt.Errorf("could not create search service: %w", err)
	}
	res, err := searchService.Cse.List().
		Cx(t.CustomSearchEngineID).
		Num(5).
		Q(query).Do()
	if err != nil {
		return nil, fmt.Errorf("could not search: %w", err)
	}

	type searchResult struct {
		URL     string `json:"url"`
		Title   string `json:"title"`
		Snippet string `json:"snippet"`
	}

	searchResults := make([]searchResult, len(res.Items))
	for i, r := range res.Items {
		searchResults[i] = searchResult{
			URL:     r.Link,
			Title:   r.Title,
			Snippet: r.Snippet,
		}
	}
	data, err := json.MarshalIndent(searchResults, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("could not encode results: %w", err)
	}

	return &tools.RunResult{
		Result: fmt.Sprintf("returning top %d results", len(searchResults)),
		Output: string(data) + "\n",
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}
