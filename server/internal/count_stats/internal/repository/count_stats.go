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

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/repository/dao"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
)

type ICountStatsRepository interface {
	DecreaseByReferenceIdAndType(ctx context.Context, countStatsType domain.CountStatsType, count int) error
	IncreaseByReferenceIdAndType(ctx context.Context, countStatsType domain.CountStatsType, delta int) error
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
		ds = append(ds, bsonx.NewD().Add("type", statsType.ToString()).Add("type", statsType.ToString()).Build())
	}
	countStats, err := r.dao.GetByFilter(ctx, query.Or(ds...))
	if err != nil {
		return nil, err
	}
	return r.toDomainCountStats(countStats), nil
}

func (r *CountStatsRepository) IncreaseByReferenceIdAndType(ctx context.Context, countStatsType domain.CountStatsType, delta int) error {
	return r.dao.IncreaseByReferenceIdAndType(ctx, countStatsType.ToString(), delta)
}

func (r *CountStatsRepository) DecreaseByReferenceIdAndType(ctx context.Context, countStatsType domain.CountStatsType, count int) error {
	return r.dao.DecreaseByReferenceIdAndType(ctx, countStatsType.ToString(), count)
}

func (r *CountStatsRepository) toDomainCountStats(stats []*dao.CountStats) []domain.CountStats {
	var countStats []domain.CountStats
	for _, stat := range stats {
		countStats = append(countStats, domain.CountStats{
			Id:    stat.ID.Hex(),
			Type:  domain.CountStatsType(stat.Type),
			Count: stat.Count,
		})
	}
	return countStats
}
