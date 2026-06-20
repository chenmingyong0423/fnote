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

package service

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/domain"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IMessageTemplateService interface {
	FindDefaultByTypeAndRecipientType(ctx context.Context, templateType domain.Type, recipientType domain.RecipientType) (*domain.MessageTemplate, error)
	FindAll(ctx context.Context) ([]domain.MessageTemplate, error)
	Create(ctx context.Context, messageTemplate domain.MessageTemplate) error
	Update(ctx context.Context, id, name, title, content string) error
	UpdateActive(ctx context.Context, id string, active bool) error
	SetDefault(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
}

var _ IMessageTemplateService = (*MessageTemplateService)(nil)

type MessageTemplateService struct {
	repo repository.IMessageTemplateRepository
}

func (s *MessageTemplateService) FindDefaultByTypeAndRecipientType(ctx context.Context, templateType domain.Type, recipientType domain.RecipientType) (*domain.MessageTemplate, error) {
	return s.repo.FindDefaultByTypeAndRecipientType(ctx, templateType, recipientType)
}

func (s *MessageTemplateService) FindAll(ctx context.Context) ([]domain.MessageTemplate, error) {
	return s.repo.FindAll(ctx)
}

func (s *MessageTemplateService) Create(ctx context.Context, messageTemplate domain.MessageTemplate) error {
	recipientType, ok := domain.RecipientForType(messageTemplate.Type)
	if !ok {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "unknown message template type")
	}
	messageTemplate.RecipientType = recipientType
	if err := validateTemplate(messageTemplate.Type, messageTemplate.Name, messageTemplate.Title, messageTemplate.Content); err != nil {
		return err
	}
	if err := s.repo.Create(ctx, messageTemplate); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return apiwrap.NewErrorResponseBody(http.StatusConflict, "message template name already exists in this type")
		}
		return err
	}
	return nil
}

func (s *MessageTemplateService) Update(ctx context.Context, id, name, title, content string) error {
	messageTemplate, err := s.findById(ctx, id)
	if err != nil {
		return err
	}
	if err = validateTemplate(messageTemplate.Type, name, title, content); err != nil {
		return err
	}
	if err = s.repo.Update(ctx, id, strings.TrimSpace(name), strings.TrimSpace(title), content); mongo.IsDuplicateKeyError(err) {
		return apiwrap.NewErrorResponseBody(http.StatusConflict, "message template name already exists in this type")
	}
	return err
}

func (s *MessageTemplateService) UpdateActive(ctx context.Context, id string, active bool) error {
	messageTemplate, err := s.findById(ctx, id)
	if err != nil {
		return err
	}
	if messageTemplate.IsDefault && !active {
		return apiwrap.NewErrorResponseBody(http.StatusConflict, "default message template cannot be disabled")
	}
	return s.repo.UpdateActive(ctx, id, active)
}

func (s *MessageTemplateService) SetDefault(ctx context.Context, id string) error {
	messageTemplate, err := s.findById(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.SetDefault(ctx, id, messageTemplate.Type)
}

func (s *MessageTemplateService) Delete(ctx context.Context, id string) error {
	messageTemplate, err := s.findById(ctx, id)
	if err != nil {
		return err
	}
	if messageTemplate.IsDefault {
		return apiwrap.NewErrorResponseBody(http.StatusConflict, "default message template cannot be deleted")
	}
	return s.repo.Delete(ctx, id)
}

func validateTemplate(templateType domain.Type, name, title, content string) error {
	if strings.TrimSpace(name) == "" {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "message template name cannot be empty")
	}
	if strings.TrimSpace(title) == "" {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "message template title cannot be empty")
	}
	if err := domain.ValidateContent(templateType, content); err != nil {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	}
	return nil
}

func (s *MessageTemplateService) findById(ctx context.Context, id string) (domain.MessageTemplate, error) {
	if _, err := bson.ObjectIDFromHex(id); err != nil {
		return domain.MessageTemplate{}, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "invalid message template id")
	}
	messageTemplate, err := s.repo.FindById(ctx, id)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.MessageTemplate{}, apiwrap.NewErrorResponseBody(http.StatusNotFound, "message template not found")
	}
	return messageTemplate, err
}

func NewMessageTemplateService(repo repository.IMessageTemplateRepository) *MessageTemplateService {
	return &MessageTemplateService{repo: repo}
}
