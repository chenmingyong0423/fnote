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

package http

import (
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	http2 "github.com/chenmingyong0423/fnote/backend/pkg/http"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func NewConfigHandler(engine *gin.Engine, serv service.IConfigService) *ConfigHandler {
	ch := &ConfigHandler{
		serv: serv,
	}

	routerGroup := engine.Group("/config")
	// 获取站长信息
	routerGroup.GET("/webmaster", ch.GetWebmasterInfo)

	return ch
}

type ConfigHandler struct {
	serv service.IConfigService
}

func (c *ConfigHandler) GetWebmasterInfo(ctx *gin.Context) {
	masterConfig, err := c.serv.GetWebmasterInfo(ctx, "webmaster")
	if err != nil {
		slog.ErrorContext(ctx, "config", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, http2.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, http2.SuccessResponse[*domain.WebMasterConfig](masterConfig))
}
