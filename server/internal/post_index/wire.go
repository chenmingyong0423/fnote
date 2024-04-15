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

package post_index

import (
	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/google/wire"
)

var ConfigProviders = wire.NewSet(web.NewPostIndexHandler, service.NewPostIndexService, service.NewBaiduService,
	wire.Bind(new(service.IPostIndexService), new(*service.PostIndexService)))

func InitPostIndexModule(cfgServ website_config.Service) Model {
	panic(wire.Build(
		ConfigProviders,
		wire.Struct(new(Model), "Svc", "Hdl"),
	))
}
