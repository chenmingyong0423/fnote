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

package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository/dao"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment) (string, error)
	FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error)
	AddReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) (string, error)
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]domain.LatestComment, error)
	FindCommentsByPostIdAndCmtStatus(ctx context.Context, postId string) ([]domain.CommentWithReplies, error)
	FindPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.AdminComment, int64, error)
	FindCommentById(ctx context.Context, id string) (*domain.Comment, error)
	UpdateCommentStatus2True(ctx context.Context, id string) error
	FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*domain.CommentReplyWithPostInfo, error)
	UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, approvalStatus bool) error
	FindCommentWithRepliesById(ctx context.Context, id string) (*domain.CommentWithReplies, error)
	DeleteCommentById(ctx context.Context, id string) error
	DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error
	CountOfToday(ctx context.Context) (int64, error)
	AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error)
}

func NewCommentRepository(dao dao.ICommentDao) *CommentRepository {
	return &CommentRepository{
		dao: dao,
	}
}

var _ ICommentRepository = (*CommentRepository)(nil)

type CommentRepository struct {
	dao dao.ICommentDao
}

func (r *CommentRepository) AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error) {
	condBuilder := query.BsonBuilder()
	if page.Keyword != "" {
		condBuilder.RegexOptions("content", fmt.Sprintf(".*%s.*", strings.TrimSpace(page.Keyword)), "i")
	}
	cond := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip(page.Skip).SetLimit(page.Size)
	if page.Field != "" && page.Order != "" {
		findOptions.SetSort(bsonx.M(page.Field, page.OrderConvertToInt()))
	} else {
		findOptions.SetSort(bsonx.M("create_time", 1))
	}
	comments, total, err := r.dao.Find(ctx, cond, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return r.toDomainAdminCommentsV2(comments), total, nil
}

func (r *CommentRepository) CountOfToday(ctx context.Context) (int64, error) {
	return r.dao.CountOfToday(ctx)
}

func (r *CommentRepository) DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error {
	objectID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}
	return r.dao.DeleteReplyByCIdAndRId(ctx, objectID, replyId)
}

func (r *CommentRepository) DeleteCommentById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *CommentRepository) FindCommentWithRepliesById(ctx context.Context, id string) (*domain.CommentWithReplies, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	comment, err := r.dao.FindCommentWithRepliesById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return r.toDomainCommentWithReplies(comment), nil
}

func (r *CommentRepository) UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, approvalStatus bool) error {
	objectID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}
	return r.dao.UpdateCommentReplyStatus(ctx, objectID, replyId, approvalStatus)
}

func (r *CommentRepository) FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*domain.CommentReplyWithPostInfo, error) {
	objectCommentID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return nil, err
	}
	commentReply, err := r.dao.FindReplyByCIdAndRId(ctx, objectCommentID, replyId)
	if err != nil {
		return nil, err
	}
	return &domain.CommentReplyWithPostInfo{
		CommentReply: domain.CommentReply{
			ReplyId:         commentReply.ReplyId,
			Content:         commentReply.Content,
			ReplyToId:       commentReply.ReplyToId,
			UserInfo:        domain.UserInfo4Reply(commentReply.UserInfo),
			RepliedUserInfo: domain.UserInfo4Reply(commentReply.RepliedUserInfo),
			ApprovalStatus:  commentReply.ApprovalStatus,
			CreatedAt:       commentReply.CreatedAt.Unix(),
		},
		PostInfo: domain.PostInfo(commentReply.PostInfo),
	}, nil
}

func (r *CommentRepository) UpdateCommentStatus2True(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdateCommentStatus2True(ctx, objectID)
}

func (r *CommentRepository) FindCommentById(ctx context.Context, id string) (*domain.Comment, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	comment, err := r.dao.FindCommentById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return r.toDomainComment(comment), nil
}

func (r *CommentRepository) FindPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.AdminComment, int64, error) {
	condBuilder := query.BsonBuilder()
	if pageDTO.Keyword != "" {
		condBuilder.RegexOptions("content", fmt.Sprintf(".*%s.*", strings.TrimSpace(pageDTO.Keyword)), "i")
	}
	cond := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip((pageDTO.PageNo - 1) * pageDTO.PageSize).SetLimit(pageDTO.PageSize)
	if pageDTO.Field != "" && pageDTO.Order != "" {
		findOptions.SetSort(bsonx.M(pageDTO.Field, pageDTO.OrderConvertToInt()))
	} else {
		findOptions.SetSort(bsonx.M("create_time", 1))
	}
	friends, total, err := r.dao.AggregationQuerySkipAndSetLimit(ctx, cond, findOptions)
	return r.toDomainAdminComments(friends), total, err
}

func (r *CommentRepository) FindCommentsByPostIdAndCmtStatus(ctx context.Context, postId string) ([]domain.CommentWithReplies, error) {
	comments, err := r.dao.FindApprovedCommentsByPostId(ctx, postId)
	if err != nil {
		return nil, err
	}
	return r.toDomainComments(comments), nil
}

func (r *CommentRepository) FineLatestCommentAndReply(ctx context.Context, cnt int) ([]domain.LatestComment, error) {
	latestComments, err := r.dao.FineLatestCommentAndReply(ctx, cnt)
	if err != nil {
		return nil, err
	}
	result := make([]domain.LatestComment, 0, len(latestComments))
	for _, latestComment := range latestComments {
		result = append(result, domain.LatestComment{
			PostInfo:  domain.PostInfo(latestComment.PostInfo),
			Name:      latestComment.Name,
			Email:     latestComment.Email,
			Content:   latestComment.Content,
			CreatedAt: latestComment.CreatedAt.Unix(),
		})
	}
	return result, nil
}

func (r *CommentRepository) AddReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) (string, error) {
	now := time.Now()
	id := uuid.NewString()
	objectID, err := primitive.ObjectIDFromHex(cmtId)
	if err != nil {
		return "", err
	}
	return id, r.dao.AddCommentReply(ctx, objectID, dao.Reply{
		ReplyId:         id,
		Content:         commentReply.Content,
		ReplyToId:       commentReply.ReplyToId,
		UserInfo:        dao.UserInfo4Reply(commentReply.UserInfo),
		RepliedUserInfo: dao.UserInfo4Reply(commentReply.RepliedUserInfo),
		ApprovalStatus:  false,
		CreatedAt:       now,
		UpdatedAt:       now,
	})
}

func (r *CommentRepository) FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error) {
	objectID, err := primitive.ObjectIDFromHex(cmtId)
	if err != nil {
		return nil, err
	}
	comment, err := r.dao.FindApprovedCommentById(ctx, objectID)
	if err != nil {
		return nil, err
	}
	commentReplies := make([]domain.CommentReply, 0, len(comment.Replies))
	for _, reply := range comment.Replies {
		commentReplies = append(commentReplies, domain.CommentReply{
			ReplyId:         reply.ReplyId,
			Content:         reply.Content,
			ReplyToId:       reply.ReplyToId,
			UserInfo:        domain.UserInfo4Reply(reply.UserInfo),
			RepliedUserInfo: domain.UserInfo4Reply(reply.RepliedUserInfo),
			ApprovalStatus:  reply.ApprovalStatus,
		})
	}
	return &domain.CommentWithReplies{
		Comment: domain.Comment{
			PostInfo: domain.PostInfo{
				PostId:    comment.PostInfo.PostId,
				PostTitle: comment.PostInfo.PostTitle,
			},
			Content: comment.Content,
			UserInfo: domain.UserInfo{
				Name:    comment.UserInfo.Name,
				Email:   comment.UserInfo.Email,
				Ip:      comment.UserInfo.Ip,
				Website: comment.UserInfo.Website,
			},
		},
		Replies: commentReplies,
	}, nil
}

func (r *CommentRepository) AddComment(ctx context.Context, comment domain.Comment) (string, error) {
	return r.dao.AddComment(ctx, &dao.Comment{
		PostInfo:       dao.PostInfo(comment.PostInfo),
		Content:        comment.Content,
		UserInfo:       dao.UserInfo4Comment(comment.UserInfo),
		Replies:        make([]dao.Reply, 0),
		ApprovalStatus: false,
	})
}

func (r *CommentRepository) toDomainComments(comments []*dao.Comment) []domain.CommentWithReplies {
	result := make([]domain.CommentWithReplies, 0, len(comments))
	for _, comment := range comments {
		replies := make([]domain.CommentReply, 0, len(comment.Replies))
		for _, commentReply := range comment.Replies {
			replies = append(replies, domain.CommentReply{
				ReplyId:         commentReply.ReplyId,
				Content:         commentReply.Content,
				ReplyToId:       commentReply.ReplyToId,
				UserInfo:        domain.UserInfo4Reply(commentReply.UserInfo),
				RepliedUserInfo: domain.UserInfo4Reply(commentReply.RepliedUserInfo),
				ApprovalStatus:  commentReply.ApprovalStatus,
				CreatedAt:       commentReply.CreatedAt.Unix(),
			})
		}
		result = append(result, domain.CommentWithReplies{
			Comment: domain.Comment{
				Id:         comment.ID.Hex(),
				PostInfo:   domain.PostInfo(comment.PostInfo),
				Content:    comment.Content,
				UserInfo:   domain.UserInfo(comment.UserInfo),
				CreateTime: comment.CreatedAt.Unix(),
			},
			Replies: replies,
		})
	}
	return result
}

func (r *CommentRepository) toDomainAdminComments(friends []dao.AdminComment) []domain.AdminComment {
	result := make([]domain.AdminComment, 0, len(friends))
	for _, friend := range friends {
		result = append(result, domain.AdminComment{
			Id:             friend.Id,
			PostInfo:       domain.PostInfo(friend.PostInfo),
			Content:        friend.Content,
			UserInfo:       domain.UserInfo4Comment(friend.UserInfo),
			ApprovalStatus: friend.ApprovalStatus,
			CreatedAt:      friend.CreateTime,
		})
	}
	return result
}

func (r *CommentRepository) toDomainAdminCommentsV2(comments []*dao.Comment) []domain.AdminComment {
	result := make([]domain.AdminComment, 0, len(comments))
	for _, comment := range comments {
		replies := make([]domain.AdminReply, 0, len(comment.Replies))
		for _, commentReply := range comment.Replies {
			replies = append(replies, domain.AdminReply{
				ReplyId:         commentReply.ReplyId,
				Content:         commentReply.Content,
				ReplyToId:       commentReply.ReplyToId,
				UserInfo:        domain.UserInfo4Reply(commentReply.UserInfo),
				RepliedUserInfo: domain.UserInfo4Reply(commentReply.RepliedUserInfo),
				ApprovalStatus:  commentReply.ApprovalStatus,
				CreatedAt:       commentReply.CreatedAt.Unix(),
				UpdatedAt:       commentReply.UpdatedAt.Unix(),
			})
		}
		result = append(result, domain.AdminComment{
			Id:             comment.ID.Hex(),
			PostInfo:       domain.PostInfo(comment.PostInfo),
			Content:        comment.Content,
			UserInfo:       domain.UserInfo4Comment(comment.UserInfo),
			Replies:        replies,
			ApprovalStatus: comment.ApprovalStatus,
			CreatedAt:      comment.CreatedAt.Unix(),
			UpdatedAt:      comment.UpdatedAt.Unix(),
		})
	}
	return result
}

func (r *CommentRepository) toDomainComment(comment *dao.Comment) *domain.Comment {
	return &domain.Comment{
		Id:             comment.ID.Hex(),
		PostInfo:       domain.PostInfo(comment.PostInfo),
		Content:        comment.Content,
		UserInfo:       domain.UserInfo(comment.UserInfo),
		ApprovalStatus: comment.ApprovalStatus,
		CreateTime:     comment.CreatedAt.Unix(),
	}
}

func (r *CommentRepository) toDomainCommentWithReplies(comment *dao.Comment) *domain.CommentWithReplies {
	replies := make([]domain.CommentReply, 0, len(comment.Replies))
	for _, commentReply := range comment.Replies {
		replies = append(replies, domain.CommentReply{
			ReplyId:         commentReply.ReplyId,
			Content:         commentReply.Content,
			ReplyToId:       commentReply.ReplyToId,
			UserInfo:        domain.UserInfo4Reply(commentReply.UserInfo),
			RepliedUserInfo: domain.UserInfo4Reply(commentReply.RepliedUserInfo),
			ApprovalStatus:  commentReply.ApprovalStatus,
			CreatedAt:       commentReply.CreatedAt.Unix(),
		})
	}
	return &domain.CommentWithReplies{
		Comment: domain.Comment{
			Id:             comment.ID.Hex(),
			PostInfo:       domain.PostInfo(comment.PostInfo),
			Content:        comment.Content,
			UserInfo:       domain.UserInfo(comment.UserInfo),
			ApprovalStatus: comment.ApprovalStatus,
			CreateTime:     comment.CreatedAt.Unix(),
		},
		Replies: replies,
	}
}
