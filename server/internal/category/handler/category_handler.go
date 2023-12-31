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
	"github.com/chenmingyong0423/fnote/backend/internal/category/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/gin-gonic/gin"
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
