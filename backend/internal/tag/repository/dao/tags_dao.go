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
	"go.mongodb.org/mongo-driver/mongo"
)

type Tags struct {
	Id         string `bson:"_id"`
	Name       string `bson:"name"`
	Route      string `bson:"route"`
	Disabled   bool   `bson:"disabled"`
	CreateTime int64  `bson:"create_time"`
	UpdateTime int64  `bson:"update_time"`
}

type ITagDao interface {
	GetTags(ctx context.Context) ([]*Tags, error)
	GetByRoute(ctx context.Context, route string) (*Tags, error)
}

var _ ITagDao = (*TagDao)(nil)

func NewTagDao(db *mongo.Database) *TagDao {
	return &TagDao{coll: mongox.NewCollection[Tags](db.Collection("tags"))}
}

type TagDao struct {
	coll *mongox.Collection[Tags]
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
