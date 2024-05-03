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

package post_visit

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostVisitProviders = wire.NewSet(web.NewPostVisitHandler, service.NewPostVisitService, repository.NewPostVisitRepository, dao.NewPostVisitDao,
	wire.Bind(new(service.IPostVisitService), new(*service.PostVisitService)),
	wire.Bind(new(repository.IPostVisitRepository), new(*repository.PostVisitRepository)),
	wire.Bind(new(dao.IPostVisitDao), new(*dao.PostVisitDao)))

func InitPostVisitModule(mongoDB *mongo.Database) *Module {
	panic(wire.Build(
		PostVisitProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
