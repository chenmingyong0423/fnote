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

	"github.com/chenmingyong0423/go-mongox"

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chenmingyong0423/go-mongox/builder/query"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type Category struct {
	mongox.Model `bson:",inline"`
	Name         string `bson:"name"`
	Route        string `bson:"route"`
	Description  string `bson:"description"`
	Sort         int64  `bson:"sort"`
	Enabled      bool   `bson:"enabled"`
	ShowInNav    bool   `bson:"show_in_nav"`
	PostCount    int64  `bson:"post_count"`
}

type ICategoryDao interface {
	GetAll(ctx context.Context) ([]*Category, error)
	GetByRoute(ctx context.Context, route string) (*Category, error)
	QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Category, int64, error)
	Create(ctx context.Context, category *Category) (string, error)
	ModifyEnabled(ctx context.Context, id primitive.ObjectID, enabled bool) error
	ModifyCategory(ctx context.Context, id primitive.ObjectID, description string) error
	DeleteById(ctx context.Context, id primitive.ObjectID) error
	GetByShowInNav(ctx context.Context) ([]*Category, error)
	ModifyCategoryNavigation(ctx context.Context, id primitive.ObjectID, showInNav bool) error
	GetById(ctx context.Context, id primitive.ObjectID) (*Category, error)
	RecoverCategory(ctx context.Context, category *Category) error
	GetEnabled(ctx context.Context) ([]*Category, error)
	IncreasePostCountByIds(ctx context.Context, categoryObjectIds []primitive.ObjectID) error
	DecreasePostCountByIds(ctx context.Context, categoryObjectIds []primitive.ObjectID) error
	FindEnabledCategories(ctx context.Context) ([]*Category, error)
}

var _ ICategoryDao = (*CategoryDao)(nil)

func NewCategoryDao(db *mongo.Database) *CategoryDao {
	return &CategoryDao{
		coll: mongox.NewCollection[Category](db.Collection("categories")),
	}
}

type CategoryDao struct {
	coll *mongox.Collection[Category]
}

func (d *CategoryDao) FindEnabledCategories(ctx context.Context) ([]*Category, error) {
	return d.coll.Finder().Filter(query.Eq("enabled", true)).Find(ctx)
}

func (d *CategoryDao) DecreasePostCountByIds(ctx context.Context, categoryObjectIds []primitive.ObjectID) error {
	updateResult, err := d.coll.Updater().Filter(query.In("_id", categoryObjectIds...)).Updates(update.BsonBuilder().Inc("post_count", -1).Set("updated_at", time.Now().Local()).Build()).UpdateMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to decrease post count by ids, ids=%+v", categoryObjectIds)
	}
	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("MatchedCount=0, decrease post count failed, ids: %+v", categoryObjectIds)
	}
	return nil
}

func (d *CategoryDao) IncreasePostCountByIds(ctx context.Context, categoryObjectIds []primitive.ObjectID) error {
	updateResult, err := d.coll.Updater().Filter(query.In("_id", categoryObjectIds...)).Updates(update.BsonBuilder().Inc("post_count", 1).Set("updated_at", time.Now().Local()).Build()).UpdateMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to increase post count by ids, ids=%+v", categoryObjectIds)
	}
	if updateResult.MatchedCount == 0 {
		return fmt.Errorf("MatchedCount=0, increase post count failed, ids: %+v", categoryObjectIds)
	}
	return nil
}

func (d *CategoryDao) GetEnabled(ctx context.Context) ([]*Category, error) {
	categories, err := d.coll.Finder().Filter(bsonx.M("enabled", true)).Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Find categories failed, enabled=true")
	}
	return categories, nil
}

func (d *CategoryDao) RecoverCategory(ctx context.Context, category *Category) error {
	_, err := d.coll.Creator().InsertOne(ctx, category)
	if err != nil {
		return errors.Wrapf(err, "Recover category failed, category: %+v", category)
	}
	return err
}

func (d *CategoryDao) GetById(ctx context.Context, id primitive.ObjectID) (*Category, error) {
	category, err := d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Get category by id failed, id=%s", id)
	}
	return category, nil
}

func (d *CategoryDao) ModifyCategoryNavigation(ctx context.Context, id primitive.ObjectID, showInNav bool) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("show_in_nav", showInNav).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "Modify category navigation failed, id=%s, showInNav=%v", id, showInNav)
	}
	if updateOne.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, Modify category navigation failed, id=%s, showInNav=%v", id, showInNav)
	}
	return nil
}

func (d *CategoryDao) GetByShowInNav(ctx context.Context) ([]*Category, error) {
	categories, err := d.coll.Finder().Filter(query.BsonBuilder().Eq("show_in_nav", true).Eq("enabled", true).Build()).Find(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "Find categories failed")
	}
	return categories, nil
}

func (d *CategoryDao) DeleteById(ctx context.Context, id primitive.ObjectID) error {
	deleteOne, err := d.coll.Deleter().Filter(bsonx.Id(id)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteOne.DeletedCount == 0 {
		return errors.New("DeletedCount=0, delete category error")
	}
	return nil
}

func (d *CategoryDao) ModifyCategory(ctx context.Context, id primitive.ObjectID, description string) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("description", description).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if updateOne.ModifiedCount == 0 {
		return errors.New("ModifiedCount=0, Modify description failed")
	}
	return nil
}

func (d *CategoryDao) ModifyEnabled(ctx context.Context, id primitive.ObjectID, enabled bool) error {
	updateOne, err := d.coll.Updater().Filter(query.Id(id)).Updates(update.BsonBuilder().Set("enabled", enabled).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return err
	}
	if updateOne.ModifiedCount == 0 {
		return errors.New("ModifiedCount=0, Modify enabled failed")
	}
	return nil
}

func (d *CategoryDao) Create(ctx context.Context, category *Category) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, category)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *CategoryDao) QuerySkipAndSetLimit(ctx context.Context, cond bson.D, findOptions *options.FindOptions) ([]*Category, int64, error) {
	finder := d.coll.Finder()
	count, err := finder.Filter(cond).Count(ctx)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Count categories failed")
	}
	categories, err := finder.Filter(cond).Find(ctx, findOptions)
	if err != nil {
		return nil, 0, errors.Wrap(err, "Find categories failed")
	}
	return categories, count, nil
}

func (d *CategoryDao) GetByRoute(ctx context.Context, route string) (*Category, error) {
	category, err := d.coll.Finder().Filter(query.Eq("route", route)).FindOne(ctx)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (d *CategoryDao) GetAll(ctx context.Context) ([]*Category, error) {
	result, err := d.coll.Finder().Filter(bsonx.M("enabled", true)).Find(ctx, options.Find().SetSort(bsonx.M("sort", 1)))
	if err != nil {
		return nil, errors.Wrap(err, "Find all categories failed failed")
	}
	return result, nil
}
