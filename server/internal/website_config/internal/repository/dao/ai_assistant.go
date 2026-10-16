package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
)

type IAIAssistantConfigDao interface {
	UpsertAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error
}

func (d *WebsiteConfigDao) UpsertAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error {
	now := time.Now().Local()
	_, err := d.coll.Updater().
		Filter(query.Eq("typ", domain.AIAssistantConfigType)).
		Updates(update.NewBuilder().
			Set("props", cfg).
			Set("updated_at", now).
			SetOnInsert("created_at", now).
			Build()).Upsert(ctx)
	if err != nil {
		// Do not include configuration values or credentials in errors.
		return fmt.Errorf("save AI assistant config: %w", err)
	}
	return nil
}
