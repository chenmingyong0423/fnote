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

package friend

import (
	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/friend/internal/web"
	msgService "github.com/chenmingyong0423/fnote/server/internal/message/service"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var FriendProviders = wire.NewSet(web.NewFriendHandler, service.NewFriendService, repository.NewFriendRepository, dao.NewFriendDao,
	wire.Bind(new(service.IFriendService), new(*service.FriendService)),
	wire.Bind(new(repository.IFriendRepository), new(*repository.FriendRepository)),
	wire.Bind(new(dao.IFriendDao), new(*dao.FriendDao)))

func InitFriendModule(mongoDB *mongo.Database, msgServ msgService.IMessageService, cfgModule *website_config.Module) *Module {
	panic(wire.Build(
		FriendProviders,
		wire.FieldsOf(new(*website_config.Module), "Svc"),
		wire.Struct(new(Module), "Svc", "Hdl"),
	))
}
