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

package ioc

import (
	handler12 "github.com/chenmingyong0423/fnote/server/internal/backup/handler"
	service13 "github.com/chenmingyong0423/fnote/server/internal/backup/service"
	service7 "github.com/chenmingyong0423/fnote/server/internal/email/service"
	service8 "github.com/chenmingyong0423/fnote/server/internal/message/service"
	"github.com/google/wire"
)

var (
	EmailProviders = wire.NewSet(service7.NewEmailService, wire.Bind(new(service7.IEmailService), new(*service7.EmailService)))

	MsgProviders = wire.NewSet(service8.NewMessageService, wire.Bind(new(service8.IMessageService), new(*service8.MessageService)))

	BackupProviders = wire.NewSet(handler12.NewBackupHandler, service13.NewBackupService,
		wire.Bind(new(service13.IBackupService), new(*service13.BackupService)),
	)
)
