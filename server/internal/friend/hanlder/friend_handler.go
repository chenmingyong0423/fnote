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
	"regexp"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"

	msgService "github.com/chenmingyong0423/fnote/server/internal/message/service"
	configServ "github.com/chenmingyong0423/fnote/server/internal/website_config/service"

	"github.com/chenmingyong0423/fnote/server/internal/friend/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

func NewFriendHandler(serv service.IFriendService, msgServ msgService.IMessageService, cfgService configServ.IWebsiteConfigService) *FriendHandler {
	return &FriendHandler{
		serv:       serv,
		msgServ:    msgServ,
		cfgService: cfgService,
	}

}

type FriendHandler struct {
	serv       service.IFriendService
	msgServ    msgService.IMessageService
	cfgService configServ.IWebsiteConfigService
}

func (h *FriendHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/friends")
	group.GET("", api.Wrap(h.GetFriends))
	group.POST("", api.WrapWithBody(h.ApplyForFriend))

	adminGroup := engine.Group("/admin/friends")
	adminGroup.GET("", api.WrapWithBody(h.AdminGetFriends))
	adminGroup.PUT("/:id", api.WrapWithBody(h.AdminUpdateFriend))
	adminGroup.DELETE("/:id", api.Wrap(h.AdminDeleteFriend))
	adminGroup.PUT("/:id/approval", api.WrapWithBody(h.AdminApproveFriend))
	adminGroup.PUT("/:id/rejection", api.WrapWithBody(h.AdminRejectFriend))
}

func (h *FriendHandler) GetFriends(ctx *gin.Context) (listVO api.ListVO[vo.FriendVO], err error) {
	friends, err := h.serv.GetFriends(ctx)
	if err != nil {
		return
	}
	listVO.List = h.toFriendVOs(friends)
	return
}

func (h *FriendHandler) toFriendVOs(friends []domain.Friend) []vo.FriendVO {
	result := make([]vo.FriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, h.toFriendVO(friend))
	}
	return result
}
func (h *FriendHandler) toFriendVO(friend domain.Friend) vo.FriendVO {
	return vo.FriendVO{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
	}
}

type FriendRequest struct {
	Name        string `json:"name" binding:"required"`
	Url         string `json:"url" binding:"required"`
	Logo        string `json:"logo" binding:"required"`
	Description string `json:"description" binding:"required,max=30"`
	Email       string `json:"email"`
}

func (h *FriendHandler) ApplyForFriend(ctx *gin.Context, req FriendRequest) (any, error) {
	if req.Email != "" {
		regExp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !regExp.MatchString(req.Email) {
			return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Email format is incorrect.")
		}
	}
	switchConfig, err := h.cfgService.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableFriendCommit {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Friend module is close.")
	}
	err = h.serv.ApplyForFriend(ctx, domain.Friend{
		Name:        req.Name,
		Url:         req.Url,
		Logo:        req.Logo,
		Description: req.Description,
		Email:       req.Email,
		Ip:          ctx.ClientIP(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key error") {
			return nil, api.NewErrorResponseBody(http.StatusTooManyRequests, "Already applied for")
		}
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

func (h *FriendHandler) AdminGetFriends(ctx *gin.Context, req request.PageRequest) (pageVO vo.PageVO[vo.AdminFriendVO], err error) {
	friends, total, err := h.serv.AdminGetFriends(ctx, dto.PageDTO{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Field:    req.Field,
		Order:    req.Order,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return
	}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.friendToAdminVO(friends)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return
}

func (h *FriendHandler) friendToAdminVO(friends []domain.Friend) []vo.AdminFriendVO {
	result := make([]vo.AdminFriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, vo.AdminFriendVO{
			Id:          friend.Id,
			Name:        friend.Name,
			Url:         friend.Url,
			Logo:        friend.Logo,
			Description: friend.Description,
			Status:      friend.Status,
			CreateTime:  friend.CreateTime,
		})
	}
	return result
}

func (h *FriendHandler) AdminUpdateFriend(ctx *gin.Context, req request.FriendReq) (any, error) {
	return nil, h.serv.AdminUpdateFriend(ctx, domain.Friend{
		Id:          ctx.Param("id"),
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
		Status:      req.Status,
	})
}

func (h *FriendHandler) AdminDeleteFriend(ctx *gin.Context) (any, error) {
	id := ctx.Param("id")
	return nil, h.serv.AdminDeleteFriend(ctx, id)
}

func (h *FriendHandler) AdminApproveFriend(ctx *gin.Context, req request.FriendApproveReq) (any, error) {
	email, err := h.serv.AdminApproveFriend(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	// 发送邮件通知朋友
	go func() {
		gErr := h.msgServ.SendEmailWithEmail(ctx, "friend-approval", email, "text/plain", req.Host)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil, nil
}

func (h *FriendHandler) AdminRejectFriend(ctx *gin.Context, req request.FriendRejectReq) (any, error) {
	email, err := h.serv.AdminRejectFriend(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	// 发送邮件通知朋友
	go func() {
		gErr := h.msgServ.SendEmailWithEmail(ctx, "friend-rejection", email, "text/plain", req.Host, req.Reason)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil, nil
}
