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
	FindMsgTplByNameAndRcpType(ctx context.Context, name domain.Name, recipientType domain.RecipientType) (*domain.MessageTemplate, error)
	FindAll(ctx context.Context) ([]domain.MessageTemplate, error)
	Update(ctx context.Context, id, title, content string) error
	UpdateActive(ctx context.Context, id string, active bool) error
}

var _ IMessageTemplateService = (*MessageTemplateService)(nil)

type MessageTemplateService struct {
	repo repository.IMessageTemplateRepository
}

func (s *MessageTemplateService) FindMsgTplByNameAndRcpType(ctx context.Context, name domain.Name, recipientType domain.RecipientType) (*domain.MessageTemplate, error) {
	return s.repo.FindMessageTemplateByNameAndRcpType(ctx, name, recipientType)
}

func (s *MessageTemplateService) FindAll(ctx context.Context) ([]domain.MessageTemplate, error) {
	return s.repo.FindAll(ctx)
}

func (s *MessageTemplateService) Update(ctx context.Context, id, title, content string) error {
	messageTemplate, err := s.findById(ctx, id)
	if err != nil {
		return err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, "message template title cannot be empty")
	}
	if err = domain.ValidateContent(messageTemplate.Name, content); err != nil {
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	}
	return s.repo.Update(ctx, id, title, content)
}

func (s *MessageTemplateService) UpdateActive(ctx context.Context, id string, active bool) error {
	if _, err := s.findById(ctx, id); err != nil {
		return err
	}
	return s.repo.UpdateActive(ctx, id, active)
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
