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

package post_draft

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft/internal/web"
	"github.com/google/wire"
	"go.mongodb.org/mongo-driver/mongo"
)

var PostDraftProviders = wire.NewSet(web.NewPostDraftHandler, service.NewPostDraftService, repository.NewPostDraftRepository, dao.NewPostDraftDao,
	wire.Bind(new(service.IPostDraftService), new(*service.PostDraftService)),
	wire.Bind(new(repository.IPostDraftRepository), new(*repository.PostDraftRepository)),
	wire.Bind(new(dao.IPostDraftDao), new(*dao.PostDraftDao)))

func InitPostDraftModule(mongoDB *mongo.Database) *Model {
	panic(wire.Build(
		PostDraftProviders,
		wire.Struct(new(Model), "Svc", "Hdl"),
	))
}
