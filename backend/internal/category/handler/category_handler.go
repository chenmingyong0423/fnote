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
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

func NewCategoryHandler(serv service.ICategoryService) *CategoryHandler {
	return &CategoryHandler{
		serv: serv,
	}
}

type CategoryHandler struct {
	serv service.ICategoryService
}

func (h *CategoryHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.GET("/categories", api.Wrap(h.GetCategoriesAndTags))
	engine.GET("/categories/:name/tags", api.Wrap(h.GetTagsByName))
	engine.GET("/menus", api.Wrap(h.GetMenus))
}

func (h *CategoryHandler) GetCategoriesAndTags(ctx *gin.Context) (listVO api.ListVO[domain.SearchCategoryVO], err error) {
	categories, err := h.serv.GetCategoriesAndTags(ctx)
	if err != nil {
		return
	}
	listVO.List = make([]domain.SearchCategoryVO, 0, len(categories))
	for _, category := range categories {
		listVO.List = append(listVO.List, domain.SearchCategoryVO{CategoryName: category.CategoryName, Tags: category.Tags})
	}
	return
}

func (h *CategoryHandler) GetMenus(ctx *gin.Context) (listVO api.ListVO[domain.MenuVO], err error) {
	menuVO, err := h.serv.GetMenus(ctx)
	if err != nil {
		return
	}
	listVO.List = menuVO
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
