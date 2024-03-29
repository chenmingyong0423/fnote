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

package ioc

import (
	"io"
	"log/slog"
	"slices"
	"strings"
	"time"

	handler6 "github.com/chenmingyong0423/fnote/server/internal/backup/handler"

	handler5 "github.com/chenmingyong0423/fnote/server/internal/count_stats/handler"

	handler4 "github.com/chenmingyong0423/fnote/server/internal/data_analysis/handler"

	"github.com/spf13/viper"

	handler3 "github.com/chenmingyong0423/fnote/server/internal/file/handler"

	handler2 "github.com/chenmingyong0423/fnote/server/internal/tag/handler"

	"github.com/chenmingyong0423/fnote/server/internal/message_template/handler"
	"github.com/chenmingyong0423/ginx/middlewares/id"
	"github.com/chenmingyong0423/ginx/middlewares/log"
	"github.com/gin-contrib/cors"

	ctgHandler "github.com/chenmingyong0423/fnote/server/internal/category/handler"
	commentHandler "github.com/chenmingyong0423/fnote/server/internal/comment/hanlder"
	friendHanlder "github.com/chenmingyong0423/fnote/server/internal/friend/hanlder"
	myValidator "github.com/chenmingyong0423/fnote/server/internal/pkg/validator"
	postHanlder "github.com/chenmingyong0423/fnote/server/internal/post/handler"
	vlHandler "github.com/chenmingyong0423/fnote/server/internal/visit_log/handler"
	cfgHandler "github.com/chenmingyong0423/fnote/server/internal/website_config/handler"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewGinEngine(fileHdr *handler3.FileHandler, ctgHdr *ctgHandler.CategoryHandler, cmtHdr *commentHandler.CommentHandler, cfgHdr *cfgHandler.WebsiteConfigHandler, frdHdr *friendHanlder.FriendHandler, postHdr *postHanlder.PostHandler, vlHdr *vlHandler.VisitLogHandler, msgTplHandler *handler.MsgTplHandler, tagsHandler *handler2.TagHandler, daHandler *handler4.DataAnalysisHandler, csHandler *handler5.CountStatsHandler, backupHandler *handler6.BackupHandler, middleware []gin.HandlerFunc, validators Validators) (*gin.Engine, error) {
	engine := gin.Default()

	// 参数校验器注册
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		for k, v := range validators {
			err := validate.RegisterValidation(k, v)
			if err != nil {
				return nil, err
			}
		}
	}

	// 中间件注册
	engine.Use(middleware...)

	// 注册路由
	{
		ctgHdr.RegisterGinRoutes(engine)
		cmtHdr.RegisterGinRoutes(engine)
		cfgHdr.RegisterGinRoutes(engine)
		frdHdr.RegisterGinRoutes(engine)
		postHdr.RegisterGinRoutes(engine)
		vlHdr.RegisterGinRoutes(engine)
		msgTplHandler.RegisterGinRoutes(engine)
		tagsHandler.RegisterGinRoutes(engine)
		fileHdr.RegisterGinRoutes(engine)
		daHandler.RegisterGinRoutes(engine)
		csHandler.RegisterGinRoutes(engine)
		backupHandler.RegisterGinRoutes(engine)
	}
	return engine, nil
}

func InitMiddlewares(writer io.Writer) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.LoggerWithWriter(writer),
		id.RequestId(),
		log.RequestLogger(*log.NewLoggerConfig(func(level string) slog.Level {
			switch level {
			case "DEBUG":
				return slog.LevelDebug
			case "INFO":
				return slog.LevelInfo
			case "WARN":
				return slog.LevelWarn
			case "ERROR":
				return slog.LevelError
			default:
				return slog.LevelInfo
			}
		}(viper.GetString("logger.level")))),
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if slices.Contains(viper.GetStringSlice("gin.allowed_origins"), "*") {
					return true
				}
				return slices.ContainsFunc(viper.GetStringSlice("gin.allowed_origins"), func(s string) bool {
					return strings.Contains(origin, s)
				})
			},
			AllowMethods: viper.GetStringSlice("gin.allowed_methods"),
			AllowHeaders: viper.GetStringSlice("gin.allowed_headers"),
			MaxAge:       12 * time.Hour,
		}),
		JwtParseMiddleware(),
	}
}

type Validators map[string]func(fl validator.FieldLevel) bool

func InitGinValidators() Validators {
	return map[string]func(fl validator.FieldLevel) bool{
		"validateEmailFormat": myValidator.ValidateEmailFormat,
	}
}
