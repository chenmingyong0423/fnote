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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/message"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/service"

	"github.com/chenmingyong0423/fnote/server/internal/post"

	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/spf13/viper"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCommentHandler(serv service.ICommentService, cfgService website_config.Service, postServ post.Service, msgServ message.Service) *CommentHandler {
	return &CommentHandler{
		serv:       serv,
		cfgService: cfgService,
		postServ:   postServ,
		msgServ:    msgServ,
	}
}

type CommentHandler struct {
	serv       service.ICommentService
	cfgService website_config.Service
	postServ   post.Service
	msgServ    message.Service
}

func (h *CommentHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/comments")
	group.GET("/latest", apiwrap.Wrap(h.GetLatestCommentAndReply))
	// 评论
	group.GET("/id/:id", apiwrap.Wrap(h.GetCommentsByPostId))
	group.POST("", apiwrap.WrapWithBody(h.AddComment))
	// 评论回复
	group.POST("/:commentId/replies", apiwrap.WrapWithBody(h.AddCommentReply))

	adminGroup := engine.Group("/admin-api/comments")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminFindCommentsWithPagination))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.AdminDeleteComment))
	adminGroup.PUT("/:id/approval", apiwrap.Wrap(h.AdminApproveComment))

	adminGroup.DELETE("/:id/replies/:rid", apiwrap.Wrap(h.AdminDeleteCommentReply))
	adminGroup.PUT("/:id/replies/:rid/approval", apiwrap.Wrap(h.AdminApproveCommentReply))
	adminGroup.PUT("/batch-approval", apiwrap.WrapWithBody(h.AdminBatchApproveComments))
	adminGroup.DELETE("/batch-approval", apiwrap.WrapWithBody(h.AdminBatchDeleteComments))
}

func (h *CommentHandler) AddComment(ctx *gin.Context, req CommentRequest) (*apiwrap.ResponseBody[api.IdVO], error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if req.Website != "" && !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableComment {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	p, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	if !p.IsCommentAllowed {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	id, err := h.serv.AddComment(ctx, domain.Comment{
		PostInfo: domain.PostInfo{
			PostId:    req.PostId,
			PostTitle: p.Title,
			PostUrl:   fmt.Sprintf("%s/posts/%s", viper.GetString("website.base_host"), req.PostId),
		},
		Content: req.Content,
		UserInfo: domain.UserInfo{
			Name:    req.UserName,
			Email:   req.Email,
			Ip:      ip,
			Website: req.Website,
		},
	})
	if err != nil {
		return nil, err
	}
	go func() {
		// todo 考虑邮件服务订阅事件发送邮件
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}

	}()
	return apiwrap.SuccessResponseWithData(api.IdVO{Id: id}), nil
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

func (h *CommentHandler) AddCommentReply(ctx *gin.Context, req ReplyRequest) (*apiwrap.ResponseBody[api.IdVO], error) {
	// 根评论的 id
	commentId := ctx.Param("commentId")
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if req.Website != "" && !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableComment {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	p, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "The postId does not exist.")
		}
		return nil, err
	}
	if !p.IsCommentAllowed {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	id, err := h.serv.AddReply(ctx, commentId, req.PostId, domain.CommentReply{
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
		return nil, err
	}
	go func() {
		// todo 考虑邮件服务订阅事件发送邮件
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponseWithData(api.IdVO{Id: id}), nil
}

func (h *CommentHandler) GetLatestCommentAndReply(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[LatestCommentVO]], error) {
	latestComments, err := h.serv.FineLatestCommentAndReply(ctx)
	if err != nil {
		return nil, err
	}
	lc := make([]LatestCommentVO, 0, len(latestComments))
	for _, latestComment := range latestComments {
		picture := ""
		if latestComment.Email != "" {
			hash := md5.Sum([]byte(strings.ToLower(latestComment.Email)))
			picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
		}
		lc = append(lc, LatestCommentVO{
			PostInfo:  PostInfo(latestComment.PostInfo),
			Name:      latestComment.Name,
			Content:   latestComment.Content,
			Picture:   picture,
			CreatedAt: latestComment.CreatedAt,
		})
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(lc)), nil
}

func (h *CommentHandler) GetCommentsByPostId(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[PostCommentVO]], error) {
	postId := ctx.Param("id")
	comments, err := h.serv.FindCommentsByPostId(ctx, postId)
	if err != nil {
		return nil, err
	}
	pc := make([]PostCommentVO, 0, len(comments))
	for _, comment := range comments {
		replies := make([]PostCommentReplyVO, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			if !reply.ApprovalStatus {
				continue
			}
			picture := ""
			if reply.UserInfo.Email != "" {
				hash := md5.Sum([]byte(strings.ToLower(reply.UserInfo.Email)))
				picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
			}
			replies = append(replies, PostCommentReplyVO{
				Id:        reply.ReplyId,
				CommentId: comment.Id,
				Content:   reply.Content,
				Name:      reply.UserInfo.Name,
				Picture:   picture,
				Website:   reply.UserInfo.Website,
				ReplyToId: reply.ReplyToId,
				ReplyTo:   reply.RepliedUserInfo.Name,
				ReplyTime: reply.CreatedAt,
			})
		}
		picture := ""
		if comment.UserInfo.Email != "" {
			hash := md5.Sum([]byte(strings.ToLower(comment.UserInfo.Email)))
			picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
		}
		pc = append(pc, PostCommentVO{
			Id:          comment.Id,
			Content:     comment.Content,
			Name:        comment.UserInfo.Name,
			Picture:     picture,
			Website:     comment.UserInfo.Website,
			CommentTime: comment.CreateTime,
			Replies:     replies,
		})
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(pc)), nil
}

func (h *CommentHandler) AdminFindCommentsWithPagination(ctx *gin.Context, req PageRequest) (*apiwrap.ResponseBody[*apiwrap.PageVO[AdminCommentVO]], error) {
	comments, total, err := h.serv.AdminFindCommentsWithPagination(ctx, domain.Page{
		Size:           req.PageSize,
		Skip:           (req.PageNo - 1) * req.PageSize,
		Sort:           req.Sort,
		ApprovalStatus: req.ApprovalStatus,
	})
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewPageVO(req.PageNo, req.PageSize, total, h.toAdminCommentVO(comments))), nil
}

func (h *CommentHandler) toAdminCommentVO(comments []domain.AdminComment) []AdminCommentVO {
	result := make([]AdminCommentVO, 0, len(comments))
	for _, comment := range comments {
		replies := make([]AdminCommentVO, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			var picture string
			if reply.UserInfo.Email != "" {
				hash := md5.Sum([]byte(strings.ToLower(reply.UserInfo.Email)))
				picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
			}
			replies = append(replies, AdminCommentVO{
				Id:        reply.ReplyId,
				PostInfo:  PostInfo(comment.PostInfo),
				ReplyToId: reply.ReplyToId,
				Content:   reply.Content,
				UserInfo: UserInfo4Comment{
					Name:    reply.UserInfo.Name,
					Email:   reply.UserInfo.Email,
					Ip:      reply.UserInfo.Ip,
					Website: reply.UserInfo.Website,
					Picture: picture,
				},
				ApprovalStatus: reply.ApprovalStatus,
				Type:           "reply",
				CreatedAt:      reply.CreatedAt,
				UpdatedAt:      reply.UpdatedAt,
			})
		}
		var picture string
		if comment.UserInfo.Email != "" {
			hash := md5.Sum([]byte(strings.ToLower(comment.UserInfo.Email)))
			picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
		}
		result = append(result, AdminCommentVO{
			Id:       comment.Id,
			PostInfo: PostInfo(comment.PostInfo),
			Content:  comment.Content,
			UserInfo: UserInfo4Comment{
				Name:    comment.UserInfo.Name,
				Email:   comment.UserInfo.Email,
				Ip:      comment.UserInfo.Ip,
				Website: comment.UserInfo.Website,
				Picture: picture,
			},
			ReplyCount:     len(replies),
			Replies:        replies,
			ApprovalStatus: comment.ApprovalStatus,
			Type:           "comment",
			CreatedAt:      comment.CreatedAt,
			UpdatedAt:      comment.UpdatedAt,
		})
	}
	return result
}

func (h *CommentHandler) AdminApproveComment(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	comment, err := h.serv.FindCommentById(ctx, commentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "Comment not found.")
		}
		return nil, err
	}
	if comment.ApprovalStatus {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Comment has been approved.")
	}
	err = h.serv.AdminApproveComment(ctx, commentId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论已通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-approval", []string{comment.UserInfo.Email}, "text/plain", comment.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminApproveCommentReply(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	replyId := ctx.Param("rid")
	commentReplyWithPostInfo, err := h.serv.FindReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "Comment reply not found.")
		}
		return nil, err
	}
	if commentReplyWithPostInfo.ApprovalStatus {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Comment reply has been approved.")
	}
	err = h.serv.AdminApproveCommentReply(ctx, commentId, replyId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论已通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-approval", []string{commentReplyWithPostInfo.UserInfo.Email}, "text/plain", commentReplyWithPostInfo.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 通知被回复的用户接收到了回复
		gErr = h.msgServ.SendEmailWithEmail(ctx, "user-comment-reply", []string{commentReplyWithPostInfo.RepliedUserInfo.Email}, "text/plain", commentReplyWithPostInfo.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDeleteComment(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	err := h.serv.DeleteCommentById(ctx, commentId)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDeleteCommentReply(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	replyId := ctx.Param("rid")
	commentReplyWithPostInfo, err := h.serv.FindReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusNotFound, "Comment reply not found.")
		}
		return nil, err
	}
	err = h.serv.DeleteReplyByCIdAndRId(ctx, commentReplyWithPostInfo.PostInfo.PostId, commentId, replyId)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminBatchApproveComments(ctx *gin.Context, req BatchApprovedCommentRequest) (*apiwrap.ResponseBody[any], error) {
	if len(req.CommentIds) == 0 && len(req.Replies) == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "CommentIds and Replies cannot be empty.")
	}
	replies := make([]domain.ReplyWithCId, 0, len(req.Replies))
	for k, v := range req.Replies {
		replies = append(replies, domain.ReplyWithCId{
			CommentId: k,
			ReplyIds:  v,
		})
	}
	approvalEmailInfos, repliedEmailInfos, err := h.serv.BatchApproveComments(ctx, req.CommentIds, replies)
	if err != nil {
		return nil, err
	}
	go func() {
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		// 通知用户评论已通过
		l.InfoContext(ctx, fmt.Sprintf("approvalEmailInfos=%v", approvalEmailInfos))
		l.InfoContext(ctx, fmt.Sprintf("repliedEmailInfos=%v", repliedEmailInfos))
		for _, approvalEmailInfo := range approvalEmailInfos {
			gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-approval", []string{approvalEmailInfo.Email}, "text/plain", approvalEmailInfo.PostUrl)
			if gErr != nil {
				l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
			}
		}

		// 通知被回复的用户接收到了回复
		for _, repliedEmailInfo := range repliedEmailInfos {
			gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-reply", []string{repliedEmailInfo.Email}, "text/plain", repliedEmailInfo.PostUrl)
			if gErr != nil {
				l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
			}
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminBatchDeleteComments(ctx *gin.Context, req BatchApprovedCommentRequest) (*apiwrap.ResponseBody[any], error) {
	if len(req.CommentIds) == 0 && len(req.Replies) == 0 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "CommentIds and Replies cannot be empty.")
	}
	replies := make([]domain.ReplyWithCId, 0, len(req.Replies))
	for k, v := range req.Replies {
		replies = append(replies, domain.ReplyWithCId{
			CommentId: k,
			ReplyIds:  v,
		})
	}
	err := h.serv.BatchDeleteComments(ctx, req.CommentIds, replies)
	if err != nil {
		return nil, err
	}

	return apiwrap.SuccessResponse(), nil
}
