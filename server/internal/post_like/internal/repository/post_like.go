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

package repository

import (
	"context"

	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/repository/dao"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPostLikeRepository interface {
	Add(ctx context.Context, postLike domain.PostLike) (string, error)
	DeleteById(ctx context.Context, id string) error
	FindByPostIdAndIp(ctx context.Context, postId string, ip string) (*domain.PostLike, error)
	CountOfToday(ctx context.Context) (int64, error)
}

var _ IPostLikeRepository = (*PostLikeRepository)(nil)

func NewPostLikeRepository(dao dao.IPostLikeDao) *PostLikeRepository {
	return &PostLikeRepository{dao: dao}
}

type PostLikeRepository struct {
	dao dao.IPostLikeDao
}

func (r *PostLikeRepository) CountOfToday(ctx context.Context) (int64, error) {
	return r.dao.CountOfToday(ctx)
}

func (r *PostLikeRepository) FindByPostIdAndIp(ctx context.Context, postId string, ip string) (*domain.PostLike, error) {
	postLike, err := r.dao.FindByPostIdAndIp(ctx, postId, ip)
	if err != nil {
		return nil, err
	}
	return r.toDomain(postLike), nil
}

func (r *PostLikeRepository) DeleteById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *PostLikeRepository) Add(ctx context.Context, postLike domain.PostLike) (string, error) {
	return r.dao.Add(ctx, &dao.PostLike{
		PostId:    postLike.PostId,
		Ip:        postLike.Ip,
		UserAgent: postLike.UserAgent,
	})
}

func (r *PostLikeRepository) toDomain(postLike *dao.PostLike) *domain.PostLike {
	return &domain.PostLike{
		Id:        postLike.ID.Hex(),
		PostId:    postLike.PostId,
		Ip:        postLike.Ip,
		UserAgent: postLike.UserAgent,
		CreatedAt: postLike.CreatedAt.UnixMilli(),
	}
}
