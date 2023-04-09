package database

type Settings struct {
	OpenAIAPIKey string `json:"openAiApiKey"`
	Model        string `json:"model"`
	Search       struct {
		GoogleCustomSearch struct {
			CustomSearchEngineID string `json:"customSearchEngineId"`
			GoogleCloudAPIKey    string `json:"googleCloudApiKey"`
		}
	}
	// Add nested struct per configurable plugin below.
}
