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
	"github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/friend/service"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

var FriendSet = wire.NewSet(NewFriendHandler, service.NewFriendService, repository.NewFriendRepository, dao.NewFriendDao,
	wire.Bind(new(service.IFriendService), new(*service.FriendService)),
	wire.Bind(new(repository.IFriendRepository), new(*repository.FriendRepository)),
	wire.Bind(new(dao.IFriendDao), new(*dao.FriendDao)))

func NewFriendHandler(serv service.IFriendService) *FriendHandler {
	return &FriendHandler{
		serv: serv,
	}

}

type FriendHandler struct {
	serv service.IFriendService
}

func (h *FriendHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.GET("/friends", h.GetFriends)
	engine.POST("/friend", h.ApplyForFriend)
}

func (h *FriendHandler) GetFriends(ctx *gin.Context) {
	vo, err := h.serv.GetFriends(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "friend", err.Error())
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
		slog.ErrorContext(ctx, "friend", err.Error())
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
		var httpCodeError *api.HttpCodeError
		if errors.As(err, &httpCodeError) {
			ctx.AbortWithStatus(int(*httpCodeError))
		} else {
			slog.ErrorContext(ctx, "friend", err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse)
}
