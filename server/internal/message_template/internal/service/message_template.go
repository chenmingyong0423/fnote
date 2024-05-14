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

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/repository"
)

type IMessageTemplateService interface {
	FindMsgTplByNameAndRcpType(ctx context.Context, name string, recipientType uint) (*domain.MessageTemplate, error)
}

var _ IMessageTemplateService = (*MessageTemplateService)(nil)

type MessageTemplateService struct {
	repo repository.IMessageTemplateRepository
}

func (s *MessageTemplateService) FindMsgTplByNameAndRcpType(ctx context.Context, name string, recipientType uint) (*domain.MessageTemplate, error) {
	return s.repo.FindMessageTemplateByNameAndRcpType(ctx, name, recipientType)
}

func NewMessageTemplateService(repo repository.IMessageTemplateRepository) *MessageTemplateService {
	return &MessageTemplateService{repo: repo}
}
