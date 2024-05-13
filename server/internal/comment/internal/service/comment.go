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
	"encoding/json"
	"log/slog"

	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/uuid"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/domain"
	"github.com/chenmingyong0423/gkit/slice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/sync/errgroup"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository"

	"github.com/pkg/errors"
)

type ICommentService interface {
	AddComment(ctx context.Context, comment domain.Comment) (string, error)
	AddReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) (string, error)
	FineLatestCommentAndReply(ctx context.Context) ([]domain.LatestComment, error)
	FindCommentsByPostId(ctx context.Context, postId string) ([]domain.CommentWithReplies, error)
	AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error)
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
	BatchDeleteComments(ctx context.Context, commentIds []string, replies []domain.ReplyWithCId) error
	FindCommentByIds(ctx context.Context, commentIds []string) ([]domain.AdminComment, error)
}

func NewCommentService(repo repository.ICommentRepository, eventBus *eventbus.EventBus) *CommentService {
	s := &CommentService{
		repo:     repo,
		eventBus: eventBus,
	}
	go s.SubscribePostDeletedEvent()
	return s
}

var _ ICommentService = (*CommentService)(nil)

type CommentService struct {
	repo     repository.ICommentRepository
	eventBus *eventbus.EventBus
}

func (s *CommentService) FindCommentByIds(ctx context.Context, commentIds []string) ([]domain.AdminComment, error) {
	var err error
	objectIDs := slice.Map(commentIds, func(i int, id string) primitive.ObjectID {
		var objectID primitive.ObjectID
		objectID, err = primitive.ObjectIDFromHex(id)
		return objectID
	})
	if err != nil {
		return nil, err
	}
	return s.repo.FindCommentByObjectIDs(ctx, objectIDs)
}

func (s *CommentService) BatchDeleteComments(ctx context.Context, commentIds []string, replies []domain.ReplyWithCId) error {
	var eg errgroup.Group
	if len(commentIds) > 0 {
		eg.Go(func() error {
			var err error
			objectIDs := slice.Map(commentIds, func(i int, id string) primitive.ObjectID {
				var objectID primitive.ObjectID
				objectID, err = primitive.ObjectIDFromHex(id)
				return objectID
			})
			if err != nil {
				return errors.Wrap(err, "primitive.ObjectIDFromHex")
			}
			err = s.repo.DeleteCommentByIds(ctx, objectIDs)
			if err != nil {
				return err
			}
			return nil
		})
	}
	for _, reply := range replies {
		eg.Go(func() error {
			err := s.repo.PullReplyByCIdAndRIds(ctx, reply.CommentId, reply.ReplyIds)
			if err != nil {
				return err
			}
			return nil
		})
	}
	err := eg.Wait()
	if err != nil {
		return err
	}
	return nil
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
	if len(commentObjectIDs) > 0 {
		eg.Go(func() error {
			// 查询未被审核的评论
			comments, err := s.repo.FindDisapprovedCommentByObjectIDs(ctx, commentObjectIDs)
			if err != nil {
				return err
			}
			commentObjectIDs = slice.Map(comments, func(_ int, c domain.AdminComment) primitive.ObjectID {
				ojbId, _ := primitive.ObjectIDFromHex(c.Id)
				return ojbId
			})
			if len(commentObjectIDs) == 0 {
				return nil
			}
			err = s.repo.UpdateCommentStatus2TrueByIds(ctx, commentObjectIDs)
			if err != nil {
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
	}
	for _, reply := range replies {
		eg.Go(func() error {
			comment, err := s.repo.FindCommentWithDisapprovedReplyByCidAndRIds(ctx, reply.CommentId, reply.ReplyIds)
			if err != nil {
				return err
			}
			reply.ReplyIds = slice.Map(comment.Replies, func(_ int, r domain.AdminReply) string {
				return r.ReplyId
			})
			if len(reply.ReplyIds) == 0 {
				return nil
			}
			err = s.repo.UpdateCReplyStatus2TrueByCidAndRIds(ctx, reply.CommentId, reply.ReplyIds)
			if err != nil {
				return err
			}

			for _, ar := range comment.Replies {
				approvalEmails = append(approvalEmails, domain.EmailInfo{
					Email:   ar.UserInfo.Email,
					PostUrl: comment.PostInfo.PostUrl,
				})
				if ar.ReplyToId != "" && ar.RepliedUserInfo.Email != "" {
					repliedEmails = append(repliedEmails, domain.EmailInfo{
						Email:   ar.RepliedUserInfo.Email,
						PostUrl: comment.PostInfo.PostUrl,
					})
				} else {
					repliedEmails = append(repliedEmails, domain.EmailInfo{
						Email:   comment.UserInfo.Email,
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

func (s *CommentService) AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error) {
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

func (s *CommentService) SubscribePostDeletedEvent() {
	type PostInfo struct {
		PostId string `json:"post_id"`
	}
	rid := uuid.NewString()
	ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
	l := slog.Default().With("X-Request-ID", rid)
	var postInfo PostInfo

	eventChan := s.eventBus.Subscribe("post-delete")
	for event := range eventChan {
		l.InfoContext(ctx, "post-delete", "payload", event.Payload)
		if payload, ok := event.Payload.(map[string]any); ok {
			marshal, err := json.Marshal(payload)
			if err != nil {
				l.ErrorContext(ctx, "post-delete: invalid payload", "payload", event.Payload)
				continue
			}
			err = json.Unmarshal(marshal, &postInfo)
			if err != nil {
				l.ErrorContext(ctx, "post-delete: json.Unmarshal failed", "payload", event.Payload, "error", err)
				continue
			}
			err = s.DeleteAllCommentByPostId(ctx, postInfo.PostId)
			if err != nil {
				l.ErrorContext(ctx, "post-delete: failed to delete all comment", "postId", postInfo.PostId, "error", err)
				continue
			}
			l.InfoContext(ctx, "post-delete: delete all comment success", "postId", postInfo.PostId)
		} else {
			l.ErrorContext(ctx, "post-delete: invalid payload", "payload", event.Payload)
		}
	}
}

func (s *CommentService) DeleteAllCommentByPostId(ctx context.Context, postId string) error {
	return s.repo.DeleteManyByPostId(ctx, postId)
}
