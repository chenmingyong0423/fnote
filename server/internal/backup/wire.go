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

package backup

import (
	"github.com/chenmingyong0423/fnote/server/internal/backup/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/backup/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var BackupProviders = wire.NewSet(web.NewBackupHandler, service.NewBackupService,
	wire.Bind(new(service.IBackupService), new(*service.BackupService)))

func InitBackupModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		BackupProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
