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
	"fmt"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/repository/dao"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
)

type IFriendRepository interface {
	FindDisplaying(ctx context.Context) ([]domain.Friend, error)
	Add(ctx context.Context, friend domain.Friend) error
	FindByUrl(ctx context.Context, url string) (domain.Friend, error)
	FindAll(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Friend, int64, error)
	UpdateById(ctx context.Context, friend domain.Friend) error
	DeleteById(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (domain.Friend, error)
	UpdateFriendApproved(ctx context.Context, id string) error
	UpdateFriendRejected(ctx context.Context, id string) error
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

func (r *FriendRepository) UpdateFriendRejected(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdateRejected(ctx, objectID)
}

func (r *FriendRepository) UpdateFriendApproved(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdateApproved(ctx, objectID)
}

func (r *FriendRepository) FindById(ctx context.Context, id string) (domain.Friend, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Friend{}, err
	}
	friend, err := r.dao.FindById(ctx, objectID)
	if err != nil {
		return domain.Friend{}, err
	}
	return r.toDomainFriend(friend), nil
}

func (r *FriendRepository) DeleteById(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.DeleteById(ctx, objectID)
}

func (r *FriendRepository) UpdateById(ctx context.Context, friend domain.Friend) error {
	id, err := primitive.ObjectIDFromHex(friend.Id)
	if err != nil {
		return err
	}
	return r.dao.UpdateById(ctx, id, dao.Friend{
		Name:        friend.Name,
		Logo:        friend.Logo,
		Description: friend.Description,
		Status:      dao.FriendStatus(friend.Status),
	})
}

func (r *FriendRepository) FindAll(ctx context.Context, pageDTO domain.PageDTO) ([]domain.Friend, int64, error) {
	condBuilder := query.BsonBuilder()
	if pageDTO.Keyword != "" {
		condBuilder.RegexOptions("name", fmt.Sprintf(".*%s.*", strings.TrimSpace(pageDTO.Keyword)), "i")
	}
	cond := condBuilder.Build()

	findOptions := options.Find()
	findOptions.SetSkip((pageDTO.PageNo - 1) * pageDTO.PageSize).SetLimit(pageDTO.PageSize)
	if pageDTO.Field != "" && pageDTO.Order != "" {
		findOptions.SetSort(bsonx.M(pageDTO.Field, pageDTO.OrderConvertToInt()))
	} else {
		findOptions.SetSort(bsonx.M("created_at", 1))
	}
	friends, total, err := r.dao.QuerySkipAndSetLimit(ctx, cond, findOptions)
	return r.toDomainFriends(friends), total, err
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
	err := r.dao.Add(ctx, &dao.Friend{
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Email:       friend.Email,
		Ip:          friend.Ip,
		Status:      dao.FriendStatusPending,
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
		Id:          friend.ID.Hex(),
		Name:        friend.Name,
		Url:         friend.Url,
		Logo:        friend.Logo,
		Description: friend.Description,
		Status:      int(friend.Status),
		Priority:    friend.Priority,
		Email:       friend.Email,
		Ip:          friend.Ip,
		CreatedAt:   friend.CreatedAt.Unix(),
	}
}
