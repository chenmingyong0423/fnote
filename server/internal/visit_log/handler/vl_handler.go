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
	"fmt"
	"log/slog"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	csServ "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/service"
	"github.com/gin-gonic/gin"
)

func NewVisitLogHandler(serv service.IVisitLogService, csServ csServ.ICountStatsService) *VisitLogHandler {
	return &VisitLogHandler{
		serv:   serv,
		csServ: csServ,
	}
}

type VisitLogHandler struct {
	serv   service.IVisitLogService
	csServ csServ.ICountStatsService
}

func (h *VisitLogHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/logs")
	routerGroup.POST("", apiwrap.WrapWithBody(h.CollectVisitLog))
}

type VisitLogReq struct {
	Url       string `json:"url" bind:"required"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Origin    string `json:"origin"`
	Referer   string `json:"referer"`
}

func (h *VisitLogHandler) CollectVisitLog(ctx *gin.Context, req VisitLogReq) (*apiwrap.ResponseBody[any], error) {
	req.Ip = ctx.ClientIP()
	req.UserAgent = ctx.GetHeader("User-Agent")
	req.Origin = ctx.GetHeader("Origin")
	req.Referer = ctx.GetHeader("Referer")
	err := h.serv.CollectVisitLog(ctx, domain.VisitHistory{Url: req.Url, Ip: req.Ip, UserAgent: req.UserAgent, Origin: req.UserAgent, Referer: req.Referer})
	if err != nil {
		return nil, err
	}
	go func() {
		gErr := h.csServ.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeWebsiteViewCount.ToString(), domain.CountStatsTypeWebsiteViewCount)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}
