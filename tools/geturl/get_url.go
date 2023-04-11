package geturl

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	gowiki "github.com/trietmn/go-wiki"

	"cuttlefish/database"
	"cuttlefish/tools"
)

type Tool struct {
}

func (t *Tool) Name() string {
	return "URL Getter"
}

func (t *Tool) Description() string {
	return "fetch the contents of a url; useful for getting details about a search result"
}

func (t *Tool) ArgumentDescriptions() map[string]string {
	return map[string]string{
		"url": "url to get",
	}
}

func (t *Tool) Instantiate(ctx context.Context, settings database.Settings) (tools.ToolInstance, error) {
	return &ToolInstance{}, nil
}

type ToolInstance struct {
}

func (t *ToolInstance) Run(ctx context.Context, args map[string]interface{}) (*tools.RunResult, error) {
	urlToFetch, ok := args["url"].(string)
	if !ok {
		return nil, fmt.Errorf("url is not a string")
	}
	content, err := fetchAndExtractMainContent(urlToFetch)
	if err != nil {
		return nil, fmt.Errorf("could not fetch url main contents: %w", err)
	}

	return &tools.RunResult{
		Result: "successfully fetched url",
		Output: strings.TrimSpace(content) + "\n",
	}, nil
}

func (t *ToolInstance) Shutdown() error {
	return nil
}

func fetchAndExtractMainContent(urlToGet string) (string, error) {
	urlParsed, err := url.Parse(urlToGet)
	if err != nil {
		return "", fmt.Errorf("could not parse URL: %w", err)
	}

	switch {
	case strings.HasSuffix(urlParsed.Hostname(), "wikipedia.org"):
		summary, err := gowiki.Summary(strings.TrimPrefix(urlParsed.Path, "/wiki/"), 5, -1, false, true)
		if err != nil {
			return "", fmt.Errorf("could not get Wikipedia summary: %w", err)
		}
		return summary, nil
	}

	resp, err := http.Get(urlToGet)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	switch {
	case strings.HasSuffix(urlParsed.Hostname(), "stackoverflow.com"):
		question := doc.Find(".question").Find(".js-post-body").Text()
		answer := doc.Find(".answer").First().Find(".js-post-body").Text()
		return fmt.Sprintf(`Question:
%s
---
Answer:
%s`, strings.TrimSpace(question), strings.TrimSpace(answer)), nil

	case doc.Find("article").Length() > 1:
		outParts := strings.Split(doc.Find("article").Text(), " ")
		if len(outParts) > 300 {
			outParts = outParts[:300]
			outParts = append(outParts, "...")
		}
		return strings.Join(outParts, " "), nil

	default:
		// Heuristics.
		var mainContent *goquery.Selection
		// Try to select the main content element using the article tag
		// Otherwise, try to select the main content element using the heuristic approach
		var maxTextLen int
		doc.Find("body *").Each(func(i int, s *goquery.Selection) {
			textLen := len(s.Text())
			if textLen > maxTextLen {
				maxTextLen = textLen
				mainContent = s
			}
			if class, exists := s.Attr("class"); exists && (strings.Contains(class, "main") || strings.Contains(class, "content")) {
				mainContent = s
			}
		})

		// If no main content element was found, return an error
		if mainContent == nil {
			return "", fmt.Errorf("could not find main content element")
		}

		outParts := strings.Split(mainContent.Text(), " ")
		if len(outParts) > 300 {
			outParts = outParts[:300]
			outParts = append(outParts, "...")
		}
		return strings.Join(outParts, " "), nil
	}
}
