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
	adminGroup := engine.Group("/admin-api")
	adminGroup.POST("/post-draft", apiwrap.WrapWithBody(h.SavePostDraft))
	adminGroup.GET("/post-draft", apiwrap.WrapWithBody(h.GetPostDraftPage))
	adminGroup.DELETE("/post-draft/:id", apiwrap.Wrap(h.DeletePostDraft))
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

func (h *PostDraftHandler) GetPostDraftPage(ctx *gin.Context, req PageRequest) (*apiwrap.ResponseBody[*apiwrap.PageVO[PostDraftBriefVO]], error) {
	postDrafts, cnt, err := h.serv.GetPostDraftPage(ctx, domain.Page{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Field:    req.Field,
		Order:    req.Order,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewPageVO(req.PageNo, req.PageSize, cnt, slice.Map(postDrafts, func(idx int, pd *domain.PostDraft) PostDraftBriefVO {
		return PostDraftBriefVO{
			Id:        pd.Id,
			Title:     pd.Title,
			CreatedAt: pd.CreatedAt,
		}
	}))), nil
}

func (h *PostDraftHandler) DeletePostDraft(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	cnt, err := h.serv.DeletePostDraftById(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, apiwrap.NewErrorResponseBody(404, "id does not exist.")
	}
	return apiwrap.SuccessResponse(), nil
}
