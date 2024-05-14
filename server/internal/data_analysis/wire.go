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

package data_analysis

import (
	"github.com/chenmingyong0423/fnote/server/internal/comment"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats"
	service2 "github.com/chenmingyong0423/fnote/server/internal/data_analysis/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/data_analysis/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/post_like"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log/service"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var DataAnalysisProviders = wire.NewSet(web.NewDataAnalysisHandler, service2.NewIpApiService,
	wire.Bind(new(service2.IIpApiService), new(*service2.IpApiService)))

func InitDataAnalysisModule(mongoDB *mongo.Database, vlServ service.IVisitLogService, countStatsModule *count_stats.Module, posLikeModule *post_like.Module, commentModule *comment.Module) *Module {
	panic(wire.Build(
		DataAnalysisProviders,
		wire.FieldsOf(new(*post_like.Module), "Svc"),
		wire.FieldsOf(new(*comment.Module), "Svc"),
		wire.FieldsOf(new(*count_stats.Module), "Svc"),
		wire.Struct(new(Module), "Hdl"),
	))
}
