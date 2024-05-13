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

package post

import (
	csServ "github.com/chenmingyong0423/fnote/server/internal/count_stats/service"
	fServ "github.com/chenmingyong0423/fnote/server/internal/file/service"
	"github.com/chenmingyong0423/fnote/server/internal/post/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/post/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/post/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/post_like"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostProviders = wire.NewSet(web.NewPostHandler, service.NewPostService, repository.NewPostRepository, dao.NewPostDao,
	wire.Bind(new(service.IPostService), new(*service.PostService)),
	wire.Bind(new(repository.IPostRepository), new(*repository.PostRepository)),
	wire.Bind(new(dao.IPostDao), new(*dao.PostDao)))

func InitPostModule(mongoDB *mongo.Database, cfgModel *website_config.Module, countStats csServ.ICountStatsService, fileService fServ.IFileService, postLikeModel *post_like.Module, eventBus *eventbus.EventBus) *Module {
	panic(wire.Build(
		PostProviders,
		wire.FieldsOf(new(*website_config.Module), "Svc"),
		wire.FieldsOf(new(*post_like.Module), "Svc"),
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
