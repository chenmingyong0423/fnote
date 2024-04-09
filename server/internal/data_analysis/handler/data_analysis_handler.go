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
	csServ "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/service"
	"github.com/gin-gonic/gin"
)

func NewDataAnalysisHandler(vlServ service.IVisitLogService, csServ csServ.ICountStatsService) *DataAnalysisHandler {
	return &DataAnalysisHandler{
		vlServ: vlServ,
		csServ: csServ,
	}
}

type DataAnalysisHandler struct {
	vlServ service.IVisitLogService
	csServ csServ.ICountStatsService
}

func (h *DataAnalysisHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/admin/data-analysis")
	routerGroup.GET("", apiwrap.Wrap(h.GetDataAnalysis))
}

func (h *DataAnalysisHandler) GetDataAnalysis(ctx *gin.Context) (*apiwrap.ResponseBody[vo.DataAnalysis], error) {
	// 查询网站统计
	websiteCountStats, err := h.csServ.GetWebsiteCountStats(ctx)
	if err != nil {
		return nil, err
	}
	result := vo.DataAnalysis{
		PostCount:      websiteCountStats.PostCount,
		CategoryCount:  websiteCountStats.CategoryCount,
		TagCount:       websiteCountStats.TagCount,
		LikeCount:      websiteCountStats.LikeCount,
		TotalViewCount: websiteCountStats.WebsiteViewCount,
		CommentCount:   websiteCountStats.CommentCount,
	}
	// 查询当日访问量
	todayViewCount, err := h.vlServ.GetTodayViewCount(ctx)
	if err != nil {
		return nil, err
	}
	result.TodayViewCount = todayViewCount
	// 查询当日实际访问用户量
	todayUserVisitCount, err := h.vlServ.GetTodayUserViewCount(ctx)
	if err != nil {
		return nil, err
	}
	result.TodayUserVisitCount = todayUserVisitCount
	return apiwrap.SuccessResponseWithData(result), nil
}
