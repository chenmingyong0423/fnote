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

package handler

import (
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
	"github.com/chenmingyong0423/fnote/server/internal/tag/service"
	"github.com/chenmingyong0423/gkit"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type TagsWithCountVO struct {
	Name  string `json:"name"`
	Route string `json:"route"`
	Count int64  `json:"count"`
}

type TagNameVO struct {
	Name string `json:"name"`
}

func NewTagHandler(serv service.ITagService) *TagHandler {
	return &TagHandler{
		serv: serv,
	}
}

type TagHandler struct {
	serv service.ITagService
}

func (h *TagHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/tags")
	group.GET("", api.Wrap(h.GetTags))
	group.GET("/route/:route", api.Wrap(h.GetTagByRoute))

	adminGroup := engine.Group("/admin/tags")
	adminGroup.GET("", api.WrapWithBody(h.AdminGetTags))
	adminGroup.GET("/select", api.Wrap(h.AdminGetSelectTags))
	adminGroup.POST("", api.WrapWithBody(h.AdminCreateTag))
	adminGroup.PUT("/enabled/:id", api.WrapWithBody(h.AdminModifyTagEnabled))
	adminGroup.DELETE("/:id", api.Wrap(h.AdminDeleteTag))
}

func (h *TagHandler) GetTags(ctx *gin.Context) (listVO api.ListVO[TagsWithCountVO], err error) {
	tags, err := h.serv.GetTags(ctx)
	if err != nil {
		return listVO, err
	}
	listVO.List = make([]TagsWithCountVO, 0, len(tags))
	for _, tag := range tags {
		listVO.List = append(listVO.List, TagsWithCountVO{
			Name:  tag.Name,
			Route: tag.Route,
			Count: tag.Count,
		})
	}
	return listVO, nil
}

func (h *TagHandler) GetTagByRoute(ctx *gin.Context) (TagNameVO, error) {
	route := ctx.Param("route")
	tag, err := h.serv.GetTagByRoute(ctx, route)
	if err != nil {
		return TagNameVO{}, err
	}
	return TagNameVO{Name: tag.Name}, nil
}

func (h *TagHandler) AdminGetTags(ctx *gin.Context, req request.PageRequest) (pageVO vo.PageVO[vo.Tag], err error) {
	tags, total, err := h.serv.AdminGetTags(ctx, dto.PageDTO{PageNo: req.PageNo, PageSize: req.PageSize, Field: req.Field, Order: req.Order, Keyword: req.Keyword})
	if err != nil {
		return
	}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.tagsToVO(tags)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return
}

func (h *TagHandler) tagsToVO(tags []domain.Tag) []vo.Tag {
	result := make([]vo.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, vo.Tag{
			Id:         tag.Id,
			Name:       tag.Name,
			Route:      tag.Route,
			PostCount:  tag.PostCount,
			Enabled:    tag.Enabled,
			CreateTime: tag.CreateTime,
			UpdateTime: tag.UpdateTime,
		})
	}
	return result
}

func (h *TagHandler) AdminCreateTag(ctx *gin.Context, req request.CreateTagRequest) (any, error) {
	err := h.serv.AdminCreateTag(ctx, domain.Tag{
		Name:    req.Name,
		Route:   req.Route,
		Enabled: req.Enabled,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, api.NewErrorResponseBody(http.StatusConflict, "tag name or route already exists")
		}
		return nil, err
	}
	return nil, nil
}

func (h *TagHandler) AdminModifyTagEnabled(ctx *gin.Context, req request.TagEnabledRequest) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.ModifyTagEnabled(ctx, id, gkit.GetValueOrDefault(req.Enabled))
}

func (h *TagHandler) AdminDeleteTag(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.DeleteTag(ctx, id)
}

func (h *TagHandler) AdminGetSelectTags(ctx *gin.Context) (api.ListVO[vo.SelectTag], error) {
	tags, err := h.serv.GetSelectTags(ctx)
	if err != nil {
		return api.ListVO[vo.SelectTag]{}, err
	}
	list := h.tagsToSelectVO(tags)
	return api.ListVO[vo.SelectTag]{List: list}, nil
}

func (h *TagHandler) tagsToSelectVO(tags []domain.Tag) []vo.SelectTag {
	result := make([]vo.SelectTag, len(tags))
	for i, tag := range tags {
		result[i] = vo.SelectTag{
			Id:    tag.Id,
			Value: tag.Name,
			Label: tag.Name,
		}
	}
	return result
}
