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

package repository

import (
	"context"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/repository/dao"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type IMessageTemplateRepository interface {
	FindDefaultByTypeAndRecipientType(ctx context.Context, templateType domain.Type, recipientType domain.RecipientType) (*domain.MessageTemplate, error)
	FindAll(ctx context.Context) ([]domain.MessageTemplate, error)
	FindById(ctx context.Context, id string) (domain.MessageTemplate, error)
	Create(ctx context.Context, messageTemplate domain.MessageTemplate) error
	Update(ctx context.Context, id, name, title, content string) error
	UpdateActive(ctx context.Context, id string, active bool) error
	SetDefault(ctx context.Context, id string, templateType domain.Type) error
	Delete(ctx context.Context, id string) error
}

var _ IMessageTemplateRepository = (*MessageTemplateRepository)(nil)

func NewMessageTemplateRepository(dao dao.IMessageTemplateDao) *MessageTemplateRepository {
	return &MessageTemplateRepository{dao: dao}
}

type MessageTemplateRepository struct {
	dao dao.IMessageTemplateDao
}

func (r *MessageTemplateRepository) FindDefaultByTypeAndRecipientType(ctx context.Context, templateType domain.Type, recipientType domain.RecipientType) (*domain.MessageTemplate, error) {
	messageTemplateByType, err := r.dao.FindDefaultByType(ctx, string(templateType), recipientType)
	if err != nil {
		return nil, err
	}
	messageTemplate := r.toDomain(messageTemplateByType)
	return &messageTemplate, nil
}

func (r *MessageTemplateRepository) FindAll(ctx context.Context) ([]domain.MessageTemplate, error) {
	templates, err := r.dao.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.MessageTemplate, 0, len(templates))
	for _, messageTemplate := range templates {
		result = append(result, r.toDomain(messageTemplate))
	}
	return result, nil
}

func (r *MessageTemplateRepository) FindById(ctx context.Context, id string) (domain.MessageTemplate, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return domain.MessageTemplate{}, err
	}
	messageTemplate, err := r.dao.FindById(ctx, objectID)
	if err != nil {
		return domain.MessageTemplate{}, err
	}
	return r.toDomain(messageTemplate), nil
}

func (r *MessageTemplateRepository) Create(ctx context.Context, messageTemplate domain.MessageTemplate) error {
	active := uint(0)
	if messageTemplate.Active {
		active = 1
	}
	return r.dao.Create(ctx, &dao.MessageTemplate{
		Type:          string(messageTemplate.Type),
		Name:          messageTemplate.Name,
		Title:         messageTemplate.Title,
		Content:       messageTemplate.Content,
		Active:        active,
		IsDefault:     false,
		RecipientType: messageTemplate.RecipientType,
	})
}

func (r *MessageTemplateRepository) Update(ctx context.Context, id, name, title, content string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.Update(ctx, objectID, name, title, content)
}

func (r *MessageTemplateRepository) SetDefault(ctx context.Context, id string, templateType domain.Type) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.SetDefault(ctx, objectID, string(templateType))
}

func (r *MessageTemplateRepository) Delete(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.Delete(ctx, objectID)
}

func (r *MessageTemplateRepository) UpdateActive(ctx context.Context, id string, active bool) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return r.dao.UpdateActive(ctx, objectID, active)
}

func (r *MessageTemplateRepository) toDomain(messageTemplate *dao.MessageTemplate) domain.MessageTemplate {
	return domain.MessageTemplate{
		Id:            messageTemplate.ID.Hex(),
		Type:          domain.Type(messageTemplate.Type),
		Name:          messageTemplate.Name,
		Title:         messageTemplate.Title,
		Content:       messageTemplate.Content,
		Active:        messageTemplate.Active == 1,
		IsDefault:     messageTemplate.IsDefault,
		RecipientType: messageTemplate.RecipientType,
		CreatedAt:     messageTemplate.CreatedAt.Unix(),
		UpdatedAt:     messageTemplate.UpdatedAt.Unix(),
	}
}
