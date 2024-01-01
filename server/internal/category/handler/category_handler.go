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

	"github.com/chenmingyong0423/fnote/backend/internal/category/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/web/request"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/web/vo"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type MenuVO struct {
	Name  string `json:"name"`
	Route string `json:"route"`
}

type CategoryWithCountVO struct {
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
	Count       int64  `json:"count"`
}

type CategoryNameVO struct {
	Name string `json:"name"`
}

func NewCategoryHandler(serv service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		serv: serv,
	}
}

type CategoryHandler struct {
	serv service.ICategoryService
}

func (h *CategoryHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/categories")
	group.GET("", api.Wrap(h.GetCategories))
	group.GET("/route/:route", api.Wrap(h.GetCategoryByRoute))
	engine.GET("/menus", api.Wrap(h.GetMenus))

	adminGroup := engine.Group("/admin/categories")
	adminGroup.GET("", api.WrapWithBody(h.AdminGetCategories))
	adminGroup.POST("", api.WrapWithBody(h.AdminCreateCategory))
	adminGroup.PUT("/disabled/:id", api.WrapWithBody(h.AdminModifyCategoryDisabled))
	adminGroup.PUT("/:id", api.WrapWithBody(h.AdminModifyCategory))
	adminGroup.DELETE("/:id", api.Wrap(h.AdminDeleteCategory))
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) (listVO api.ListVO[CategoryWithCountVO], err error) {
	categoriesWithCount, err := h.serv.GetCategories(ctx)
	if err != nil {
		return
	}
	result := make([]CategoryWithCountVO, len(categoriesWithCount))
	for i, category := range categoriesWithCount {
		result[i] = CategoryWithCountVO{
			Name:        category.Name,
			Route:       category.Route,
			Description: category.Description,
			Count:       category.Count,
		}
	}
	listVO.List = result
	return
}

func (h *CategoryHandler) GetMenus(ctx *gin.Context) (listVO api.ListVO[MenuVO], err error) {
	menus, err := h.serv.GetMenus(ctx)
	if err != nil {
		return
	}
	menuVOs := make([]MenuVO, len(menus))
	for i, menu := range menus {
		menuVOs[i] = MenuVO{Name: menu.Name, Route: menu.Route}
	}
	listVO.List = menuVOs
	return
}

func (h *CategoryHandler) GetCategoryByRoute(ctx *gin.Context) (CategoryNameVO, error) {
	route := ctx.Param("route")
	category, err := h.serv.GetCategoryByRoute(ctx, route)
	if err != nil {
		return CategoryNameVO{}, err
	}
	return CategoryNameVO{Name: category.Name}, nil
}

func (h *CategoryHandler) AdminGetCategories(ctx *gin.Context, req request.PageRequest) (pageVO vo.PageVO[vo.Category], err error) {
	categories, total, err := h.serv.AdminGetCategories(ctx, dto.PageDTO{PageNo: req.PageNo, PageSize: req.PageSize, Field: req.Field, Order: req.Order, Keyword: req.Keyword})
	if err != nil {
		return vo.PageVO[vo.Category]{}, err
	}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.categoriesToVO(categories)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return
}

func (h *CategoryHandler) categoriesToVO(categories []domain.Category) []vo.Category {
	categoryVOs := make([]vo.Category, len(categories))
	for i, category := range categories {
		categoryVOs[i] = vo.Category{
			Id:          category.Id,
			Name:        category.Name,
			Route:       category.Route,
			Disabled:    category.Disabled,
			Description: category.Description,
			CreateTime:  category.CreateTime,
			UpdateTime:  category.UpdateTime,
		}
	}
	return categoryVOs
}

func (h *CategoryHandler) AdminCreateCategory(ctx *gin.Context, req request.CreateCategoryRequest) (any, error) {
	err := h.serv.AdminCreateCategory(ctx, domain.Category{
		Name:        req.Name,
		Route:       req.Route,
		Description: req.Description,
		Disabled:    req.Disabled,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, api.NewErrorResponseBody(http.StatusConflict, "category name or route already exists")
		}
		return nil, err
	}
	return nil, nil
}

func (h *CategoryHandler) AdminModifyCategoryDisabled(ctx *gin.Context, req request.CategoryDisabledRequest) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.ModifyCategoryDisabled(ctx, id, req.Disabled)
}

func (h *CategoryHandler) AdminModifyCategory(ctx *gin.Context, req request.UpdateCategoryRequest) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.ModifyCategory(ctx, id, req.Description)
}

func (h *CategoryHandler) AdminDeleteCategory(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.DeleteCategory(ctx, id)
}
