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

package file

import (
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/file/internal/web"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var FileProviders = wire.NewSet(web.NewFileHandler, service.NewFileService, repository.NewFileRepository, dao.NewFileDao,
	wire.Bind(new(service.IFileService), new(*service.FileService)),
	wire.Bind(new(repository.IFileRepository), new(*repository.FileRepository)),
	wire.Bind(new(dao.IFileDao), new(*dao.FileDao)))

func InitFileModule(mongoDB *mongo.Database, eventBus *eventbus.EventBus) *Module {
	panic(wire.Build(
		FileProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
