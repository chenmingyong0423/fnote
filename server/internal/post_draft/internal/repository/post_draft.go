// Copyright 2024 chenmingyong0423

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
	"time"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/chenmingyong0423/gkit/slice"

	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository/dao"
)

type IPostDraftRepository interface {
	Save(ctx context.Context, postDraft domain.PostDraft) (string, error)
	GetById(ctx context.Context, id string) (*domain.PostDraft, error)
	DeleteById(ctx context.Context, id string) (int64, error)
}

var _ IPostDraftRepository = (*PostDraftRepository)(nil)

func NewPostDraftRepository(dao dao.IPostDraftDao) *PostDraftRepository {
	return &PostDraftRepository{dao: dao}
}

type PostDraftRepository struct {
	dao dao.IPostDraftDao
}

func (r *PostDraftRepository) DeleteById(ctx context.Context, id string) (int64, error) {
	return r.dao.DeleteById(ctx, id)
}

func (r *PostDraftRepository) GetById(ctx context.Context, id string) (*domain.PostDraft, error) {
	postDraft, err := r.dao.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(postDraft), nil
}

func (r *PostDraftRepository) Save(ctx context.Context, postDraft domain.PostDraft) (string, error) {
	var (
		createdAt time.Time
	)

	if postDraft.CreatedAt != 0 {
		createdAt = time.Unix(postDraft.CreatedAt, 0).Local()
	}

	if postDraft.Id == "" {
		postDraft.Id = uuidx.RearrangeUUID4()
	}

	categories := slice.Map(postDraft.Categories, func(idx int, c domain.Category4PostDraft) dao.Category4PostDraft {
		return dao.Category4PostDraft{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	tags := slice.Map(postDraft.Tags, func(idx int, t domain.Tag4PostDraft) dao.Tag4PostDraft {
		return dao.Tag4PostDraft{
			Id:   t.Id,
			Name: t.Name,
		}
	})

	return r.dao.Save(ctx, &dao.PostDraft{
		ID:               postDraft.Id,
		CreatedAt:        createdAt,
		Author:           postDraft.Author,
		Title:            postDraft.Title,
		Summary:          postDraft.Summary,
		Content:          postDraft.Content,
		CoverImg:         postDraft.CoverImg,
		Categories:       categories,
		Tags:             tags,
		IsDisplayed:      postDraft.IsDisplayed,
		StickyWeight:     postDraft.StickyWeight,
		MetaDescription:  postDraft.MetaDescription,
		MetaKeywords:     postDraft.MetaKeywords,
		WordCount:        postDraft.WordCount,
		IsCommentAllowed: postDraft.IsCommentAllowed,
	})
}

func (r *PostDraftRepository) toDomain(postDraft *dao.PostDraft) *domain.PostDraft {
	categories := slice.Map(postDraft.Categories, func(idx int, c dao.Category4PostDraft) domain.Category4PostDraft {
		return domain.Category4PostDraft{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	tags := slice.Map(postDraft.Tags, func(idx int, t dao.Tag4PostDraft) domain.Tag4PostDraft {
		return domain.Tag4PostDraft{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	return &domain.PostDraft{
		Id:               postDraft.ID,
		Author:           postDraft.Author,
		Title:            postDraft.Title,
		Summary:          postDraft.Summary,
		CoverImg:         postDraft.CoverImg,
		Categories:       categories,
		Tags:             tags,
		StickyWeight:     postDraft.StickyWeight,
		Content:          postDraft.Content,
		MetaDescription:  postDraft.MetaDescription,
		MetaKeywords:     postDraft.MetaKeywords,
		WordCount:        postDraft.WordCount,
		IsDisplayed:      postDraft.IsDisplayed,
		IsCommentAllowed: postDraft.IsCommentAllowed,
		CreatedAt:        postDraft.CreatedAt.Unix(),
	}
}
