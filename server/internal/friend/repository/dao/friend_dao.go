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
	"fmt"
	"time"

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox/builder/query"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/bsonx"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Friend struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Url         string             `bson:"url"`
	Logo        string             `bson:"logo"`
	Description string             `bson:"description"`
	Email       string             `bson:"email"`
	Show        bool               `bson:"show"`
	Priority    int                `bson:"priority"`
	Ip          string             `bson:"ip"`
	// 表示是否已经通过审核
	Accepted   bool  `bson:"accepted"`
	CreateTime int64 `bson:"create_time"`
	UpdateTime int64 `bson:"update_time"`
}

type IFriendDao interface {
	FindDisplaying(ctx context.Context) ([]*Friend, error)
	Add(ctx context.Context, friend Friend) error
	FindByUrl(ctx context.Context, url string) (*Friend, error)
	QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Friend, int64, error)
	UpdateById(ctx context.Context, objectID primitive.ObjectID, friend Friend) error
	DeleteById(ctx context.Context, objectID primitive.ObjectID) error
	FindById(ctx context.Context, objectID primitive.ObjectID) (*Friend, error)
	UpdateAccept(ctx context.Context, objectID primitive.ObjectID) error
}

var _ IFriendDao = (*FriendDao)(nil)

func NewFriendDao(db *mongo.Database) *FriendDao {
	return &FriendDao{
		coll: mongox.NewCollection[Friend](db.Collection("friends")),
	}
}

type FriendDao struct {
	coll *mongox.Collection[Friend]
}

func (d *FriendDao) UpdateAccept(ctx context.Context, objectID primitive.ObjectID) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(update.BsonBuilder().Set("accepted", true).Set("show", true).Set("update_time", time.Now().Unix()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update the document from friends, id=%s", objectID.Hex())
	}
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("fails to update the document from friends, id=%s", objectID.Hex())
	}
	return nil
}

func (d *FriendDao) FindById(ctx context.Context, objectID primitive.ObjectID) (*Friend, error) {
	friend, err := d.coll.Finder().Filter(query.Id(objectID)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from friends, id=%s", objectID.Hex())
	}
	return friend, nil
}

func (d *FriendDao) DeleteById(ctx context.Context, objectID primitive.ObjectID) error {
	deleteOne, err := d.coll.Deleter().Filter(query.Id(objectID)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to delete the document from friends, id=%s", objectID.Hex())
	}
	if deleteOne.DeletedCount == 0 {
		return fmt.Errorf("fails to delete the document from friends, id=%s", objectID.Hex())
	}
	return nil
}

func (d *FriendDao) UpdateById(ctx context.Context, objectID primitive.ObjectID, friend Friend) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(objectID)).Updates(
		update.BsonBuilder().Set("name", friend.Name).Set("logo", friend.Logo).Set("description", friend.Description).Set("show", friend.Show).Set("update_time", time.Now().Unix()).Build(),
	).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to update the document from friends, id=%s, friend=%v", objectID.Hex(), friend)
	}
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("fails to update the document from friends, id=%s, friend=%v", objectID.Hex(), friend)
	}
	return nil
}

func (d *FriendDao) QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Friend, int64, error) {
	count, err := d.coll.Finder().Filter(cond).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fails to count the documents from friends, cond=%v", cond)
	}
	friends, err := d.coll.Finder().Filter(cond).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "fails to find the documents from friends, cond=%v, findOptions=%v", cond, findOptions)
	}
	return friends, count, nil
}

func (d *FriendDao) FindByUrl(ctx context.Context, url string) (*Friend, error) {
	friend, err := d.coll.Finder().Filter(bsonx.M("url", url)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the document from friends, url=%s", url)
	}
	return friend, nil
}

func (d *FriendDao) Add(ctx context.Context, friend Friend) error {
	result, err := d.coll.Creator().InsertOne(ctx, friend)
	if err != nil {
		return errors.Wrapf(err, "fails to insert into friends, friend=%v", friend)
	}
	if result.InsertedID == nil {
		return errors.Wrapf(err, "InsertedID=nil, fails to insert into friends, friend=%v", friend)
	}
	return nil
}

func (d *FriendDao) FindDisplaying(ctx context.Context) ([]*Friend, error) {
	friends, err := d.coll.Finder().Filter(query.BsonBuilder().Eq("show", true).Eq("accepted", true).Build()).Find(ctx, options.Find().SetSort(bsonx.M("create_time", 1)))
	if err != nil {
		return nil, errors.Wrapf(err, "fails to find the documents from friends")
	}
	return friends, nil
}
