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

package category

import (
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/category/internal/web"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var CategoryProviders = wire.NewSet(web.NewCategoryHandler, service.NewCategoryService, repository.NewCategoryRepository, dao.NewCategoryDao,
	wire.Bind(new(service.ICategoryService), new(*service.CategoryService)),
	wire.Bind(new(repository.ICategoryRepository), new(*repository.CategoryRepository)),
	wire.Bind(new(dao.ICategoryDao), new(*dao.CategoryDao)))

func InitCategoryModule(mongoDB *mongo.Database, eventBus *eventbus.EventBus) *Module {
	panic(wire.Build(
		CategoryProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
