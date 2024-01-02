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

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/update"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tags struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Route      string             `bson:"route"`
	Disabled   bool               `bson:"disabled"`
	CreateTime int64              `bson:"create_time"`
	UpdateTime int64              `bson:"update_time"`
}

type ITagDao interface {
	GetTags(ctx context.Context) ([]*Tags, error)
	GetByRoute(ctx context.Context, route string) (*Tags, error)
	QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Tags, int64, error)
	Create(ctx context.Context, tag *Tags) (string, error)
	ModifyDisabled(ctx context.Context, id primitive.ObjectID, disabled bool) error
	GetById(ctx context.Context, id primitive.ObjectID) (*Tags, error)
	DeleteById(ctx context.Context, id primitive.ObjectID) error
	RecoverTag(ctx context.Context, tag Tags) error
}

var _ ITagDao = (*TagDao)(nil)

func NewTagDao(db *mongo.Database) *TagDao {
	return &TagDao{coll: mongox.NewCollection[Tags](db.Collection("tags"))}
}

type TagDao struct {
	coll *mongox.Collection[Tags]
}

func (d *TagDao) RecoverTag(ctx context.Context, tag Tags) error {
	_, err := d.coll.Creator().InsertOne(ctx, tag)
	if err != nil {
		return errors.Wrapf(err, "Recover tag failed, tag: %+v", tag)
	}
	return err
}

func (d *TagDao) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	deleteOne, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteOne.DeletedCount == 0 {
		return fmt.Errorf("DeletedCount=0, Delete tag failed, id: %s", id.Hex())
	}
	return nil
}

func (d *TagDao) GetById(ctx context.Context, id primitive.ObjectID) (*Tags, error) {
	tag, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Get tag by id failed, id: %s", id.Hex())
	}
	return tag, nil
}

func (d *TagDao) ModifyDisabled(ctx context.Context, id primitive.ObjectID, disabled bool) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.Set(bsonx.D(bsonx.E("disabled", disabled), bsonx.E("update_time", time.Now().Unix())))).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, Modify tag disabled failed, id: %s", id.Hex())
	}
	return nil
}

func (d *TagDao) Create(ctx context.Context, tag *Tags) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, *tag)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *TagDao) QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Tags, int64, error) {
	finder := d.coll.Finder()
	count, err := finder.Filter(cond).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Count tags failed, cond: %+v", cond)
	}
	categories, err := finder.Filter(cond).Options(findOptions).Find(ctx)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "Query tags failed, cond: %+v, findOptions: %+v", cond, findOptions)
	}
	return categories, count, nil
}

func (d *TagDao) GetByRoute(ctx context.Context, route string) (*Tags, error) {
	return d.coll.Finder().Filter(query.Eq("route", route)).FindOne(ctx)
}

func (d *TagDao) GetTags(ctx context.Context) ([]*Tags, error) {
	tags, err := d.coll.Finder().Filter(query.Eq("disabled", true)).Find(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}
