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

package service

import (
	"context"
	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/domain"
	"github.com/chenmingyong0423/gkit/slice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"slices"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository"

	"github.com/pkg/errors"
)

type ICommentService interface {
	AddComment(ctx context.Context, comment domain.Comment) (string, error)
	AddReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) (string, error)
	FineLatestCommentAndReply(ctx context.Context) ([]domain.LatestComment, error)
	FindCommentsByPostId(ctx context.Context, postId string) ([]domain.CommentWithReplies, error)
	AdminGetComments(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error)
	AdminApproveComment(ctx context.Context, id string) error
	FindCommentById(ctx context.Context, id string) (*domain.Comment, error)
	FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*domain.CommentReplyWithPostInfo, error)
	AdminApproveCommentReply(ctx context.Context, commentId string, replyId string) error
	FindCommentWithRepliesById(ctx context.Context, id string) (*domain.CommentWithReplies, error)
	DeleteCommentById(ctx context.Context, id string) error
	DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error
	UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, approvalStatus bool) error
	FindCommentCountOfToday(ctx context.Context) (int64, error)
	BatchApproveComments(ctx context.Context, commentIds []string, replies []domain.ReplyWithCId) ([]domain.EmailInfo, []domain.EmailInfo, error)
}

func NewCommentService(repo repository.ICommentRepository) *CommentService {
	return &CommentService{
		repo: repo,
	}
}

var _ ICommentService = (*CommentService)(nil)

type CommentService struct {
	repo repository.ICommentRepository
}

func (s *CommentService) BatchApproveComments(ctx context.Context, commentIds []string, replies []domain.ReplyWithCId) ([]domain.EmailInfo, []domain.EmailInfo, error) {
	approvalEmails := make([]domain.EmailInfo, 0, len(commentIds)+len(replies))
	// 评论被回复的邮件
	repliedEmails := make([]domain.EmailInfo, 0)
	commentObjectIDs := make([]primitive.ObjectID, len(commentIds))
	for i, commentId := range commentIds {
		objectID, err := primitive.ObjectIDFromHex(commentId)
		if err != nil {
			return nil, nil, err
		}
		commentObjectIDs[i] = objectID
	}
	var eg errgroup.Group
	eg.Go(func() error {
		approvedIds, err := s.repo.UpdateCommentStatus2TrueByIds(ctx, commentObjectIDs)
		if err != nil {
			return err
		}
		comments, err := s.repo.FindCommentByObjectIDs(ctx, approvedIds)
		if err != nil {
			slog.Error("comment-approval", "message", "approve successfully but not send email, ids=%v", approvedIds)
			return err
		}
		for _, comment := range comments {
			approvalEmails = append(approvalEmails, domain.EmailInfo{
				Email:   comment.UserInfo.Email,
				PostUrl: comment.PostInfo.PostUrl,
			})
		}
		return nil
	})
	for _, reply := range replies {
		eg.Go(func() error {
			err := s.repo.UpdateCReplyStatus2TrueByCidAndRIds(ctx, reply.CommentId, reply.ReplyIds)
			if err != nil {
				return err
			}
			comment, err := s.repo.FindCommentWithRepliesById(ctx, reply.CommentId)
			if err != nil {
				slog.Error("comment-reply-approval", "message", "approve successfully but not send email, ids=%v", reply.ReplyIds)
				return err
			}
			approvalReplies := slice.FilterMap(comment.Replies, func(_ int, c domain.CommentReply) (domain.CommentReply, bool) {
				if slices.Contains(reply.ReplyIds, c.ReplyId) {
					return c, true
				}
				return domain.CommentReply{}, false
			})
			for _, ar := range approvalReplies {
				approvalEmails = append(approvalEmails, domain.EmailInfo{
					Email:   ar.UserInfo.Email,
					PostUrl: comment.PostInfo.PostUrl,
				})
				if ar.ReplyToId != "" && ar.RepliedUserInfo.Email != "" {
					repliedEmails = append(repliedEmails, domain.EmailInfo{
						Email:   ar.RepliedUserInfo.Email,
						PostUrl: comment.PostInfo.PostUrl,
					})
				}
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return approvalEmails, repliedEmails, err
	}
	return approvalEmails, repliedEmails, nil
}

func (s *CommentService) FindCommentCountOfToday(ctx context.Context) (int64, error) {
	return s.repo.CountOfToday(ctx)
}

func (s *CommentService) UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, approvalStatus bool) error {
	return s.repo.UpdateCommentReplyStatus(ctx, commentId, replyId, approvalStatus)
}

func (s *CommentService) DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error {
	return s.repo.DeleteReplyByCIdAndRId(ctx, commentId, replyId)
}

func (s *CommentService) DeleteCommentById(ctx context.Context, id string) error {
	return s.repo.DeleteCommentById(ctx, id)
}

func (s *CommentService) FindCommentWithRepliesById(ctx context.Context, id string) (*domain.CommentWithReplies, error) {
	return s.repo.FindCommentWithRepliesById(ctx, id)
}

func (s *CommentService) AdminApproveCommentReply(ctx context.Context, commentId string, replyId string) error {
	return s.repo.UpdateCommentReplyStatus(ctx, commentId, replyId, true)
}

func (s *CommentService) FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*domain.CommentReplyWithPostInfo, error) {
	return s.repo.FindReplyByCIdAndRId(ctx, commentId, replyId)
}

func (s *CommentService) FindCommentById(ctx context.Context, id string) (*domain.Comment, error) {
	return s.repo.FindCommentById(ctx, id)
}

func (s *CommentService) AdminApproveComment(ctx context.Context, id string) error {
	return s.repo.UpdateCommentStatus2True(ctx, id)
}

func (s *CommentService) AdminGetComments(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error) {
	return s.repo.AdminFindCommentsWithPagination(ctx, page)
}

func (s *CommentService) FindCommentsByPostId(ctx context.Context, postId string) ([]domain.CommentWithReplies, error) {
	return s.repo.FindCommentsByPostIdAndCmtStatus(ctx, postId)
}

func (s *CommentService) FineLatestCommentAndReply(ctx context.Context) ([]domain.LatestComment, error) {
	// todo 默认查找最新的前 5 条，后续可能考虑动态配置
	return s.repo.FineLatestCommentAndReply(ctx, 5)
}

func (s *CommentService) AddReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) (string, error) {
	// todo 待优化，直接查询评论信息
	commentWithReplies, err := s.repo.FindApprovedCommentById(ctx, cmtId)
	if err != nil {
		return "", err
	}
	if commentWithReplies.PostInfo.PostId != postId {
		return "", errors.New("PostId is invalid.")
	}
	commentReply.RepliedUserInfo = domain.UserInfo4Reply{
		Name:  commentWithReplies.UserInfo.Name,
		Email: commentWithReplies.UserInfo.Email,
		Ip:    commentWithReplies.UserInfo.Ip,
	}
	if commentReply.ReplyToId != "" {
		isExist := false
		for _, reply := range commentWithReplies.Replies {
			if reply.ReplyId == commentReply.ReplyToId && reply.ApprovalStatus {
				commentReply.RepliedUserInfo.Name, commentReply.RepliedUserInfo.Email, commentReply.RepliedUserInfo.Website, commentReply.RepliedUserInfo.Ip = reply.UserInfo.Name, reply.UserInfo.Email, reply.UserInfo.Website, reply.UserInfo.Ip
				isExist = true
				break
			}
		}
		if !isExist {
			return "", errors.New("The replyToId does not exist.")
		}
	}
	return s.repo.AddReply(ctx, cmtId, commentReply)
}

func (s *CommentService) AddComment(ctx context.Context, comment domain.Comment) (string, error) {
	return s.repo.AddComment(ctx, comment)
}
