package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"github.com/pkg/errors"
)

type IConfigCheckStateDao interface {
	GetConfigCheckStates(ctx context.Context) ([]*ConfigCheckState, error)
	UpsertConfigCheckState(ctx context.Context, state domain.ConfigCheckState) error
}

type ConfigCheckState struct {
	mongox.Model  `bson:",inline"`
	CheckKey      string     `bson:"check_key"`
	Status        string     `bson:"status"`
	IgnoredReason string     `bson:"ignored_reason,omitempty"`
	SnoozedUntil  *time.Time `bson:"snoozed_until,omitempty"`
}

func (d *WebsiteConfigDao) GetConfigCheckStates(ctx context.Context) ([]*ConfigCheckState, error) {
	states, err := d.configCheckStateColl.Finder().Find(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fails to find config check states")
	}
	return states, nil
}

func (d *WebsiteConfigDao) UpsertConfigCheckState(ctx context.Context, state domain.ConfigCheckState) error {
	now := time.Now().Local()
	updateResult, err := d.configCheckStateColl.Updater().
		Filter(query.Eq("check_key", state.CheckKey)).
		Updates(update.NewBuilder().
			Set("check_key", state.CheckKey).
			Set("status", state.Status).
			Set("ignored_reason", state.IgnoredReason).
			Set("snoozed_until", state.SnoozedUntil).
			Set("updated_at", now).
			SetOnInsert("created_at", now).
			Build()).
		Upsert(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to upsert config check state, state=%v", state)
	}
	if updateResult.ModifiedCount == 0 && updateResult.UpsertedCount == 0 {
		return fmt.Errorf("ModifiedCount=0 && UpsertedCount=0, fails to upsert config check state, state=%v", state)
	}
	return nil
}
