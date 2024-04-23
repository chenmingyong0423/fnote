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
	adminGroup := engine.Group("/admin")
	adminGroup.POST("/post-draft", apiwrap.WrapWithBody(h.SavePostDraft))
}

func (h *PostDraftHandler) SavePostDraft(ctx *gin.Context, req PostDraftRequest) (*apiwrap.ResponseBody[any], error) {
	categories := slice.Map(req.Categories, func(idx int, s Category4Post) domain.Category4PostDraft {
		return domain.Category4PostDraft{
			Id:   s.Id,
			Name: s.Name,
		}
	})
	tags := slice.Map(req.Tags, func(idx int, s Tag4Post) domain.Tag4PostDraft {
		return domain.Tag4PostDraft{
			Id:   s.Id,
			Name: s.Name,
		}
	})
	return apiwrap.SuccessResponse(), h.serv.SavePostDraft(ctx, domain.PostDraft{
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
}
