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
	"time"

	"github.com/chenmingyong0423/go-mongox/builder/aggregation"

	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/query"

	"github.com/chenmingyong0423/go-mongox"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type VisitHistory struct {
	Id        string    `bson:"_id"`
	Url       string    `bson:"url"`
	Ip        string    `bson:"ip"`
	UserAgent string    `bson:"user_agent"`
	Origin    string    `bson:"origin"`
	Referer   string    `bson:"referer"`
	CreatedAt time.Time `bson:"created_at"`
}

type TendencyData struct {
	Date      time.Time `bson:"_id"`
	ViewCount int64     `bson:"view_count"`
}

type IVisitLogDao interface {
	Add(ctx context.Context, visitHistory *VisitHistory) error
	CountOfToday(ctx context.Context) (int64, error)
	CountOfTodayByIp(ctx context.Context) (int64, error)
	GetViewTendencyStats4PV(ctx context.Context, days int) ([]*TendencyData, error)
	GetViewTendencyStats4UV(ctx context.Context, days int) ([]*TendencyData, error)
	GetByDate(ctx context.Context, start time.Time, end time.Time) ([]*VisitHistory, error)
}

var _ IVisitLogDao = (*VisitLogDao)(nil)

type VisitLogDao struct {
	coll *mongox.Collection[VisitHistory]
}

func (d *VisitLogDao) GetByDate(ctx context.Context, start time.Time, end time.Time) ([]*VisitHistory, error) {
	return d.coll.Finder().
		Filter(query.BsonBuilder().
			Gte("created_at", start).
			Lte("created_at", end).
			Build()).
		Find(ctx)
}

func (d *VisitLogDao) GetViewTendencyStats4UV(ctx context.Context, days int) ([]*TendencyData, error) {
	daysAgo := time.Now().Local().AddDate(0, 0, -days).Truncate(24 * time.Hour)
	pipeline := aggregation.StageBsonBuilder().
		Match(query.Gte("created_at", daysAgo)).
		Set(bsonx.M("dayStart", bsonx.M("$dateTrunc", bsonx.NewD().Add("date", "$created_at").Add("unit", "day").Build()))).
		Group("$dayStart", bsonx.E("ips", bsonx.M("$addToSet", "$ip"))).
		Project(aggregation.BsonBuilder().AddKeyValues("_id", "$_id").Size("view_count", "$ips").Build()).
		Sort(bsonx.M("_id", 1)).
		Build()
	var result []*TendencyData
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *VisitLogDao) GetViewTendencyStats4PV(ctx context.Context, days int) ([]*TendencyData, error) {
	daysAgo := time.Now().Local().AddDate(0, 0, -days).Truncate(24 * time.Hour)
	pipeline := aggregation.StageBsonBuilder().
		Match(query.Gte("created_at", daysAgo)).
		Set(bsonx.M("dayStart", bsonx.M("$dateTrunc", bsonx.NewD().Add("date", "$created_at").Add("unit", "day").Build()))).
		Group("$dayStart", aggregation.Sum("view_count", 1)...).
		Sort(bsonx.M("_id", 1)).
		Build()
	var result []*TendencyData
	err := d.coll.Aggregator().Pipeline(pipeline).AggregateWithParse(ctx, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (d *VisitLogDao) CountOfTodayByIp(ctx context.Context) (int64, error) {
	startOfDayUnix, endOfDayUnix := d.getBeginSecondsAndEnd()
	distinct, err := d.coll.Collection().Distinct(ctx, "ip", query.BsonBuilder().Gte("created_at", startOfDayUnix).Lte("created_at", endOfDayUnix).Build())
	if err != nil {
		return 0, errors.Wrap(err, "fails to find the count of today from visit_logs")
	}
	return int64(len(distinct)), nil
}

func (d *VisitLogDao) getBeginSecondsAndEnd() (time.Time, time.Time) {
	now := time.Now().Local()
	// 获取当日0点的时间
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 获取当日23:59:59的时间
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	return startOfDay, endOfDay
}

func (d *VisitLogDao) CountOfToday(ctx context.Context) (int64, error) {
	startOfDayUnix, endOfDayUnix := d.getBeginSecondsAndEnd()
	count, err := d.coll.Finder().Filter(query.BsonBuilder().Gte("created_at", startOfDayUnix).Lte("created_at", endOfDayUnix).Build()).Count(ctx)
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
