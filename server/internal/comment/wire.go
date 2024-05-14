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

package comment

import (
	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/comment/internal/web"
	msgService "github.com/chenmingyong0423/fnote/server/internal/message/service"
	"github.com/chenmingyong0423/fnote/server/internal/post"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/chenmingyong0423/go-eventbus"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var CommentProviders = wire.NewSet(web.NewCommentHandler, service.NewCommentService, repository.NewCommentRepository, dao.NewCommentDao,
	wire.Bind(new(service.ICommentService), new(*service.CommentService)),
	wire.Bind(new(repository.ICommentRepository), new(*repository.CommentRepository)),
	wire.Bind(new(dao.ICommentDao), new(*dao.CommentDao)))

func InitCommentModule(mongoDB *mongo.Database, msgServ msgService.IMessageService, cfgModule *website_config.Module, postModule *post.Module, eventBus *eventbus.EventBus) *Module {
	panic(wire.Build(
		CommentProviders,
		wire.FieldsOf(new(*website_config.Module), "Svc"),
		wire.FieldsOf(new(*post.Module), "Svc"),
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
