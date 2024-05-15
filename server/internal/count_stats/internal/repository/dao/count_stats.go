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

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chenmingyong0423/go-mongox/builder/update"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type CountStats struct {
	mongox.Model `bson:",inline"`
	Type         string `bson:"type"`
	Count        int64  `bson:"count"`
}

type ICountStatsDao interface {
	Create(ctx context.Context, countStats *CountStats) (string, error)
	DeleteByReferenceIdAndType(ctx context.Context, statsType string) error
	DecreaseByReferenceIdAndType(ctx context.Context, statsType string, count int) error
	IncreaseByReferenceIdAndType(ctx context.Context, statsType string, delta int) error
	GetByFilter(ctx context.Context, filter bson.D) ([]*CountStats, error)
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

func (d *CountStatsDao) GetByFilter(ctx context.Context, filter bson.D) ([]*CountStats, error) {
	countStats, err := d.coll.Finder().Filter(filter).Find(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "get count stats by filter error, filter=%v", filter)
	}
	return countStats, nil
}

func (d *CountStatsDao) IncreaseByReferenceIdAndType(ctx context.Context, statsType string, delta int) error {
	oneResult, err := d.coll.Updater().Filter(query.Eq("type", statsType)).Updates(update.BsonBuilder().Inc("count", delta).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "iucrease count stats error, type=%s", statsType)
	}
	if oneResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, iucrease count stats error, type=%s", statsType)
	}
	return nil
}

func (d *CountStatsDao) DecreaseByReferenceIdAndType(ctx context.Context, statsType string, count int) error {
	oneResult, err := d.coll.Updater().Filter(query.Eq("type", statsType)).Updates(update.BsonBuilder().Inc("count", -count).Set("updated_at", time.Now().Local()).Build()).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "decrease count stats error, type=%s", statsType)
	}
	if oneResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, decrease count stats error, type=%s", statsType)
	}
	return nil
}

func (d *CountStatsDao) DeleteByReferenceIdAndType(ctx context.Context, statsType string) error {
	deleteOne, err := d.coll.Deleter().Filter(query.Eq("type", statsType)).DeleteOne(ctx)
	if err != nil {
		return err
	}
	if deleteOne.DeletedCount == 0 {
		return errors.New("DeletedCount=0, delete count stats error")
	}
	return nil
}

func (d *CountStatsDao) Create(ctx context.Context, countStats *CountStats) (string, error) {
	oneResult, err := d.coll.Creator().InsertOne(ctx, countStats)
	if err != nil {
		return "", err
	}
	return oneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
