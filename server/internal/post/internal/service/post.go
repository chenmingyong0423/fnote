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
	"fmt"
	"log/slog"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/google/uuid"

	"github.com/chenmingyong0423/fnote/server/internal/post/internal/domain"

	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/post/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	"github.com/chenmingyong0423/gkit/slice"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
)

type IPostService interface {
	GetLatestPosts(ctx context.Context, count int64) ([]*domain.Post, error)
	GetPosts(ctx context.Context, pageRequest *domain.PostRequest) ([]*domain.Post, int64, error)
	GetPunishedPostById(ctx context.Context, id string) (*domain.Post, error)
	IncreaseVisitCount(ctx context.Context, id string) error
	AdminGetPosts(ctx context.Context, page domain.Page) ([]*domain.Post, int64, error)
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

func NewPostService(repo repository.IPostRepository, cfgService website_config.Service, eventBus *eventbus.EventBus) *PostService {
	s := &PostService{
		repo:       repo,
		cfgService: cfgService,
		eventBus:   eventBus,
	}
	go s.subscribeCommentEvent()
	return s
}

type PostService struct {
	repo       repository.IPostRepository
	cfgService website_config.Service
	eventBus   *eventbus.EventBus
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
	var (
		marshal []byte
		err     error
	)
	if isNewPost {
		marshal, err = s.marshalPostEvent(savedPost)
		if err != nil {
			return err
		}
	} else {
		marshal, err = s.marshalUpdatePostEvent(originalPost, savedPost)
		if err != nil {
			return err
		}
	}

	err = s.repo.SavePost(ctx, savedPost)
	if err != nil {
		return err
	}

	go func() {
		if isNewPost {
			s.eventBus.Publish("post", eventbus.Event{Payload: marshal})
		} else {
			s.eventBus.Publish("post", eventbus.Event{Payload: marshal})
		}
	}()
	return nil
}

func (s *PostService) marshalUpdatePostEvent(post *domain.Post, savedPost *domain.Post) ([]byte, error) {
	// categories
	// 被移除的
	removedCategories := slice.DiffFunc(post.Categories, savedPost.Categories, func(srcItem, dstItem domain.Category4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	var removedCategoryIds []string
	if len(removedCategories) > 0 {
		removedCategoryIds = slice.Map[domain.Category4Post, string](removedCategories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
	}
	// 被添加的
	addedCategories := slice.DiffFunc(savedPost.Categories, post.Categories, func(srcItem, dstItem domain.Category4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	var addedCategoryIds []string
	if len(addedCategories) > 0 {
		addedCategoryIds = slice.Map[domain.Category4Post, string](addedCategories, func(_ int, c domain.Category4Post) string {
			return c.Id
		})
	}
	// tags
	// 被移除的
	removedTags := slice.DiffFunc(post.Tags, savedPost.Tags, func(srcItem, dstItem domain.Tag4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	var removedTagIds []string
	if len(removedTags) > 0 {
		removedTagIds = slice.Map[domain.Tag4Post, string](removedTags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
	}
	// 被添加的
	addedTags := slice.DiffFunc(savedPost.Tags, post.Tags, func(srcItem, dstItem domain.Tag4Post) bool {
		return srcItem.Id == dstItem.Id
	})
	var addedTagIds []string
	if len(addedTags) > 0 {
		addedTagIds = slice.Map[domain.Tag4Post, string](addedTags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		})
	}

	postEvent := domain.PostEvent{
		PostId:            savedPost.Id,
		AddedCategoryId:   addedCategoryIds,
		DeletedCategoryId: removedCategoryIds,
		AddedTagId:        addedTagIds,
		DeletedTagId:      removedTagIds,
		NewFileId:         strings.Split(savedPost.CoverImg[8:], ".")[0],
		OldFileId:         strings.Split(post.CoverImg[1:], ".")[0],
		Type:              "update",
	}
	return json.Marshal(postEvent)
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

	postInfo := domain.PostEvent{
		PostId: id,
		DeletedCategoryId: slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
			return c.Id
		}),
		DeletedTagId: slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		}),
		OldFileId: strings.Split(post.CoverImg[1:], ".")[0],
		Type:      "delete",
	}
	marshal, err := json.Marshal(postInfo)
	if err != nil {
		return err
	}

	err = s.repo.DeletePost(ctx, id)
	if err != nil {
		return err
	}

	s.eventBus.Publish("post", eventbus.Event{Payload: marshal})
	return nil
}

func (s *PostService) AddPost(ctx context.Context, post *domain.Post) error {
	if post.Id == "" {
		post.Id = uuidx.RearrangeUUID4()
	}

	marshal, err := s.marshalPostEvent(post)
	if err != nil {
		return err
	}

	err = s.repo.AddPost(ctx, post)
	if err != nil {
		return err
	}
	s.eventBus.Publish("post", eventbus.Event{Payload: marshal})
	return nil
}

func (s *PostService) marshalPostEvent(post *domain.Post) ([]byte, error) {
	postInfo := domain.PostEvent{
		PostId: post.Id,
		AddedCategoryId: slice.Map[domain.Category4Post, string](post.Categories, func(_ int, c domain.Category4Post) string {
			return c.Id
		}),
		AddedTagId: slice.Map[domain.Tag4Post, string](post.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Id
		}),
		NewFileId: strings.Split(post.CoverImg[8:], ".")[0],
		Type:      "create",
	}
	marshal, err := json.Marshal(postInfo)
	if err != nil {
		return nil, err
	}
	return marshal, nil
}

func (s *PostService) AdminGetPosts(ctx context.Context, page domain.Page) ([]*domain.Post, int64, error) {
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

func (s *PostService) subscribeCommentEvent() {
	eventChan := s.eventBus.Subscribe("comment")
	type contextKey string
	for event := range eventChan {
		rid := uuid.NewString()
		var key contextKey = "X-Request-ID"
		ctx := context.WithValue(context.Background(), key, rid)
		l := slog.Default().With("X-Request-ID", rid)
		l.InfoContext(ctx, "Post: comment event", "payload", string(event.Payload))
		var e domain.CommentEvent
		err := jsoniter.Unmarshal(event.Payload, &e)
		if err != nil {
			l.ErrorContext(ctx, "Post: comment event: failed to unmarshal", "error", err)
			continue
		}
		switch e.Type {
		case "create":
			err = s.repo.IncreaseCommentCount(ctx, e.PostId)
			if err != nil {
				l.ErrorContext(ctx, "Post: comment event: failed to increase the count of comment in post", "count", 1, "error", err)
				continue
			}
		case "delete":
			err = s.repo.DecreaseCommentCount(ctx, e.PostId, e.Count)
			if err != nil {
				l.ErrorContext(ctx, "Post: comment event: failed to increase the count of comment in post", "count", e.Count, "error", err)
				continue
			}
		}
		l.InfoContext(ctx, "Post: comment event: handle successfully")
	}
}
