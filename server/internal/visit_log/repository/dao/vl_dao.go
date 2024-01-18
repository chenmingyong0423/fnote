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
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"time"

	"github.com/chenmingyong0423/go-mongox"

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
	Add(ctx context.Context, visitHistory *VisitHistory) error
	CountOfToday(ctx context.Context) (int64, error)
	CountOfTodayByIp(ctx context.Context) (int64, error)
}

var _ IVisitLogDao = (*VisitLogDao)(nil)

type VisitLogDao struct {
	coll *mongox.Collection[VisitHistory]
}

func (d *VisitLogDao) CountOfTodayByIp(ctx context.Context) (int64, error) {
	startOfDayUnix, endOfDayUnix := d.getBeginSecondsAndEndSeconds()
	distinct, err := d.coll.Collection().Distinct(ctx, "ip", query.BsonBuilder().Gte("create_time", startOfDayUnix).Lte("create_time", endOfDayUnix).Build())
	if err != nil {
		return 0, errors.Wrap(err, "fails to find the count of today from visit_logs")
	}
	return int64(len(distinct)), nil
}

func (d *VisitLogDao) getBeginSecondsAndEndSeconds() (int64, int64) {
	now := time.Now()
	// 获取当日0点的时间
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取当日23:59:59的时间
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return startOfDay.Unix(), endOfDay.Unix()
}

func (d *VisitLogDao) CountOfToday(ctx context.Context) (int64, error) {
	startOfDayUnix, endOfDayUnix := d.getBeginSecondsAndEndSeconds()
	count, err := d.coll.Finder().Filter(query.BsonBuilder().Gte("create_time", startOfDayUnix).Lte("create_time", endOfDayUnix).Build()).Count(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "fails to find the count of today from visit_logs")
	}
	return count, nil
}

func (d *VisitLogDao) Add(ctx context.Context, visitHistory *VisitHistory) error {
	result, err := d.coll.Creator().InsertOne(ctx, visitHistory)
	if err != nil {
		return errors.Wrapf(err, "fails to insert info visit_logs, visitHistory=%v", visitHistory)
	}
	if result == nil {
		return errors.Wrapf(err, "result=nil, fails to insert info visit_logs, visitHistory=%v", visitHistory)
	}
	return nil
}

func NewVisitLogDao(db *mongo.Database) *VisitLogDao {
	return &VisitLogDao{coll: mongox.NewCollection[VisitHistory](db.Collection("visit_logs"))}
}
