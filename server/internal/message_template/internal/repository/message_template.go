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
)

type IMessageTemplateRepository interface {
	FindMessageTemplateByNameAndRcpType(ctx context.Context, name string, recipientType uint) (*domain.MessageTemplate, error)
}

var _ IMessageTemplateRepository = (*MessageTemplateRepository)(nil)

func NewMessageTemplateRepository(dao dao.IMessageTemplateDao) *MessageTemplateRepository {
	return &MessageTemplateRepository{dao: dao}
}

type MessageTemplateRepository struct {
	dao dao.IMessageTemplateDao
}

func (r *MessageTemplateRepository) FindMessageTemplateByNameAndRcpType(ctx context.Context, name string, recipientType uint) (*domain.MessageTemplate, error) {
	MessageTemplateByName, err := r.dao.FindMsgTplByName(ctx, name, recipientType)
	if err != nil {
		return nil, err
	}
	return &domain.MessageTemplate{
		Name:    MessageTemplateByName.Name,
		Title:   MessageTemplateByName.Title,
		Content: MessageTemplateByName.Content,
	}, nil
}
