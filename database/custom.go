package database

type Settings struct {
	OpenAIAPIKey string         `json:"openAiApiKey"`
	Model        string         `json:"model"`
	Search       SearchSettings `json:"search"`
	Python       PythonSettings `json:"python"`
}

type SearchSettings struct {
	GoogleCustomSearch GoogleCustomSearchSettings `json:"googleCustomSearch"`
}

type GoogleCustomSearchSettings struct {
	CustomSearchEngineID string `json:"customSearchEngineId"`
	GoogleCloudAPIKey    string `json:"googleCloudApiKey"`
}

type PythonSettings struct {
	InterpreterPath string `json:"interpreterPath"`
}
