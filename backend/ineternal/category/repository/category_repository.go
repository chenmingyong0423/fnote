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
	"github.com/chenmingyong0423/fnote/backend/ineternal/category/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
)

type ICategoryRepository interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
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

func (r *CategoryRepository) GetAll(ctx context.Context) ([]domain.Category, error) {
	categories, err := r.dao.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Category, 0, len(categories))
	for _, category := range categories {
		result = append(result, domain.Category{Name: category.Name, Tags: category.Tags})
	}
	return result, nil
}
