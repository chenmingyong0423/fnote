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

package web

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/message"

	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/service"

	"github.com/spf13/viper"

	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

func NewFriendHandler(serv service.IFriendService, msgServ message.Service, cfgService website_config.Service) *FriendHandler {
	return &FriendHandler{
		serv:       serv,
		msgServ:    msgServ,
		cfgService: cfgService,
	}

}

type FriendHandler struct {
	serv       service.IFriendService
	msgServ    message.Service
	cfgService website_config.Service
}

func (h *FriendHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/friends")
	group.GET("", apiwrap.Wrap(h.GetFriends))
	group.GET("/summary", apiwrap.Wrap(h.GetFriendSummary))
	group.POST("", apiwrap.WrapWithBody(h.ApplyForFriend))

	adminGroup := engine.Group("/admin-api/friends")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminGetFriends))
	adminGroup.PUT("/:id", apiwrap.WrapWithBody(h.AdminUpdateFriend))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.AdminDeleteFriend))
	adminGroup.PUT("/:id/approval", apiwrap.Wrap(h.AdminApproveFriend))
	adminGroup.PUT("/:id/rejection", apiwrap.WrapWithBody(h.AdminRejectFriend))
}

func (h *FriendHandler) GetFriends(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[FriendVO]], error) {
	friends, err := h.serv.GetFriends(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toFriendVOs(friends))), nil
}

func (h *FriendHandler) toFriendVOs(friends []domain.Friend) []FriendVO {
	result := make([]FriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, h.toFriendVO(friend))
	}
	return result
}
func (h *FriendHandler) toFriendVO(friend domain.Friend) FriendVO {
	return FriendVO{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
	}
}

func (h *FriendHandler) ApplyForFriend(ctx *gin.Context, req FriendRequest) (*apiwrap.ResponseBody[any], error) {
	if req.Email != "" {
		regExp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !regExp.MatchString(req.Email) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Email format is incorrect.")
		}
	}
	switchConfig, err := h.cfgService.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableFriendCommit {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "Friend module is close.")
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
			return nil, apiwrap.NewErrorResponseBody(http.StatusTooManyRequests, "Already applied for")
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

	return apiwrap.SuccessResponse(), nil
}

func (h *FriendHandler) AdminGetFriends(ctx *gin.Context, req request.PageRequest) (*apiwrap.ResponseBody[vo.PageVO[AdminFriendVO]], error) {
	friends, total, err := h.serv.AdminGetFriends(ctx, dto.PageDTO{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Field:    req.Field,
		Order:    req.Order,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	pageVO := vo.PageVO[AdminFriendVO]{}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.friendToAdminVO(friends)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *FriendHandler) friendToAdminVO(friends []domain.Friend) []AdminFriendVO {
	result := make([]AdminFriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, AdminFriendVO{
			Id:          friend.Id,
			Name:        friend.Name,
			Url:         friend.Url,
			Logo:        friend.Logo,
			Description: friend.Description,
			Status:      friend.Status,
			CreateTime:  friend.CreatedAt,
		})
	}
	return result
}

func (h *FriendHandler) AdminUpdateFriend(ctx *gin.Context, req request.FriendReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AdminUpdateFriend(ctx, domain.Friend{
		Id:          ctx.Param("id"),
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
		Status:      req.Status,
	})
}

func (h *FriendHandler) AdminDeleteFriend(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	id := ctx.Param("id")
	return apiwrap.SuccessResponse(), h.serv.AdminDeleteFriend(ctx, id)
}

func (h *FriendHandler) AdminApproveFriend(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	email, err := h.serv.AdminApproveFriend(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	// 发送邮件通知朋友
	go func() {
		gErr := h.msgServ.SendEmailWithEmail(ctx, "friend-approval", []string{email}, "text/plain", viper.GetString("website.base_host")+"/friend")
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *FriendHandler) AdminRejectFriend(ctx *gin.Context, req request.FriendRejectReq) (*apiwrap.ResponseBody[any], error) {
	email, err := h.serv.AdminRejectFriend(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	// 发送邮件通知朋友
	go func() {
		gErr := h.msgServ.SendEmailWithEmail(ctx, "friend-rejection", []string{email}, "text/plain", viper.GetString("website.base_host")+"/friend", req.Reason)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *FriendHandler) GetFriendSummary(ctx *gin.Context) (*apiwrap.ResponseBody[FriendSummaryVO], error) {
	friendConfig, err := h.cfgService.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(FriendSummaryVO{Introduction: friendConfig.Introduction}), nil
}
