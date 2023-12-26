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

	"github.com/chenmingyong0423/fnote/backend/internal/count_stats/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/tag/repository"
	"github.com/chenmingyong0423/gkit/slice"
)

func NewTagService(repo repository.ITagRepository, countStatsService service.ICountStatsService) *TagService {
	return &TagService{
		repo:              repo,
		countStatsService: countStatsService,
	}
}

var _ ITagService = (*TagService)(nil)

type ITagService interface {
	GetTags(ctx context.Context) ([]domain.TagWithCount, error)
}

type TagService struct {
	repo              repository.ITagRepository
	countStatsService service.ICountStatsService
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
	tagCounts, err := s.countStatsService.GetByReferenceIdAndType(ctx, ids, domain.CountStatsTypePostCountInTag)
	if err != nil {
		return nil, err
	}

	tagWithCounts := slice.Map[domain.CountStats, domain.TagWithCount](tagCounts, func(_ int, s domain.CountStats) domain.TagWithCount {
		return domain.TagWithCount{
			Name:  tagMap[s.ReferenceId].Name,
			Count: s.Count,
		}
	})
	return tagWithCounts, nil
}
