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

package dao

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Friend struct {
	Id          string `bson:"_id"`
	Name        string `bson:"name"`
	Url         string `bson:"url"`
	Logo        string `bson:"logo"`
	Description string `bson:"description"`
	Email       string `bson:"email"`
	// 0 hiding，1 pending，2 showing
	Status     int   `bson:"status"`
	Priority   int   `bson:"priority"`
	CreateTime int64 `bson:"create_time"`
	UpdateTime int64 `bson:"update_time"`
}

type IFriendDao interface {
	FindDisplaying(ctx context.Context) ([]*Friend, error)
	Add(ctx context.Context, friend Friend) error
	FindByUrl(ctx context.Context, url string) (*Friend, error)
}

var _ IFriendDao = (*FriendDao)(nil)

func NewFriendDao(coll *mongo.Collection) *FriendDao {
	return &FriendDao{
		coll: coll,
	}
}

type FriendDao struct {
	coll *mongo.Collection
}

func (d *FriendDao) FindByUrl(ctx context.Context, url string) (*Friend, error) {
	friend := new(Friend)
	if err := d.coll.FindOne(ctx, bson.M{"url": url}).Decode(friend); err != nil {
		return nil, errors.Wrapf(err, "Fails to find the document from %s, url=%s", d.coll.Name(), url)
	}
	return friend, nil
}

func (d *FriendDao) Add(ctx context.Context, friend Friend) error {
	result, err := d.coll.InsertOne(ctx, friend)
	if err != nil {
		return errors.Wrapf(err, "Fails to insert into %s, friend=%v", d.coll.Name(), friend)
	}
	if result.InsertedID == nil {
		return errors.Wrapf(err, "InsertedID=nil, fails to insert into %s, friend=%v", d.coll.Name(), friend)
	}
	return nil
}

func (d *FriendDao) FindDisplaying(ctx context.Context) ([]*Friend, error) {
	cursor, err := d.coll.Find(ctx, bson.M{"status": 2})
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find the documents from %s", d.coll.Name())
	}
	defer cursor.Close(ctx)
	friends := make([]*Friend, 0)
	err = cursor.All(ctx, &friends)
	if err != nil {
		return nil, errors.Wrap(err, "Fails to decode the result")
	}
	return friends, nil
}
