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
	"fmt"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/post/repository/dao"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
)

type IPostRepository interface {
	GetLatest5Posts(ctx context.Context) ([]*domain.Post, error)
	QueryPostsPage(ctx context.Context, postsQueryCondition domain.PostsQueryCondition) ([]*domain.Post, int64, error)
	GetPostById(ctx context.Context, id string) (*domain.Post, error)
	IncreaseVisitCount(ctx context.Context, id string) error
	HadLikePost(ctx context.Context, id string, ip string) (bool, error)
	AddLike(ctx context.Context, id string, ip string) error
	DeleteLike(ctx context.Context, id string, ip string) error
	IncreaseCommentCount(ctx context.Context, id string) error
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

func (r *PostRepository) IncreaseCommentCount(ctx context.Context, id string) error {
	return r.dao.IncreaseFieldById(ctx, id, "comment_count")
}

func (r *PostRepository) DeleteLike(ctx context.Context, id string, ip string) error {
	err := r.dao.DeleteLike(ctx, id, ip)
	if err != nil {
		return errors.WithMessage(err, "r.dao.DeleteLike failed")
	}
	return nil
}

func (r *PostRepository) AddLike(ctx context.Context, id string, ip string) error {
	err := r.dao.AddLike(ctx, id, ip)
	if err != nil {
		return errors.WithMessage(err, "r.dao.AddLike failed")
	}
	return nil
}

func (r *PostRepository) HadLikePost(ctx context.Context, id string, ip string) (bool, error) {
	_, err := r.dao.FindByIdAndIp(ctx, id, ip)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, errors.WithMessage(err, "r.dao.FindByIdAndIp")
	}
	return true, nil
}

func (r *PostRepository) IncreaseVisitCount(ctx context.Context, id string) error {
	return r.dao.IncreaseFieldById(ctx, id, "visit_count")
}

func (r *PostRepository) GetPostById(ctx context.Context, id string) (*domain.Post, error) {
	post, err := r.dao.GetPostById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.GetPostById failed")
	}
	return r.daoPostToDomainPost(post), nil
}

func (r *PostRepository) QueryPostsPage(ctx context.Context, postsQueryCondition domain.PostsQueryCondition) ([]*domain.Post, int64, error) {
	con := bson.D{}
	if postsQueryCondition.Category != nil {
		con = append(con, bson.E{Key: "category", Value: *postsQueryCondition.Category})
	}
	if postsQueryCondition.Tag != nil {
		con = append(con, bson.E{Key: "tags", Value: *postsQueryCondition.Tag})
	}
	if postsQueryCondition.Search != nil {
		con = append(con, bson.E{Key: "title", Value: primitive.Regex{
			Pattern: fmt.Sprintf(".*%s.*", strings.TrimSpace(*postsQueryCondition.Search)),
		}})
	}

	findOptions := options.Find()
	findOptions.SetSkip(postsQueryCondition.Skip).SetLimit(postsQueryCondition.Size)
	if postsQueryCondition.Sort != nil {
		findOptions.SetSort(bson.E{Key: postsQueryCondition.Sort.Filed, Value: orderConvertToInt(postsQueryCondition.Sort.Order)})
	} else {
		findOptions.SetSort(bson.E{Key: "create_time", Value: -1})
	}

	posts, cnt, err := r.dao.QueryPostsPage(ctx, con, findOptions)
	if err != nil {
		return nil, 0, errors.WithMessage(err, "r.dao.QueryPostsPage failed")
	}
	return r.toDomainPosts(posts), cnt, nil
}

func orderConvertToInt(order string) int {
	switch order {
	case "ASC":
		return 1
	case "DESC":
		return -1
	default:
		return -1
	}
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
	return &domain.Post{PrimaryPost: domain.PrimaryPost{Sug: post.Sug, Author: post.Author, Title: post.Title, Summary: post.Summary, CoverImg: post.CoverImg, Category: post.Category, Tags: post.Tags, LikeCount: post.LikeCount, CommentCount: post.CommentCount, VisitCount: post.VisitCount, Priority: post.Priority, CreateTime: post.CreateTime}, ExtraPost: domain.ExtraPost{Content: post.Content, MetaDescription: post.MetaDescription, MetaKeywords: post.MetaKeywords, UpdateTime: post.UpdateTime}, IsCommentAllowed: post.IsCommentAllowed, Likes: post.Likes}
}
