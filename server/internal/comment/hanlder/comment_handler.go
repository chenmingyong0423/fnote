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
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/vo"

	"github.com/chenmingyong0423/fnote/server/internal/comment/service"
	msgService "github.com/chenmingyong0423/fnote/server/internal/message/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	postServ "github.com/chenmingyong0423/fnote/server/internal/post/service"
	configServ "github.com/chenmingyong0423/fnote/server/internal/website_config/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCommentHandler(serv service.ICommentService, cfgService configServ.IWebsiteConfigService, postServ postServ.IPostService, msgServ msgService.IMessageService) *CommentHandler {
	return &CommentHandler{
		serv:       serv,
		cfgService: cfgService,
		postServ:   postServ,
		msgServ:    msgServ,
	}
}

type CommentHandler struct {
	serv       service.ICommentService
	cfgService configServ.IWebsiteConfigService
	postServ   postServ.IPostService
	msgServ    msgService.IMessageService
}

type CommentRequest struct {
	PostId   string `json:"postId" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,validateEmailFormat"`
	Website  string `json:"website"`
	Content  string `json:"content" binding:"required,max=200"`
}

func (h *CommentHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/comments")
	group.GET("/id/:id", api.Wrap(h.GetCommentsByPostId))
	group.POST("", api.WrapWithBody(h.AddComment))
	group.POST("/:commentId/replies", api.WrapWithBody(h.AddCommentReply))
	group.GET("/latest", api.Wrap(h.GetLatestCommentAndReply))
}

func (h *CommentHandler) AddComment(ctx *gin.Context, req CommentRequest) (vo api.IdVO, err error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return vo, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return vo, api.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return
	}
	if !switchConfig.EnableComment {
		return vo, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	post, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		return
	}
	if !post.IsCommentAllowed {
		return vo, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	vo.Id, err = h.serv.AddComment(ctx, domain.Comment{
		PostInfo: domain.PostInfo4Comment{
			PostId:    req.PostId,
			PostTitle: post.Title,
		},
		Content: req.Content,
		UserInfo: domain.UserInfo4Comment{
			Name:    req.UserName,
			Email:   req.Email,
			Ip:      ip,
			Website: req.Website,
		},
	})
	if err != nil {
		return
	}
	go func() {
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return
}

type ReplyRequest struct {
	PostId string `json:"postId" binding:"required"`
	// 如果是对某个回复进行回复，则是某个回复的 id
	ReplyToId string `json:"replyToId"`
	UserName  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,validateEmailFormat"`
	Website   string `json:"website"`
	Content   string `json:"content" binding:"required,max=200"`
}

func (h *CommentHandler) AddCommentReply(ctx *gin.Context, req ReplyRequest) (vo api.IdVO, err error) {
	// 根评论的 id
	commentId := ctx.Param("commentId")
	ip := ctx.ClientIP()
	if ip == "" {
		return vo, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return vo, api.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return
	}
	if !switchConfig.EnableComment {
		return vo, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	post, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return vo, api.NewErrorResponseBody(http.StatusForbidden, "The postId does not exist.")
		}
		return
	}
	if !post.IsCommentAllowed {
		return vo, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	vo.Id, err = h.serv.AddCommentReply(ctx, commentId, req.PostId, domain.CommentReply{
		Content:   req.Content,
		ReplyToId: req.ReplyToId,
		UserInfo: domain.UserInfo4Reply{
			Name:    req.UserName,
			Email:   req.Email,
			Website: req.Website,
			Ip:      ip,
		},
	})
	if err != nil {
		return
	}
	go func() {
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return
}

func (h *CommentHandler) GetLatestCommentAndReply(ctx *gin.Context) (result api.ListVO[vo.LatestCommentVO], err error) {
	latestComments, err := h.serv.FineLatestCommentAndReply(ctx)
	if err != nil {
		return
	}
	lc := make([]vo.LatestCommentVO, 0, len(latestComments))
	for _, latestComment := range latestComments {
		lc = append(lc, vo.LatestCommentVO{
			PostInfo4Comment: vo.PostInfo4Comment(latestComment.PostInfo4Comment),
			Name:             latestComment.Name,
			Content:          latestComment.Content,
			CreateTime:       latestComment.CreateTime,
		})
	}
	result.List = lc
	return
}

func (h *CommentHandler) GetCommentsByPostId(ctx *gin.Context) (listVO api.ListVO[vo.PostCommentVO], err error) {
	postId := ctx.Param("id")
	comments, err := h.serv.FindCommentsByPostId(ctx, postId)
	if err != nil {
		return
	}
	pc := make([]vo.PostCommentVO, 0, len(comments))
	for _, comment := range comments {
		replies := make([]vo.PostCommentReplyVO, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			if reply.Status != domain.CommentStatusApproved {
				continue
			}
			replies = append(replies, vo.PostCommentReplyVO{
				Id:        reply.ReplyId,
				CommentId: comment.Id,
				Content:   reply.Content,
				Name:      reply.UserInfo.Name,
				Email:     reply.UserInfo.Email,
				Website:   reply.UserInfo.Website,
				ReplyToId: reply.ReplyToId,
				ReplyTo:   reply.RepliedUserInfo.Name,
				ReplyTime: reply.CreateTime,
			})
		}
		pc = append(pc, vo.PostCommentVO{
			Id:          comment.Id,
			Content:     comment.Content,
			Name:        comment.UserInfo.Name,
			Email:       comment.UserInfo.Email,
			Website:     comment.UserInfo.Website,
			CommentTime: comment.CreateTime,
			Replies:     replies,
		})
	}
	listVO.List = pc
	return
}
