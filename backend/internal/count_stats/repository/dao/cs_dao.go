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

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type CountStats struct {
	Id          string `bson:"_id"`
	Type        string `bson:"type"`
	ReferenceId string `bson:"reference_id"`
	Count       int64  `bson:"count"`
	CreateTime  int64  `bson:"create_time"`
	UpdateTime  int64  `bson:"update_time"`
}

type ICountStatsDao interface {
	GetByReferenceIdAndType(ctx context.Context, referenceIds []string, statsType string) ([]*CountStats, error)
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

func (d *CountStatsDao) GetByReferenceIdAndType(ctx context.Context, referenceIds []string, statsType string) ([]*CountStats, error) {
	countStats, err := d.coll.Finder().Filter(query.BsonBuilder().InString("reference_id", referenceIds...).Eq("type", statsType).Build()).Find(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get count stats by reference id and type error")
	}

	return countStats, nil
}
