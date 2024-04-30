// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repository

import (
	"context"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
)

type ICountStatsRepository interface {
	GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error)
	Create(ctx context.Context, countStats domain.CountStats) (string, error)
	DeleteByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error
	DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType, count int) error
	IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error
	GetWebsiteCountStats(ctx context.Context, countStatsTypes []domain.CountStatsType) ([]domain.CountStats, error)
}

var _ ICountStatsRepository = (*CountStatsRepository)(nil)

func NewCountStatsRepository(dao dao.ICountStatsDao) *CountStatsRepository {
	return &CountStatsRepository{
		dao: dao,
	}
}

type CountStatsRepository struct {
	dao dao.ICountStatsDao
}

func (r *CountStatsRepository) GetWebsiteCountStats(ctx context.Context, countStatsTypes []domain.CountStatsType) ([]domain.CountStats, error) {
	ds := make([]any, 0, len(countStatsTypes))
	for _, statsType := range countStatsTypes {
		ds = append(ds, bsonx.NewD().Add("reference_id", statsType.ToString()).Add("type", statsType.ToString()).Build())
	}
	countStats, err := r.dao.GetByFilter(ctx, query.Or(ds...))
	if err != nil {
		return nil, err
	}
	return r.toDomainCountStats(countStats), nil
}

func (r *CountStatsRepository) IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error {
	return r.dao.IncreaseByReferenceIdAndType(ctx, referenceId, countStatsType.ToString())
}

func (r *CountStatsRepository) DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType, count int) error {
	return r.dao.DecreaseByReferenceIdAndType(ctx, referenceId, countStatsType.ToString(), count)
}

func (r *CountStatsRepository) IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return r.dao.IncreaseByReferenceIdsAndType(ctx, ids, countStatsType.ToString())
}

func (r *CountStatsRepository) DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return r.dao.DecreaseByReferenceIdsAndType(ctx, ids, countStatsType.ToString())
}

func (r *CountStatsRepository) DeleteByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error {
	return r.dao.DeleteByReferenceIdAndType(ctx, referenceId, countStatsType.ToString())
}

func (r *CountStatsRepository) Create(ctx context.Context, countStats domain.CountStats) (string, error) {
	unix := time.Now().Unix()
	return r.dao.Create(ctx, &dao.CountStats{
		Type:        countStats.Type.ToString(),
		ReferenceId: countStats.ReferenceId,
		CreateTime:  unix,
		UpdateTime:  unix,
	})
}

func (r *CountStatsRepository) GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error) {
	countStats, err := r.dao.GetByReferenceIdAndType(ctx, referenceIds, countStatsType.ToString())
	if err != nil {
		return nil, err
	}

	return r.toDomainCountStats(countStats), nil
}

func (r *CountStatsRepository) toDomainCountStats(stats []*dao.CountStats) []domain.CountStats {
	var countStats []domain.CountStats
	for _, stat := range stats {
		countStats = append(countStats, domain.CountStats{
			Id:          stat.Id.Hex(),
			Type:        domain.CountStatsType(stat.Type),
			ReferenceId: stat.ReferenceId,
			Count:       stat.Count,
		})
	}
	return countStats
}
