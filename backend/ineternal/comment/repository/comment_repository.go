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
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/types"
	"github.com/google/uuid"
	"time"
)

type ICommentRepository interface {
	AddComment(ctx context.Context, comment domain.Comment) (any, error)
	FindApprovedCommentById(ctx context.Context, cmtId string) (*domain.CommentWithReplies, error)
	AddCommentReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) error
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

func (r *CommentRepository) AddCommentReply(ctx context.Context, cmtId string, commentReply domain.CommentReply) error {
	unix := time.Now().Unix()
	return r.dao.AddCommentReply(ctx, cmtId, dao.CommentReply{
		CommentReply: types.CommentReply{
			ReplyId:         uuid.NewString(),
			Content:         commentReply.Content,
			ReplyToId:       commentReply.ReplyToId,
			UserInfo:        commentReply.UserInfo,
			RepliedUserInfo: commentReply.RepliedUserInfo,
		},
		Status:     domain.CommentStatusPending,
		CreateTime: unix,
		UpdateTime: unix,
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
			CommentReply: types.CommentReply{
				ReplyId:         reply.ReplyId,
				Content:         reply.Content,
				ReplyToId:       reply.ReplyToId,
				UserInfo:        reply.UserInfo,
				RepliedUserInfo: reply.RepliedUserInfo,
			},
			Status: reply.Status,
		})
	}
	return &domain.CommentWithReplies{
		Comment: domain.Comment{
			Comment: types.Comment{
				PostInfo: types.PostInfo4Comment{
					PostId:    comment.PostInfo.PostId,
					PostTitle: comment.PostInfo.PostTitle,
				},
				Content: comment.Content,
				UserInfo: types.UserInfo4Comment{
					Name:    comment.UserInfo.Name,
					Email:   comment.UserInfo.Email,
					Ip:      comment.UserInfo.Ip,
					Website: comment.UserInfo.Website,
				},
			},
		},
		Replies: commentReplies,
	}, nil
}

func (r *CommentRepository) AddComment(ctx context.Context, comment domain.Comment) (any, error) {
	unix := time.Now().Unix()
	return r.dao.AddComment(ctx, dao.Comment{
		Id:         uuid.NewString(),
		Comment:    comment.Comment,
		Replies:    make([]dao.CommentReply, 0),
		Status:     domain.CommentStatusPending,
		CreateTime: unix,
		UpdateTime: unix,
	})
}
