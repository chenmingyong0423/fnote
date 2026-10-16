package domain

import "errors"

const AIAssistantConfigType = "ai_assistant"

var ErrInvalidAIConfig = errors.New("AI 助手配置不合法")

type AIAssistantConfig struct {
	Enabled          bool         `bson:"enabled"`
	AssistantName    string       `bson:"assistant_name"`
	SystemPrompt     string       `bson:"system_prompt"`
	EnabledTools     []string     `bson:"enabled_tools"`
	ActiveProviderId string       `bson:"active_provider_id"`
	Providers        []AIProvider `bson:"providers"`
}

type AIProvider struct {
	Id              string `bson:"id"`
	Name            string `bson:"name"`
	ApiType         string `bson:"api_type"`
	BaseUrl         string `bson:"base_url"`
	Model           string `bson:"model"`
	MaxOutputTokens *int   `bson:"max_output_tokens,omitempty"`
	ReasoningEffort string `bson:"reasoning_effort,omitempty"`
	ThinkingMode    string `bson:"thinking_mode,omitempty"`
	EncryptedApiKey string `bson:"encrypted_api_key" json:"-"`
	// ApiKey is accepted on updates only and must never be persisted or serialized.
	ApiKey string `bson:"-" json:"-"`
}

type AIProviderTestResult struct {
	Success    bool
	Message    string
	Reply      string
	DurationMs int64
}
