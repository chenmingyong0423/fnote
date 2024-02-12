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

package service

import (
	"context"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/repository"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
)

type ICountStatsService interface {
	GetByReferenceIdsAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error)
	Create(ctx context.Context, countStats domain.CountStats) error
	DeleteByReferenceIdAndType(ctx context.Context, referenceId string, statsType domain.CountStatsType) error
	DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error
	IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error
	GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error)
}

var _ ICountStatsService = (*CountStatsService)(nil)

func NewCountStatsService(repo repository.ICountStatsRepository) *CountStatsService {
	return &CountStatsService{
		repo: repo,
	}
}

type CountStatsService struct {
	repo repository.ICountStatsRepository
}

func (s *CountStatsService) GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error) {
	var result = new(domain.WebsiteCountStats)
	countStatsSlice, err := s.repo.GetWebsiteCountStats(ctx, []domain.CountStatsType{
		domain.CountStatsTypePostCountInWebsite,
		domain.CountStatsTypeCategoryCount,
		domain.CountStatsTypeTagCount,
		domain.CountStatsTypeCommentCount,
		domain.CountStatsTypeLikeCount,
		domain.CountStatsTypeWebsiteViewCount,
	})
	if err != nil {
		return *result, err
	}
	for _, countStats := range countStatsSlice {
		result.SetCountByType(countStats.Type, countStats.Count)
	}
	return *result, nil
}

func (s *CountStatsService) IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error {
	return s.repo.IncreaseByReferenceIdAndType(ctx, referenceId, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error {
	return s.repo.DecreaseByReferenceIdAndType(ctx, referenceId, countStatsType)
}

func (s *CountStatsService) IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.IncreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.DecreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DeleteByReferenceIdAndType(ctx context.Context, referenceId string, statsType domain.CountStatsType) error {
	return s.repo.DeleteByReferenceIdAndType(ctx, referenceId, statsType)
}

func (s *CountStatsService) Create(ctx context.Context, countStats domain.CountStats) error {
	err := countStats.Type.Valid()
	if err != nil {
		return err
	}
	_, err = s.repo.Create(ctx, countStats)
	if err != nil {
		return err
	}
	return nil
}

func (s *CountStatsService) GetByReferenceIdsAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error) {
	return s.repo.GetByReferenceIdAndType(ctx, referenceIds, countStatsType)
}
