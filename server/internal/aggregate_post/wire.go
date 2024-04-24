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

package aggregate_post

import (
	"github.com/chenmingyong0423/fnote/server/internal/aggregate_post/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/aggregate_post/internal/web"
	postServ "github.com/chenmingyong0423/fnote/server/internal/post/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft"
	"github.com/google/wire"
)

var AggregatePostProviders = wire.NewSet(web.NewAggregatePostHandler, service.NewAggregatePostService,
	wire.Bind(new(service.IAggregatePostService), new(*service.AggregatePostService)),
)

func InitAggregatePostModule(postServ postServ.IPostService, postDraftModel *post_draft.Model) *Model {
	panic(wire.Build(
		wire.FieldsOf(new(*post_draft.Model), "Svc"),
		AggregatePostProviders,
		wire.Struct(new(Model), "Svc", "Hdl"),
	))
}
