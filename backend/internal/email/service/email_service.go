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

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

type IEmailService interface {
	SendEmail(ctx context.Context, email domain.Email) error
}

var (
	_ IEmailService = (*EmailService)(nil)
)

type EmailService struct {
}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (s *EmailService) SendEmail(ctx context.Context, email domain.Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", email.Account)
	m.SetHeader("To", email.To...)
	m.SetHeader("Subject", email.Subject)
	m.SetBody(email.ContentType, email.Body)
	dialer := gomail.NewDialer(email.Host, email.Port, email.Account, email.Password)
	err := dialer.DialAndSend(m)
	if err != nil {
		return errors.Wrap(err, "dialer.DialAndSend failed")
	}
	return nil
}
