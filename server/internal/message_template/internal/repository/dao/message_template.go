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

	"github.com/chenmingyong0423/go-mongox"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

type MessageTemplate struct {
	mongox.Model `bson:",inline"`
	Name         string `bson:"name"`
	Title        string `bson:"title"`
	Content      string `bson:"content"`
	// 0 未激活，1 激活
	Active uint `bson:"active"`
	// 0 webmaster 站长， 1 user 用户
	RecipientType uint `bson:"recipient_type"`
}

type IMessageTemplateDao interface {
	FindMsgTplByName(ctx context.Context, name string, recipientType uint) (*MessageTemplate, error)
}

var _ IMessageTemplateDao = (*MessageTemplateDao)(nil)

func NewMessageTemplateDao(db *mongo.Database) *MessageTemplateDao {
	return &MessageTemplateDao{coll: mongox.NewCollection[MessageTemplate](db.Collection("message_templates"))}
}

type MessageTemplateDao struct {
	coll *mongox.Collection[MessageTemplate]
}

func (d *MessageTemplateDao) FindMsgTplByName(ctx context.Context, name string, recipientType uint) (*MessageTemplate, error) {
	msgTpl, err := d.coll.Finder().Filter(
		query.BsonBuilder().Eq("name", name).Eq("active", 1).Eq("recipient_type", recipientType).Build(),
	).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find a docment from message_template, name=%s, recipient_type=%d", name, recipientType)
	}
	return msgTpl, nil
}
