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

package count_stats

import (
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats/internal/web"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var CountStatsProviders = wire.NewSet(web.NewCountStatsHandler, service.NewCountStatsService, repository.NewCountStatsRepository, dao.NewCountStatsDao,
	wire.Bind(new(service.ICountStatsService), new(*service.CountStatsService)),
	wire.Bind(new(repository.ICountStatsRepository), new(*repository.CountStatsRepository)),
	wire.Bind(new(dao.ICountStatsDao), new(*dao.CountStatsDao)))

func InitCountStatsModule(mongoDB *mongo.Database, eventbus *eventbus.EventBus) *Module {
	panic(wire.Build(
		CountStatsProviders,
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
