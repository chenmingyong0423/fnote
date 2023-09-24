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
	"log/slog"
	"net/http"

	"github.com/chenmingyong0423/fnote/backend/internal/config/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

func NewConfigHandler(serv service.IConfigService) *ConfigHandler {
	return &ConfigHandler{
		serv: serv,
	}
}

type ConfigHandler struct {
	serv service.IConfigService
}

func (h *ConfigHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/configs")
	// 获取站长信息
	routerGroup.GET("/webmaster", h.GetWebmasterInfo)
}

func (c *ConfigHandler) GetWebmasterInfo(ctx *gin.Context) {
	masterConfigVO, err := c.serv.GetWebmasterInfo(ctx, "webmaster")
	if err != nil {
		slog.ErrorContext(ctx, "config", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponseWithData[*domain.WebMasterConfigVO](masterConfigVO))
}
