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
	"time"

	"github.com/chenmingyong0423/fnote/backend/internal/comment/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"

	"github.com/google/uuid"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment) (any, error)
	FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error)
	AddCommentReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) error
	FineLatestCommentAndReply(ctx context.Context, cnt int) ([]domain.LatestComment, error)
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

func (r *CommentRepository) FineLatestCommentAndReply(ctx context.Context, cnt int) ([]domain.LatestComment, error) {
	latestComments, err := r.dao.FineLatestCommentAndReply(ctx, cnt)
	if err != nil {
		return nil, err
	}
	result := make([]domain.LatestComment, 0, len(latestComments))
	for _, latestComment := range latestComments {
		result = append(result, domain.LatestComment{
			PostInfo4Comment: domain.PostInfo4Comment(latestComment.PostInfo4Comment),
			Name:             latestComment.Name,
			Content:          latestComment.Content,
			CreateTime:       latestComment.CreateTime,
		})
	}
	return result, nil
}

func (r *CommentRepository) AddCommentReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) error {
	unix := time.Now().Unix()
	return r.dao.AddCommentReply(ctx, cmtId, dao.CommentReply{
		ReplyId:         uuid.NewString(),
		Content:         commentReply.Content,
		ReplyToId:       commentReply.ReplyToId,
		UserInfo:        dao.UserInfo4Reply(commentReply.UserInfo),
		RepliedUserInfo: dao.UserInfo4Reply(commentReply.RepliedUserInfo),
		Status:          dao.CommentStatusPending,
		CreateTime:      unix,
		UpdateTime:      unix,
	})
}

func (r *CommentRepository) FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error) {
	comment, err := r.dao.FindCommentById(ctx, cmtId)
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

			Status: domain.CommentStatus(reply.Status),
		})
	}
	return &domain.CommentWithReplies{
		Comment: domain.Comment{
			PostInfo: domain.PostInfo4Comment{
				PostId:    comment.PostInfo.PostId,
				PostTitle: comment.PostInfo.PostTitle,
			},
			Content: comment.Content,
			UserInfo: domain.UserInfo4Comment{
				Name:    comment.UserInfo.Name,
				Email:   comment.UserInfo.Email,
				Ip:      comment.UserInfo.Ip,
				Website: comment.UserInfo.Website,
			},
		},
		Replies: commentReplies,
	}, nil
}

func (r *CommentRepository) AddComment(ctx context.Context, comment domain.Comment) (any, error) {
	unix := time.Now().Unix()
	return r.dao.AddComment(ctx, dao.Comment{
		Id:         uuid.NewString(),
		PostInfo:   dao.PostInfo4Comment(comment.PostInfo),
		Content:    comment.Content,
		UserInfo:   dao.UserInfo4Comment(comment.UserInfo),
		Replies:    make([]dao.CommentReply, 0),
		Status:     dao.CommentStatusPending,
		CreateTime: unix,
		UpdateTime: unix,
	})
}
