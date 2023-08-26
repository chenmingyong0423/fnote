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
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/post/repository"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log/slog"
	"slices"
)

type IPostService interface {
	GetHomePosts(ctx context.Context) (api.ListVO[*domain.SummaryPostVO], error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) (*api.PageVO[*domain.SummaryPostVO], error)
	GetPostBySug(ctx context.Context, sug string) (*domain.DetailPostVO, error)
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository) *PostService {
	return &PostService{
		repo: repo,
	}
}

type PostService struct {
	repo repository.IPostRepository
}

func (s *PostService) GetPostBySug(ctx context.Context, sug string) (*domain.DetailPostVO, error) {
	post, err := s.repo.GetPostBySug(ctx, sug)
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.GetPostBySug failed")
	}

	// increase visits
	go func() {
		gErr := s.repo.IncreaseVisits(ctx, post.Sug)
		if gErr != nil {
			slog.WarnContext(ctx, "post", err)
		}
	}()
	postVO := &domain.DetailPostVO{PrimaryPost: post.PrimaryPost, ExtraPost: post.ExtraPost}
	postVO.IsLiked = slices.Contains(post.Likes, ctx.(*gin.Context).ClientIP())

	return postVO, nil
}

func (s *PostService) GetPosts(ctx context.Context, pageRequest *domain.PostRequest) (*api.PageVO[*domain.SummaryPostVO], error) {
	pageVO := &api.PageVO[*domain.SummaryPostVO]{Page: pageRequest.Page}

	posts, cnt, err := s.repo.QueryPostsPage(ctx, domain.PostsQueryCondition{Size: pageRequest.PageSize, Skip: (pageRequest.PageNo - 1) * pageRequest.PageSize, Search: pageRequest.Search, Sort: pageRequest.Sort, Category: pageRequest.Category, Tag: pageRequest.Tag})
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
