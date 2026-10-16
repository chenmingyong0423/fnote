package repository

import (
	"context"
	"fmt"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IAIAssistantConfigRepository interface {
	GetAIAssistantConfig(ctx context.Context) (*domain.AIAssistantConfig, error)
	UpsertAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error
}

func (r *WebsiteConfigRepository) GetAIAssistantConfig(ctx context.Context) (*domain.AIAssistantConfig, error) {
	stored, err := r.dao.FindByTyp(ctx, domain.AIAssistantConfigType)
	if err != nil {
		return nil, err
	}
	data, err := bson.Marshal(stored.Props)
	if err != nil {
		return nil, fmt.Errorf("encode AI assistant config: %w", err)
	}
	var cfg domain.AIAssistantConfig
	if err = bson.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode AI assistant config: %w", err)
	}
	return &cfg, nil
}

func (r *WebsiteConfigRepository) UpsertAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error {
	return r.dao.UpsertAIAssistantConfig(ctx, cfg)
}
