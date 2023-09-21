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

	configServ "github.com/chenmingyong0423/fnote/backend/internal/config/service"
	emailServ "github.com/chenmingyong0423/fnote/backend/internal/email/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/google/wire"
)

type IMessageService interface {
	SendEmail(ctx context.Context, subject string, body string, email string, contentType string) error
	SendEmailToWebmaster(ctx context.Context, subject string, body string, contentType string) error
}

var (
	_      IMessageService = (*MessageService)(nil)
	MsgSet                 = wire.NewSet(NewMessageService, wire.Bind(new(IMessageService), new(*MessageService)))
)

func NewMessageService(configServ configServ.IConfigService, emailServ emailServ.IEmailService) *MessageService {
	return &MessageService{
		configServ: configServ,
		emailServ:  emailServ,
	}
}

type MessageService struct {
	configServ configServ.IConfigService
	emailServ  emailServ.IEmailService
}

func (s *MessageService) SendEmailToWebmaster(ctx context.Context, subject string, body string, contentType string) error {
	return s.sendEmail(ctx, subject, body, contentType, "")
}

func (s *MessageService) sendEmail(ctx context.Context, subject, body, contentType, email string) error {
	emailCfg, err := s.configServ.GetEmailConfig(ctx)
	if err != nil {
		return err
	}
	webNMasterCfg, err := s.configServ.GetWebmasterInfo(ctx, "webmaster")
	if err != nil {
		return err
	}
	if email == "" {
		email = emailCfg.Email
	}
	return s.emailServ.SendEmail(ctx, domain.Email{
		Host:        emailCfg.Host,
		Port:        emailCfg.Port,
		Account:     emailCfg.Account,
		Password:    emailCfg.Password,
		Name:        webNMasterCfg.Name,
		To:          []string{email},
		Subject:     subject,
		Body:        body,
		ContentType: contentType,
	})
}

func (s *MessageService) SendEmail(ctx context.Context, subject string, body string, email string, contentType string) error {
	return s.sendEmail(ctx, subject, body, contentType, email)
}
