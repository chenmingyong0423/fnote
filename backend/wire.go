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
	commentHandler "github.com/chenmingyong0423/fnote/backend/ineternal/comment/hanlder"
	cfgHandler "github.com/chenmingyong0423/fnote/backend/ineternal/config/handler"
	emailServ "github.com/chenmingyong0423/fnote/backend/ineternal/email/service"
	friendHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/friend/hanlder"
	msgServ "github.com/chenmingyong0423/fnote/backend/ineternal/message/service"
	postHanlder "github.com/chenmingyong0423/fnote/backend/ineternal/post/handler"
	vlHandler "github.com/chenmingyong0423/fnote/backend/ineternal/visit_log/handler"
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
		emailServ.EmailSet,
		msgServ.MsgSet,
		ioc.NewMongoDB,
		ioc.NewGinEngine,
	))
}
