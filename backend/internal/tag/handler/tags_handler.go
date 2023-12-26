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
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/tag/service"
	"github.com/gin-gonic/gin"
)

type TagsWithCountVO struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
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
			Count: tag.Count,
		})
	}
	return listVO, nil
}
