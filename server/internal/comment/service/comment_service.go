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

	"github.com/chenmingyong0423/fnote/server/internal/comment/repository"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/pkg/errors"
)

type ICommentService interface {
	AddComment(ctx context.Context, comment domain.Comment) (string, error)
	AddCommentReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) (string, error)
	FineLatestCommentAndReply(ctx context.Context) ([]domain.LatestComment, error)
	FindCommentsByPostId(ctx context.Context, postId string) ([]domain.CommentWithReplies, error)
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

func (s *CommentService) FindCommentsByPostId(ctx context.Context, postId string) ([]domain.CommentWithReplies, error) {
	return s.repo.FindCommentsByPostIdAndCmtStatus(ctx, postId, domain.CommentStatusApproved)
}

func (s *CommentService) FineLatestCommentAndReply(ctx context.Context) ([]domain.LatestComment, error) {
	// todo 默认查找最新的前 5 条，后续可能考虑动态配置
	return s.repo.FineLatestCommentAndReply(ctx, 5)
}

func (s *CommentService) AddCommentReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) (string, error) {
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
			if reply.ReplyId == commentReply.ReplyToId && reply.Status == domain.CommentStatusApproved {
				commentReply.RepliedUserInfo.Name, commentReply.RepliedUserInfo.Email, commentReply.RepliedUserInfo.Website, commentReply.RepliedUserInfo.Ip = reply.UserInfo.Name, reply.UserInfo.Email, reply.UserInfo.Website, reply.UserInfo.Ip
				isExist = true
				break
			}
		}
		if !isExist {
			return "", errors.New("The replyToId does not exist.")
		}
	}
	return s.repo.AddCommentReply(ctx, cmtId, commentReply)
}

func (s *CommentService) AddComment(ctx context.Context, comment domain.Comment) (string, error) {
	return s.repo.AddComment(ctx, comment)
}
