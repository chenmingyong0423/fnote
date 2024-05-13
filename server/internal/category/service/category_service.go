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
	"fmt"
	"log/slog"
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/service"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/category/repository"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategories(ctx context.Context) ([]domain.CategoryWithCount, error)
	GetMenus(ctx context.Context) ([]domain.Category, error)
	GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error)
	AdminGetCategories(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Category, int64, error)
	AdminCreateCategory(ctx context.Context, category domain.Category) error
	ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error
	ModifyCategory(ctx context.Context, id string, description string) error
	DeleteCategory(ctx context.Context, id string) error
	ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error
	AdminGetSelectCategories(ctx context.Context) ([]domain.Category, error)
}

var _ ICategoryService = (*CategoryService)(nil)

func NewCategoryService(repo repository.ICategoryRepository, countStatsService count_stats.Service, configService website_config.Service) *CategoryService {
	return &CategoryService{
		countStatsService: countStatsService,
		configService:     configService,
		repo:              repo,
	}
}

type CategoryService struct {
	configService     website_config.Service
	countStatsService service.ICountStatsService
	repo              repository.ICategoryRepository
}

func (s *CategoryService) AdminGetSelectCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetSelectCategories(ctx)
}

func (s *CategoryService) ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error {
	return s.repo.ModifyCategoryNavigation(ctx, id, showInNav)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	category, err := s.repo.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "category not found")
		}
		return err
	}
	err = s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	// 删除分类时，同时删除该分类下的文章数量统计数据
	err = s.countStatsService.DeleteByReferenceIdAndType(ctx, id, domain.CountStatsTypePostCountInCategory)
	if err != nil {
		nErr := s.repo.RecoverCategory(ctx, category)
		if nErr != nil {
			return nErr
		}
		return err
	}
	go func() {
		// 更新分类数量
		gErr := s.countStatsService.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCategoryCount.ToString(), domain.CountStatsTypeCategoryCount, 1)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil
}

func (s *CategoryService) ModifyCategory(ctx context.Context, id string, description string) error {
	return s.repo.ModifyCategory(ctx, id, description)
}

func (s *CategoryService) ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error {
	return s.repo.ModifyCategoryEnabled(ctx, id, enabled)
}

func (s *CategoryService) AdminCreateCategory(ctx context.Context, category domain.Category) error {
	id, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return err
	}
	// 创建分类时，同时创建分类的统计数据
	err = s.countStatsService.Create(ctx, domain.CountStats{
		Type:        domain.CountStatsTypePostCountInCategory,
		ReferenceId: id,
	})
	if err != nil {
		gErr := s.DeleteCategory(ctx, id)
		if gErr != nil {
			return gErr
		}
		return err
	}
	go func() {
		// 更新分类数量
		gErr := s.countStatsService.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCategoryCount.ToString(), domain.CountStatsTypeCategoryCount)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil
}

func (s *CategoryService) AdminGetCategories(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Category, int64, error) {
	categories, total, err := s.QueryCategoriesPage(ctx, pageDTO)
	if err != nil {
		return nil, 0, err
	}
	if len(categories) > 0 {
		// 获取分类的统计数据
		ids := slice.Map[domain.Category, string](categories, func(_ int, s domain.Category) string {
			return s.Id
		})
		countStats, fErr := s.countStatsService.GetByReferenceIdsAndType(ctx, ids, domain.CountStatsTypePostCountInCategory)
		if err != nil {
			return nil, 0, fErr
		}
		countStatsMap := slice.IndexStructsByKey[domain.CountStats, string](countStats, func(countStats domain.CountStats) string {
			return countStats.ReferenceId
		})
		for i, category := range categories {
			if cs, ok := countStatsMap[category.Id]; ok {
				categories[i].PostCount = cs.Count
			}
		}
	}
	return categories, total, nil
}

func (s *CategoryService) GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error) {
	return s.repo.GetCategoryByRoute(ctx, route)
}

func (s *CategoryService) GetMenus(ctx context.Context) (menuVO []domain.Category, err error) {
	categories, err := s.repo.GetNavigations(ctx)
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

	categoryCounts, err := s.countStatsService.GetByReferenceIdsAndType(ctx, ids, domain.CountStatsTypePostCountInCategory)
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

func (s *CategoryService) QueryCategoriesPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Category, int64, error) {
	return s.repo.QueryCategoriesPage(ctx, pageDTO)
}
