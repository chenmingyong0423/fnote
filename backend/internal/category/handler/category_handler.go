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

type CategoryAndTagWithCountVO struct {
	Categories []CategoryWithCountVO `json:"categories"`
	Tags       []TagWithCountVO      `json:"tags"`
}

type CategoryWithCountVO struct {
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
	Count       int64  `json:"count"`
}

type TagWithCountVO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
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
	group.GET("", api.Wrap(h.GetCategoriesAndTags))
	group.GET("/:name/tags", api.Wrap(h.GetTagsByName))
	engine.GET("/menus", api.Wrap(h.GetMenus))
}

func (h *CategoryHandler) GetCategoriesAndTags(ctx *gin.Context) (VO CategoryAndTagWithCountVO, err error) {
	categoryAndTagWithCount, err := h.serv.GetCategoriesAndTags(ctx)
	if err != nil {
		return
	}
	VO.Categories = func() []CategoryWithCountVO {
		result := make([]CategoryWithCountVO, len(categoryAndTagWithCount.Categories))
		for i, category := range categoryAndTagWithCount.Categories {
			result[i] = CategoryWithCountVO{
				Name:        category.Name,
				Route:       category.Route,
				Description: category.Description,
				Count:       category.Count,
			}
		}
		return result
	}()
	VO.Tags = func() []TagWithCountVO {
		result := make([]TagWithCountVO, len(categoryAndTagWithCount.Tags))
		for i, tag := range categoryAndTagWithCount.Tags {
			result[i] = TagWithCountVO{
				Name:  tag.Name,
				Count: tag.Count,
			}
		}
		return result
	}()
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

func (h *CategoryHandler) GetTagsByName(ctx *gin.Context) (listVO api.ListVO[string], err error) {
	name := ctx.Param("name")
	tags, err := h.serv.GetTagsByName(ctx, name)
	if err != nil {
		return
	}
	listVO.List = append(listVO.List, tags...)
	return
}
