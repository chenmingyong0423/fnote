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
	"fmt"
	"strings"
	"time"

	"github.com/chenmingyong0423/go-mongox"

	"github.com/chenmingyong0423/fnote/server/internal/category/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/repository/dao"
	"github.com/chenmingyong0423/gkit/slice"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICategoryRepository interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error)
	QueryCategoriesPage(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Category, int64, error)
	CreateCategory(ctx context.Context, category domain.Category) (string, error)
	ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error
	ModifyCategory(ctx context.Context, id string, description string) error
	DeleteCategory(ctx context.Context, id string) error
	GetNavigations(ctx context.Context) ([]domain.Category, error)
	ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error
	GetCategoryById(ctx context.Context, id string) (domain.Category, error)
	RecoverCategory(ctx context.Context, category domain.Category) error
	GetSelectCategories(ctx context.Context) ([]domain.Category, error)
	IncreasePostCountByIds(ctx context.Context, categoryIds []string) error
	DecreasePostCountByIds(ctx context.Context, categoryIds []string) error
	FindEnabledCategories(ctx context.Context) ([]domain.Category, error)
}

var _ ICategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(dao dao.ICategoryDao) *CategoryRepository {
	return &CategoryRepository{
		dao: dao,
	}
}

type CategoryRepository struct {
	dao dao.ICategoryDao
}

func (r *CategoryRepository) FindEnabledCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.FindEnabledCategories(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainCategories(categories), nil
}

func (r *CategoryRepository) DecreasePostCountByIds(ctx context.Context, categoryIds []string) (err error) {
	categoryObjectIds := slice.Map(categoryIds, func(_ int, id string) (ojbId primitive.ObjectID) {
		if err != nil {
			return ojbId
		}
		ojbId, err = primitive.ObjectIDFromHex(id)
		return ojbId
	})
	if err != nil {
		return
	}
	return r.dao.DecreasePostCountByIds(ctx, categoryObjectIds)
}

func (r *CategoryRepository) IncreasePostCountByIds(ctx context.Context, categoryIds []string) (err error) {
	categoryObjectIds := slice.Map(categoryIds, func(_ int, id string) (ojbId primitive.ObjectID) {
		if err != nil {
			return ojbId
		}
		ojbId, err = primitive.ObjectIDFromHex(id)
		return ojbId
	})
	if err != nil {
		return
	}
	return r.dao.IncreasePostCountByIds(ctx, categoryObjectIds)
}

func (r *CategoryRepository) GetSelectCategories(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.GetEnabled(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainCategories(categories), nil
}

func (r *CategoryRepository) RecoverCategory(ctx context.Context, category domain.Category) error {
	return r.dao.RecoverCategory(ctx, &dao.Category{
		Name:        category.Name,
		Route:       category.Route,
		Description: category.Description,
		Sort:        category.Sort,
		Enabled:     category.Enabled,
		ShowInNav:   category.ShowInNav,
		Model: mongox.Model{
			CreatedAt: time.Unix(category.CreatedAt, 0),
			UpdatedAt: time.Unix(category.UpdatedAt, 0),
		},
	})
}

func (r *CategoryRepository) GetCategoryById(ctx context.Context, id string) (t domain.Category, err error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	tag, err := r.dao.GetById(ctx, objId)
	if err != nil {
		return
	}
	return r.toDomainCategory(tag), nil
}

func (r *CategoryRepository) ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyCategoryNavigation(ctx, objectID, showInNav)
}

func (r *CategoryRepository) GetNavigations(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.GetByShowInNav(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainCategories(categories), nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *CategoryRepository) ModifyCategory(ctx context.Context, id string, description string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyCategory(ctx, objectID, description)
}

func (r *CategoryRepository) ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyEnabled(ctx, objectID, enabled)
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category domain.Category) (string, error) {
	return r.dao.Create(ctx, &dao.Category{Name: category.Name, Route: category.Route, Description: category.Description, ShowInNav: category.ShowInNav, Enabled: category.Enabled})
}

func (r *CategoryRepository) QueryCategoriesPage(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Category, int64, error) {
	condBuilder := query.BsonBuilder()
	if pageDTO.Keyword != "" {
		condBuilder.RegexOptions("name", fmt.Sprintf(".*%s.*", strings.TrimSpace(pageDTO.Keyword)), "i")
	}
	cond := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip((pageDTO.PageNo - 1) * pageDTO.PageSize).SetLimit(pageDTO.PageSize)
	if pageDTO.Field != "" && pageDTO.Order != "" {
		findOptions.SetSort(bsonx.M(pageDTO.Field, pageDTO.OrderConvertToInt()))
	} else {
		findOptions.SetSort(bsonx.M("created_at", -1))
	}
	categories, total, err := r.dao.QuerySkipAndSetLimit(ctx, cond, findOptions)
	return r.toDomainCategories(categories), total, err
}

func (r *CategoryRepository) toDomainCategories(categories []*dao.Category) []domain.Category {
	result := make([]domain.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, r.toDomainCategory(category))
	}
	return result
}

func (r *CategoryRepository) GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error) {
	category, err := r.dao.GetByRoute(ctx, route)
	if err != nil {
		return domain.Category{}, err
	}
	return r.toDomainCategory(category), nil
}

func (r *CategoryRepository) toDomainCategory(category *dao.Category) domain.Category {
	return domain.Category{Id: category.ID.Hex(), Name: category.Name, Route: category.Route, Description: category.Description, Enabled: category.Enabled, ShowInNav: category.ShowInNav, PostCount: category.PostCount, Sort: category.Sort, CreatedAt: category.CreatedAt.Unix(), UpdatedAt: category.UpdatedAt.Unix()}
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, r.toDomainCategory(category))
	}
	return result, nil
}
