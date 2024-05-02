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
	"fmt"
	"log/slog"
	"net/http"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/fnote/server/internal/tag/repository"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewTagService(repo repository.ITagRepository, countStatsService service.ICountStatsService) *TagService {
	return &TagService{
		repo:              repo,
		countStatsService: countStatsService,
	}
}

type ITagService interface {
	GetTags(ctx context.Context) ([]domain.TagWithCount, error)
	GetTagByRoute(ctx context.Context, route string) (domain.Tag, error)
	AdminGetTags(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error)
	AdminCreateTag(ctx context.Context, tag domain.Tag) error
	ModifyTagEnabled(ctx context.Context, id string, enabled bool) error
	DeleteTag(ctx context.Context, id string) error
	GetSelectTags(ctx context.Context) ([]domain.Tag, error)
}

var _ ITagService = (*TagService)(nil)

type TagService struct {
	repo              repository.ITagRepository
	countStatsService service.ICountStatsService
}

func (s *TagService) GetSelectTags(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.GetSelectTags(ctx)
}

func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	tag, err := s.repo.GetTagById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "tag not found")
		}
		return err
	}
	err = s.repo.DeleteTagById(ctx, id)
	if err != nil {
		return err
	}
	// 删除分类时，同时删除该分类下文章数量的统计数据
	err = s.countStatsService.DeleteByReferenceIdAndType(ctx, id, domain.CountStatsTypePostCountInTag)
	if err != nil {
		nErr := s.repo.RecoverTag(ctx, tag)
		if nErr != nil {
			return nErr
		}
		return err
	}
	go func() {
		// 分类数量 -1
		gErr := s.countStatsService.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeTagCount.ToString(), domain.CountStatsTypeTagCount, 1)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil
}

func (s *TagService) ModifyTagEnabled(ctx context.Context, id string, enabled bool) error {
	return s.repo.ModifyTagEnabled(ctx, id, enabled)
}

func (s *TagService) AdminCreateTag(ctx context.Context, tag domain.Tag) error {
	id, err := s.repo.CreateTag(ctx, tag)
	if err != nil {
		return err
	}
	// 创建标签时，同时创建标签的统计数据
	err = s.countStatsService.Create(ctx, domain.CountStats{
		Type:        domain.CountStatsTypePostCountInTag,
		ReferenceId: id,
	})
	if err != nil {
		gErr := s.DeleteTag(ctx, id)
		if gErr != nil {
			return gErr
		}
		return err
	}
	go func() {
		// 分类数量 +1
		gErr := s.countStatsService.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeTagCount.ToString(), domain.CountStatsTypeTagCount)
		if gErr != nil {
			l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
			l.WarnContext(ctx, fmt.Sprintf("%+v", gErr))
		}
	}()
	return nil
}

func (s *TagService) AdminGetTags(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error) {
	tags, total, err := s.QueryTagsPage(ctx, pageDTO)
	if err != nil {
		return nil, 0, err
	}
	// 查询标签的统计数据
	if len(tags) > 0 {
		ids := slice.Map[domain.Tag, string](tags, func(_ int, s domain.Tag) string {
			return s.Id
		})
		countStats, fErr := s.countStatsService.GetByReferenceIdsAndType(ctx, ids, domain.CountStatsTypePostCountInTag)
		if err != nil {
			return nil, 0, fErr
		}
		countStatsMap := slice.IndexStructsByKey[domain.CountStats, string](countStats, func(countStats domain.CountStats) string {
			return countStats.ReferenceId
		})
		for i, tag := range tags {
			if cs, ok := countStatsMap[tag.Id]; ok {
				tags[i].PostCount = cs.Count
			}
		}
	}
	return tags, total, nil
}

func (s *TagService) GetTagByRoute(ctx context.Context, route string) (domain.Tag, error) {
	return s.repo.GetTagByRoute(ctx, route)
}

func (s *TagService) GetTags(ctx context.Context) ([]domain.TagWithCount, error) {
	tags, err := s.repo.GetTags(ctx)
	if err != nil {
		return nil, err
	}
	if len(tags) == 0 {
		return nil, nil
	}
	ids := slice.Map[domain.Tag, string](tags, func(_ int, s domain.Tag) string {
		return s.Id
	})
	tagMap := slice.IndexStructsByKey[domain.Tag, string](tags, func(tag domain.Tag) string {
		return tag.Id
	})
	tagCounts, err := s.countStatsService.GetByReferenceIdsAndType(ctx, ids, domain.CountStatsTypePostCountInTag)
	if err != nil {
		return nil, err
	}

	tagWithCounts := slice.Map[domain.CountStats, domain.TagWithCount](tagCounts, func(_ int, s domain.CountStats) domain.TagWithCount {
		return domain.TagWithCount{
			Name:  tagMap[s.ReferenceId].Name,
			Route: tagMap[s.ReferenceId].Route,
			Count: s.Count,
		}
	})
	return tagWithCounts, nil
}

func (s *TagService) QueryTagsPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error) {
	return s.repo.QueryTagsPage(ctx, pageDTO)
}
