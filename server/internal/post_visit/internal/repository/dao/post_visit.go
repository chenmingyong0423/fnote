// Copyright 2024 chenmingyong0423

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
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/chenmingyong0423/go-mongox"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostVisit struct {
	Id        primitive.ObjectID `bson:"_id"`
	PostId    string             `bson:"post_id"`
	Ip        string             `bson:"ip"`
	UserAgent string             `bson:"user_agent"`
	Origin    string             `bson:"origin"`
	Referer   string             `bson:"referer"`
	StayTime  int64              `bson:"stay_time"`
	VisitAt   time.Time          `bson:"visit_at"`
}

type IPostVisitDao interface {
	Insert(ctx context.Context, postVisit *PostVisit) (string, error)
}

var _ IPostVisitDao = (*PostVisitDao)(nil)

func NewPostVisitDao(db *mongo.Database) *PostVisitDao {
	return &PostVisitDao{coll: mongox.NewCollection[PostVisit](db.Collection("post_visits"))}
}

type PostVisitDao struct {
	coll *mongox.Collection[PostVisit]
}

func (d *PostVisitDao) Insert(ctx context.Context, postVisit *PostVisit) (string, error) {
	insertOneResult, err := d.coll.Creator().InsertOne(ctx, postVisit)
	if err != nil {
		return "", errors.Wrapf(err, "failed to insert post visit, postVisit: %+v", postVisit)
	}
	return insertOneResult.InsertedID.(primitive.ObjectID).Hex(), nil
}
