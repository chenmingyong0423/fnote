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
	"strings"
	"time"

	"github.com/chenmingyong0423/gkit/slice"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post/repository/dao"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPostRepository interface {
	GetLatest5Posts(ctx context.Context, count int64) ([]*domain.Post, error)
	QueryPostsPage(ctx context.Context, postsQueryCondition domain.PostsQueryCondition) ([]*domain.Post, int64, error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	IncreaseVisitCount(ctx context.Context, id string) error
	HadLikePost(ctx context.Context, id string, ip string) (bool, error)
	AddLike(ctx context.Context, id string, ip string) error
	DeleteLike(ctx context.Context, id string, ip string) error
	IncreaseCommentCount(ctx context.Context, id string) error
	QueryAdminPostsPage(ctx context.Context, postsQueryDTO dto.PostsQueryDTO) ([]*domain.Post, int64, error)
	AddPost(ctx context.Context, post domain.Post) error
	DeletePost(ctx context.Context, id string) error
	FindPostById(ctx context.Context, id string) (*domain.Post, error)
	DecreaseCommentCount(ctx context.Context, postId string, cnt int) error
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

func (r *PostRepository) DecreaseCommentCount(ctx context.Context, postId string, cnt int) error {
	return r.dao.DecreaseByField(ctx, postId, "comment_count", cnt)
}

func (r *PostRepository) FindPostById(ctx context.Context, id string) (*domain.Post, error) {
	post, err := r.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.daoPostToDomainPost(post), nil
}

func (r *PostRepository) DeletePost(ctx context.Context, id string) error {
	return r.dao.DeleteById(ctx, id)
}

func (r *PostRepository) AddPost(ctx context.Context, post domain.Post) error {
	unix := time.Now().Unix()
	categories := make([]dao.Category4Post, 0, len(post.Categories))
	for _, category := range post.Categories {
		categories = append(categories, dao.Category4Post{
			Id:   category.Id,
			Name: category.Name,
		})
	}
	tags := make([]dao.Tag4Post, 0, len(post.Tags))
	for _, tag := range post.Tags {
		tags = append(tags, dao.Tag4Post{
			Id:   tag.Id,
			Name: tag.Name,
		})
	}
	err := r.dao.AddPost(ctx, &dao.Post{
		Id:               post.Id,
		Author:           post.Author,
		Title:            post.Title,
		Summary:          post.Summary,
		Content:          post.Content,
		CoverImg:         post.CoverImg,
		Categories:       categories,
		Tags:             tags,
		Status:           domain.PostStatus(post.Status),
		Likes:            make([]string, 0),
		LikeCount:        0,
		CommentCount:     0,
		VisitCount:       0,
		StickyWeight:     post.StickyWeight,
		MetaDescription:  post.MetaDescription,
		MetaKeywords:     post.MetaKeywords,
		WordCount:        0,
		IsCommentAllowed: post.IsCommentAllowed,
		CreateTime:       unix,
		UpdateTime:       unix,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *PostRepository) QueryAdminPostsPage(ctx context.Context, postsQueryDTO dto.PostsQueryDTO) ([]*domain.Post, int64, error) {
	condBuilder := query.BsonBuilder()
	if postsQueryDTO.Keyword != "" {
		condBuilder.RegexOptions("title", fmt.Sprintf(".*%s.*", strings.TrimSpace(postsQueryDTO.Keyword)), "i")
	}
	con := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip(postsQueryDTO.Skip).SetLimit(postsQueryDTO.Size)
	if postsQueryDTO.Field != "" && postsQueryDTO.Order != "" {
		findOptions.SetSort(bsonx.M(postsQueryDTO.Field, orderConvertToInt(postsQueryDTO.Order)))
	} else {
		findOptions.SetSort(bsonx.M("create_time", -1))
	}

	posts, cnt, err := r.dao.QueryPostsPage(ctx, con, findOptions)
	if err != nil {
		return nil, 0, err
	}
	return r.toDomainPosts(posts), cnt, nil
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
	return r.dao.AddLike(ctx, id, ip)
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

func (r *PostRepository) GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error) {
	post, err := r.dao.GetPunishedPostById(ctx, id)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.GetPunishedPostById failed")
	}
	return r.daoPostToDomainPost(post), nil
}

func (r *PostRepository) QueryPostsPage(ctx context.Context, postsQueryCondition domain.PostsQueryCondition) ([]*domain.Post, int64, error) {
	condBuilder := query.BsonBuilder().Eq("status", dao.PostStatusPunished)
	if postsQueryCondition.Categories != nil && len(postsQueryCondition.Categories) > 0 {
		condBuilder.Eq("categories.name", postsQueryCondition.Categories[0])
	}
	if postsQueryCondition.Tags != nil && len(postsQueryCondition.Tags) > 0 {
		condBuilder.Eq("tags.name", postsQueryCondition.Tags[0])
	}
	if postsQueryCondition.Keyword != nil && *postsQueryCondition.Keyword != "" {
		condBuilder.RegexOptions("title", fmt.Sprintf(".*%s.*", strings.TrimSpace(*postsQueryCondition.Keyword)), "i")
	}
	con := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip(postsQueryCondition.Skip).SetLimit(postsQueryCondition.Size)
	if postsQueryCondition.Sorting.Field != nil && postsQueryCondition.Sorting.Order != nil {
		findOptions.SetSort(bsonx.M(*postsQueryCondition.Sorting.Field, orderConvertToInt(*postsQueryCondition.Sorting.Order)))
	} else {
		findOptions.SetSort(bsonx.M("create_time", -1))
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

func (r *PostRepository) GetLatest5Posts(ctx context.Context, count int64) ([]*domain.Post, error) {
	posts, err := r.dao.GetFrontPosts(ctx, count)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.GetFrontPosts failed")
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
	categories := slice.Map[dao.Category4Post, domain.Category4Post](post.Categories, func(_ int, c dao.Category4Post) domain.Category4Post {
		return domain.Category4Post{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	tags := slice.Map[dao.Tag4Post, domain.Tag4Post](post.Tags, func(_ int, t dao.Tag4Post) domain.Tag4Post {
		return domain.Tag4Post{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	return &domain.Post{PrimaryPost: domain.PrimaryPost{Id: post.Id, Author: post.Author, Title: post.Title, Summary: post.Summary, CoverImg: post.CoverImg, Categories: categories, Tags: tags, LikeCount: post.LikeCount, CommentCount: post.CommentCount, VisitCount: post.VisitCount, StickyWeight: post.StickyWeight, CreateTime: post.CreateTime}, ExtraPost: domain.ExtraPost{Content: post.Content, MetaDescription: post.MetaDescription, MetaKeywords: post.MetaKeywords, WordCount: post.WordCount, UpdateTime: post.UpdateTime}, IsCommentAllowed: post.IsCommentAllowed, Likes: post.Likes}
}
