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

	"github.com/chenmingyong0423/fnote/backend/internal/comment/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/comment/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/internal/comment/service"
	configServ "github.com/chenmingyong0423/fnote/backend/internal/config/service"
	msgService "github.com/chenmingyong0423/fnote/backend/internal/message/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/types"
	postServ "github.com/chenmingyong0423/fnote/backend/internal/post/service"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentSet = wire.NewSet(NewCommentHandler, service.NewCommentService, repository.NewCommentRepository, dao.NewCommentDao,
	wire.Bind(new(service.ICommentService), new(*service.CommentService)),
	wire.Bind(new(repository.ICommentRepository), new(*repository.CommentRepository)),
	wire.Bind(new(dao.ICommentDao), new(*dao.CommentDao)))

func NewCommentHandler(serv service.ICommentService, cfgService configServ.IConfigService, postServ postServ.IPostService, msgServ msgService.IMessageService) *CommentHandler {
	return &CommentHandler{
		serv:       serv,
		cfgService: cfgService,
		postServ:   postServ,
		msgServ:    msgServ,
	}
}

type CommentHandler struct {
	serv       service.ICommentService
	cfgService configServ.IConfigService
	postServ   postServ.IPostService
	msgServ    msgService.IMessageService
}

func (h *CommentHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/comment")
	group.POST("", h.AddComment)
	group.POST("/:commentId/reply", h.AddCommentReply)
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
		slog.ErrorContext(ctx, "comment", err)
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
		slog.ErrorContext(ctx, "comment", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !switchConfig.Status {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	post, err := h.postServ.InternalGetPunishedPostById(ctx, req.PostId)
	if err != nil {
		slog.ErrorContext(ctx, "post", err)
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
		slog.ErrorContext(ctx, "comment", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	go func() {
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			slog.WarnContext(ctx, "comment", gErr)
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "文章评论通知", "您好，您在文章有新的评论，详情请前往后台进行查看。", "text/plain")
		if gErr != nil {
			slog.WarnContext(ctx, "message", gErr)
		}
	}()
	ctx.JSON(http.StatusOK, api.SuccessResponse)
}

func (h *CommentHandler) AddCommentReply(ctx *gin.Context) {
	type ReplyRequest struct {
		PostId string `json:"postId" binding:"required"`
		// 如果是对某个回复进行回复，则是某个回复的 id
		ReplyToId string `json:"replyToId"`
		UserName  string `json:"username" binding:"required"`
		Email     string `json:"email" binding:"required,validateEmailFormat"`
		Website   string `json:"website"`
		Content   string `json:"content" binding:"required,max=200"`
	}
	req := new(ReplyRequest)
	err := ctx.BindJSON(req)
	if err != nil {
		slog.ErrorContext(ctx, "reply", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// 根评论的 id
	commentId := ctx.Param("commentId")
	ip := ctx.ClientIP()
	if ip == "" {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switchConfig, err := h.cfgService.GetSwitchStatusByTyp(ctx, "comment")
	if err != nil {
		slog.ErrorContext(ctx, "reply", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !switchConfig.Status {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	post, err := h.postServ.InternalGetPunishedPostById(ctx, req.PostId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			slog.ErrorContext(ctx, "post", err)
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		slog.ErrorContext(ctx, "post", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if !post.IsCommentAllowed {
		ctx.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = h.serv.AddCommentReply(ctx, commentId, req.PostId, domain.CommentReply{
		CommentReply: types.CommentReply{
			Content:   req.Content,
			ReplyToId: req.ReplyToId,
			UserInfo: types.UserInfo4Reply{
				Name:    req.UserName,
				Email:   req.Email,
				Website: req.Website,
				Ip:      ip,
			},
		}})
	if err != nil {
		var httpCodeError *api.HttpCodeError
		if errors.As(err, &httpCodeError) {
			ctx.AbortWithStatus(int(*httpCodeError))
		} else {
			slog.ErrorContext(ctx, "reply", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}
	go func() {
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			slog.WarnContext(ctx, "comment", gErr)
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "文章评论通知", "您好，您在文章有新的评论，详情请前往后台进行查看。", "text/plain")
		if gErr != nil {
			slog.WarnContext(ctx, "message", gErr)
		}
	}()
	ctx.JSON(http.StatusOK, api.SuccessResponse)
}
