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

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/backend/internal/tag/repository/dao"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
)

type ITagRepository interface {
	GetTags(ctx context.Context) ([]domain.Tag, error)
	GetTagByRoute(ctx context.Context, route string) (domain.Tag, error)
	QueryTagsPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error)
	CreateTag(ctx context.Context, tag domain.Tag) (string, error)
	ModifyTagDisabled(ctx context.Context, id string, disabled bool) error
	GetTagById(ctx context.Context, id string) (domain.Tag, error)
	DeleteTagById(ctx context.Context, id string) error
	RecoverTag(ctx context.Context, tag domain.Tag) error
}

var _ ITagRepository = (*TagRepository)(nil)

func NewTagRepository(dao dao.ITagDao) *TagRepository {
	return &TagRepository{dao: dao}
}

type TagRepository struct {
	dao dao.ITagDao
}

func (r *TagRepository) RecoverTag(ctx context.Context, tag domain.Tag) error {
	id, err := primitive.ObjectIDFromHex(tag.Id)
	if err != nil {
		return err
	}
	return r.dao.RecoverTag(ctx, dao.Tags{
		Id:         id,
		Name:       tag.Name,
		Route:      tag.Route,
		Disabled:   tag.Disabled,
		CreateTime: tag.CreateTime,
		UpdateTime: tag.UpdateTime,
	})
}

func (r *TagRepository) DeleteTagById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *TagRepository) GetTagById(ctx context.Context, id string) (t domain.Tag, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	tag, err := r.dao.GetById(ctx, objectID)
	if err != nil {
		return
	}
	return r.toDomainTag(tag), nil
}

func (r *TagRepository) ModifyTagDisabled(ctx context.Context, id string, disabled bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyDisabled(ctx, objectID, disabled)
}

func (r *TagRepository) CreateTag(ctx context.Context, tag domain.Tag) (string, error) {
	now := time.Now().Unix()
	return r.dao.Create(ctx, &dao.Tags{Name: tag.Name, Route: tag.Route, Disabled: tag.Disabled, CreateTime: now, UpdateTime: now})

}

func (r *TagRepository) QueryTagsPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error) {
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
		findOptions.SetSort(bsonx.M("create_time", -1))
	}
	categories, total, err := r.dao.QuerySkipAndSetLimit(ctx, cond, findOptions)
	return r.toDomainTags(categories), total, err
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
		Id:         tag.Id.Hex(),
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
