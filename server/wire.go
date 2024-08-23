// Copyright 2023 chenmingyong0423

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

package main

import (
	"github.com/chenmingyong0423/fnote/server/internal/aggregate_post"
	"github.com/chenmingyong0423/fnote/server/internal/backup"
	"github.com/chenmingyong0423/fnote/server/internal/category"
	"github.com/chenmingyong0423/fnote/server/internal/comment"
	"github.com/chenmingyong0423/fnote/server/internal/count_stats"
	"github.com/chenmingyong0423/fnote/server/internal/data_analysis"
	"github.com/chenmingyong0423/fnote/server/internal/email"
	"github.com/chenmingyong0423/fnote/server/internal/file"
	"github.com/chenmingyong0423/fnote/server/internal/friend"
	"github.com/chenmingyong0423/fnote/server/internal/global"
	"github.com/chenmingyong0423/fnote/server/internal/ioc"
	"github.com/chenmingyong0423/fnote/server/internal/message"
	"github.com/chenmingyong0423/fnote/server/internal/message_template"
	"github.com/chenmingyong0423/fnote/server/internal/post"
	"github.com/chenmingyong0423/fnote/server/internal/post_asset"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft"
	"github.com/chenmingyong0423/fnote/server/internal/post_index"
	"github.com/chenmingyong0423/fnote/server/internal/post_like"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit"
	"github.com/chenmingyong0423/fnote/server/internal/tag"
	"github.com/chenmingyong0423/fnote/server/internal/visit_log"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func initializeApp() (*gin.Engine, error) {
	panic(wire.Build(
		ioc.NewEventBus,
		ioc.InitLogger,
		ioc.NewMongoDB,
		ioc.InitMiddlewares,
		ioc.InitGinValidators,
		ioc.NewGinEngine,
		global.IsWebsiteInitializedFn,

		website_config.InitWebsiteConfigModule,
		wire.FieldsOf(new(*website_config.Module), "Hdl"),
		post_index.InitPostIndexModule,
		wire.FieldsOf(new(*post_index.Module), "Hdl"),
		post_draft.InitPostDraftModule,
		wire.FieldsOf(new(*post_draft.Module), "Hdl"),
		aggregate_post.InitAggregatePostModule,
		wire.FieldsOf(new(*aggregate_post.Module), "Hdl"),
		post_like.InitPostLikeModule,
		wire.FieldsOf(new(*post_like.Module), "Hdl"),
		data_analysis.InitDataAnalysisModule,
		wire.FieldsOf(new(*data_analysis.Module), "Hdl"),
		post.InitPostModule,
		wire.FieldsOf(new(*post.Module), "Hdl"),
		comment.InitCommentModule,
		wire.FieldsOf(new(*comment.Module), "Hdl"),
		post_visit.InitPostVisitModule,
		wire.FieldsOf(new(*post_visit.Module), "Hdl"),
		friend.InitFriendModule,
		wire.FieldsOf(new(*friend.Module), "Hdl"),
		count_stats.InitCountStatsModule,
		wire.FieldsOf(new(*count_stats.Module), "Hdl"),
		file.InitFileModule,
		wire.FieldsOf(new(*file.Module), "Hdl"),
		category.InitCategoryModule,
		wire.FieldsOf(new(*category.Module), "Hdl"),
		tag.InitTagModule,
		wire.FieldsOf(new(*tag.Module), "Hdl"),
		message_template.InitMessageTemplateModule,
		wire.FieldsOf(new(*message_template.Module), "Hdl"),
		visit_log.InitVisitLogModule,
		wire.FieldsOf(new(*visit_log.Module), "Hdl"),
		message.InitMessageModule,
		email.InitEmailModule,
		backup.InitBackupModule,
		wire.FieldsOf(new(*backup.Module), "Hdl"),
		post_asset.InitPostAssetModule,
		wire.FieldsOf(new(*post_asset.Module), "Hdl"),
	))
}
