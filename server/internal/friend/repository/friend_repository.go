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

package repository

import (
	"context"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/friend/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/pkg/errors"
)

type IFriendRepository interface {
	FindDisplaying(ctx context.Context) ([]domain.Friend, error)
	Add(ctx context.Context, friend domain.Friend) error
	FindByUrl(ctx context.Context, url string) (domain.Friend, error)
}

var _ IFriendRepository = (*FriendRepository)(nil)

func NewFriendRepository(dao dao.IFriendDao) *FriendRepository {
	return &FriendRepository{
		dao: dao,
	}
}

type FriendRepository struct {
	dao dao.IFriendDao
}

func (r *FriendRepository) FindByUrl(ctx context.Context, url string) (friend domain.Friend, err error) {
	friendDao, err := r.dao.FindByUrl(ctx, url)
	if err != nil {
		return friend, err
	}
	friend = r.toDomainFriend(friendDao)
	return
}

func (r *FriendRepository) Add(ctx context.Context, friend domain.Friend) error {
	unix := time.Now().Unix()
	err := r.dao.Add(ctx, dao.Friend{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Email:       friend.Email,
		Status:      dao.FriendStatusHidden,
		Ip:          friend.Ip,
		Accepted:    false,
		CreateTime:  unix,
		UpdateTime:  unix,
	})
	if err != nil {
		return errors.WithMessage(err, "r.dao.Add failed")
	}
	return nil
}

func (r *FriendRepository) FindDisplaying(ctx context.Context) ([]domain.Friend, error) {
	friends, err := r.dao.FindDisplaying(ctx)
	if err != nil {
		return nil, err
	}
	return r.toDomainFriends(friends), nil
}

func (r *FriendRepository) toDomainFriends(friends []*dao.Friend) []domain.Friend {
	results := make([]domain.Friend, 0, len(friends))
	for _, friend := range friends {
		results = append(results, r.toDomainFriend(friend))
	}
	return results
}

func (r *FriendRepository) toDomainFriend(friend *dao.Friend) domain.Friend {
	return domain.Friend{
		Id:          friend.Id.Hex(),
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Status:      friend.Status,
		Priority:    friend.Priority,
	}
}
