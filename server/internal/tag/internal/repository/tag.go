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

	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/repository/dao"
	"github.com/chenmingyong0423/gkit/slice"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITagRepository interface {
	GetTags(ctx context.Context) ([]domain.Tag, error)
	GetTagByRoute(ctx context.Context, route string) (domain.Tag, error)
	QueryTagsPage(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Tag, int64, error)
	CreateTag(ctx context.Context, tag domain.Tag) (string, error)
	ModifyTagEnabled(ctx context.Context, id string, enabled bool) error
	GetTagById(ctx context.Context, id string) (domain.Tag, error)
	DeleteTagById(ctx context.Context, id string) error
	GetSelectTags(ctx context.Context) ([]domain.Tag, error)
	IncreasePostCountByIds(ctx context.Context, tagIds []string) error
	DecreasePostCountByIds(ctx context.Context, tagIds []string) error
	FindEnabledTags(ctx context.Context) ([]domain.Tag, error)
}

var _ ITagRepository = (*TagRepository)(nil)

func NewTagRepository(dao dao.ITagDao) *TagRepository {
	return &TagRepository{dao: dao}
}

type TagRepository struct {
	dao dao.ITagDao
}

func (r *TagRepository) FindEnabledTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := r.dao.FindEnabledTags(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainTags(tags), nil
}

func (r *TagRepository) DecreasePostCountByIds(ctx context.Context, tagIds []string) (err error) {
	tagObjectIds := slice.Map(tagIds, func(_ int, t string) (objId primitive.ObjectID) {
		if err != nil {
			return objId
		}
		objId, err = primitive.ObjectIDFromHex(t)
		return objId
	})
	if err != nil {
		return err
	}
	return r.dao.DecreasePostCountByIds(ctx, tagObjectIds)
}

func (r *TagRepository) IncreasePostCountByIds(ctx context.Context, tagIds []string) (err error) {
	tagObjectIds := slice.Map(tagIds, func(_ int, t string) (objId primitive.ObjectID) {
		if err != nil {
			return objId
		}
		objId, err = primitive.ObjectIDFromHex(t)
		return objId
	})
	if err != nil {
		return err
	}
	return r.dao.IncreasePostCountByIds(ctx, tagObjectIds)
}

func (r *TagRepository) GetSelectTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := r.dao.GetEnabled(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainTags(tags), nil
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

func (r *TagRepository) ModifyTagEnabled(ctx context.Context, id string, enabled bool) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.ModifyEnabled(ctx, objectID, enabled)
}

func (r *TagRepository) CreateTag(ctx context.Context, tag domain.Tag) (string, error) {
	return r.dao.Create(ctx, &dao.Tags{Name: tag.Name, Route: tag.Route, Enabled: tag.Enabled})
}

func (r *TagRepository) QueryTagsPage(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Tag, int64, error) {
	condBuilder := query.NewBuilder()
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
		Id:        tag.ID.Hex(),
		Name:      tag.Name,
		Route:     tag.Route,
		Enabled:   tag.Enabled,
		PostCount: tag.PostCount,
		CreatedAt: tag.CreatedAt.Unix(),
		UpdatedAt: tag.UpdatedAt.Unix(),
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
