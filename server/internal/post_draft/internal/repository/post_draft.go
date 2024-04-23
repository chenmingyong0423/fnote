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

	"github.com/chenmingyong0423/gkit/slice"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository/dao"
	"github.com/chenmingyong0423/go-mongox"
)

type IPostDraftRepository interface {
	Save(ctx context.Context, postDraft domain.PostDraft) error
}

var _ IPostDraftRepository = (*PostDraftRepository)(nil)

func NewPostDraftRepository(dao dao.IPostDraftDao) *PostDraftRepository {
	return &PostDraftRepository{dao: dao}
}

type PostDraftRepository struct {
	dao dao.IPostDraftDao
}

func (r *PostDraftRepository) Save(ctx context.Context, postDraft domain.PostDraft) error {
	var (
		objectID  primitive.ObjectID
		err       error
		createdAt time.Time
	)
	if postDraft.Id != "" {
		objectID, err = primitive.ObjectIDFromHex(postDraft.Id)
		if err != nil {
			return err
		}
	}

	if postDraft.CreatedAt != 0 {
		createdAt = time.Unix(postDraft.CreatedAt, 0).Local()
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

	return r.dao.Save(ctx, dao.PostDraft{
		Model: mongox.Model{
			ID:        objectID,
			CreatedAt: createdAt,
		},
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
