package repository

import (
	"context"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
)

type IConfigCheckStateRepository interface {
	GetConfigCheckStates(ctx context.Context) ([]domain.ConfigCheckState, error)
	UpsertConfigCheckState(ctx context.Context, state domain.ConfigCheckState) error
}

func (r *WebsiteConfigRepository) GetConfigCheckStates(ctx context.Context) ([]domain.ConfigCheckState, error) {
	states, err := r.dao.GetConfigCheckStates(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ConfigCheckState, 0, len(states))
	for _, state := range states {
		result = append(result, domain.ConfigCheckState{
			CheckKey:      state.CheckKey,
			Status:        state.Status,
			IgnoredReason: state.IgnoredReason,
			SnoozedUntil:  state.SnoozedUntil,
			UpdatedAt:     state.UpdatedAt,
		})
	}
	return result, nil
}

func (r *WebsiteConfigRepository) UpsertConfigCheckState(ctx context.Context, state domain.ConfigCheckState) error {
	return r.dao.UpsertConfigCheckState(ctx, state)
}
