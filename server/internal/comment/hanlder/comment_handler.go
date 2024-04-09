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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	csService "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"
	"github.com/spf13/viper"

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

func NewCommentHandler(serv service.ICommentService, cfgService configServ.IWebsiteConfigService, postServ postServ.IPostService, msgServ msgService.IMessageService, statsServ csService.ICountStatsService) *CommentHandler {
	return &CommentHandler{
		serv:       serv,
		cfgService: cfgService,
		postServ:   postServ,
		msgServ:    msgServ,
		statsServ:  statsServ,
	}
}

type CommentHandler struct {
	serv       service.ICommentService
	cfgService configServ.IWebsiteConfigService
	postServ   postServ.IPostService
	msgServ    msgService.IMessageService
	statsServ  csService.ICountStatsService
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
	group.GET("/latest", apiwrap.Wrap(h.GetLatestCommentAndReply))
	// 评论
	group.GET("/id/:id", apiwrap.Wrap(h.GetCommentsByPostId))
	group.POST("", apiwrap.WrapWithBody(h.AddComment))
	// 评论回复
	group.POST("/:commentId/replies", apiwrap.WrapWithBody(h.AddCommentReply))

	adminGroup := engine.Group("/admin/comments")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminGetComments))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.AdminDeleteComment))
	adminGroup.PUT("/:id/status", apiwrap.WrapWithBody(h.AdminUpdateCommentStatus))
	adminGroup.PUT("/:id/approval", apiwrap.Wrap(h.AdminApproveComment))
	adminGroup.PUT("/:id/disapproval", apiwrap.WrapWithBody(h.AdminDisapproveComment))

	adminGroup.DELETE("/:id/replies/:rid", apiwrap.Wrap(h.AdminDeleteCommentReply))
	adminGroup.PUT("/:id/replies/:rid/status", apiwrap.WrapWithBody(h.AdminUpdateReplyStatus))
	adminGroup.PUT("/:id/replies/:rid/approval", apiwrap.Wrap(h.AdminApproveCommentReply))
	adminGroup.PUT("/:id/replies/:rid/disapproval", apiwrap.WrapWithBody(h.AdminDisapproveCommentReply))
}

func (h *CommentHandler) AddComment(ctx *gin.Context, req CommentRequest) (*apiwrap.ResponseBody[api.IdVO], error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if req.Website != "" && !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableComment {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	post, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		return nil, err
	}
	if !post.IsCommentAllowed {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	id, err := h.serv.AddComment(ctx, domain.Comment{
		PostInfo: domain.PostInfo{
			PostId:    req.PostId,
			PostTitle: post.Title,
			PostUrl:   fmt.Sprintf("%s/posts/%s", ctx.Request.Host, req.PostId),
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
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 统计评论数
		gErr = h.statsServ.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount.ToString(), domain.CountStatsTypeCommentCount)
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
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	if req.Website != "" && !strings.HasPrefix(req.Website, "http://") && !strings.HasPrefix(req.Website, "https://") {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "website format is invalid.")
	}
	switchConfig, err := h.cfgService.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !switchConfig.EnableComment {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
	}
	post, err := h.postServ.GetPunishedPostById(ctx, req.PostId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusForbidden, "The postId does not exist.")
		}
		return nil, err
	}
	if !post.IsCommentAllowed {
		return nil, api.NewErrorResponseBody(http.StatusForbidden, "Comment module is closed.")
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
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.IncreaseVisitCount(ctx, req.PostId)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = h.msgServ.SendEmailToWebmaster(ctx, "comment", "text/plain")
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 统计评论数
		gErr = h.statsServ.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount.ToString(), domain.CountStatsTypeCommentCount)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponseWithData(api.IdVO{Id: id}), nil
}

func (h *CommentHandler) GetLatestCommentAndReply(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.LatestCommentVO]], error) {
	latestComments, err := h.serv.FineLatestCommentAndReply(ctx)
	if err != nil {
		return nil, err
	}
	lc := make([]vo.LatestCommentVO, 0, len(latestComments))
	for _, latestComment := range latestComments {
		picture := ""
		if latestComment.Email != "" {
			hash := md5.Sum([]byte(strings.ToLower(latestComment.Email)))
			picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
		}
		lc = append(lc, vo.LatestCommentVO{
			PostInfo:   vo.PostInfo(latestComment.PostInfo),
			Name:       latestComment.Name,
			Content:    latestComment.Content,
			Picture:    picture,
			CreateTime: latestComment.CreateTime,
		})
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(lc)), nil
}

func (h *CommentHandler) GetCommentsByPostId(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.PostCommentVO]], error) {
	postId := ctx.Param("id")
	comments, err := h.serv.FindCommentsByPostId(ctx, postId)
	if err != nil {
		return nil, err
	}
	pc := make([]vo.PostCommentVO, 0, len(comments))
	for _, comment := range comments {
		replies := make([]vo.PostCommentReplyVO, 0, len(comment.Replies))
		for _, reply := range comment.Replies {
			if reply.Status != domain.CommentStatusApproved {
				continue
			}
			picture := ""
			if reply.UserInfo.Email != "" {
				hash := md5.Sum([]byte(strings.ToLower(reply.UserInfo.Email)))
				picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
			}
			replies = append(replies, vo.PostCommentReplyVO{
				Id:        reply.ReplyId,
				CommentId: comment.Id,
				Content:   reply.Content,
				Name:      reply.UserInfo.Name,
				Picture:   picture,
				Website:   reply.UserInfo.Website,
				ReplyToId: reply.ReplyToId,
				ReplyTo:   reply.RepliedUserInfo.Name,
				ReplyTime: reply.CreateTime,
			})
		}
		picture := ""
		if comment.UserInfo.Email != "" {
			hash := md5.Sum([]byte(strings.ToLower(comment.UserInfo.Email)))
			picture = viper.GetString("gravatar.api") + hex.EncodeToString(hash[:])
		}
		pc = append(pc, vo.PostCommentVO{
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

func (h *CommentHandler) AdminGetComments(ctx *gin.Context, req request.PageRequest) (*apiwrap.ResponseBody[apiwrap.PageVO[vo.AdminCommentVO]], error) {
	friends, total, err := h.serv.AdminGetComments(ctx, dto.PageDTO{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Field:    req.Field,
		Order:    req.Order,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	pageVO := apiwrap.PageVO[vo.AdminCommentVO]{}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.toAdminCommentVO(friends)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *CommentHandler) toAdminCommentVO(friends []domain.AdminComment) []vo.AdminCommentVO {
	result := make([]vo.AdminCommentVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, vo.AdminCommentVO{
			Id:         friend.Id,
			PostInfo:   vo.PostInfo(friend.PostInfo),
			Content:    friend.Content,
			UserInfo:   vo.AdminUserInfoVO(friend.UserInfo),
			Fid:        friend.Fid,
			Type:       friend.Type,
			Status:     friend.Status,
			CreateTime: friend.CreateTime,
		})
	}
	return result
}

func (h *CommentHandler) AdminApproveComment(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	comment, err := h.serv.FindCommentById(ctx, commentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment not found.")
		}
		return nil, err
	}
	if comment.IsApproved() {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Comment has been approved.")
	}
	err = h.serv.AdminApproveComment(ctx, commentId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论已通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-approval", comment.UserInfo.Email, "text/plain", comment.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDisapproveComment(ctx *gin.Context, req request.DisapproveCommentRequest) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	comment, err := h.serv.FindCommentById(ctx, commentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment not found.")
		}
		return nil, err
	}
	if comment.IsDisapproved() {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Comment has been disapproved.")
	}

	err = h.serv.AdminDisapproveComment(ctx, commentId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论未通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-disapproval", comment.UserInfo.Email, "text/plain", comment.PostInfo.PostUrl, req.Reason)
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
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment reply not found.")
		}
		return nil, err
	}
	if commentReplyWithPostInfo.IsApproved() {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Comment reply has been approved.")
	}
	err = h.serv.AdminApproveCommentReply(ctx, commentId, replyId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论已通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-approval", commentReplyWithPostInfo.UserInfo.Email, "text/plain", commentReplyWithPostInfo.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 通知被回复的用户接收到了回复
		gErr = h.msgServ.SendEmailWithEmail(ctx, "user-comment-reply", commentReplyWithPostInfo.RepliedUserInfo.Email, "text/plain", commentReplyWithPostInfo.PostInfo.PostUrl)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDisapproveCommentReply(ctx *gin.Context, req request.DisapproveCommentRequest) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	replyId := ctx.Param("rid")
	commentReplyWithPostInfo, err := h.serv.FindReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment reply not found.")
		}
		return nil, err
	}
	if commentReplyWithPostInfo.IsDisapproved() {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Comment reply has been disapproved.")
	}
	err = h.serv.AdminDisapproveCommentReply(ctx, commentId, replyId)
	if err != nil {
		return nil, err
	}
	go func() {
		// 通知用户评论未通过
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.msgServ.SendEmailWithEmail(ctx, "user-comment-disapproval", commentReplyWithPostInfo.UserInfo.Email, "text/plain", commentReplyWithPostInfo.PostInfo.PostUrl, req.Reason)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDeleteComment(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	commentWithReplies, err := h.serv.FindCommentWithRepliesById(ctx, commentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment not found.")
		}
		return nil, err
	}
	err = h.serv.DeleteCommentById(ctx, commentId)
	if err != nil {
		return nil, err
	}
	// 统计删除的评论数
	go func() {
		cnt := 1 + len(commentWithReplies.Replies)
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.DecreaseCommentCount(ctx, commentWithReplies.PostInfo.PostId, cnt)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 减少评论数
		gErr = h.statsServ.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount.ToString(), domain.CountStatsTypeCommentCount)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminDeleteCommentReply(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	replyId := ctx.Param("rid")
	commentReplyWithPostInfo, err := h.serv.FindReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusNotFound, "Comment reply not found.")
		}
		return nil, err
	}
	err = h.serv.DeleteReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		return nil, err
	}

	go func() {
		l := slog.Default().With("X-Request-ID", ctx.GetString("X-Request-ID"))
		gErr := h.postServ.DecreaseCommentCount(ctx, commentReplyWithPostInfo.PostInfo.PostId, 1)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 减少评论数
		gErr = h.statsServ.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount.ToString(), domain.CountStatsTypeCommentCount)
		if gErr != nil {
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminUpdateCommentStatus(ctx *gin.Context, req request.CommentStatusRequest) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	err := h.serv.UpdateCommentStatus(ctx, commentId, domain.CommentStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *CommentHandler) AdminUpdateReplyStatus(ctx *gin.Context, req request.CommentStatusRequest) (*apiwrap.ResponseBody[any], error) {
	commentId := ctx.Param("id")
	replyId := ctx.Param("rid")
	err := h.serv.UpdateCommentReplyStatus(ctx, commentId, replyId, domain.CommentStatus(req.Status))
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}
