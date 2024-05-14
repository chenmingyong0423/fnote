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

package tag

import (
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/tag/internal/web"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var TagProviders = wire.NewSet(web.NewTagHandler, service.NewTagService, repository.NewTagRepository, dao.NewTagDao,
	wire.Bind(new(service.ITagService), new(*service.TagService)),
	wire.Bind(new(repository.ITagRepository), new(*repository.TagRepository)),
	wire.Bind(new(dao.ITagDao), new(*dao.TagDao)))

func InitTagModule(mongoDB *mongo.Database, eventBus *eventbus.EventBus) *Module {
	panic(wire.Build(
		TagProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
