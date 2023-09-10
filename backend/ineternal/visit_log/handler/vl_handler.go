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
	configServ "github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func NewVisitLogHandler(engine *gin.Engine, serv service.IVisitLogService, cfgServ configServ.IConfigService) *VisitLogHandler {
	h := &VisitLogHandler{
		serv:    serv,
		cfgServ: cfgServ,
	}

	routerGroup := engine.Group("/log")
	routerGroup.POST("", h.CollectVisitLog)
	return h
}

type VisitLogHandler struct {
	serv    service.IVisitLogService
	cfgServ configServ.IConfigService
}

func (h *VisitLogHandler) CollectVisitLog(ctx *gin.Context) {
	type VisitLogReq struct {
		Url       string `json:"url" bind:"required"`
		Ip        string `json:"ip"`
		UserAgent string `json:"user_agent"`
		Origin    string `json:"origin"`
		Referer   string `json:"referer"`
	}
	req := new(VisitLogReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		slog.ErrorContext(ctx, "visitLog", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	req.Ip = ctx.ClientIP()
	req.UserAgent = ctx.GetHeader("User-Agent")
	req.Origin = ctx.GetHeader("Origin")
	req.Referer = ctx.GetHeader("Referer")
	err = h.serv.CollectVisitLog(ctx, domain.VisitHistory{Url: req.Url, Ip: req.Ip, UserAgent: req.UserAgent, Origin: req.UserAgent, Referer: req.Referer})
	if err != nil {
		slog.ErrorContext(ctx, "visitLog", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	go func() {
		gErr := h.cfgServ.IncreaseWebsiteViews(ctx)
		if gErr != nil {
			slog.WarnContext(ctx, "config", gErr)
		}
	}()

	ctx.JSON(http.StatusOK, api.SuccessResponse)
}
