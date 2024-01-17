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
	GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error)
	Create(ctx context.Context, countStats domain.CountStats) error
	DeleteByReferenceId(ctx context.Context, referenceId string) error
	IncreaseByReferenceIds(ctx context.Context, ids []string) error
	DecreaseByReferenceIds(ctx context.Context, ids []string) error
	DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
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

func (s *CountStatsService) IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.IncreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.DecreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIds(ctx context.Context, ids []string) error {
	return s.repo.DecreaseByReferenceIds(ctx, ids)
}

func (s *CountStatsService) IncreaseByReferenceIds(ctx context.Context, ids []string) error {
	return s.repo.IncreaseByReferenceIds(ctx, ids)
}

func (s *CountStatsService) DeleteByReferenceId(ctx context.Context, referenceId string) error {
	return s.repo.DeleteByReferenceId(ctx, referenceId)
}

func (s *CountStatsService) Create(ctx context.Context, countStats domain.CountStats) error {
	_, err := s.repo.Create(ctx, countStats)
	if err != nil {
		return err
	}
	return nil
}

func (s *CountStatsService) GetByReferenceIdAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error) {
	return s.repo.GetByReferenceIdAndType(ctx, referenceIds, countStatsType)
}
