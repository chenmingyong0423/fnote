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
	"github.com/chenmingyong0423/fnote/backend/ineternal/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/post/repository/dao"
	"github.com/pkg/errors"
)

type IPostRepository interface {
	GetLatest5Posts(ctx context.Context) ([]*domain.Post, error)
}

var _ IPostRepository = (*PostRepository)(nil)

func NewPostRepository(dao dao.IPostDao) *PostRepository {
	return &PostRepository{
		dao: dao,
	}
}

type PostRepository struct {
	dao dao.IPostDao
}

func (r *PostRepository) GetLatest5Posts(ctx context.Context) ([]*domain.Post, error) {
	posts, err := r.dao.GetLatest5Posts(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.GetLatest5Posts failed")
	}
	return r.toDomainPosts(posts), nil
}
func (r *PostRepository) toDomainPosts(posts []*dao.Post) []*domain.Post {
	result := make([]*domain.Post, 0, len(posts))
	for _, post := range posts {
		result = append(result, r.daoPostToDomainPost(post))
	}
	return result
}

func (r *PostRepository) daoPostToDomainPost(post *dao.Post) *domain.Post {
	return &domain.Post{BasePost: domain.BasePost{Sug: post.Sug, Author: post.Author, Title: post.Title, Summary: post.Summary, CoverImg: post.CoverImg, Category: post.Category, Tags: post.Tags, LikeCount: post.LikeCount, Comments: post.Comments, Visits: post.Visits, Priority: post.Priority, CreateTime: post.CreateTime}}
}
