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

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/web/dto"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chenmingyong0423/fnote/backend/internal/category/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
)

type ICategoryRepository interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error)
	QueryCategoriesPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Category, int64, error)
	CreateCategory(ctx context.Context, category domain.Category) (string, error)
	ModifyCategoryDisabled(ctx context.Context, id string, disabled bool) error
	ModifyCategory(ctx context.Context, id string, description string) error
	DeleteCategory(ctx context.Context, id string) error
	GetNavigations(ctx context.Context) ([]domain.Category, error)
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

func (r *CategoryRepository) ModifyCategoryDisabled(ctx context.Context, id string, disabled bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyDisabled(ctx, objectID, disabled)
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category domain.Category) (string, error) {
	now := time.Now().Unix()
	return r.dao.Create(ctx, &dao.Category{Name: category.Name, Route: category.Route, Description: category.Description, CreateTime: now, UpdateTime: now})
}

func (r *CategoryRepository) QueryCategoriesPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Category, int64, error) {
	condBuilder := query.BsonBuilder()
	if pageDTO.Keyword != "" {
		condBuilder.RegexOptions("title", fmt.Sprintf(".*%s.*", strings.TrimSpace(pageDTO.Keyword)), "i")
	}
	cond := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip((pageDTO.PageNo - 1) * pageDTO.PageSize).SetLimit(pageDTO.PageSize)
	if pageDTO.Field != "" && pageDTO.Order != "" {
		findOptions.SetSort(bsonx.M(pageDTO.Field, pageDTO.OrderConvertToInt()))
	} else {
		findOptions.SetSort(bsonx.M("create_time", -1))
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
	return domain.Category{Id: category.Id.Hex(), Name: category.Name, Route: category.Route, Description: category.Description, Disabled: category.Disabled, CreateTime: category.CreateTime, UpdateTime: category.UpdateTime}
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
