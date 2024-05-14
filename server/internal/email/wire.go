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

package email

import (
	"github.com/chenmingyong0423/fnote/server/internal/email/internal/service"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var EmailProviders = wire.NewSet(service.NewEmailService, wire.Bind(new(service.IEmailService), new(*service.EmailService)))

func InitEmailModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		EmailProviders,
		wire.Struct(new(Module), "Svc"),
	))
}
