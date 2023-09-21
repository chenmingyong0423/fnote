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

// Package types 存放公共属性的包
package types

import (
	"github.com/gin-gonic/gin"
)

type CommentReply struct {
	ReplyId string `bson:"reply_id"`
	// 回复内容
	Content string `bson:"content"`
	// 被回复的回复 Id
	ReplyToId string `bson:"reply_to_id"`
	// 用户信息
	UserInfo UserInfo4Reply `bson:"user_info"`
	// 被回复用户的信息
	RepliedUserInfo UserInfo4Reply `bson:"replied_user_info"`
}
type UserInfo4Reply UserInfo4Comment

type Comment struct {
	// 文章信息
	PostInfo PostInfo4Comment `bson:"post_info"`
	// 评论的内容
	Content string `bson:"content"`
	// 用户信息
	UserInfo UserInfo4Comment `bson:"user_info"`
}

type PostInfo4Comment struct {
	// 文章 ID
	PostId string `bson:"post_id"`
	// 文章标题字段
	PostTitle string `bson:"post_title"`
}

type UserInfo4Comment struct {
	Name    string `bson:"name"`
	Email   string `bson:"email"`
	Ip      string `bson:"ip"`
	Website string `bson:"website"`
}

type GinRoutes interface {
	RegisterGinRoutes(engine *gin.Engine)
}
