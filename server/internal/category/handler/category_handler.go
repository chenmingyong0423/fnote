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

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit"

	"github.com/chenmingyong0423/fnote/server/internal/category/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
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
	group.GET("", apiwrap.Wrap(h.GetCategories))
	group.GET("/route/:route", apiwrap.Wrap(h.GetCategoryByRoute))
	engine.GET("/menus", apiwrap.Wrap(h.GetMenus))

	adminGroup := engine.Group("/admin-api/categories")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminGetCategories))
	adminGroup.GET("/select", apiwrap.Wrap(h.AdminGetSelectCategories))
	adminGroup.POST("", apiwrap.WrapWithBody(h.AdminCreateCategory))
	adminGroup.PUT("/:id/enabled", apiwrap.WrapWithBody(h.AdminModifyCategoryEnabled))
	adminGroup.PUT("/:id", apiwrap.WrapWithBody(h.AdminModifyCategory))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.AdminDeleteCategory))
	adminGroup.PUT("/:id/navigation", apiwrap.WrapWithBody(h.AdminModifyCategoryNavigation))
}

func (h *CategoryHandler) GetCategories(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[CategoryWithCountVO]], error) {
	categoriesWithCount, err := h.serv.GetCategories(ctx)
	if err != nil {
		return nil, err
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
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(result)), nil
}

func (h *CategoryHandler) GetMenus(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[MenuVO]], error) {
	menus, err := h.serv.GetMenus(ctx)
	if err != nil {
		return nil, err
	}
	menuVOs := make([]MenuVO, len(menus))
	for i, menu := range menus {
		menuVOs[i] = MenuVO{Name: menu.Name, Route: menu.Route}
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(menuVOs)), nil
}

func (h *CategoryHandler) GetCategoryByRoute(ctx *gin.Context) (*apiwrap.ResponseBody[CategoryNameVO], error) {
	route := ctx.Param("route")
	category, err := h.serv.GetCategoryByRoute(ctx, route)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(CategoryNameVO{Name: category.Name}), nil
}

func (h *CategoryHandler) AdminGetCategories(ctx *gin.Context, req request.PageRequest) (*apiwrap.ResponseBody[vo.PageVO[vo.Category]], error) {
	categories, total, err := h.serv.AdminGetCategories(ctx, dto.PageDTO{PageNo: req.PageNo, PageSize: req.PageSize, Field: req.Field, Order: req.Order, Keyword: req.Keyword})
	if err != nil {
		return nil, err
	}
	pageVO := vo.PageVO[vo.Category]{}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.categoriesToVO(categories)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *CategoryHandler) categoriesToVO(categories []domain.Category) []vo.Category {
	categoryVOs := make([]vo.Category, len(categories))
	for i, category := range categories {
		categoryVOs[i] = vo.Category{
			Id:          category.Id,
			Name:        category.Name,
			Route:       category.Route,
			Enabled:     category.Enabled,
			ShowInNav:   category.ShowInNav,
			Description: category.Description,
			PostCount:   category.PostCount,
			CreateTime:  category.CreateTime,
			UpdateTime:  category.UpdateTime,
		}
	}
	return categoryVOs
}

func (h *CategoryHandler) AdminCreateCategory(ctx *gin.Context, req request.CreateCategoryRequest) (*apiwrap.ResponseBody[any], error) {
	err := h.serv.AdminCreateCategory(ctx, domain.Category{
		Name:        req.Name,
		Route:       req.Route,
		Description: req.Description,
		ShowInNav:   req.ShowInNav,
		Enabled:     req.Enabled,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusConflict, "category name or route already exists")
		}
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *CategoryHandler) AdminModifyCategoryEnabled(ctx *gin.Context, req request.CategoryEnabledRequest) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.ModifyCategoryEnabled(ctx, id, gkit.GetValueOrDefault(req.Enabled))
}

func (h *CategoryHandler) AdminModifyCategory(ctx *gin.Context, req request.UpdateCategoryRequest) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.ModifyCategory(ctx, id, req.Description)
}

func (h *CategoryHandler) AdminDeleteCategory(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.DeleteCategory(ctx, id)
}

func (h *CategoryHandler) AdminModifyCategoryNavigation(ctx *gin.Context, req request.CategoryNavRequest) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.ModifyCategoryNavigation(ctx, id, gkit.GetValueOrDefault(req.ShowInNav))
}

func (h *CategoryHandler) AdminGetSelectCategories(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.SelectCategory]], error) {
	categories, err := h.serv.AdminGetSelectCategories(ctx)
	if err != nil {
		return nil, err
	}
	list := h.categoriesToSelectVO(categories)
	return apiwrap.SuccessResponseWithData(apiwrap.ListVO[vo.SelectCategory]{List: list}), nil
}

func (h *CategoryHandler) categoriesToSelectVO(categories []domain.Category) []vo.SelectCategory {
	result := make([]vo.SelectCategory, len(categories))
	for i, category := range categories {
		result[i] = vo.SelectCategory{Id: category.Id, Value: category.Name, Label: category.Name}
	}
	return result
}
