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
	"fmt"
	"log/slog"
	"sync"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/post/repository"
	"github.com/pkg/errors"
)

type IPostService interface {
	GetHomePosts(ctx context.Context) ([]*domain.Post, error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	AddLike(ctx context.Context, id string, ip string) error
	DeleteLike(ctx context.Context, id string, ip string) error
	IncreaseVisitCount(ctx context.Context, id string) error
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

type PostService struct {
	repo  repository.IPostRepository
	ipMap sync.Map
}

func (s *PostService) IncreaseVisitCount(ctx context.Context, id string) error {
	return s.repo.IncreaseCommentCount(ctx, id)
}

func (s *PostService) DeleteLike(ctx context.Context, id string, ip string) error {
	// 先判断是否已经点过赞
	had, err := s.repo.HadLikePost(ctx, id, ip)
	if err != nil {
		return errors.WithMessage(err, "s.repo.HadLikePost failed")
	}
	if !had {
		return nil
	}
	_, isExist := s.ipMap.LoadOrStore(ip, struct{}{})
	if !isExist {
		defer s.ipMap.Delete(ip)
		err = s.repo.DeleteLike(ctx, id, ip)
		if err != nil {
			return errors.WithMessage(err, "s.repo.DeleteLike")
		}
	}
	return nil
}

func (s *PostService) AddLike(ctx context.Context, id string, ip string) error {
	// 先判断是否已经点过赞
	had, err := s.repo.HadLikePost(ctx, id, ip)
	if err != nil {
		return errors.WithMessage(err, "s.repo.HadLikePost failed")
	}
	if had {
		return nil
	}
	_, isExist := s.ipMap.LoadOrStore(ip, struct{}{})
	if !isExist {
		defer s.ipMap.Delete(ip)
		err = s.repo.AddLike(ctx, id, ip)
		if err != nil {
			return errors.WithMessage(err, "s.repo.AddLike")
		}
	}
	return nil
}

func (s *PostService) GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.repo.GetPunishedPostById(ctx, id)
	if err != nil {
		return nil, err
	}
	// increase visits
	go func() {
		gErr := s.repo.IncreaseVisitCount(ctx, post.Sug)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error) {
	return s.repo.QueryPostsPage(ctx, domain.PostsQueryCondition{Size: pageRequest.PageSize, Skip: (pageRequest.PageNo - 1) * pageRequest.PageSize, Search: pageRequest.Search, Sorting: api.Sorting{
		Field: pageRequest.Sorting.Field,
		Order: pageRequest.Sorting.Order,
	}, Category: pageRequest.Category, Tags: pageRequest.Tags})

}

func (s *PostService) GetHomePosts(ctx context.Context) ([]*domain.Post, error) {
	return s.repo.GetLatest5Posts(ctx)
}
