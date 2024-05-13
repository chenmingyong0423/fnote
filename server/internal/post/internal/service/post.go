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
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"

	"github.com/chenmingyong0423/go-eventbus"

	domain2 "github.com/chenmingyong0423/fnote/server/internal/post/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/post/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	service3 "github.com/chenmingyong0423/fnote/server/internal/file/service"
	"github.com/chenmingyong0423/gkit/slice"

	service2 "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
)

type IPostService interface {
	GetLatestPosts(ctx context.Context, count int64) ([]*domain.Post, error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	IncreaseVisitCount(ctx context.Context, id string) error
	AdminGetPosts(ctx context.Context, page domain2.Page) ([]*domain.Post, int64, error)
	AddPost(ctx context.Context, post *domain.Post) error
	DeletePost(ctx context.Context, id string) error
	DecreaseCommentCount(ctx context.Context, postId string, cnt int) error
	AdminGetPostById(ctx context.Context, id string) (*domain.Post, error)
	UpdatePostIsDisplayed(ctx context.Context, id string, isDisplayed bool) error
	UpdatePostIsCommentAllowed(ctx context.Context, id string, isCommentAllowed bool) error
	SavePost(ctx context.Context, originalPost *domain.Post, savedPost *domain.Post, isNewPost bool) error
	IncreasePostLikeCount(ctx context.Context, postId string) error
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository, cfgService website_config.Service, countStats service2.ICountStatsService, fileService service3.IFileService, eventBus *eventbus.EventBus) *PostService {
	return &PostService{
		repo:        repo,
		cfgService:  cfgService,
		countStats:  countStats,
		fileService: fileService,
		eventBus:    eventBus,
	}
}

type PostService struct {
	repo        repository.IPostRepository
	cfgService  website_config.Service
	countStats  service2.ICountStatsService
	fileService service3.IFileService
	eventBus    *eventbus.EventBus
}

func (s *PostService) IncreasePostLikeCount(ctx context.Context, postId string) error {
	return s.repo.IncreasePostLikeCount(ctx, postId)
}

func (s *PostService) UpdatePostIsCommentAllowed(ctx context.Context, id string, isCommentAllowed bool) error {
	return s.repo.UpdatePostIsCommentAllowedById(ctx, id, isCommentAllowed)
}

func (s *PostService) UpdatePostIsDisplayed(ctx context.Context, id string, isDisplayed bool) error {
	return s.repo.UpdatePostIsDisplayedById(ctx, id, isDisplayed)
}

func (s *PostService) SavePost(ctx context.Context, originalPost *domain.Post, savedPost *domain.Post, isNewPost bool) error {
	err := s.repo.SavePost(ctx, savedPost)
	if err != nil {
		return err
	}
	go func() {
		if isNewPost {
			s.addPostCallback(ctx, savedPost)
		} else {
			s.updatePostCallback(ctx, originalPost, savedPost)
		}
	}()
	return nil
}

func (s *PostService) AdminGetPostById(ctx context.Context, id string) (*domain.Post, error) {
	return s.repo.FindPostById(ctx, id)
}

func (s *PostService) DecreaseCommentCount(ctx context.Context, postId string, cnt int) error {
	return s.repo.DecreaseCommentCount(ctx, postId, cnt)
}

func (s *PostService) DeletePost(ctx context.Context, id string) error {
	post, err := s.repo.FindPostById(ctx, id)
	if err != nil {
		return err
	}
	err = s.repo.DeletePost(ctx, id)
	if err != nil {
		return err
	}
	go func() {
		// 网站文章数-1
		gErr := s.countStats.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCountInWebsite.ToString(), domain.CountStatsTypePostCountInWebsite, 1)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 对应的分类和标签文章数-1
		categoryIds := slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
		if len(categoryIds) > 0 {
			gErr = s.countStats.DecreaseByReferenceIdsAndType(ctx, categoryIds, domain.CountStatsTypePostCountInCategory)
			if gErr != nil {
				l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
				l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
			}
		}

		tagIds := slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
		if len(tagIds) > 0 {
			gErr = s.countStats.DecreaseByReferenceIdsAndType(ctx, tagIds, domain.CountStatsTypePostCountInTag)
			if gErr != nil {
				l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
				l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
			}
		}
		// 删除封面文件索引
		fid, gErr := hex.DecodeString(strings.Split(post.CoverImg[1:], ".")[0])
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = s.fileService.DeleteIndexFileMeta(ctx, fid, id, "post")
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// todo 评论等数据待删除
	}()
	return nil
}

func (s *PostService) AddPost(ctx context.Context, post *domain.Post) error {
	if post.Id == "" {
		post.Id = uuidx.RearrangeUUID4()
	}
	err := s.repo.AddPost(ctx, post)
	if err != nil {
		return err
	}
	go func() {
		s.addPostCallback(ctx, post)
	}()
	return nil
}

// addPostCallback 添加文章后的回调
func (s *PostService) addPostCallback(ctx context.Context, post *domain.Post) {
	// 网站文章数+1
	gErr := s.countStats.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCountInWebsite.ToString(), domain.CountStatsTypePostCountInWebsite)
	if gErr != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
	}
	// 对应的分类和标签文章数+1
	categoryIds := slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
		return c.Id
	})
	if len(categoryIds) > 0 {
		gErr = s.countStats.IncreaseByReferenceIdsAndType(ctx, categoryIds, domain.CountStatsTypePostCountInCategory)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}
	tagIds := slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
		return t.Id
	})
	if len(tagIds) > 0 {
		gErr = s.countStats.IncreaseByReferenceIdsAndType(ctx, tagIds, domain.CountStatsTypePostCountInTag)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}
	// 封面文件索引
	fileId, gErr := hex.DecodeString(strings.Split(post.CoverImg[8:], ".")[0])
	if gErr != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
	}
	gErr = s.fileService.IndexFileMeta(ctx, fileId, post.Id, "post")
	if gErr != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
	}
}

func (s *PostService) AdminGetPosts(ctx context.Context, page domain2.Page) ([]*domain.Post, int64, error) {
	return s.repo.QueryAdminPostsPage(ctx, page)
}

func (s *PostService) IncreaseVisitCount(ctx context.Context, id string) error {
	return s.repo.IncreaseCommentCount(ctx, id)
}

func (s *PostService) GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error) {
	post, err := s.repo.GetPunishedPostById(ctx, id)
	if err != nil {
		return nil, err
	}
	// increase visits
	go func() {
		gErr := s.repo.IncreaseVisitCount(ctx, post.Id)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return post, nil
}

func (s *PostService) GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error) {
	return s.repo.QueryPostsPage(ctx, domain.PostsQueryCondition{Size: pageRequest.PageSize, Skip: (pageRequest.PageNo - 1) * pageRequest.PageSize, Keyword: pageRequest.Keyword, Sorting: api.Sorting{
		Field: pageRequest.Sorting.Field,
		Order: pageRequest.Sorting.Order,
	}, Categories: pageRequest.Categories, Tags: pageRequest.Tags})

}

func (s *PostService) GetLatestPosts(ctx context.Context, count int64) ([]*domain.Post, error) {
	return s.repo.GetLatest5Posts(ctx, count)
}

func (s *PostService) updatePostCallback(ctx context.Context, post *domain.Post, savedPost *domain.Post) {
	// categories
	// 被移除的
	removedCategories := slice.DiffFunc(post.Categories, savedPost.Categories, func(srcItem, dstItem domain.Category4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	if len(removedCategories) > 0 {
		removedCategoryIds := slice.Map[domain.Category4Post, string](removedCategories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
		err := s.countStats.DecreaseByReferenceIdsAndType(ctx, removedCategoryIds, domain.CountStatsTypePostCountInCategory)
		if err != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", err))
		}
	}
	// 被添加的
	addedCategories := slice.DiffFunc(savedPost.Categories, post.Categories, func(srcItem, dstItem domain.Category4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	if len(addedCategories) > 0 {
		addedCategoryIds := slice.Map[domain.Category4Post, string](addedCategories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
		err := s.countStats.IncreaseByReferenceIdsAndType(ctx, addedCategoryIds, domain.CountStatsTypePostCountInCategory)
		if err != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", err))
		}
	}
	// tags
	// 被移除的
	removedTags := slice.DiffFunc(post.Tags, savedPost.Tags, func(srcItem, dstItem domain.Tag4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	if len(removedTags) > 0 {
		removedTagIds := slice.Map[domain.Tag4Post, string](removedTags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
		err := s.countStats.DecreaseByReferenceIdsAndType(ctx, removedTagIds, domain.CountStatsTypePostCountInTag)
		if err != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", err))
		}
	}
	// 被添加的
	addedTags := slice.DiffFunc(savedPost.Tags, post.Tags, func(srcItem, dstItem domain.Tag4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	if len(addedTags) > 0 {
		addedTagIds := slice.Map[domain.Tag4Post, string](addedTags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
		err := s.countStats.IncreaseByReferenceIdsAndType(ctx, addedTagIds, domain.CountStatsTypePostCountInTag)
		if err != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", err))
		}
	}
}
