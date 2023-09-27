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
	"net/http"

	"github.com/chenmingyong0423/fnote/backend/internal/comment/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICommentService interface {
	AddComment(ctx context.Context, comment domain.Comment) error
	AddCommentReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) error
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
	// 默认查找最新的前 5 条，后续可能考虑动态配置
	cmts, err := s.repo.FineLatestCommentAndReply(ctx, 5)
	if err != nil {
		return nil, err
	}
	return cmts, err
}

func (s *CommentService) AddCommentReply(ctx context.Context, cmtId string, postId string, commentReply domain.CommentReply) error {
	commentWithReplies, err := s.repo.FindApprovedCommentById(ctx, cmtId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return api.NewHttpCodeError(http.StatusBadRequest)
		}
		return err
	}
	if commentWithReplies.PostInfo.PostId != postId {
		return api.NewHttpCodeError(http.StatusBadRequest)
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
			return errors.New("the replyToId does not exist.")
		}
	}
	return s.repo.AddCommentReply(ctx, cmtId, commentReply)
}

func (s *CommentService) AddComment(ctx context.Context, comment domain.Comment) error {
	_, err := s.repo.AddComment(ctx, comment)
	return err
}