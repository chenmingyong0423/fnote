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

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/chenmingyong0423/go-mongox/bsonx"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type CountStats struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Type        string             `bson:"type"`
	ReferenceId string             `bson:"reference_id"`
	Count       int64              `bson:"count"`
	CreateTime  int64              `bson:"create_time"`
	UpdateTime  int64              `bson:"update_time"`
}

type ICountStatsDao interface {
	GetByReferenceIdAndType(ctx context.Context, referenceIds []string, statsType string) ([]*CountStats, error)
	Create(ctx context.Context, countStats CountStats) (string, error)
	DeleteByReferenceId(ctx context.Context, referenceId string) error
	IncreaseByReferenceIds(ctx context.Context, ids []string) error
	DecreaseByReferenceIds(ctx context.Context, ids []string) error
}

var _ ICountStatsDao = (*CountStatsDao)(nil)

func NewCountStatsDao(db *mongo.Database) *CountStatsDao {
	return &CountStatsDao{
		coll: mongox.NewCollection[CountStats](db.Collection("count_stats")),
	}
}

type CountStatsDao struct {
	coll *mongox.Collection[CountStats]
}

func (d *CountStatsDao) DecreaseByReferenceIds(ctx context.Context, ids []string) error {
	manyResult, err := d.coll.Updater().Filter(query.In[string]("reference_id", ids...)).Updates(
		update.Inc(bsonx.M("count", -1)),
	).UpdateMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "decrease count stats error, ids=%v", ids)
	}
	if manyResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, decrease count stats error, ids=%v", ids)
	}
	return nil
}

func (d *CountStatsDao) IncreaseByReferenceIds(ctx context.Context, ids []string) error {
	manyResult, err := d.coll.Updater().Filter(query.In[string]("reference_id", ids...)).Updates(
		update.Inc(bsonx.M("count", 1)),
	).UpdateMany(ctx)
	if err != nil {
		return errors.Wrapf(err, "increase count stats error, ids=%v", ids)
	}
	if manyResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, increase count stats error, ids=%v", ids)
	}
	return nil
}

func (d *CountStatsDao) DeleteByReferenceId(ctx context.Context, referenceId string) error {
	deleteOne, err := d.coll.Deleter().Filter(bsonx.M("reference_id", referenceId)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteOne.DeletedCount == 0 {
		return errors.New("DeletedCount=0, delete count stats error")
	}
	return nil
}

func (d *CountStatsDao) Create(ctx context.Context, countStats CountStats) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, countStats)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (d *CountStatsDao) GetByReferenceIdAndType(ctx context.Context, referenceIds []string, statsType string) ([]*CountStats, error) {
	countStats, err := d.coll.Finder().Filter(query.BsonBuilder().InString("reference_id", referenceIds...).Eq("type", statsType).Build()).Find(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get count stats by reference id and type error")
	}

	return countStats, nil
}
