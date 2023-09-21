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
	"log/slog"
	"sync"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/post/repository"
	"github.com/pkg/errors"
)

type IPostService interface {
	GetHomePosts(ctx context.Context) (api.ListVO[*domain.SummaryPostVO], error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) (*api.PageVO[*domain.SummaryPostVO], error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	AddLike(ctx context.Context, id string, ip string) error
	DeleteLike(ctx context.Context, id string, ip string) error
	IncreaseVisitCount(ctx context.Context, id string) error
	InternalGetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
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

func (s *PostService) InternalGetPunishedPostById(ctx context.Context, id string) (*domain.Post, error) {
	return s.repo.GetPunishedPostById(ctx, id)
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
		err := s.repo.DeleteLike(ctx, id, ip)
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
		err := s.repo.AddLike(ctx, id, ip)
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
			slog.WarnContext(ctx, "post", err)
		}
	}()

	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, pageRequest *domain.PostRequest) (*api.PageVO[*domain.SummaryPostVO], error) {
	pageVO := &api.PageVO[*domain.SummaryPostVO]{Page: pageRequest.Page}

	posts, cnt, err := s.repo.QueryPostsPage(ctx, domain.PostsQueryCondition{Size: pageRequest.PageSize, Skip: (pageRequest.PageNo - 1) * pageRequest.PageSize, Search: pageRequest.Search, Sorting: api.Sorting{
		Filed: pageRequest.Sorting.Filed,
		Order: pageRequest.Sorting.Order,
	}, Category: pageRequest.Category, Tag: pageRequest.Tag})
	if err != nil {
		return pageVO, errors.WithMessage(err, "s.repo.QueryPostsPage failed")
	}

	pageVO.List = s.postsToPostVOs(posts)
	pageVO.SetTotalCountAndCalculateTotalPages(cnt)

	return pageVO, nil
}

func (s *PostService) GetHomePosts(ctx context.Context) (api.ListVO[*domain.SummaryPostVO], error) {
	listVO := api.ListVO[*domain.SummaryPostVO]{}
	posts, err := s.repo.GetLatest5Posts(ctx)
	if err != nil {
		return listVO, err
	}
	listVO.List = s.postsToPostVOs(posts)
	return listVO, nil
}

func (s *PostService) postsToPostVOs(posts []*domain.Post) []*domain.SummaryPostVO {
	postVOs := make([]*domain.SummaryPostVO, 0, len(posts))
	for _, post := range posts {
		postVOs = append(postVOs, &domain.SummaryPostVO{PrimaryPost: post.PrimaryPost})
	}
	return postVOs
}
