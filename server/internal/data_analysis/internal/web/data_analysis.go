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

package web

import (
	"fmt"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/visit_log"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats"

	"github.com/chenmingyong0423/fnote/server/internal/comment"
	service2 "github.com/chenmingyong0423/fnote/server/internal/data_analysis/internal/service"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/post_like"
	"github.com/gin-gonic/gin"
)

func NewDataAnalysisHandler(vlServ visit_log.Service, csServ count_stats.Service, postLikeServ post_like.Service, commentServ comment.Service, ipAPiServ service2.IIpApiService) *DataAnalysisHandler {
	return &DataAnalysisHandler{
		vlServ:       vlServ,
		csServ:       csServ,
		postLikeServ: postLikeServ,
		commentServ:  commentServ,
		ipAPiServ:    ipAPiServ,
	}
}

type DataAnalysisHandler struct {
	vlServ       visit_log.Service
	csServ       count_stats.Service
	postLikeServ post_like.Service
	commentServ  comment.Service
	ipAPiServ    service2.IIpApiService
}

func (h *DataAnalysisHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/admin-api/data-analysis")
	routerGroup.GET("/traffic/today", apiwrap.Wrap(h.GetTodayTrafficStats))
	routerGroup.GET("/traffic", apiwrap.Wrap(h.GetWebsiteCountStats))
	routerGroup.GET("/content", apiwrap.Wrap(h.GetWebsiteContentStats))
	routerGroup.GET("/tendency", apiwrap.Wrap(h.GetTendencyStats))
	routerGroup.GET("/user-distribution", apiwrap.Wrap(h.GetUserDistributionStats))
}

func (h *DataAnalysisHandler) GetTodayTrafficStats(ctx *gin.Context) (*apiwrap.ResponseBody[TodayTrafficStatsVO], error) {
	// 查询当日访问量
	todayViewCount, err := h.vlServ.GetTodayViewCount(ctx)
	if err != nil {
		return nil, err
	}
	// 查询当日实际访问用户量
	userViewCount, err := h.vlServ.GetTodayUserViewCount(ctx)
	if err != nil {
		return nil, err
	}

	commentCount, err := h.commentServ.FindCommentCountOfToday(ctx)
	if err != nil {
		return nil, err
	}

	likeCount, err := h.postLikeServ.FindLikeCountToday(ctx)
	if err != nil {
		return nil, err
	}

	return apiwrap.SuccessResponseWithData(TodayTrafficStatsVO{
		ViewCount:     todayViewCount,
		UserViewCount: userViewCount,
		CommentCount:  commentCount,
		LikeCount:     likeCount,
	}), nil
}

func (h *DataAnalysisHandler) GetWebsiteCountStats(ctx *gin.Context) (*apiwrap.ResponseBody[TrafficStatsVO], error) {
	// 查询网站统计
	websiteCountStats, err := h.csServ.GetWebsiteCountStats(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(TrafficStatsVO{
		ViewCount:    websiteCountStats.WebsiteViewCount,
		CommentCount: websiteCountStats.CommentCount,
		LikeCount:    websiteCountStats.LikeCount,
	}), nil
}

func (h *DataAnalysisHandler) GetWebsiteContentStats(ctx *gin.Context) (*apiwrap.ResponseBody[ContentStatsVO], error) {
	// 查询网站统计
	websiteCountStats, err := h.csServ.GetWebsiteCountStats(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(ContentStatsVO{
		PostCount:     websiteCountStats.PostCount,
		CategoryCount: websiteCountStats.CategoryCount,
		TagCount:      websiteCountStats.TagCount,
	}), nil
}

func (h *DataAnalysisHandler) GetTendencyStats(ctx *gin.Context) (*apiwrap.ResponseBody[TendencyDataVO], error) {
	var (
		period = ctx.Query("period")
		days   = 7
	)
	if period == "month" {
		days = 30
	}
	tendencyData4PV, err := h.vlServ.GetViewTendencyStats4PV(ctx, days)
	if err != nil {
		return nil, err
	}
	tendencyData4UV, err := h.vlServ.GetViewTendencyStats4UV(ctx, days)
	if err != nil {
		return nil, err
	}

	return apiwrap.SuccessResponseWithData(TendencyDataVO{
		PV: h.tdToVO(tendencyData4PV),
		UV: h.tdToVO(tendencyData4UV),
	}), nil
}

func (h *DataAnalysisHandler) tdToVO(data []visit_log.TendencyData) []TendencyData {
	voList := make([]TendencyData, 0, len(data))
	for _, td := range data {
		voList = append(voList, TendencyData{
			Timestamp: td.Timestamp,
			ViewCount: td.ViewCount,
		})
	}
	return voList
}

func (h *DataAnalysisHandler) GetUserDistributionStats(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[UserDistributionVO]], error) {
	var (
		now   = time.Now()
		err   error
		start = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		end   = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
	)
	startParam := ctx.Query("start")
	endParam := ctx.Query("end")
	if startParam != "" {
		start, err = time.ParseInLocation(time.DateTime, startParam, time.Local)
		if err != nil {
			return nil, apiwrap.NewErrorResponseBody(400, "invalid date")
		}
		start = start.Local()
	}
	if endParam != "" {
		end, err = time.ParseInLocation(time.DateTime, endParam, time.Local)
		if err != nil {
			return nil, apiwrap.NewErrorResponseBody(400, "invalid date")
		}
		end = end.Local()
	}
	ips, err := h.vlServ.GetIpsByDate(ctx, start, end)
	if err != nil {
		return nil, err
	}
	var result []UserDistributionVO
	if len(ips) != 0 {
		userInfos, err := h.ipAPiServ.BatchGetLocation(ctx, ips)
		if err != nil {
			return nil, err
		}
		mp := make(map[string]int64, len(userInfos))
		for _, userInfo := range userInfos {
			mp[fmt.Sprintf("%s-%s", withDefault(userInfo.Country, "未知"), withDefault(userInfo.City, "未知"))]++
		}
		result = make([]UserDistributionVO, 0, len(mp))
		for k, v := range mp {
			result = append(result, UserDistributionVO{
				UserCount: v,
				Location:  k,
			})
		}
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(result)), nil
}

func withDefault(src string, dft string) string {
	if src == "" {
		return dft
	}
	return src
}
