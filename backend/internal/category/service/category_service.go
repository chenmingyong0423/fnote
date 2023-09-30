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
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategoriesAndTags(ctx context.Context) ([]domain.Category, error)
	GetMenus(ctx context.Context) ([]domain.MenuVO, error)
	GetTagsByName(ctx context.Context, name string) ([]string, error)
}

var _ ICategoryService = (*CategoryService)(nil)

func NewCategoryService(repo repository.ICategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

type CategoryService struct {
	repo repository.ICategoryRepository
}

func (s *CategoryService) GetTagsByName(ctx context.Context, name string) ([]string, error) {
	return s.repo.GetTagsByName(ctx, name)
}

func (s *CategoryService) GetMenus(ctx context.Context) (menuVO []domain.MenuVO, err error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	menuVO = make([]domain.MenuVO, 0, len(categories))
	for _, category := range categories {
		menuVO = append(menuVO, domain.MenuVO{Menu: domain.Menu{CategoryName: category.CategoryName, Route: category.Route}})
	}
	return menuVO, nil
}

func (s *CategoryService) GetCategoriesAndTags(ctx context.Context) ([]domain.Category, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return categories, err
	}
	return categories, nil
}
