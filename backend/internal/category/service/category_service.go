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

	"github.com/chenmingyong0423/fnote/backend/internal/count_stats/service"
	"github.com/chenmingyong0423/gkit/slice"
	"golang.org/x/sync/errgroup"

	"github.com/chenmingyong0423/fnote/backend/internal/category/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategoriesAndTags(ctx context.Context) (domain.CategoryAndTagWithCount, error)
	GetMenus(ctx context.Context) ([]domain.Category, error)
	GetTagsByName(ctx context.Context, name string) ([]string, error)
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

func (s *CategoryService) GetTagsByName(ctx context.Context, name string) ([]string, error) {
	return s.repo.GetTagsByName(ctx, name)
}

func (s *CategoryService) GetMenus(ctx context.Context) (menuVO []domain.Category, err error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategoriesAndTags(ctx context.Context) (domain.CategoryAndTagWithCount, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return domain.CategoryAndTagWithCount{}, err
	}
	if len(categories) == 0 {
		return domain.CategoryAndTagWithCount{}, nil
	}
	// 将所有分类的 id 转换为 string 数组
	ids := slice.Map[domain.Category, string](categories, func(_ int, s domain.Category) string {
		return s.Id
	})
	categoryMap := slice.IndexStructsByKey[domain.Category, string](categories, func(category domain.Category) string {
		return category.Id
	})
	// 获取所有分类下的标签
	tags := slice.CombineAndDeduplicateNestedSlices[domain.Category, string](categories, func(_ int, s domain.Category) []string {
		return s.Tags
	})

	categoryCounts := make([]domain.CountStats, 0, len(ids))
	tagCounts := make([]domain.CountStats, 0, len(tags))
	group := &errgroup.Group{}
	group.Go(func() error {
		categoryCounts, err = s.countStatsService.GetByReferenceIdAndType(ctx, ids, domain.CountStatsTypePostCountInCategory)
		if err != nil {
			return err
		}
		return nil
	})
	group.Go(func() error {
		if len(tags) > 0 {
			tagCounts, err = s.countStatsService.GetByReferenceIdAndType(ctx, tags, domain.CountStatsTypePostCountInTag)
			if err != nil {
				return err
			}
		}
		return nil
	})
	err = group.Wait()
	if err != nil {
		return domain.CategoryAndTagWithCount{}, err
	}

	tagWithCounts := slice.Map[domain.CountStats, domain.TagWithCount](tagCounts, func(idx int, s domain.CountStats) domain.TagWithCount {
		return domain.TagWithCount{
			Name:  s.ReferenceId,
			Count: s.Count,
		}
	})

	categoryWithCounts := slice.Map[domain.CountStats, domain.CategoryWithCount](categoryCounts, func(idx int, s domain.CountStats) domain.CategoryWithCount {
		return domain.CategoryWithCount{
			Name:        categoryMap[s.ReferenceId].Name,
			Route:       categoryMap[s.ReferenceId].Route,
			Description: categoryMap[s.ReferenceId].Description,
			Count:       s.Count,
		}
	})

	return domain.CategoryAndTagWithCount{
		Categories: categoryWithCounts,
		Tags:       tagWithCounts,
	}, nil
}
