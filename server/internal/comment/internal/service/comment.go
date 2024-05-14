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
	"net/http"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/mongo"

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
	DeleteReplyByCIdAndRId(ctx context.Context, postId string, commentId string, replyId string) error
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
	go s.subscribePostEvent()
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
	comments, err := s.FindCommentByIds(ctx, commentIds)
	if err != nil {
		return err
	}
	var eg errgroup.Group
	if len(comments) > 0 {
		eg.Go(func() error {
			var gErr error
			objectIDs := slice.Map(commentIds, func(i int, id string) primitive.ObjectID {
				var objectID primitive.ObjectID
				objectID, gErr = primitive.ObjectIDFromHex(id)
				return objectID
			})
			if gErr != nil {
				return errors.Wrap(gErr, "primitive.ObjectIDFromHex")
			}
			gErr = s.repo.DeleteCommentByIds(ctx, objectIDs)
			if gErr != nil {
				return gErr
			}
			return nil
		})
	}
	for _, reply := range replies {
		eg.Go(func() error {
			gErr := s.repo.PullReplyByCIdAndRIds(ctx, reply.CommentId, reply.ReplyIds)
			if gErr != nil {
				return gErr
			}
			return nil
		})
	}
	err = eg.Wait()
	if err != nil {
		return err
	}

	go func() {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		commentEvents := make([]domain.CommentEvent, 0, len(comments)+len(replies))
		if len(comments) > 0 {
			commentEvents = append(commentEvents, slice.Map(comments, func(_ int, c domain.AdminComment) domain.CommentEvent {
				return domain.CommentEvent{
					PostId:    c.PostInfo.PostId,
					CommentId: c.Id,
					RepliesId: slice.Map(c.Replies, func(_ int, r domain.AdminReply) string {
						return r.ReplyId
					}),
					Count: 1 + len(c.Replies),
					Type:  "delete",
				}
			})...)
		}
		if len(replies) > 0 {
			comments2, gErr := s.FindCommentByIds(ctx, slice.Map(replies, func(_ int, r domain.ReplyWithCId) string {
				return r.CommentId
			}))
			if gErr != nil {
				l.ErrorContext(ctx, "BatchDeleteComments: comment event: failed to FindCommentById", "err", gErr)
			} else {
				commentMp := make(map[string]domain.AdminComment, len(comments2))
				for _, c := range comments2 {
					commentMp[c.Id] = c
				}
				for _, cr := range replies {
					comment := commentMp[cr.CommentId]
					commentEvents = append(commentEvents, domain.CommentEvent{
						PostId:    comment.PostInfo.PostId,
						CommentId: cr.CommentId,
						RepliesId: cr.ReplyIds,
						Count:     len(cr.ReplyIds),
						Type:      "delete",
					})
				}
			}
		}
		for _, commentEvent := range commentEvents {
			marshal, fErr := json.Marshal(commentEvent)
			if fErr != nil {
				l.ErrorContext(ctx, "BatchDeleteComments: comment event: failed to marshal comment event", "err", fErr)
				continue
			}
			s.eventBus.Publish("comment", eventbus.Event{Payload: marshal})
		}
	}()
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

func (s *CommentService) DeleteReplyByCIdAndRId(ctx context.Context, postId string, commentId string, replyId string) error {
	commentEvent := domain.CommentEvent{
		PostId:    postId,
		CommentId: commentId,
		RepliesId: []string{replyId},
		Count:     1,
		Type:      "delete",
	}
	marshal, err := json.Marshal(&commentEvent)
	if err != nil {
		return err
	}
	err = s.repo.DeleteReplyByCIdAndRId(ctx, commentId, replyId)
	if err != nil {
		return err
	}
	s.eventBus.Publish("comment", eventbus.Event{Payload: marshal})
	return nil
}

func (s *CommentService) DeleteCommentById(ctx context.Context, commentId string) error {
	commentWithReplies, err := s.repo.FindCommentWithRepliesById(ctx, commentId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "Comment not found.")
		}
		return err
	}

	commentEvent := domain.CommentEvent{
		PostId:    commentWithReplies.PostInfo.PostId,
		CommentId: commentId,
		RepliesId: slice.Map(commentWithReplies.Replies, func(_ int, r domain.CommentReply) string {
			return r.ReplyId
		}),
		Count: 1 + len(commentWithReplies.Replies),
		Type:  "delete",
	}
	marshal, err := json.Marshal(&commentEvent)
	if err != nil {
		return err
	}
	err = s.repo.DeleteCommentById(ctx, commentId)
	if err != nil {
		return err
	}
	s.eventBus.Publish("comment", eventbus.Event{Payload: marshal})
	return nil
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
	rid, err := s.repo.AddReply(ctx, cmtId, commentReply)
	if err != nil {
		return "", err
	}
	commentEvent := domain.CommentEvent{
		PostId:    postId,
		CommentId: cmtId,
		RepliesId: []string{rid},
		Count:     1,
		Type:      "addition",
	}
	marshal, err := json.Marshal(&commentEvent)
	if err != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "AddReply: comment event: failed to marshal comment event")
		return rid, nil
	}
	s.eventBus.Publish("comment", eventbus.Event{Payload: marshal})
	return rid, nil
}

func (s *CommentService) AddComment(ctx context.Context, comment domain.Comment) (string, error) {
	commentId, err := s.repo.AddComment(ctx, comment)
	if err != nil {
		return "", err
	}
	commentEvent := domain.CommentEvent{
		PostId:    comment.PostInfo.PostId,
		CommentId: commentId,
		RepliesId: nil,
		Count:     1,
		Type:      "addition",
	}
	marshal, err := json.Marshal(&commentEvent)
	if err != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.ErrorContext(ctx, "AddComment: comment event: failed to marshal comment event")
		return commentId, nil
	}
	s.eventBus.Publish("comment", eventbus.Event{Payload: marshal})
	return commentId, nil
}

func (s *CommentService) subscribePostEvent() {
	eventChan := s.eventBus.Subscribe("post")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "Comment: post event", "payload", string(event.Payload))
		var e domain.PostEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "Comment: post event: failed to json.Unmarshal", "err", err)
			continue
		}
		switch e.Type {
		case "delete":
			err = s.DeleteAllCommentByPostId(ctx, e.PostId)
			if err != nil {
				l.ErrorContext(ctx, "Comment: post event: failed to delete all comment", "postId", e.PostId, "error", err)
				continue
			}
		}
		l.InfoContext(ctx, "Comment: post event: handle successfully")
	}
}

func (s *CommentService) DeleteAllCommentByPostId(ctx context.Context, postId string) error {
	return s.repo.DeleteManyByPostId(ctx, postId)
}
