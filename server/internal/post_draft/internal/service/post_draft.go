// Copyright 2024 chenmingyong0423

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
	"errors"
	"regexp"

	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository"
	"github.com/chenmingyong0423/gkit/uuidx"
	"github.com/chenmingyong0423/go-eventbus"
	jsoniter "github.com/json-iterator/go"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IPostDraftService interface {
	SavePostDraft(ctx context.Context, postDraft domain.PostDraft) (string, error)
	GetPostDraftById(ctx context.Context, id string) (*domain.PostDraft, error)
	DeletePostDraftById(ctx context.Context, id string) (int64, error)
	GetPostDraftPage(ctx context.Context, page domain.Page) ([]*domain.PostDraft, int64, error)
}

var _ IPostDraftService = (*PostDraftService)(nil)

func NewPostDraftService(repo repository.IPostDraftRepository, eventBus *eventbus.EventBus) *PostDraftService {
	return &PostDraftService{
		repo:     repo,
		eventBus: eventBus,
	}
}

type PostDraftService struct {
	repo     repository.IPostDraftRepository
	eventBus *eventbus.EventBus
}

var staticFileIDPattern = regexp.MustCompile(`/static/([[:xdigit:]]+)(?:\.[^/?#)\s]+)?`)

func (s *PostDraftService) GetPostDraftPage(ctx context.Context, page domain.Page) ([]*domain.PostDraft, int64, error) {
	return s.repo.GetPostDraftPage(ctx, domain.PageQuery{
		Size:    page.PageSize,
		Skip:    (page.PageNo - 1) * page.PageSize,
		Keyword: page.Keyword,
		Field:   page.Field,
		Order:   page.OrderConvertToInt(),
	})
}

func (s *PostDraftService) DeletePostDraftById(ctx context.Context, id string) (int64, error) {
	draft, err := s.repo.GetById(ctx, id)
	if err != nil {
		return 0, err
	}
	count, err := s.repo.DeleteById(ctx, id)
	if err != nil || count == 0 {
		return count, err
	}
	s.publishFileUsageEvent(id, nil, staticFileIDs(draft.CoverImg, draft.Content))
	return count, nil
}

func (s *PostDraftService) GetPostDraftById(ctx context.Context, id string) (*domain.PostDraft, error) {
	return s.repo.GetById(ctx, id)
}

func (s *PostDraftService) SavePostDraft(ctx context.Context, postDraft domain.PostDraft) (string, error) {
	var oldFileIDs []string
	if postDraft.Id == "" {
		postDraft.Id = uuidx.RearrangeUUID4()
	} else {
		existing, err := s.repo.GetById(ctx, postDraft.Id)
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return "", err
		}
		if existing != nil {
			oldFileIDs = staticFileIDs(existing.CoverImg, existing.Content)
		}
	}
	id, err := s.repo.Save(ctx, postDraft)
	if err != nil {
		return "", err
	}
	newFileIDs := staticFileIDs(postDraft.CoverImg, postDraft.Content)
	s.publishFileUsageEvent(id, diffFileIDs(newFileIDs, oldFileIDs), diffFileIDs(oldFileIDs, newFileIDs))
	return id, nil
}

func (s *PostDraftService) publishFileUsageEvent(entityID string, added, deleted []string) {
	payload, err := jsoniter.Marshal(domain.FileUsageEvent{EntityId: entityID, AddedFileIds: added, DeletedFileIds: deleted})
	if err == nil {
		s.eventBus.Publish("post-draft", eventbus.Event{Payload: payload})
	}
}

func staticFileIDs(values ...string) []string {
	seen := make(map[string]struct{})
	result := make([]string, 0)
	for _, value := range values {
		for _, match := range staticFileIDPattern.FindAllStringSubmatch(value, -1) {
			if _, ok := seen[match[1]]; !ok {
				seen[match[1]] = struct{}{}
				result = append(result, match[1])
			}
		}
	}
	return result
}

func diffFileIDs(source, target []string) []string {
	targetSet := make(map[string]struct{}, len(target))
	for _, id := range target {
		targetSet[id] = struct{}{}
	}
	result := make([]string, 0)
	for _, id := range source {
		if _, ok := targetSet[id]; !ok {
			result = append(result, id)
		}
	}
	return result
}
