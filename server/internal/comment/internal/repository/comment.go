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
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chenmingyong0423/go-mongox/types"

	"github.com/chenmingyong0423/go-mongox/builder/aggregation"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository/dao"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/google/uuid"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment) (string, error)
	FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error)
	AddReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) (string, error)
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]domain.LatestComment, error)
	FindCommentsByPostIdAndCmtStatus(ctx context.Context, postId string) ([]domain.CommentWithReplies, error)
	FindCommentById(ctx context.Context, id string) (*domain.Comment, error)
	UpdateCommentStatus2True(ctx context.Context, id string) error
	FindReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) (*domain.CommentReplyWithPostInfo, error)
	UpdateCommentReplyStatus(ctx context.Context, commentId string, replyId string, approvalStatus bool) error
	FindCommentWithRepliesById(ctx context.Context, id string) (*domain.CommentWithReplies, error)
	DeleteCommentById(ctx context.Context, id string) error
	DeleteReplyByCIdAndRId(ctx context.Context, commentId string, replyId string) error
	CountOfToday(ctx context.Context) (int64, error)
	AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error)
	UpdateCommentStatus2TrueByIds(ctx context.Context, ids []primitive.ObjectID) error
	FindCommentByObjectIDs(ctx context.Context, ids []primitive.ObjectID) ([]domain.AdminComment, error)
	UpdateCReplyStatus2TrueByCidAndRIds(ctx context.Context, commentId string, replyIds []string) error
	FindCommentWithDisapprovedReplyByCidAndRIds(ctx context.Context, commentId string, replyIds []string) (*domain.AdminComment, error)
	FindDisapprovedCommentByObjectIDs(ctx context.Context, commentObjectIDs []primitive.ObjectID) ([]domain.AdminComment, error)
	DeleteCommentByIds(ctx context.Context, ids []primitive.ObjectID) error
	PullReplyByCIdAndRIds(ctx context.Context, commentId string, replyIds []string) error
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

func (r *CommentRepository) PullReplyByCIdAndRIds(ctx context.Context, commentId string, replyIds []string) error {
	objectID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}
	return r.dao.PullReplyByCIdAndRIds(ctx, objectID, replyIds)
}

func (r *CommentRepository) DeleteCommentByIds(ctx context.Context, ids []primitive.ObjectID) error {
	return r.dao.DeleteByIds(ctx, ids)
}

func (r *CommentRepository) FindDisapprovedCommentByObjectIDs(ctx context.Context, commentObjectIDs []primitive.ObjectID) ([]domain.AdminComment, error) {
	comments, err := r.dao.FindDisapprovedCommentByObjectIDs(ctx, commentObjectIDs)
	if err != nil {
		return nil, err
	}
	return r.toDomainAdminComments(comments), nil
}

func (r *CommentRepository) FindCommentWithDisapprovedReplyByCidAndRIds(ctx context.Context, commentId string, replyIds []string) (*domain.AdminComment, error) {
	commentObjID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return nil, err
	}
	comment, err := r.dao.FindWithDisapprovedReplyByCidAndRIds(ctx, commentObjID, replyIds)
	if err != nil {
		return nil, err
	}
	adminComment := r.toDomainAdminComment(comment)
	return &adminComment, nil
}

func (r *CommentRepository) UpdateCReplyStatus2TrueByCidAndRIds(ctx context.Context, commentId string, replyIds []string) error {
	objectID, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return err
	}
	return r.dao.UpdateCReplyStatus2TrueByCidAndRIds(ctx, objectID, replyIds)
}

func (r *CommentRepository) FindCommentByObjectIDs(ctx context.Context, ids []primitive.ObjectID) ([]domain.AdminComment, error) {
	comments, err := r.dao.FindByObjectIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return r.toDomainAdminComments(comments), nil
}

func (r *CommentRepository) UpdateCommentStatus2TrueByIds(ctx context.Context, ids []primitive.ObjectID) error {
	return r.dao.UpdateCommentStatus2TrueByIds(ctx, ids)
}

func (r *CommentRepository) AdminFindCommentsWithPagination(ctx context.Context, page domain.Page) ([]domain.AdminComment, int64, error) {
	var (
		comments []*dao.Comment
		total    int64
		err      error
	)
	sortBson := page.SortToBson()

	if page.ApprovalStatus == nil {
		comments, total, err = r.find(ctx, page, sortBson)
		if err != nil {
			return nil, 0, err
		}
	} else {
		// 如果是根据关键字 approvalStatus 查询的话，使用聚合操作
		approvalStatus := *page.ApprovalStatus
		pipelineBuilder := aggregation.StageBsonBuilder()
		if approvalStatus {
			pipelineBuilder.Match(bsonx.M("approval_status", true)).
				Project(
					aggregation.BsonBuilder().
						AddKeyValues("content", 1).
						AddKeyValues("post_info", 1).
						AddKeyValues("user_info", 1).
						AddKeyValues("approval_status", 1).
						AddKeyValues("created_at", 1).
						Filter("replies", "$replies", aggregation.EqWithoutKey("$$reply.approval_status", true), &types.FilterOptions{As: "reply"}).Build(),
				)
		} else {
			pipelineBuilder.Match(
				aggregation.OrWithoutKey(
					bsonx.M("approval_status", false),
					bsonx.NewD().
						Add("approval_status", true).
						Add("replies.approval_status", false).
						Build()),
			).Project(
				aggregation.BsonBuilder().
					AddKeyValues("content", 1).
					AddKeyValues("post_info", 1).
					AddKeyValues("user_info", 1).
					AddKeyValues("approval_status", 1).
					AddKeyValues("created_at", 1).
					Cond("replies",
						aggregation.EqWithoutKey("$approval_status", true),
						aggregation.FilterWithoutKey("$replies", aggregation.EqWithoutKey("$$reply.approval_status", false), &types.FilterOptions{As: "reply"}),
						"$replies",
					).Build(),
			)
		}
		if len(sortBson) > 0 {
			pipelineBuilder.Sort(sortBson)
		}
		pipelineBuilder.Skip(page.Skip).Limit(page.Size)
		comments, err = r.dao.FindByAggregation(ctx, pipelineBuilder.Build())
		if err != nil {
			return nil, 0, err
		}
	}
	// 对 replies 进行排序
	for _, comment := range comments {
		if len(comment.Replies) > 0 {
			sort.Slice(comment.Replies, func(i, j int) bool {
				return comment.Replies[i].CreatedAt.After(comment.Replies[j].CreatedAt)
			})
		}
	}
	return r.toDomainAdminComments(comments), total, nil
}

func (r *CommentRepository) find(ctx context.Context, page domain.Page, sortBson bson.D) ([]*dao.Comment, int64, error) {
	findOptions := options.Find()
	findOptions.SetSkip(page.Skip).SetLimit(page.Size)
	if len(sortBson) > 0 {
		findOptions.SetSort(sortBson)
	} else {
		findOptions.SetSort(bsonx.M("created_at", -1))
	}
	comments, total, err := r.dao.Find(ctx, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
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

func (r *CommentRepository) toDomainAdminComment(comment *dao.Comment) domain.AdminComment {
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
	return domain.AdminComment{
		Id:             comment.ID.Hex(),
		PostInfo:       domain.PostInfo(comment.PostInfo),
		Content:        comment.Content,
		UserInfo:       domain.UserInfo4Comment(comment.UserInfo),
		Replies:        replies,
		ApprovalStatus: comment.ApprovalStatus,
		CreatedAt:      comment.CreatedAt.Unix(),
		UpdatedAt:      comment.UpdatedAt.Unix(),
	}
}

func (r *CommentRepository) toDomainAdminComments(comments []*dao.Comment) []domain.AdminComment {
	result := make([]domain.AdminComment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, r.toDomainAdminComment(comment))
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
