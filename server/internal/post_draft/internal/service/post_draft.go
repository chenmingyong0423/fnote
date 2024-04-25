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

	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository"
)

type IPostDraftService interface {
	SavePostDraft(ctx context.Context, postDraft domain.PostDraft) (string, error)
	GetPostDraftById(ctx context.Context, id string) (*domain.PostDraft, error)
	DeletePostDraftById(ctx context.Context, id string) (int64, error)
	GetPostDraftPage(ctx context.Context, page domain.Page) ([]*domain.PostDraft, int64, error)
}

var _ IPostDraftService = (*PostDraftService)(nil)

func NewPostDraftService(repo repository.IPostDraftRepository) *PostDraftService {
	return &PostDraftService{
		repo: repo,
	}
}

type PostDraftService struct {
	repo repository.IPostDraftRepository
}

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
	return s.repo.DeleteById(ctx, id)
}

func (s *PostDraftService) GetPostDraftById(ctx context.Context, id string) (*domain.PostDraft, error) {
	return s.repo.GetById(ctx, id)
}

func (s *PostDraftService) SavePostDraft(ctx context.Context, postDraft domain.PostDraft) (string, error) {
	return s.repo.Save(ctx, postDraft)
}
