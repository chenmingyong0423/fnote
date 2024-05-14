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
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/category/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/repository"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit/slice"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type ICategoryService interface {
	GetCategories(ctx context.Context) ([]domain.CategoryWithCount, error)
	GetMenus(ctx context.Context) ([]domain.Category, error)
	GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error)
	AdminGetCategories(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Category, int64, error)
	AdminCreateCategory(ctx context.Context, category domain.Category) error
	ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error
	ModifyCategory(ctx context.Context, id string, description string) error
	DeleteCategory(ctx context.Context, id string) error
	ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error
	AdminGetSelectCategories(ctx context.Context) ([]domain.Category, error)
}

var _ ICategoryService = (*CategoryService)(nil)

func NewCategoryService(repo repository.ICategoryRepository, eventbus *eventbus.EventBus) *CategoryService {
	s := &CategoryService{
		repo:     repo,
		eventBus: eventbus,
	}
	go s.subscribePostEvent()
	return s
}

type CategoryService struct {
	repo     repository.ICategoryRepository
	eventBus *eventbus.EventBus
}

func (s *CategoryService) AdminGetSelectCategories(ctx context.Context) ([]domain.Category, error) {
	return s.repo.GetSelectCategories(ctx)
}

func (s *CategoryService) ModifyCategoryNavigation(ctx context.Context, id string, showInNav bool) error {
	return s.repo.ModifyCategoryNavigation(ctx, id, showInNav)
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string) error {
	_, err := s.repo.GetCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "category not found")
		}
		return err
	}
	var categoryEvent = domain.CategoryEvent{CategoryId: id, Type: "delete"}
	marshal, err := json.Marshal(&categoryEvent)
	if err != nil {
		return err
	}
	err = s.repo.DeleteCategory(ctx, id)
	if err != nil {
		return err
	}
	go s.eventBus.Publish("category", eventbus.Event{Payload: marshal})
	return nil
}

func (s *CategoryService) ModifyCategory(ctx context.Context, id string, description string) error {
	return s.repo.ModifyCategory(ctx, id, description)
}

func (s *CategoryService) ModifyCategoryEnabled(ctx context.Context, id string, enabled bool) error {
	return s.repo.ModifyCategoryEnabled(ctx, id, enabled)
}

func (s *CategoryService) AdminCreateCategory(ctx context.Context, category domain.Category) error {
	id, err := s.repo.CreateCategory(ctx, category)
	if err != nil {
		return err
	}
	var categoryEvent = domain.CategoryEvent{CategoryId: id, Type: "create"}
	marshal, err := json.Marshal(&categoryEvent)
	if err != nil {
		return err
	}
	if err != nil {
		err2 := s.DeleteCategory(ctx, id)
		if err2 != nil {
			return err2
		}
		return err
	}

	go s.eventBus.Publish("category", eventbus.Event{Payload: marshal})
	return nil
}

func (s *CategoryService) AdminGetCategories(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Category, int64, error) {
	categories, total, err := s.QueryCategoriesPage(ctx, pageDTO)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (s *CategoryService) GetCategoryByRoute(ctx context.Context, route string) (domain.Category, error) {
	return s.repo.GetCategoryByRoute(ctx, route)
}

func (s *CategoryService) GetMenus(ctx context.Context) (menuVO []domain.Category, err error) {
	categories, err := s.repo.GetNavigations(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	return categories, nil
}

func (s *CategoryService) GetCategories(ctx context.Context) ([]domain.CategoryWithCount, error) {
	categories, err := s.repo.GetAll(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNilDocument) {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, nil
	}
	return slice.Map(categories, func(_ int, c domain.Category) domain.CategoryWithCount {
		return domain.CategoryWithCount{
			Name:        c.Name,
			Route:       c.Route,
			Description: c.Description,
			Count:       c.PostCount,
		}
	}), nil
}

func (s *CategoryService) QueryCategoriesPage(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Category, int64, error) {
	return s.repo.QueryCategoriesPage(ctx, pageDTO)
}

func (s *CategoryService) subscribePostEvent() {
	eventChan := s.eventBus.Subscribe("post")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "Category: post event", "payload", string(event.Payload))
		var e domain.PostEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "Category: post event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			// 对应分类的文章数量 +1
			if len(e.AddedCategoryId) > 0 {
				err = s.repo.IncreasePostCountByIds(ctx, e.AddedCategoryId)
				if err != nil {
					l.ErrorContext(ctx, "Category: post event: failed to increase post count", "categoryIds", e.AddedCategoryId, "error", err)
					continue
				}
			}
		case "delete":
			// 分类的文章数量 -1
			if len(e.DeletedCategoryId) > 0 {
				err = s.repo.DecreasePostCountByIds(ctx, e.DeletedCategoryId)
				if err != nil {
					l.ErrorContext(ctx, "Category: post event: failed to decrease post count", "categoryIds", e.DeletedCategoryId, "error", err)
					continue
				}
			}
		case "update":
			// 对应分类的文章数量 +1
			if len(e.AddedCategoryId) > 0 {
				err = s.repo.IncreasePostCountByIds(ctx, e.AddedCategoryId)
				if err != nil {
					l.ErrorContext(ctx, "Category: post event: failed to increase post count", "categoryIds", e.AddedCategoryId, "error", err)
					continue
				}
			}
			// 分类的文章数量 -1
			if len(e.DeletedCategoryId) > 0 {
				err = s.repo.DecreasePostCountByIds(ctx, e.DeletedCategoryId)
				if err != nil {
					l.ErrorContext(ctx, "Category: post event: failed to decrease post count", "categoryIds", e.DeletedCategoryId, "error", err)
					continue
				}
			}
		}
		l.InfoContext(ctx, "Category: post event: handle successfully")
	}
}
