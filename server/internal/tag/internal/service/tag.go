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
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/repository"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewTagService(repo repository.ITagRepository, eventBus *eventbus.EventBus) *TagService {
	s := &TagService{
		repo:     repo,
		eventBus: eventBus,
	}
	go s.subscribePostEvent()
	return s
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
	repo     repository.ITagRepository
	eventBus *eventbus.EventBus
}

func (s *TagService) GetSelectTags(ctx context.Context) ([]domain.Tag, error) {
	return s.repo.GetSelectTags(ctx)
}

func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	_, err := s.repo.GetTagById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "tag not found")
		}
		return err
	}
	tagEvent := domain.TagEvent{
		TagId: id,
		Type:  "delete",
	}
	marshal, err := jsoniter.Marshal(tagEvent)
	if err != nil {
		return err
	}
	err = s.repo.DeleteTagById(ctx, id)
	if err != nil {
		return err
	}
	s.eventBus.Publish("tag", eventbus.Event{Payload: marshal})
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
	tagEvent := domain.TagEvent{
		TagId: id,
		Type:  "create",
	}
	marshal, err := jsoniter.Marshal(tagEvent)
	if err != nil {
		l := slog.Default().With("X-Request-ID", ctx.(*gin.Context).GetString("X-Request-ID"))
		l.WarnContext(ctx, "AdminCreateTag: tag event: failed to jsoniter.Marshal", "error", err)
		return nil
	}
	s.eventBus.Publish("tag", eventbus.Event{Payload: marshal})
	return nil
}

func (s *TagService) AdminGetTags(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error) {
	tags, total, err := s.QueryTagsPage(ctx, pageDTO)
	if err != nil {
		return nil, 0, err
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
	return slice.Map(tags, func(_ int, t domain.Tag) domain.TagWithCount {
		return domain.TagWithCount{
			Name:  t.Name,
			Route: t.Route,
			Count: t.PostCount,
		}
	}), nil
}

func (s *TagService) QueryTagsPage(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Tag, int64, error) {
	return s.repo.QueryTagsPage(ctx, pageDTO)
}

func (s *TagService) subscribePostEvent() {
	eventChan := s.eventBus.Subscribe("post")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "Tag: post event", "payload", string(event.Payload))
		var e domain.PostEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "Tag: post event: failed to json.Unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			// 对应标签的文章数量 +1
			if len(e.AddedTagId) > 0 {
				err = s.repo.IncreasePostCountByIds(ctx, e.AddedTagId)
				if err != nil {
					l.ErrorContext(ctx, "Tag: post event: failed to increase the count of post in tag", "error", err)
					continue
				}
			}
		case "delete":
			if len(e.DeletedTagId) > 0 {
				// 对应标签的文章数量 -1
				err = s.repo.DecreasePostCountByIds(ctx, e.DeletedTagId)
				if err != nil {
					l.ErrorContext(ctx, "Tag: post event: failed to decrease the count of post in tag", "error", err)
					continue
				}
			}
		case "update":
			if len(e.AddedTagId) > 0 {
				err = s.repo.IncreasePostCountByIds(ctx, e.AddedTagId)
				if err != nil {
					l.ErrorContext(ctx, "Tag: post event: failed to increase the count of post in tag", "error", err)
					if len(e.DeletedTagId) == 0 {
						continue
					}
				}
			}
			if len(e.DeletedTagId) > 0 {
				err = s.repo.DecreasePostCountByIds(ctx, e.DeletedTagId)
				if err != nil {
					l.ErrorContext(ctx, "Tag: post event: failed to decrease the count of post in tag", "error", err)
					continue
				}
			}
		}
		l.InfoContext(ctx, "Tag: post event: handle successfully")
	}
}
