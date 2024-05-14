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

package domain

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

type Page struct {
	Size           int64
	Skip           int64
	Sort           string
	ApprovalStatus *bool
}

func (p *Page) SortToBson() bson.D {
	sort := p.Sort
	if p.Sort == "" {
		return nil
	}
	split := strings.Split(sort, ",")
	var sortBson bson.D
	for _, s := range split {
		if strings.HasPrefix(s, "+") {
			sortBson = append(sortBson, bson.E{Key: strings.TrimLeft(s, "-"), Value: 1})
		} else {
			sortBson = append(sortBson, bson.E{Key: s, Value: -1})
		}
	}
	return sortBson
}

type AdminComment struct {
	Id string
	// 文章信息
	PostInfo PostInfo
	// 评论的内容
	Content string
	// 用户信息
	UserInfo UserInfo4Comment

	// 该评论下的所有回复的内容
	Replies        []AdminReply
	ApprovalStatus bool
	// 评论时间
	CreatedAt int64
	// 修改时间
	UpdatedAt int64
}

type PostInfo struct {
	// 文章 ID
	PostId string
	// 文章标题字段
	PostTitle string
	// 文章链接
	PostUrl string
}

type UserInfo4Reply UserInfo

type UserInfo4Comment UserInfo

type UserInfo struct {
	Name    string
	Email   string
	Ip      string
	Website string
}

type AdminReply struct {
	ReplyId string
	// 回复内容
	Content string
	// 被回复的回复 Id
	ReplyToId string
	// 用户信息
	UserInfo UserInfo4Reply
	// 被回复用户的信息
	RepliedUserInfo UserInfo4Reply
	ApprovalStatus  bool
	// 回复时间
	CreatedAt int64
	// 修改时间
	UpdatedAt int64
}

type LatestComment struct {
	PostInfo
	Name      string
	Content   string
	Email     string
	CreatedAt int64
}

type CommentWithReplies struct {
	Comment
	Replies []CommentReply
}

type Comment struct {
	Id string
	// 文章信息
	PostInfo PostInfo
	// 评论的内容
	Content string
	// 用户信息
	UserInfo       UserInfo
	ApprovalStatus bool
	CreateTime     int64
}

type CommentReply struct {
	ReplyId string
	// 回复内容
	Content string
	// 被回复的回复 Id
	ReplyToId string
	// 用户信息
	UserInfo UserInfo4Reply
	// 被回复用户的信息
	RepliedUserInfo UserInfo4Reply
	ApprovalStatus  bool
	CreatedAt       int64
}

type CommentReplyWithPostInfo struct {
	CommentReply
	PostInfo PostInfo
}

type EmailInfo struct {
	Email   string
	PostUrl string
}

type ReplyWithCId struct {
	CommentId string
	ReplyIds  []string
}

type CommentEvent struct {
	PostId    string   `json:"post_id"`
	CommentId string   `json:"comment_id"`
	RepliesId []string `json:"replies_id"`
	Count     int      `json:"count"`
	Type      string   `json:"type"`
}
