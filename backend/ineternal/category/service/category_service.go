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
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/result"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategoriesAndTagsByQueryCond(ctx context.Context) (result.ListVO, error)
	GetMenus(ctx context.Context) (result.ListVO, error)
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

func (s *CategoryService) GetMenus(ctx context.Context) (result.ListVO, error) {
	var listVO result.ListVO
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return listVO, errors.WithMessage(err, "s.repo.GetAll failed")
	}
	listVO.List = make([]any, 0, len(categories))
	for _, category := range categories {
		listVO.List = append(listVO.List, domain.MenuCategoryVO{Name: category.Name, Route: category.Route})
	}
	return listVO, nil
}

func (s *CategoryService) GetCategoriesAndTagsByQueryCond(ctx context.Context) (result.ListVO, error) {
	var listVO result.ListVO
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return listVO, errors.WithMessage(err, "s.repo.GetAll failed")
	}
	listVO.List = make([]any, 0, len(categories))
	for _, category := range categories {
		listVO.List = append(listVO.List, domain.SearchCategoryVO{Name: category.Name, Tags: category.Tags})
	}
	return listVO, nil
}
