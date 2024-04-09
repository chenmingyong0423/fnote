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

package handler

import (
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/gin-gonic/gin"
)

func NewCountStatsHandler(serv service.ICountStatsService) *CountStatsHandler {
	return &CountStatsHandler{
		serv: serv,
	}
}

type CountStatsHandler struct {
	serv service.ICountStatsService
}

func (h *CountStatsHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/stats")
	routerGroup.GET("", apiwrap.Wrap(h.GetWebsiteCountStats))
}

func (h *CountStatsHandler) GetWebsiteCountStats(ctx *gin.Context) (*apiwrap.ResponseBody[vo.WebsiteCountStatsVO], error) {
	// 查询网站的统计数据
	websiteCountStats, err := h.serv.GetWebsiteCountStats(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(vo.WebsiteCountStatsVO{
		PostCount:        websiteCountStats.PostCount,
		CategoryCount:    websiteCountStats.CategoryCount,
		TagCount:         websiteCountStats.TagCount,
		CommentCount:     websiteCountStats.CommentCount,
		LikeCount:        websiteCountStats.LikeCount,
		WebsiteViewCount: websiteCountStats.WebsiteViewCount,
	}), nil
}
