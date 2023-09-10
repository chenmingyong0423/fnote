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
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitHistory struct {
	Id         string `bson:"_id"`
	Url        string `bson:"url"`
	Ip         string `bson:"ip"`
	UserAgent  string `bson:"user_agent"`
	Origin     string `bson:"origin"`
	Referer    string `bson:"referer"`
	CreateTime int64  `bson:"create_time"`
}

type IVisitLogDao interface {
	Add(ctx context.Context, visitHistory VisitHistory) error
}

var _ IVisitLogDao = (*VisitLogDao)(nil)

type VisitLogDao struct {
	coll *mongo.Collection
}

func (d *VisitLogDao) Add(ctx context.Context, visitHistory VisitHistory) error {
	result, err := d.coll.InsertOne(ctx, visitHistory)
	if err != nil {
		return errors.Wrapf(err, "fails to insert info %s, visitHistory=%v", d.coll.Name(), visitHistory)
	}
	if result == nil {
		return errors.Wrapf(err, "result=nil, fails to insert info %s, visitHistory=%v", d.coll.Name(), visitHistory)
	}
	return nil
}

func NewVisitLogDao(coll *mongo.Collection) *VisitLogDao {
	return &VisitLogDao{coll: coll}
}
