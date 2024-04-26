// Copyright 2024 chenmingyong0423

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

	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type IPostLikeService interface {
	Add(ctx context.Context, postLike domain.PostLike) (string, error)
	DeleteById(ctx context.Context, id string) error
	GetLikeStatus(ctx context.Context, id string, ip string) (bool, error)
}

var _ IPostLikeService = (*PostLikeService)(nil)

func NewPostLikeService(repo repository.IPostLikeRepository) *PostLikeService {
	return &PostLikeService{
		repo: repo,
	}
}

type PostLikeService struct {
	repo repository.IPostLikeRepository
}

func (s *PostLikeService) GetLikeStatus(ctx context.Context, postId string, ip string) (bool, error) {
	postLike, err := s.repo.FindByPostIdAndIp(ctx, postId, ip)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return false, err
	}
	return postLike != nil, nil
}

func (s *PostLikeService) Add(ctx context.Context, postLike domain.PostLike) (string, error) {
	return s.repo.Add(ctx, postLike)
}

func (s *PostLikeService) DeleteById(ctx context.Context, id string) error {
	return s.repo.DeleteById(ctx, id)
}
