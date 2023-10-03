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
	"fmt"
	"log/slog"
	"net/http"

	configServ "github.com/chenmingyong0423/fnote/backend/internal/config/service"
	msgService "github.com/chenmingyong0423/fnote/backend/internal/message/service"

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
	group := engine.Group("/friends")
	group.GET("", api.Wrap(h.GetFriends))
	group.POST("", api.WrapWithBody(h.ApplyForFriend))
}

func (h *FriendHandler) GetFriends(ctx *gin.Context) (listVO api.ListVO[domain.FriendVO], err error) {
	friends, err := h.serv.GetFriends(ctx)
	if err != nil {
		return
	}
	listVO.List = h.toFriendVOs(friends)
	return
}
func (h *FriendHandler) toFriendVOs(friends []domain.Friend) []domain.FriendVO {
	result := make([]domain.FriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, h.toFriendVO(friend))
	}
	return result
}
func (h *FriendHandler) toFriendVO(friend domain.Friend) domain.FriendVO {
	return domain.FriendVO{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Priority:    friend.Priority,
	}
}

type FriendRequest struct {
	Name        string `json:"name" binding:"required"`
	Url         string `json:"url" binding:"required"`
	Logo        string `json:"logo" binding:"required"`
	Description string `json:"description" binding:"required,max=20"`
	Email       string `json:"email" binding:"required,validateEmailFormat"`
}

func (h *FriendHandler) ApplyForFriend(ctx *gin.Context, req FriendRequest) (any, error) {
	switchConfig, err := h.cfgService.GetSwitchStatusByTyp(ctx, "friend")
	if err != nil {
		return nil, err
	}
	if !switchConfig.Status {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Friend module is close.")
	}
	err = h.serv.ApplyForFriend(ctx, domain.Friend{
		Name:        req.Name,
		Url:         req.Url,
		Logo:        req.Logo,
		Description: req.Description,
		Email:       req.Email,
	})
	if err != nil {
		return nil, err
	}

	// 发送邮件
	go func() {
		gErr := h.msgServ.SendEmailToWebmaster(ctx, "friend", "text/plain")
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()

	return nil, nil
}
