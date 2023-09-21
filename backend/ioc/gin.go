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
	"strings"
	"time"

	ctgHandler "github.com/chenmingyong0423/fnote/backend/internal/category/handler"
	commentHandler "github.com/chenmingyong0423/fnote/backend/internal/comment/hanlder"
	cfgHandler "github.com/chenmingyong0423/fnote/backend/internal/config/handler"
	friendHanlder "github.com/chenmingyong0423/fnote/backend/internal/friend/hanlder"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/middleware"
	myValidator "github.com/chenmingyong0423/fnote/backend/internal/pkg/validator"
	postHanlder "github.com/chenmingyong0423/fnote/backend/internal/post/handler"
	vlHandler "github.com/chenmingyong0423/fnote/backend/internal/visit_log/handler"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewGinEngine(ctgHdr *ctgHandler.CategoryHandler, cmtHdr *commentHandler.CommentHandler, cfgHdr *cfgHandler.ConfigHandler, frdHdr *friendHanlder.FriendHandler, postHdr *postHanlder.PostHandler, vlHdr *vlHandler.VisitLogHandler) (*gin.Engine, error) {
	engine := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validateEmailFormat", myValidator.ValidateEmailFormat)
		if err != nil {
			return nil, err
		}
	}

	engine.Use(middleware.RequestId())
	engine.Use(middleware.Logger())

	engine.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 你的开发环境
				return true
			}
			return strings.Contains(origin, "chenmingyong.cn")
		},
		MaxAge: 12 * time.Hour,
	}))

	// 注册路由
	{
		ctgHdr.RegisterGinRoutes(engine)
		cmtHdr.RegisterGinRoutes(engine)
		cfgHdr.RegisterGinRoutes(engine)
		frdHdr.RegisterGinRoutes(engine)
		postHdr.RegisterGinRoutes(engine)
		vlHdr.RegisterGinRoutes(engine)
	}
	return engine, nil
}
