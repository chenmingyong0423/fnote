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
	ctgHandler "github.com/chenmingyong0423/fnote/backend/ineternal/category/handler"
	ctgRepo "github.com/chenmingyong0423/fnote/backend/ineternal/category/repository"
	ctgDao "github.com/chenmingyong0423/fnote/backend/ineternal/category/repository/dao"
	ctgServ "github.com/chenmingyong0423/fnote/backend/ineternal/category/service"
	commentHandler "github.com/chenmingyong0423/fnote/backend/ineternal/comment/hanlder"
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/repository/dao"
	"github.com/chenmingyong0423/fnote/backend/ineternal/comment/service"
	cfgHandler "github.com/chenmingyong0423/fnote/backend/ineternal/config/handler"
	repository2 "github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	dao2 "github.com/chenmingyong0423/fnote/backend/ineternal/config/repository/dao"
	service2 "github.com/chenmingyong0423/fnote/backend/ineternal/config/service"
	emailServ "github.com/chenmingyong0423/fnote/backend/ineternal/email/service"
	friendHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/friend/hanlder"
	repository3 "github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository"
	dao3 "github.com/chenmingyong0423/fnote/backend/ineternal/friend/repository/dao"
	service3 "github.com/chenmingyong0423/fnote/backend/ineternal/friend/service"
	msgServ "github.com/chenmingyong0423/fnote/backend/ineternal/message/service"
	postHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/post/handler"
	repository4 "github.com/chenmingyong0423/fnote/backend/ineternal/post/repository"
	dao4 "github.com/chenmingyong0423/fnote/backend/ineternal/post/repository/dao"
	service4 "github.com/chenmingyong0423/fnote/backend/ineternal/post/service"
	vlHandler "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/handler"
	repository5 "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/repository"
	dao5 "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/repository/dao"
	service5 "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/service"
	"github.com/chenmingyong0423/fnote/backend/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func initializeApp(username ioc.Username, password ioc.Password) (*gin.Engine, error) {
	panic(wire.Build(
		ctgHandler.CategorySet,
		commentHandler.CommentSet,
		cfgHandler.ConfigSet,
		friendHanlder.FriendSet,
		postHanlder.PostSet,
		vlHandler.VlSet,
		emailServ.NewEmailService,
		msgServ.NewMessageService,
		ioc.NewMongoDB,
		ioc.NewGinEngine,
		// 绑定接口
		wire.Bind(new(ctgServ.ICategoryService), new(*ctgServ.CategoryService)),
		wire.Bind(new(ctgRepo.ICategoryRepository), new(*ctgRepo.CategoryRepository)),
		wire.Bind(new(ctgDao.ICategoryDao), new(*ctgDao.CategoryDao)),

		wire.Bind(new(service.ICommentService), new(*service.CommentService)),
		wire.Bind(new(repository.ICommentRepository), new(*repository.CommentRepository)),
		wire.Bind(new(dao.ICommentDao), new(*dao.CommentDao)),

		wire.Bind(new(service2.IConfigService), new(*service2.ConfigService)),
		wire.Bind(new(repository2.IConfigRepository), new(*repository2.ConfigRepository)),
		wire.Bind(new(dao2.IConfigDao), new(*dao2.ConfigDao)),

		wire.Bind(new(emailServ.IEmailService), new(*emailServ.EmailService)),

		wire.Bind(new(service3.IFriendService), new(*service3.FriendService)),
		wire.Bind(new(repository3.IFriendRepository), new(*repository3.FriendRepository)),
		wire.Bind(new(dao3.IFriendDao), new(*dao3.FriendDao)),

		wire.Bind(new(msgServ.IMessageService), new(*msgServ.MessageService)),

		wire.Bind(new(service4.IPostService), new(*service4.PostService)),
		wire.Bind(new(repository4.IPostRepository), new(*repository4.PostRepository)),
		wire.Bind(new(dao4.IPostDao), new(*dao4.PostDao)),

		wire.Bind(new(service5.IVisitLogService), new(*service5.VisitLogService)),
		wire.Bind(new(repository5.IVisitLogRepository), new(*repository5.VisitLogRepository)),
		wire.Bind(new(dao5.IVisitLogDao), new(*dao5.VisitLogDao)),
	))
}
