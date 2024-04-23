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

package website_config

import (
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var ConfigProviders = wire.NewSet(web.NewWebsiteConfigHandler, service.NewWebsiteConfigService, repository.NewWebsiteConfigRepository, dao.NewWebsiteConfigDao,
	wire.Bind(new(service.IWebsiteConfigService), new(*service.WebsiteConfigService)),
	wire.Bind(new(repository.IWebsiteConfigRepository), new(*repository.WebsiteConfigRepository)),
	wire.Bind(new(dao.IWebsiteConfigDao), new(*dao.WebsiteConfigDao)))

func InitWebsiteConfigModule(mongoDB *mongo.Database) *Model {
	panic(wire.Build(
		ConfigProviders,
		wire.Struct(new(Model), "Svc", "Hdl"),
	))
}
