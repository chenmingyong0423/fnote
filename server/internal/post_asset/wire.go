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

package post_asset

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostAssetProviders = wire.NewSet(web.NewPostAssetHandler, service.NewPostAssetService, repository.NewPostAssetRepository, dao.NewPostAssetDao,
	wire.Bind(new(service.IPostAssetService), new(*service.PostAssetService)),
	wire.Bind(new(repository.IPostAssetRepository), new(*repository.PostAssetRepository)),
	wire.Bind(new(dao.IPostAssetDao), new(*dao.PostAssetDao)),

	service.NewAssetFolderService, repository.NewAssetFolderRepository, dao.NewAssetFolderDao,
	wire.Bind(new(service.IAssetFolderService), new(*service.AssetFolderService)),
	wire.Bind(new(repository.IAssetFolderRepository), new(*repository.AssetFolderRepository)),
	wire.Bind(new(dao.IAssetFolderDao), new(*dao.AssetFolderDao)),
)

func InitPostAssetModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		PostAssetProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
