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

package main

import (
	"context"
	"errors"
	emailCfg "github.com/chenmingyong0423/fnote/backend/ineternal/email/service"
	friendHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/friend/hanlder"
	friendRepo "github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository"
	friendDao "github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository/dao"
	friendServ "github.com/chenmingyong0423/fnote/backend/ineternal/friend/service"
	myValidator "github.com/chenmingyong0423/fnote/backend/ineternal/pkg/validator"
	vlHandler "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/handler"
	vlRepo "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/repository"
	vlLogDao "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/repository/dao"
	vlServ "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/service"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"os"
	"strings"
	"time"

	ctgHandler "github.com/chenmingyong0423/fnote/backend/ineternal/category/handler"
	ctgRepo "github.com/chenmingyong0423/fnote/backend/ineternal/category/repository"
	ctgDao "github.com/chenmingyong0423/fnote/backend/ineternal/category/repository/dao"
	ctgService "github.com/chenmingyong0423/fnote/backend/ineternal/category/service"
	cHandler "github.com/chenmingyong0423/fnote/backend/ineternal/config/handler"
	cRepository "github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	cDao "github.com/chenmingyong0423/fnote/backend/ineternal/config/repository/dao"
	cService "github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	postHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/post/handler"
	postRepo "github.com/chenmingyong0423/fnote/backend/ineternal/post/repository"
	postDao "github.com/chenmingyong0423/fnote/backend/ineternal/post/repository/dao"
	postServ "github.com/chenmingyong0423/fnote/backend/ineternal/post/service"

	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/middleware"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if len(os.Args) < 3 {
		panic(errors.New("missing parameters"))
	}
	username := os.Args[1]
	password := os.Args[2]
	db := initDb(username, password)

	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("validateEmailFormat", myValidator.ValidateEmailFormat)
		if err != nil {
			panic(err)
		}
	}

	r.Use(middleware.RequestId())
	r.Use(middleware.Logger())

	r.Use(cors.New(cors.Config{
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
	cfgServ := cService.NewConfigService(cRepository.NewConfigRepository(cDao.NewConfigDao(db.Collection("configs"))))
	cHandler.NewConfigHandler(r, cfgServ)
	ctgHandler.NewCategoryHandler(r, ctgService.NewCategoryService(ctgRepo.NewCategoryRepository(ctgDao.NewCategoryDao(db.Collection("categories")))))
	postHanlder.NewPostHandler(r, postServ.NewPostService(postRepo.NewPostRepository(postDao.NewPostDao(db.Collection("posts")))))
	vlHandler.NewVisitLogHandler(r, vlServ.NewVisitLogService(vlRepo.NewVisitLogRepository(vlLogDao.NewVisitLogDao(db.Collection("visit_logs")))), cfgServ)
	emailServ := emailCfg.NewEmailService()
	friendHanlder.NewFriendHandler(r, friendServ.NewFriendService(friendRepo.NewFriendRepository(friendDao.NewFriendDao(db.Collection("friends"))), emailServ, cfgServ))
	err := r.Run()
	if err != nil {
		panic(err)
	}
}

func initDb(username, password string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(options.Credential{
		Username:   username,
		Password:   password,
		AuthSource: "fnote",
	}).SetDirect(true))
	if err != nil {
		panic(err)
	}
	return client.Database("fnote")
}
