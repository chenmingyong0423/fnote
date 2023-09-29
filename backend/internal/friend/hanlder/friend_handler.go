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

package hanlder

import (
	"log/slog"
	"net/http"

	configServ "github.com/chenmingyong0423/fnote/backend/internal/config/service"
	msgService "github.com/chenmingyong0423/fnote/backend/internal/message/service"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/log"

	"github.com/chenmingyong0423/fnote/backend/internal/friend/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

func NewFriendHandler(serv service.IFriendService, msgServ msgService.IMessageService, cfgService configServ.IConfigService) *FriendHandler {
	return &FriendHandler{
		serv:       serv,
		msgServ:    msgServ,
		cfgService: cfgService,
	}

}

type FriendHandler struct {
	serv       service.IFriendService
	msgServ    msgService.IMessageService
	cfgService configServ.IConfigService
}

func (h *FriendHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.GET("/friends", h.GetFriends)
	engine.POST("/friends", h.ApplyForFriend)
}

func (h *FriendHandler) GetFriends(ctx *gin.Context) {
	vo, err := h.serv.GetFriends(ctx)
	if err != nil {
		log.ErrorWithStack(ctx, "friend", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponseWithData(vo))
}

func (h *FriendHandler) ApplyForFriend(ctx *gin.Context) {
	type FriendRequest struct {
		Name        string `json:"name" binding:"required"`
		Url         string `json:"url" binding:"required"`
		Logo        string `json:"logo" binding:"required"`
		Description string `json:"description" binding:"required"`
		Email       string `json:"email" binding:"required,validateEmailFormat"`
	}
	req := new(FriendRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		log.ErrorWithStack(ctx, "friend", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	switchConfig, err := h.cfgService.GetSwitchStatusByTyp(ctx, "friend")
	if err != nil {
		log.ErrorWithStack(ctx, "config", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !switchConfig.Status {
		slog.WarnContext(ctx, "config", "Friend module is close.")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = h.serv.ApplyForFriend(ctx, domain.Friend{
		Name:        req.Name,
		Url:         req.Url,
		Logo:        req.Logo,
		Description: req.Description,
		Email:       req.Email,
	})
	if err != nil {
		log.ErrorWithStack(ctx, "friend", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// 发送邮件
	go func() {
		gErr := h.msgServ.SendEmailToWebmaster(ctx, "friend", "text/plain")
		if gErr != nil {
			log.WarnWithStack(ctx, "message", gErr)
		}
	}()

	ctx.JSON(http.StatusOK, api.SuccessResponse)
}
