// Package aimodel constructs Eino models from provider connection settings.
// It has no dependencies on HTTP handlers, configuration persistence or visitors.
package aimodel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/eino-ext/components/model/agenticopenai"
	"github.com/cloudwego/eino/components/model"
)

type Config struct {
	APIType         string
	BaseURL         string
	APIKey          string
	Model           string
	MaxOutputTokens *int
	ReasoningEffort string
	ThinkingMode    string
}

// New uses Eino's Responses and OpenAI-compatible Chat Completions components.
// DeepSeek's additional parameters are supplied through Eino's ExtraFields.
func New(ctx context.Context, cfg Config, client *http.Client) (model.AgenticModel, error) {
	switch cfg.APIType {
	case "openai_responses":
		store, retries := false, 0
		config := &agenticopenai.ResponsesConfig{
			BaseURL: cfg.BaseURL, APIKey: cfg.APIKey, Model: cfg.Model,
			HTTPClient: client, MaxTokens: cfg.MaxOutputTokens, Store: &store, MaxRetries: &retries,
		}
		if cfg.ReasoningEffort != "" {
			config.ExtraFields = map[string]any{"reasoning": map[string]string{"effort": cfg.ReasoningEffort}}
		}
		return agenticopenai.NewResponsesModel(ctx, config)
	case "openai_chat_completions":
		extra := map[string]any{}
		if cfg.MaxOutputTokens != nil {
			extra["max_tokens"] = *cfg.MaxOutputTokens
		}
		if cfg.ThinkingMode != "" {
			extra["thinking"] = map[string]string{"type": cfg.ThinkingMode}
		}
		if cfg.ReasoningEffort != "" {
			extra["reasoning_effort"] = cfg.ReasoningEffort
		}
		return agenticopenai.NewChatModel(ctx, &agenticopenai.ChatConfig{
			BaseURL: cfg.BaseURL, APIKey: cfg.APIKey, Model: cfg.Model,
			HTTPClient: client, ExtraFields: extra,
		})
	default:
		return nil, fmt.Errorf("unsupported AI API type")
	}
}
