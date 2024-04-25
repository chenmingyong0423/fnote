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

package post_like

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_like/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostLikeProviders = wire.NewSet(web.NewPostLikeHandler, service.NewPostLikeService, repository.NewPostLikeRepository, dao.NewPostLikeDao,
	wire.Bind(new(service.IPostLikeService), new(*service.PostLikeService)),
	wire.Bind(new(repository.IPostLikeRepository), new(*repository.PostLikeRepository)),
	wire.Bind(new(dao.IPostLikeDao), new(*dao.PostLikeDao)))

func InitPostLikeModule(mongoDB *mongo.Database) *Model {
	panic(wire.Build(
		PostLikeProviders,
		wire.Struct(new(Model), "Svc", "Hdl"),
	))
}
