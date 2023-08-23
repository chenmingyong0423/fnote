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
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/result"
	"github.com/chenmingyong0423/fnote/backend/ineternal/post/repository"
)

type IPostService interface {
	GetHomePosts(ctx context.Context) (result.ListVO[*domain.PostVO], error)
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

func (s *PostService) GetHomePosts(ctx context.Context) (result.ListVO[*domain.PostVO], error) {
	listVO := result.ListVO[*domain.PostVO]{}
	posts, err := s.repo.GetLatest5Posts(ctx)
	if err != nil {
		return listVO, err
	}
	listVO.List = s.postsToPostVOs(posts)
	return listVO, nil
}

func (s *PostService) postsToPostVOs(posts []*domain.Post) []*domain.PostVO {
	postVOs := make([]*domain.PostVO, 0, len(posts))
	for _, post := range posts {
		postVOs = append(postVOs, &domain.PostVO{BasePost: post.BasePost})
	}
	return postVOs
}
