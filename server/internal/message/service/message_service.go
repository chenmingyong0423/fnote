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

	configServ "github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/service"

	emailServ "github.com/chenmingyong0423/fnote/server/internal/email/service"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
)

type IMessageService interface {
	SendEmailWithEmail(ctx context.Context, msgTplName, email, contentType string, args ...any) error
	SendEmailToWebmaster(ctx context.Context, msgTplName, contentType string) error
}

var (
	_ IMessageService = (*MessageService)(nil)
)

func NewMessageService(configServ configServ.IWebsiteConfigService, emailServ emailServ.IEmailService, msgTplService service.IMsgTplService) *MessageService {
	return &MessageService{
		configServ:    configServ,
		emailServ:     emailServ,
		msgTplService: msgTplService,
	}
}

type MessageService struct {
	configServ    configServ.IWebsiteConfigService
	emailServ     emailServ.IEmailService
	msgTplService service.IMsgTplService
}

func (s *MessageService) SendEmailToWebmaster(ctx context.Context, msgTplName, contentType string) error {
	return s.sendEmail(ctx, msgTplName, contentType, 0, "")
}

func (s *MessageService) sendEmail(ctx context.Context, msgTplName, contentType string, recipientType uint, email string, args ...any) error {
	emailCfg, err := s.configServ.GetEmailConfig(ctx)
	if err != nil {
		return err
	}
	webNMasterCfg, err := s.configServ.GetWebSiteConfig(ctx)
	if err != nil {
		return err
	}
	if email == "" {
		email = emailCfg.Email
	}
	msgTpl, err := s.msgTplService.FindMsgTplByNameAndRcpType(ctx, msgTplName, recipientType)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		msgTpl.FormatContent(args...)
	}
	return s.emailServ.SendEmail(ctx, domain.Email{
		Host:        emailCfg.Host,
		Port:        emailCfg.Port,
		Username:    emailCfg.Username,
		Password:    emailCfg.Password,
		Name:        webNMasterCfg.WebsiteName,
		To:          []string{email},
		Subject:     msgTpl.Title,
		Body:        msgTpl.Content,
		ContentType: contentType,
	})
}

func (s *MessageService) SendEmailWithEmail(ctx context.Context, msgTplName string, email, contentType string, args ...any) error {
	return s.sendEmail(ctx, msgTplName, contentType, 1, email, args...)
}
