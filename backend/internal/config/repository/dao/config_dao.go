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

	"github.com/chenmingyong0423/go-mongox"
	"github.com/chenmingyong0423/go-mongox/bsonx"
	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config defines for the MongoDB Collection "config"
type Config struct {
	Id         string `bson:"_id"`
	Props      any    `bson:"props"`
	Typ        string `bson:"typ"`
	CreateTime int64  `bson:"create_time"`
	UpdateTime int64  `bson:"update_time"`
}

type IConfigDao interface {
	FindByTyp(ctx context.Context, typ string) (*Config, error)
	Increase(ctx context.Context, field string) error
}

func NewConfigDao(db *mongo.Database) *ConfigDao {
	return &ConfigDao{
		coll: mongox.NewCollection[Config](db.Collection("configs")),
	}
}

var _ IConfigDao = (*ConfigDao)(nil)

type ConfigDao struct {
	coll *mongox.Collection[Config]
}

func (d *ConfigDao) Increase(ctx context.Context, field string) error {
	field = fmt.Sprintf("props.%s", field)
	updateResult, err := d.coll.Updater().Filter(bsonx.M("typ", "webmaster")).Updates(update.Inc(bsonx.M(field, 1))).UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "fails to increase %s", field)
	}
	if updateResult.ModifiedCount == 0 {
		return fmt.Errorf("ModifiedCount=0, fails to increase %s", field)
	}
	return nil
}

func (d *ConfigDao) FindByTyp(ctx context.Context, typ string) (*Config, error) {
	config, err := d.coll.Finder().Filter(bsonx.M("typ", typ)).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Find config failed, typ=%s", typ)
	}
	return config, nil
}
