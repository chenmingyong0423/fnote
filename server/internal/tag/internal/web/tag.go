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

package web

import (
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/service"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
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
	group.GET("", apiwrap.Wrap(h.GetTags))
	group.GET("/route/:route", apiwrap.Wrap(h.GetTagByRoute))

	adminGroup := engine.Group("/admin-api/tags")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminGetTags))
	adminGroup.GET("/select", apiwrap.Wrap(h.AdminGetSelectTags))
	adminGroup.POST("", apiwrap.WrapWithBody(h.AdminCreateTag))
	adminGroup.PUT("/:id/enabled", apiwrap.WrapWithBody(h.AdminModifyTagEnabled))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.AdminDeleteTag))
}

func (h *TagHandler) GetTags(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[TagsWithCountVO]], error) {
	tags, err := h.serv.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	listVO := apiwrap.NewListVO(make([]TagsWithCountVO, 0, len(tags)))
	for _, tag := range tags {
		listVO.List = append(listVO.List, TagsWithCountVO{
			Name:  tag.Name,
			Route: tag.Route,
			Count: tag.Count,
		})
	}
	return apiwrap.SuccessResponseWithData(listVO), nil
}

func (h *TagHandler) GetTagByRoute(ctx *gin.Context) (*apiwrap.ResponseBody[TagNameVO], error) {
	route := ctx.Param("route")
	tag, err := h.serv.GetTagByRoute(ctx, route)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(TagNameVO{Name: tag.Name}), nil
}

func (h *TagHandler) AdminGetTags(ctx *gin.Context, req request.PageRequest) (*apiwrap.ResponseBody[vo.PageVO[vo.Tag]], error) {
	tags, total, err := h.serv.AdminGetTags(ctx, dto.PageDTO{PageNo: req.PageNo, PageSize: req.PageSize, Field: req.Field, Order: req.Order, Keyword: req.Keyword})
	if err != nil {
		return nil, err
	}
	pageVO := vo.PageVO[vo.Tag]{}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.tagsToVO(tags)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
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
			CreateTime: tag.CreatedAt,
			UpdateTime: tag.UpdatedAt,
		})
	}
	return result
}

func (h *TagHandler) AdminCreateTag(ctx *gin.Context, req request.CreateTagRequest) (*apiwrap.ResponseBody[any], error) {
	err := h.serv.AdminCreateTag(ctx, domain.Tag{
		Name:    req.Name,
		Route:   req.Route,
		Enabled: req.Enabled,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusConflict, "tag name or route already exists")
		}
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *TagHandler) AdminModifyTagEnabled(ctx *gin.Context, req request.TagEnabledRequest) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.ModifyTagEnabled(ctx, id, gkit.GetValueOrDefault(req.Enabled))
}

func (h *TagHandler) AdminDeleteTag(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.DeleteTag(ctx, id)
}

func (h *TagHandler) AdminGetSelectTags(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.SelectTag]], error) {
	tags, err := h.serv.GetSelectTags(ctx)
	if err != nil {
		return nil, err
	}
	list := h.tagsToSelectVO(tags)
	return apiwrap.SuccessResponseWithData(apiwrap.ListVO[vo.SelectTag]{List: list}), nil
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
