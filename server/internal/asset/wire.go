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

package asset

import (
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/asset/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var AssetProviders = wire.NewSet(web.NewAssetHandler, service.NewAssetService, repository.NewAssetRepository, dao.NewAssetDao,
	wire.Bind(new(service.IAssetService), new(*service.AssetService)),
	wire.Bind(new(repository.IAssetRepository), new(*repository.AssetRepository)),
	wire.Bind(new(dao.IAssetDao), new(*dao.AssetDao)),

	repository.NewAssetFolderRepository, dao.NewAssetFolderDao,
	wire.Bind(new(repository.IAssetFolderRepository), new(*repository.AssetFolderRepository)),
	wire.Bind(new(dao.IAssetFolderDao), new(*dao.AssetFolderDao)),
)

func InitAssetModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		AssetProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
