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
	"sync"

	"github.com/chenmingyong0423/gkit/slice"

	service3 "github.com/chenmingyong0423/fnote/server/internal/file/service"

	service2 "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/service"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post/repository"
	"github.com/pkg/errors"
)

type IPostService interface {
	GetLatestPosts(ctx context.Context, count int64) ([]*domain.Post, error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	AddLike(ctx context.Context, id string, ip string) error
	DeleteLike(ctx context.Context, id string, ip string) error
	IncreaseVisitCount(ctx context.Context, id string) error
	AdminGetPosts(ctx context.Context, pageDTO dto.PageDTO) ([]*domain.Post, int64, error)
	AddPost(ctx context.Context, post domain.Post) error
	DeletePost(ctx context.Context, id string) error
}

var _ IPostService = (*PostService)(nil)

func NewPostService(repo repository.IPostRepository, cfgService service.IWebsiteConfigService, countStats service2.ICountStatsService, fileService service3.IFileService) *PostService {
	return &PostService{
		repo:        repo,
		cfgService:  cfgService,
		countStats:  countStats,
		fileService: fileService,
	}
}

type PostService struct {
	repo        repository.IPostRepository
	cfgService  service.IWebsiteConfigService
	countStats  service2.ICountStatsService
	fileService service3.IFileService
	ipMap       sync.Map
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
		gErr := s.cfgService.DecreaseWebsitePostCount(ctx)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 对应的分类和标签文章数-1
		referenceIds := make([]string, 0, len(post.Categories)+len(post.Tags))
		categoryIds := slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
		referenceIds = append(referenceIds, categoryIds...)
		tagIds := slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
		referenceIds = append(referenceIds, tagIds...)
		gErr = s.countStats.DecreaseByReferenceIds(ctx, referenceIds)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
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

func (s *PostService) AddPost(ctx context.Context, post domain.Post) error {
	if post.Id == "" {
		post.Id = uuidx.RearrangeUUID4()
	}
	err := s.repo.AddPost(ctx, post)
	if err != nil {
		return err
	}
	go func() {
		// 网站文章数+1
		gErr := s.cfgService.IncreaseWebsitePostCount(ctx)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 对应的分类和标签文章数+1
		referenceIds := make([]string, 0, len(post.Categories)+len(post.Tags))
		categoryIds := slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
		referenceIds = append(referenceIds, categoryIds...)
		tagIds := slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
		referenceIds = append(referenceIds, tagIds...)
		for _, tag := range post.Tags {
			referenceIds = append(referenceIds, tag.Id)
		}
		gErr = s.countStats.IncreaseByReferenceIds(ctx, referenceIds)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		// 封面文件索引
		fileId, gErr := hex.DecodeString(strings.Split(post.CoverImg[1:], ".")[0])
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
		gErr = s.fileService.IndexFileMeta(ctx, fileId, post.Id, "post")
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}

	}()
	return nil
}

func (s *PostService) AdminGetPosts(ctx context.Context, pageDTO dto.PageDTO) ([]*domain.Post, int64, error) {
	return s.repo.QueryAdminPostsPage(ctx, dto.PostsQueryDTO{Size: pageDTO.PageSize, Skip: (pageDTO.PageNo - 1) * pageDTO.PageSize, Keyword: pageDTO.Keyword})
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
