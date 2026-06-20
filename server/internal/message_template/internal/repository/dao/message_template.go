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

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/domain"
	"github.com/chenmingyong0423/go-mongox/v2/bsonx"
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/pkg/errors"
)

type MessageTemplate struct {
	mongox.Model `bson:",inline"`
	Name         string `bson:"name"`
	Title        string `bson:"title"`
	Content      string `bson:"content"`
	// 0 未激活，1 激活
	Active uint `bson:"active"`
	// 0 webmaster 站长， 1 user 用户
	RecipientType domain.RecipientType `bson:"recipient_type"`
}

type IMessageTemplateDao interface {
	FindMsgTplByName(ctx context.Context, name string, recipientType domain.RecipientType) (*MessageTemplate, error)
	FindAll(ctx context.Context) ([]*MessageTemplate, error)
	FindById(ctx context.Context, id bson.ObjectID) (*MessageTemplate, error)
	Update(ctx context.Context, id bson.ObjectID, title, content string) error
	UpdateActive(ctx context.Context, id bson.ObjectID, active bool) error
}

var _ IMessageTemplateDao = (*MessageTemplateDao)(nil)

func NewMessageTemplateDao(db *mongox.Database) *MessageTemplateDao {
	return &MessageTemplateDao{coll: mongox.NewCollection[MessageTemplate](db, "message_templates")}
}

type MessageTemplateDao struct {
	coll *mongox.Collection[MessageTemplate]
}

func (d *MessageTemplateDao) FindMsgTplByName(ctx context.Context, name string, recipientType domain.RecipientType) (*MessageTemplate, error) {
	msgTpl, err := d.coll.Finder().Filter(
		query.NewBuilder().Eq("name", name).Eq("active", 1).Eq("recipient_type", recipientType).Build(),
	).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "Fails to find a docment from message_template, name=%s, recipient_type=%d", name, recipientType)
	}
	return msgTpl, nil
}

func (d *MessageTemplateDao) FindAll(ctx context.Context) ([]*MessageTemplate, error) {
	return d.coll.Finder().Find(ctx, options.Find().SetSort(bsonx.M("created_at", 1)))
}

func (d *MessageTemplateDao) FindById(ctx context.Context, id bson.ObjectID) (*MessageTemplate, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

func (d *MessageTemplateDao) Update(ctx context.Context, id bson.ObjectID, title, content string) error {
	result, err := d.coll.Updater().
		Filter(query.Id(id)).
		Updates(update.NewBuilder().Set("title", title).Set("content", content).Set("updated_at", time.Now().Local()).Build()).
		UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to update message template, id=%s", id.Hex())
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("message template not found, id=%s", id.Hex())
	}
	return nil
}

func (d *MessageTemplateDao) UpdateActive(ctx context.Context, id bson.ObjectID, active bool) error {
	activeValue := uint(0)
	if active {
		activeValue = 1
	}
	result, err := d.coll.Updater().
		Filter(query.Id(id)).
		Updates(update.NewBuilder().Set("active", activeValue).Set("updated_at", time.Now().Local()).Build()).
		UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to update message template active status, id=%s", id.Hex())
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("message template not found, id=%s", id.Hex())
	}
	return nil
}
