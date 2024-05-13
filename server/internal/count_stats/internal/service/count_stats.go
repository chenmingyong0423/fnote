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
	"encoding/json"
	"log/slog"

	"github.com/google/uuid"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/domain"
	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/repository"
)

type ICountStatsService interface {
	GetByReferenceIdsAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error)
	Create(ctx context.Context, countStats domain.CountStats) error
	DeleteByReferenceIdAndType(ctx context.Context, referenceId string, statsType domain.CountStatsType) error
	DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error
	DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType, count int) error
	IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error
	GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error)
}

var _ ICountStatsService = (*CountStatsService)(nil)

func NewCountStatsService(repo repository.ICountStatsRepository, eventbus *eventbus.EventBus) *CountStatsService {
	s := &CountStatsService{
		repo:     repo,
		eventBus: eventbus,
	}
	go s.SubscribePostDeletedEvent()
	go s.SubscribePostAddedEvent()
	go s.SubscribePostUpdatedEvent()
	go s.SubscribePostLikedEvent()
	return s
}

type CountStatsService struct {
	repo     repository.ICountStatsRepository
	eventBus *eventbus.EventBus
}

func (s *CountStatsService) GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error) {
	var result = new(domain.WebsiteCountStats)
	countStatsSlice, err := s.repo.GetWebsiteCountStats(ctx, []domain.CountStatsType{
		domain.CountStatsTypePostCountInWebsite,
		domain.CountStatsTypeCategoryCount,
		domain.CountStatsTypeTagCount,
		domain.CountStatsTypeCommentCount,
		domain.CountStatsTypeLikeCount,
		domain.CountStatsTypeWebsiteViewCount,
	})
	if err != nil {
		return *result, err
	}
	for _, countStats := range countStatsSlice {
		result.SetCountByType(countStats.Type, countStats.Count)
	}
	return *result, nil
}

func (s *CountStatsService) IncreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType) error {
	return s.repo.IncreaseByReferenceIdAndType(ctx, referenceId, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIdAndType(ctx context.Context, referenceId string, countStatsType domain.CountStatsType, count int) error {
	return s.repo.DecreaseByReferenceIdAndType(ctx, referenceId, countStatsType, count)
}

func (s *CountStatsService) IncreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.IncreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DecreaseByReferenceIdsAndType(ctx context.Context, ids []string, countStatsType domain.CountStatsType) error {
	return s.repo.DecreaseByReferenceIdsAndType(ctx, ids, countStatsType)
}

func (s *CountStatsService) DeleteByReferenceIdAndType(ctx context.Context, referenceId string, statsType domain.CountStatsType) error {
	return s.repo.DeleteByReferenceIdAndType(ctx, referenceId, statsType)
}

func (s *CountStatsService) Create(ctx context.Context, countStats domain.CountStats) error {
	err := countStats.Type.Valid()
	if err != nil {
		return err
	}
	_, err = s.repo.Create(ctx, countStats)
	if err != nil {
		return err
	}
	return nil
}

func (s *CountStatsService) GetByReferenceIdsAndType(ctx context.Context, referenceIds []string, countStatsType domain.CountStatsType) ([]domain.CountStats, error) {
	return s.repo.GetByReferenceIdAndType(ctx, referenceIds, countStatsType)
}

func (s *CountStatsService) SubscribePostDeletedEvent() {
	eventChan := s.eventBus.Subscribe("post-delete")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-delete", "payload", string(event.Payload))
		var postEvent domain.PostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-delete: failed to json.Unmarshal", "err", err)
			continue
		}
		// todo 后面可以考虑使用事务
		{
			// 网站文章数 -1
			err = s.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCountInWebsite.ToString(), domain.CountStatsTypePostCountInWebsite, 1)
			if err != nil {
				l.ErrorContext(ctx, "post-delete: failed to decrease the count of post in website", "count", 1, "err", err)
				continue
			}
			// 文章对应的分类和标签文章数 -1
			if len(postEvent.CategoryId) > 0 {
				err = s.DecreaseByReferenceIdsAndType(ctx, postEvent.CategoryId, domain.CountStatsTypePostCountInCategory)
				if err != nil {
					l.ErrorContext(ctx, "post-delete: failed to decrease the count of post in category", "count", 1, "err", err)
				}
			}
			if len(postEvent.TagId) > 0 {
				err = s.DecreaseByReferenceIdsAndType(ctx, postEvent.TagId, domain.CountStatsTypePostCountInTag)
				if err != nil {
					l.ErrorContext(ctx, "post-delete: failed to decrease the count of post in tag", "count", 1, "err", err)
				}
			}
		}
		l.InfoContext(ctx, "post-delete: successfully decrease the count of post in website, category and tag", "count", 1)
	}
}

func (s *CountStatsService) SubscribePostAddedEvent() {
	eventChan := s.eventBus.Subscribe("post-addition")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-addition", "payload", string(event.Payload))
		var postEvent domain.PostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-addition: failed to json.Unmarshal", "err", err)
			continue
		}
		// todo 后面可以考虑使用事务
		{
			// 网站文章数 +1
			err = s.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCountInWebsite.ToString(), domain.CountStatsTypePostCountInWebsite)
			if err != nil {
				l.ErrorContext(ctx, "post-addition: failed to increase the count of post in website", "count", 1, "err", err)
				continue
			}
			// 文章对应的分类和标签文章数 +1
			if len(postEvent.CategoryId) > 0 {
				err = s.IncreaseByReferenceIdsAndType(ctx, postEvent.CategoryId, domain.CountStatsTypePostCountInCategory)
				if err != nil {
					l.ErrorContext(ctx, "post-addition: failed to increase the count of post in category", "count", 1, "err", err)
				}
			}
			if len(postEvent.TagId) > 0 {
				err = s.IncreaseByReferenceIdsAndType(ctx, postEvent.TagId, domain.CountStatsTypePostCountInTag)
				if err != nil {
					l.ErrorContext(ctx, "post-addition: failed to increase the count of post in tag", "count", 1, "err", err)
				}
			}
		}
		l.InfoContext(ctx, "post-addition: successfully increase the count of post in website, category and tag", "count", 1)
	}
}

func (s *CountStatsService) SubscribePostUpdatedEvent() {
	eventChan := s.eventBus.Subscribe("post-update")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-update", "payload", string(event.Payload))
		var postEvent domain.UpdatedPostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-update: failed to json.Unmarshal", "err", err)
			continue
		}
		// todo 后面可以考虑使用事务
		{
			// 文章对应的分类和标签文章数 +1
			if len(postEvent.AddedCategoryId) > 0 {
				err = s.IncreaseByReferenceIdsAndType(ctx, postEvent.AddedCategoryId, domain.CountStatsTypePostCountInCategory)
				if err != nil {
					l.ErrorContext(ctx, "post-update: failed to increase the count of post in category", "categoryId", postEvent.AddedCategoryId, "count", 1, "err", err)
				}
			}
			if len(postEvent.DeletedCategoryId) > 0 {
				err = s.DecreaseByReferenceIdsAndType(ctx, postEvent.DeletedCategoryId, domain.CountStatsTypePostCountInCategory)
				if err != nil {
					l.ErrorContext(ctx, "post-update: failed to decrease the count of post in category", "categoryId", postEvent.DeletedCategoryId, "count", 1, "err", err)
				}
			}
			if len(postEvent.AddedTagId) > 0 {
				err = s.IncreaseByReferenceIdsAndType(ctx, postEvent.AddedTagId, domain.CountStatsTypePostCountInTag)
				if err != nil {
					l.ErrorContext(ctx, "post-update: failed to increase the count of post in tag", "tagId", postEvent.AddedTagId, "count", 1, "err", err)
				}
			}
			if len(postEvent.DeletedTagId) > 0 {
				err = s.DecreaseByReferenceIdsAndType(ctx, postEvent.AddedTagId, domain.CountStatsTypePostCountInTag)
				if err != nil {
					l.ErrorContext(ctx, "post-update: failed to decrease the count of post in tag", "tagId", postEvent.DeletedTagId, "count", 1, "err", err)
				}
			}
		}
		l.InfoContext(ctx, "post-update: successfully update the count of post in website, category and tag", "postEvent", postEvent)
	}
}

func (s *CountStatsService) SubscribePostLikedEvent() {
	eventChan := s.eventBus.Subscribe("post-like")
	for event := range eventChan {
		rid := uuid.NewString()
		ctx := context.WithValue(context.Background(), "X-Request-ID", rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "post-like", "payload", string(event.Payload))
		var postEvent domain.LikePostEvent
		err := json.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "post-like: failed to json.Unmarshal", "err", err)
			continue
		}

		// 点赞数+1
		err = s.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeLikeCount.ToString(), domain.CountStatsTypeLikeCount)
		if err != nil {
			l.ErrorContext(ctx, "post-like: failed to increase the count of like in website", "count", 1, "err", err)
			continue
		}
		l.InfoContext(ctx, "post-like: successfully increase the count of like in website", "count", 1)
	}
}
