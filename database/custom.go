package database

type Settings struct {
	OpenAIAPIKey string         `json:"openAiApiKey"`
	Model        string         `json:"model"`
	Search       SearchSettings `json:"search"`
}

type SearchSettings struct {
	GoogleCustomSearch GoogleCustomSearchSettings `json:"googleCustomSearch"`
}

type GoogleCustomSearchSettings struct {
	CustomSearchEngineID string `json:"customSearchEngineId"`
	GoogleCloudAPIKey    string `json:"googleCloudApiKey"`
}
