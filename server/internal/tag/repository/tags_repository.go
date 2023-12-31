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

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/tag/repository/dao"
)

type ITagRepository interface {
	GetTags(ctx context.Context) ([]domain.Tag, error)
	GetTagByRoute(ctx context.Context, route string) (domain.Tag, error)
}

var _ ITagRepository = (*TagRepository)(nil)

func NewTagRepository(dao dao.ITagDao) *TagRepository {
	return &TagRepository{dao: dao}
}

type TagRepository struct {
	dao dao.ITagDao
}

func (r *TagRepository) GetTagByRoute(ctx context.Context, route string) (domain.Tag, error) {
	tag, err := r.dao.GetByRoute(ctx, route)
	if err != nil {
		return domain.Tag{}, err
	}
	return r.toDomainTag(tag), nil
}

func (r *TagRepository) toDomainTag(tag *dao.Tags) domain.Tag {
	return domain.Tag{
		Id:         tag.Id,
		Name:       tag.Name,
		Route:      tag.Route,
		Disabled:   tag.Disabled,
		CreateTime: tag.CreateTime,
		UpdateTime: tag.UpdateTime,
	}
}

func (r *TagRepository) GetTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := r.dao.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainTags(tags), nil
}

func (r *TagRepository) toDomainTags(tags []*dao.Tags) []domain.Tag {
	var domainTags []domain.Tag
	for _, tag := range tags {
		domainTags = append(domainTags, r.toDomainTag(tag))
	}
	return domainTags
}
