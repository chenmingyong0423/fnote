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

	"github.com/chenmingyong0423/fnote/backend/internal/count_stats/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
)

type ICountStatsRepository interface {
	GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error)
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

func (r *CountStatsRepository) GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error) {
	countStats, err := r.dao.GetByReferenceIdAndType(ctx, referenceIds, string(countStatsType))
	if err != nil {
		return nil, err
	}

	return r.toDomainCountStats(countStats), nil
}

func (r *CountStatsRepository) toDomainCountStats(stats []*dao.CountStats) []domain.CountStats {
	var countStats []domain.CountStats
	for _, stat := range stats {
		countStats = append(countStats, domain.CountStats{
			Id:          stat.Id,
			Type:        stat.Type,
			ReferenceId: stat.ReferenceId,
			Count:       stat.Count,
		})
	}
	return countStats
}
