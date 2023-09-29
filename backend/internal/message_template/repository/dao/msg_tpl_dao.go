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
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageTemplate struct {
	Id         string `bson:"_id"`
	Name       string `bson:"name"`
	Content    string `bson:"content"`
	Active     uint   `bson:"active"`
	CreateTime int64  `bson:"create_time"`
	UpdateTime int64  `bson:"update_time"`
}

type IMsgTplDao interface {
}

var _ IMsgTplDao = (*MsgTplDao)(nil)

func NewMsgTplDao(db *mongo.Database) *MsgTplDao {
	return &MsgTplDao{coll: db.Collection("message_template")}
}

type MsgTplDao struct {
	coll *mongo.Collection
}
