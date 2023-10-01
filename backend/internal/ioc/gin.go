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

	"github.com/chenmingyong0423/fnote/backend/internal/message_template/handler"

	"github.com/chenmingyong0423/ginx/middlewares/id"
	"github.com/chenmingyong0423/ginx/middlewares/log"

	ctgHandler "github.com/chenmingyong0423/fnote/backend/internal/category/handler"
	commentHandler "github.com/chenmingyong0423/fnote/backend/internal/comment/hanlder"
	cfgHandler "github.com/chenmingyong0423/fnote/backend/internal/config/handler"
	friendHanlder "github.com/chenmingyong0423/fnote/backend/internal/friend/hanlder"
	myValidator "github.com/chenmingyong0423/fnote/backend/internal/pkg/validator"
	postHanlder "github.com/chenmingyong0423/fnote/backend/internal/post/handler"
	vlHandler "github.com/chenmingyong0423/fnote/backend/internal/visit_log/handler"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewGinEngine(ctgHdr *ctgHandler.CategoryHandler, cmtHdr *commentHandler.CommentHandler, cfgHdr *cfgHandler.ConfigHandler, frdHdr *friendHanlder.FriendHandler, postHdr *postHanlder.PostHandler, vlHdr *vlHandler.VisitLogHandler, msgTplHandler *handler.MsgTplHandler, middleware []gin.HandlerFunc, validators Validators) (*gin.Engine, error) {
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
	}
	return engine, nil
}

func InitMiddlewares(cfg *Config, writer io.Writer) []gin.HandlerFunc {
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
		}(cfg.Logger.Level))),
		cors.New(cors.Config{
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if slices.Contains(cfg.Gin.AllowedOrigins, "*") {
					return true
				}
				return slices.ContainsFunc(cfg.Gin.AllowedOrigins, func(s string) bool {
					return strings.Contains(origin, s)
				})
			},
			MaxAge: 12 * time.Hour,
		}),
	}
}

type Validators map[string]func(fl validator.FieldLevel) bool

func InitGinValidators() Validators {
	return map[string]func(fl validator.FieldLevel) bool{
		"validateEmailFormat": myValidator.ValidateEmailFormat,
	}
}
