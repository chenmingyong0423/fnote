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
	"github.com/chenmingyong0423/fnote/backend/ineternal/category/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategoriesAndTags(ctx context.Context) (api.ListVO[domain.SearchCategoryVO], error)
	GetMenus(ctx context.Context) (api.ListVO[domain.MenuVO], error)
	GetTagsByName(ctx context.Context, name string) (api.ListVO[string], error)
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

func (s *CategoryService) GetTagsByName(ctx context.Context, name string) (api.ListVO[string], error) {
	var listVO api.ListVO[string]
	tags, err := s.repo.GetTagsByName(ctx, name)
	if err != nil {
		return listVO, errors.WithMessage(err, "s.repo.GetTagsByName failed")
	}
	for _, tag := range tags {
		listVO.List = append(listVO.List, tag)
	}
	return listVO, nil
}

func (s *CategoryService) GetMenus(ctx context.Context) (api.ListVO[domain.MenuVO], error) {
	var listVO api.ListVO[domain.MenuVO]
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return listVO, errors.WithMessage(err, "s.repo.GetAll failed")
	}
	listVO.List = make([]domain.MenuVO, 0, len(categories))
	for _, category := range categories {
		listVO.List = append(listVO.List, domain.MenuVO{Menu: domain.Menu{CategoryName: category.CategoryName, Route: category.Route}})
	}
	return listVO, nil
}

func (s *CategoryService) GetCategoriesAndTags(ctx context.Context) (api.ListVO[domain.SearchCategoryVO], error) {
	var listVO api.ListVO[domain.SearchCategoryVO]
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return listVO, errors.WithMessage(err, "s.repo.GetAll failed")
	}
	listVO.List = make([]domain.SearchCategoryVO, 0, len(categories))
	for _, category := range categories {
		listVO.List = append(listVO.List, domain.SearchCategoryVO{CategoryName: category.CategoryName, Tags: category.Tags})
	}
	return listVO, nil
}