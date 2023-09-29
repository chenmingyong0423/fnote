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
	"fmt"

	"github.com/chenmingyong0423/fnote/backend/internal/friend/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type IFriendService interface {
	GetFriends(ctx context.Context) (api.ListVO[*domain.FriendVO], error)
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
	if f, err := s.repo.FindByUrl(ctx, friend.Url); !errors.Is(err, mongo.ErrNoDocuments) {
		if err != nil {
			return err
		}
		return fmt.Errorf("the friend had already applied for, friend=%v", f)
	}
	err := s.repo.Add(ctx, friend)
	if err != nil {
		return errors.WithMessage(err, "s.repo.Add failed")
	}
	return nil
}

func (s *FriendService) GetFriends(ctx context.Context) (api.ListVO[*domain.FriendVO], error) {
	vo := api.ListVO[*domain.FriendVO]{}
	friends, err := s.repo.FindDisplaying(ctx)
	if err != nil {
		return vo, errors.WithMessage(err, "s.repo.FindDisplaying failed")
	}
	vo.List = s.toFriendVOs(friends)
	return vo, nil
}

func (s *FriendService) toFriendVOs(friends []*domain.Friend) []*domain.FriendVO {
	result := make([]*domain.FriendVO, 0, len(friends))
	for _, friend := range friends {
		result = append(result, s.toFriendVO(friend))
	}
	return result
}

func (s *FriendService) toFriendVO(friend *domain.Friend) *domain.FriendVO {
	return &domain.FriendVO{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Status:      friend.Status,
		Priority:    friend.Priority,
	}
}
