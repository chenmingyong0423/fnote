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
	"net/http"

	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/repository"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"
	"github.com/pkg/errors"
)

type IFriendService interface {
	GetFriends(ctx context.Context) ([]domain.Friend, error)
	ApplyForFriend(ctx context.Context, friend domain.Friend) error
	AdminGetFriends(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Friend, int64, error)
	AdminUpdateFriend(ctx context.Context, friend domain.Friend) error
	AdminDeleteFriend(ctx context.Context, id string) error
	AdminApproveFriend(ctx context.Context, id string) (string, error)
	AdminRejectFriend(ctx context.Context, id string) (string, error)
}

var _ IFriendService = (*FriendService)(nil)

func NewFriendService(repo repository.IFriendRepository) *FriendService {
	return &FriendService{
		repo: repo,
	}
}

type FriendService struct {
	repo repository.IFriendRepository
}

func (s *FriendService) AdminRejectFriend(ctx context.Context, id string) (string, error) {
	friend, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", apiwrap.NewErrorResponseBody(http.StatusNotFound, "friend not found")
		}
		return "", err
	}
	if friend.IsRejected() {
		return "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, "friend already rejected")
	}
	err = s.repo.UpdateFriendRejected(ctx, id)
	if err != nil {
		return "", err
	}
	return friend.Email, nil
}

func (s *FriendService) AdminApproveFriend(ctx context.Context, id string) (string, error) {
	friend, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", apiwrap.NewErrorResponseBody(http.StatusNotFound, "friend not found")
		}
		return "", err
	}
	if friend.IsApproved() {
		return "", apiwrap.NewErrorResponseBody(http.StatusBadRequest, "friend already accepted")
	}
	err = s.repo.UpdateFriendApproved(ctx, id)
	if err != nil {
		return "", err
	}
	return friend.Email, nil
}

func (s *FriendService) AdminDeleteFriend(ctx context.Context, id string) error {
	_, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "friend not found")
		}
		return errors.WithMessage(err, "s.repo.FindById failed")
	}
	return s.repo.DeleteById(ctx, id)
}

func (s *FriendService) AdminUpdateFriend(ctx context.Context, friend domain.Friend) error {
	_, err := s.repo.FindById(ctx, friend.Id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return apiwrap.NewErrorResponseBody(http.StatusNotFound, "friend not found")
		}
		return errors.WithMessage(err, "s.repo.FindById failed")
	}
	return s.repo.UpdateById(ctx, friend)
}

func (s *FriendService) AdminGetFriends(ctx context.Context, pageDTO dto.PageDTO) ([]domain.Friend, int64, error) {
	return s.repo.FindAll(ctx, pageDTO)
}

func (s *FriendService) ApplyForFriend(ctx context.Context, friend domain.Friend) error {
	err := s.repo.Add(ctx, friend)
	if err != nil {
		return errors.WithMessage(err, "s.repo.Add failed")
	}
	return nil
}

func (s *FriendService) GetFriends(ctx context.Context) ([]domain.Friend, error) {
	return s.repo.FindDisplaying(ctx)
}
