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
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/service"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPostDraftHandler(serv service.IPostDraftService) *PostDraftHandler {
	return &PostDraftHandler{
		serv: serv,
	}
}

type PostDraftHandler struct {
	serv service.IPostDraftService
}

func (h *PostDraftHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api")
	adminGroup.POST("/post-draft", apiwrap.WrapWithBody(h.SavePostDraft))
	adminGroup.GET("/post-draft/:id", apiwrap.Wrap(h.GetPostDraftById))
}

func (h *PostDraftHandler) SavePostDraft(ctx *gin.Context, req PostDraftRequest) (*apiwrap.ResponseBody[map[string]string], error) {
	categories := slice.Map(req.Categories, func(idx int, s Category4PostDraft) domain.Category4PostDraft {
		return domain.Category4PostDraft{
			Id:   s.Id,
			Name: s.Name,
		}
	})
	tags := slice.Map(req.Tags, func(idx int, s Tag4PostDraft) domain.Tag4PostDraft {
		return domain.Tag4PostDraft{
			Id:   s.Id,
			Name: s.Name,
		}
	})
	id, err := h.serv.SavePostDraft(ctx, domain.PostDraft{
		Id:               req.Id,
		Author:           req.Author,
		Title:            req.Title,
		Summary:          req.Summary,
		CoverImg:         req.CoverImg,
		Categories:       categories,
		Tags:             tags,
		LikeCount:        req.WordCount,
		StickyWeight:     req.StickyWeight,
		Content:          req.Content,
		MetaDescription:  req.MetaDescription,
		MetaKeywords:     req.MetaKeywords,
		WordCount:        req.WordCount,
		IsDisplayed:      req.IsDisplayed,
		IsCommentAllowed: req.IsCommentAllowed,
		CreatedAt:        req.CreatedAt,
	})
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(map[string]string{
		"id": id,
	}), nil
}

func (h *PostDraftHandler) GetPostDraftById(ctx *gin.Context) (*apiwrap.ResponseBody[*PostDraftVO], error) {
	postDraft, err := h.serv.GetPostDraftById(ctx, ctx.Param("id"))
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if postDraft == nil {
		return nil, apiwrap.NewErrorResponseBody(404, "post draft not found")
	}
	return apiwrap.SuccessResponseWithData(h.toVO(postDraft)), nil
}

func (h *PostDraftHandler) toVO(postDraft *domain.PostDraft) *PostDraftVO {
	categories := slice.Map(postDraft.Categories, func(idx int, s domain.Category4PostDraft) Category4PostDraft {
		return Category4PostDraft{
			Id:   s.Id,
			Name: s.Name,
		}
	})
	tags := slice.Map(postDraft.Tags, func(idx int, t domain.Tag4PostDraft) Tag4PostDraft {
		return Tag4PostDraft{
			Id:   t.Id,
			Name: t.Name,
		}
	})
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
	}
}
