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
	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/service"
	service2 "github.com/chenmingyong0423/fnote/server/internal/website_config/service"
	"github.com/gin-gonic/gin"
)

func NewDataAnalysisHandler(vlServ service.IVisitLogService, cfgServ service2.IWebsiteConfigService) *DataAnalysisHandler {
	return &DataAnalysisHandler{
		vlServ:  vlServ,
		cfgServ: cfgServ,
	}
}

type DataAnalysisHandler struct {
	vlServ  service.IVisitLogService
	cfgServ service2.IWebsiteConfigService
}

func (h *DataAnalysisHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/admin/data-analysis")
	routerGroup.GET("", api.Wrap(h.GetDataAnalysis))
}

func (h *DataAnalysisHandler) GetDataAnalysis(ctx *gin.Context) (result vo.DataAnalysis, err error) {
	webSiteConfig, err := h.cfgServ.GetWebSiteConfig(ctx)
	if err != nil {
		return
	}
	// 查询当日访问量
	todayViewCount, err := h.vlServ.GetTodayViewCount(ctx)
	if err != nil {
		return
	}
	// 查询当日实际访问用户量
	todayUserVisitCount, err := h.vlServ.GetTodayUserViewCount(ctx)
	if err != nil {
		return
	}
	result.PostCount = webSiteConfig.PostCount
	result.CategoryCount = webSiteConfig.CategoryCount
	result.TotalViewCount = webSiteConfig.ViewCount
	result.TodayViewCount = todayViewCount
	result.TodayUserVisitCount = todayUserVisitCount
	result.CommentCount = 0
	result.LikeCount = 0
	return
}
