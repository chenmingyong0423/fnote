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
	"log/slog"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/domain"
	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/repository"
)

type ICountStatsService interface {
	GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error)
}

var _ ICountStatsService = (*CountStatsService)(nil)

func NewCountStatsService(repo repository.ICountStatsRepository, eventbus *eventbus.EventBus) *CountStatsService {
	s := &CountStatsService{
		repo:     repo,
		eventBus: eventbus,
	}
	go s.subscribePostEvent()
	go s.subscribePostLikedEvent()
	go s.subscribeCategoryEvent()
	go s.subscribeCommentEvent()
	go s.subscribeWebsiteVisitEvent()
	go s.subscribeTagEvent()
	return s
}

type CountStatsService struct {
	repo     repository.ICountStatsRepository
	eventBus *eventbus.EventBus
}

func (s *CountStatsService) GetWebsiteCountStats(ctx context.Context) (domain.WebsiteCountStats, error) {
	var result = new(domain.WebsiteCountStats)
	countStatsSlice, err := s.repo.GetWebsiteCountStats(ctx, []domain.CountStatsType{
		domain.CountStatsTypePostCount,
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

func (s *CountStatsService) subscribePostLikedEvent() {
	eventChan := s.eventBus.Subscribe("post-like")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "CountStats post-like event", "payload", string(event.Payload))
		var postEvent domain.LikePostEvent
		err := jsoniter.Unmarshal(event.Payload, &postEvent)
		if err != nil {
			l.ErrorContext(ctx, "CountStats post-like event: failed to unmarshal", "error", err)
			continue
		}

		// 点赞数+1
		err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeLikeCount, 1)
		if err != nil {
			l.ErrorContext(ctx, "CountStats post-like event: failed to increase the count of like in website", "count", 1, "error", err)
			continue
		}
		l.InfoContext(ctx, "CountStats post-like event: handle successfully")
	}
}

func (s *CountStatsService) subscribeCommentEvent() {
	eventChan := s.eventBus.Subscribe("comment")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "CountStats: comment event", "payload", string(event.Payload))
		var e domain.CommentEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "CountStats: comment event: failed to unmarshal", "error", err)
			continue
		}

		switch e.Type {
		case "create":
			err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount, e.Count)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: comment event: failed to increase the count of comment", "count", e.Count, "error", err)
				continue
			}
		case "delete":
			err = s.repo.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount, e.Count)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: comment event: failed to decrease the count of comment", "count", e.Count, "error", err)
				continue
			}
		}
		l.InfoContext(ctx, "CountStats: comment event: handle successfully ")
	}
}

func (s *CountStatsService) subscribeWebsiteVisitEvent() {
	eventChan := s.eventBus.Subscribe("website visit")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "CountStats: website visit event", "payload", string(event.Payload))
		var e domain.CommentEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "CountStats: website visit event: failed to unmarshal", "error", err)
			continue
		}
		err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeWebsiteViewCount, 1)
		if err != nil {
			l.ErrorContext(ctx, "CountStats: website visit event: failed to increase the count of website visit", "count", 1, "error", err)
			continue
		}
		l.InfoContext(ctx, "CountStats: website visit event: handle successfully")
	}
}

func (s *CountStatsService) subscribeTagEvent() {
	eventChan := s.eventBus.Subscribe("tag")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With(slog.Any("X-Request-ID", rid))
		l.InfoContext(ctx, "CountStats: tag event", "payload", string(event.Payload))
		var e domain.TagEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "CountStats: tag event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeTagCount, 1)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: tag event: failed to increase the count of tag", "count", 1, "error", err)
				continue
			}
		case "delete":
			err = s.repo.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeTagCount, 1)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: tag event: failed to decrease the count of tag", "count", 1, "error", err)
				continue
			}
		}
		l.InfoContext(ctx, "CountStats: tag event: handle successfully")
	}
}

func (s *CountStatsService) subscribePostEvent() {
	eventChan := s.eventBus.Subscribe("post")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "CountStats: post", "payload", string(event.Payload))
		var e domain.PostEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "CountStats: post event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			{
				// 网站文章数 +1
				err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCount, 1)
				if err != nil {
					l.ErrorContext(ctx, "CountStats: post event: failed to increase the count of post", "count", 1, "error", err)
				}
			}
		case "update":
		case "delete":
			// 网站文章数 -1
			err = s.repo.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypePostCount, 1)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: post event: failed to decrease the count of post", "count", 1, "error", err)
				continue
			}
			// 删除评论数
			err = s.repo.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCommentCount, e.CommentCount)
			if err != nil {
				l.ErrorContext(ctx, "CountStats: post event: failed to decrease the count of comment", "count", e.CommentCount, "error", err)
				continue
			}

		}
		l.InfoContext(ctx, "CountStats: post: handle successfully")
	}
}

func (s *CountStatsService) subscribeCategoryEvent() {
	eventChan := s.eventBus.Subscribe("category")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "CountStats category event", "payload", string(event.Payload))
		var e domain.CategoryEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "CountStats category event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			err = s.repo.IncreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCategoryCount, 1)
			if err != nil {
				l.ErrorContext(ctx, "CountStats category event: failed to increase the count of category", "count", 1, "error", err)
				continue
			}
		case "delete":
			err = s.repo.DecreaseByReferenceIdAndType(ctx, domain.CountStatsTypeCategoryCount, 1)
			if err != nil {
				l.ErrorContext(ctx, "CountStats category event: failed to decrease the count of category", "count", 1, "error", err)
				continue
			}
		}
		l.InfoContext(ctx, "CountStats category event: handle successfully")
	}
}
