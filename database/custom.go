package database

type Settings struct {
	OpenAIAPIKey string `json:"openAiApiKey"`
	Model        string `json:"model"`
	// Add nested struct per configurable plugin below.
}
