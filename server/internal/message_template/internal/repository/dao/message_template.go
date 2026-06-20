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
	"github.com/chenmingyong0423/go-mongox/v2/builder/query"
	"github.com/chenmingyong0423/go-mongox/v2/builder/update"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/pkg/errors"
)

type MessageTemplate struct {
	mongox.Model `bson:",inline"`
	Type         string `bson:"type"`
	Name         string `bson:"name"`
	Title        string `bson:"title"`
	Content      string `bson:"content"`
	IsDefault    bool   `bson:"is_default"`
	// 0 未激活，1 激活
	Active uint `bson:"active"`
	// 0 webmaster 站长， 1 user 用户
	RecipientType domain.RecipientType `bson:"recipient_type"`
}

type IMessageTemplateDao interface {
	FindDefaultByType(ctx context.Context, templateType string, recipientType domain.RecipientType) (*MessageTemplate, error)
	FindAll(ctx context.Context) ([]*MessageTemplate, error)
	FindById(ctx context.Context, id bson.ObjectID) (*MessageTemplate, error)
	Create(ctx context.Context, messageTemplate *MessageTemplate) error
	Update(ctx context.Context, id bson.ObjectID, name, title, content string) error
	UpdateActive(ctx context.Context, id bson.ObjectID, active bool) error
	SetDefault(ctx context.Context, id bson.ObjectID, templateType string) error
	Delete(ctx context.Context, id bson.ObjectID) error
}

var _ IMessageTemplateDao = (*MessageTemplateDao)(nil)

func NewMessageTemplateDao(db *mongox.Database) *MessageTemplateDao {
	return &MessageTemplateDao{db: db, coll: mongox.NewCollection[MessageTemplate](db, "message_templates")}
}

type MessageTemplateDao struct {
	db   *mongox.Database
	coll *mongox.Collection[MessageTemplate]
}

func (d *MessageTemplateDao) FindDefaultByType(ctx context.Context, templateType string, recipientType domain.RecipientType) (*MessageTemplate, error) {
	msgTpl, err := d.coll.Finder().Filter(
		query.NewBuilder().Eq("type", templateType).Eq("is_default", true).Eq("active", 1).Eq("recipient_type", recipientType).Build(),
	).FindOne(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find default message template, type=%s, recipient_type=%d", templateType, recipientType)
	}
	return msgTpl, nil
}

func (d *MessageTemplateDao) FindAll(ctx context.Context) ([]*MessageTemplate, error) {
	return d.coll.Finder().Find(ctx, options.Find().SetSort(bson.D{
		{Key: "type", Value: 1},
		{Key: "is_default", Value: -1},
		{Key: "created_at", Value: 1},
	}))
}

func (d *MessageTemplateDao) FindById(ctx context.Context, id bson.ObjectID) (*MessageTemplate, error) {
	return d.coll.Finder().Filter(query.Id(id)).FindOne(ctx)
}

func (d *MessageTemplateDao) Create(ctx context.Context, messageTemplate *MessageTemplate) error {
	_, err := d.coll.Creator().InsertOne(ctx, messageTemplate)
	return errors.Wrap(err, "failed to create message template")
}

func (d *MessageTemplateDao) Update(ctx context.Context, id bson.ObjectID, name, title, content string) error {
	result, err := d.coll.Updater().
		Filter(query.Id(id)).
		Updates(update.NewBuilder().Set("name", name).Set("title", title).Set("content", content).Set("updated_at", time.Now().Local()).Build()).
		UpdateOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to update message template, id=%s", id.Hex())
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("message template not found, id=%s", id.Hex())
	}
	return nil
}

func (d *MessageTemplateDao) SetDefault(ctx context.Context, id bson.ObjectID, templateType string) error {
	session, err := d.db.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (any, error) {
		_, txErr := d.coll.Updater().
			Filter(query.Eq("type", templateType)).
			Updates(update.NewBuilder().Set("is_default", false).Set("updated_at", time.Now().Local()).Build()).
			UpdateMany(txCtx)
		if txErr != nil {
			return nil, txErr
		}
		result, txErr := d.coll.Updater().
			Filter(query.NewBuilder().Id(id).Eq("type", templateType).Build()).
			Updates(update.NewBuilder().Set("is_default", true).Set("active", uint(1)).Set("updated_at", time.Now().Local()).Build()).
			UpdateOne(txCtx)
		if txErr != nil {
			return nil, txErr
		}
		if result.MatchedCount == 0 {
			return nil, fmt.Errorf("message template not found, id=%s", id.Hex())
		}
		return nil, nil
	})
	return errors.Wrap(err, "failed to set default message template")
}

func (d *MessageTemplateDao) Delete(ctx context.Context, id bson.ObjectID) error {
	result, err := d.coll.Deleter().Filter(query.Id(id)).DeleteOne(ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to delete message template, id=%s", id.Hex())
	}
	if result.DeletedCount == 0 {
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
