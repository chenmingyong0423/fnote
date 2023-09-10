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
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/service"
	configServ "github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	msgService "github.com/chenmingyong0423/fnote/backend/ineternal/message/service"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/types"
	postServ "github.com/chenmingyong0423/fnote/backend/ineternal/post/service"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func NewCommentHandler(engine *gin.Engine, serv service.ICommentService, cfgService configServ.IConfigService, postServ postServ.IPostService, msgServ msgService.IMessageService) *CommentHandler {
	ch := &CommentHandler{
		serv:       serv,
		cfgService: cfgService,
		postServ:   postServ,
		msgServ:    msgServ,
	}
	group := engine.Group("/comment")
	group.POST("", ch.AddComment)
	return ch
}

type CommentHandler struct {
	serv       service.ICommentService
	cfgService configServ.IConfigService
	postServ   postServ.IPostService
	msgServ    msgService.IMessageService
}

func (h *CommentHandler) AddComment(ctx *gin.Context) {
	type CommentRequest struct {
		PostId   string `json:"postId" binding:"required"`
		UserName string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,validateEmailFormat"`
		Website  string `json:"website"`
		Content  string `json:"content" binding:"required,max=200"`
	}
	req := new(CommentRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		slog.ErrorContext(ctx, "comment", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ip := ctx.ClientIP()
	if ip == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	switchConfig, err := h.cfgService.GetSwitchStatusByTyp(ctx, "comment")
	if err != nil {
		slog.ErrorContext(ctx, "comment", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !switchConfig.Status {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	post, err := h.postServ.InternalGetPostById(ctx, req.PostId)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !post.IsCommentAllowed {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = h.serv.AddComment(ctx, domain.Comment{Comment: types.Comment{
		PostInfo: types.PostInfo4Comment{
			PostId:    req.PostId,
			PostTitle: post.Title,
		},
		Content: req.Content,
		UserInfo: types.UserInfo4Comment{
			Name:    req.UserName,
			Email:   req.Email,
			Ip:      ip,
			Website: req.Website,
		},
	}})
	if err != nil {
		slog.ErrorContext(ctx, "comment", err.Error())
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	go func() {
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			slog.WarnContext(ctx, "comment", gErr.Error())
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "文章评论通知", "您好，您在文章有新的评论，详情请前往后台进行查看。", "text/plain")
		if gErr != nil {
			slog.WarnContext(ctx, "message", gErr.Error())
		}
	}()
	ctx.JSON(http.StatusOK, api.SuccessResponse)
}
