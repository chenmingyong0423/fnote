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

	"github.com/chenmingyong0423/fnote/server/internal/friend/repository"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/pkg/errors"
)

type IFriendService interface {
	GetFriends(ctx context.Context) ([]domain.Friend, error)
	ApplyForFriend(ctx context.Context, friend domain.Friend) error
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
