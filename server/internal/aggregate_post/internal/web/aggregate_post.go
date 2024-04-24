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

package web

import (
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	postServ "github.com/chenmingyong0423/fnote/server/internal/post/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewAggregatePostHandler(postServ postServ.IPostService, postDraftServ post_draft.Service) *AggregatePostHandler {
	return &AggregatePostHandler{
		postServ:      postServ,
		postDraftServ: postDraftServ,
	}
}

type AggregatePostHandler struct {
	postServ      postServ.IPostService
	postDraftServ post_draft.Service
}

func (h *AggregatePostHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api")
	adminGroup.GET("/post-draft/:id", apiwrap.Wrap(h.GetPostDraftById))
}

func (h *AggregatePostHandler) GetPostDraftById(ctx *gin.Context) (*apiwrap.ResponseBody[*PostDraftVO], error) {
	id := ctx.Param("id")
	postDraftVO, err := h.getPostDraftById(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if postDraftVO == nil {
		// 查询文章是否存在
		return nil, apiwrap.NewErrorResponseBody(404, "post draft not found")
	}
	return apiwrap.SuccessResponseWithData(postDraftVO), nil
}

func (h *AggregatePostHandler) getPostDraftById(ctx *gin.Context, id string) (*PostDraftVO, error) {
	postDraft, err := h.postDraftServ.GetPostDraftById(ctx, id)
	if err != nil {
		return nil, err
	}
	var categories []Category4PostDraft
	for _, category := range postDraft.Categories {
		categories = append(categories, Category4PostDraft{
			Id:   category.Id,
			Name: category.Name,
		})
	}
	var tags []Tag4PostDraft
	for _, tag := range postDraft.Tags {
		tags = append(tags, Tag4PostDraft{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}
	return &PostDraftVO{
		Id:               postDraft.Id,
		Author:           postDraft.Author,
		Title:            postDraft.Title,
		Summary:          postDraft.Summary,
		Content:          postDraft.Content,
		CoverImg:         postDraft.CoverImg,
		Categories:       categories,
		Tags:             tags,
		StickyWeight:     postDraft.StickyWeight,
		IsDisplayed:      postDraft.IsDisplayed,
		MetaDescription:  postDraft.MetaDescription,
		MetaKeywords:     postDraft.MetaKeywords,
		WordCount:        postDraft.WordCount,
		IsCommentAllowed: postDraft.IsCommentAllowed,
		CreatedAt:        postDraft.CreatedAt,
	}, nil
}
