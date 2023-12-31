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

	"github.com/chenmingyong0423/fnote/backend/internal/category/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/count_stats/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategories(ctx context.Context) ([]domain.CategoryWithCount, error)
	GetMenus(ctx context.Context) ([]domain.Category, error)
	GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error)
}

var _ ICategoryService = (*CategoryService)(nil)

func NewCategoryService(repo repository.ICategoryRepository, countStatsService service.ICountStatsService) *CategoryService {
	return &CategoryService{
		countStatsService: countStatsService,
		repo:              repo,
	}
}

type CategoryService struct {
	countStatsService service.ICountStatsService
	repo              repository.ICategoryRepository
}

func (s *CategoryService) GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error) {
	return s.repo.GetCategoryByRoute(ctx, route)
}

func (s *CategoryService) GetMenus(ctx context.Context) (menuVO []domain.Category, err error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]domain.CategoryWithCount, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, nil
	}
	// 将所有分类的 id 转换为 string 数组
	ids := slice.Map[domain.Category, string](categories, func(_ int, s domain.Category) string {
		return s.Id
	})
	categoryMap := slice.IndexStructsByKey[domain.Category, string](categories, func(category domain.Category) string {
		return category.Id
	})

	categoryCounts, err := s.countStatsService.GetByReferenceIdAndType(ctx, ids, domain.CountStatsTypePostCountInCategory)
	if err != nil {
		return nil, err
	}

	categoryWithCounts := slice.Map[domain.CountStats, domain.CategoryWithCount](categoryCounts, func(idx int, s domain.CountStats) domain.CategoryWithCount {
		return domain.CategoryWithCount{
			Name:        categoryMap[s.ReferenceId].Name,
			Route:       categoryMap[s.ReferenceId].Route,
			Description: categoryMap[s.ReferenceId].Description,
			Count:       s.Count,
		}
	})

	return categoryWithCounts, nil
}
