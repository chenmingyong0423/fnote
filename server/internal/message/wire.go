// Copyright 2024 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build wireinject

package message

import (
	emailServ "github.com/chenmingyong0423/fnote/server/internal/email/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/message/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/message_template"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/google/wire"
)

var MessageProviders = wire.NewSet(service.NewMessageService,
	wire.Bind(new(service.IMessageService), new(*service.MessageService)),
)

func InitMessageModule(emailServ emailServ.IEmailService, messageTemplateModule *message_template.Module, websiteConfigModule *website_config.Module) *Module {
	panic(wire.Build(
		MessageProviders,
		wire.FieldsOf(new(*message_template.Module), "Svc"),
		wire.FieldsOf(new(*website_config.Module), "Svc"),
		wire.Struct(new(Module), "Svc"),
	))
}
