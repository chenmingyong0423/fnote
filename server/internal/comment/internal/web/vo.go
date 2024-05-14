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

package web

type LatestCommentVO struct {
	PostInfo
	Picture   string `json:"picture"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type PostInfo struct {
	// 文章 ID
	PostId string `json:"post_id"`
	// 文章标题字段
	PostTitle string `json:"post_title"`
	// 文章链接
	PostUrl string `json:"post_url"`
}

type PostCommentVO struct {
	Id string `json:"id"`
	// 评论的内容
	Content string `json:"content"`
	// 评论的用户
	Name        string               `json:"username"`
	Picture     string               `json:"picture"`
	Website     string               `json:"website"`
	CommentTime int64                `json:"comment_time"`
	Replies     []PostCommentReplyVO `json:"replies,omitempty"`
}

type PostCommentReplyVO struct {
	Id        string `json:"id"`
	CommentId string `json:"comment_id"`
	// 回复的内容
	Content string `json:"content"`
	// 回复的用户
	Name    string `json:"name"`
	Picture string `json:"picture"`
	Website string `json:"website"`
	// 被回复的回复 Id
	ReplyToId string `json:"reply_to_id"`
	// 被回复的用户
	ReplyTo string `json:"reply_to"`
	// 回复时间
	ReplyTime int64 `json:"reply_time"`
}

type AdminCommentVO struct {
	Id string `json:"id"`
	// 文章信息
	PostInfo PostInfo `json:"post_info"`
	// 评论的内容
	Content string `json:"content"`
	// 用户信息
	UserInfo UserInfo4Comment `json:"user_info"`

	ReplyCount int `json:"reply_count"`

	// 该评论下的所有回复的内容
	Replies        []AdminCommentVO `json:"replies,omitempty"`
	ApprovalStatus bool             `json:"approval_status"`
	Type           string           `json:"type"`
	// 评论时间
	CreatedAt int64 `json:"created_at"`
	// 修改时间
	UpdatedAt int64 `json:"updated_at"`

	// 被回复的回复 Id
	ReplyToId string `json:"reply_to_id"`
}

type AdminCommentReplyVO struct {
	ReplyId string `json:"reply_id"`
	// 回复内容
	Content string `json:"content"`
	// 被回复的回复 Id
	ReplyToId string `json:"reply_to_id"`
	// 用户信息
	UserInfo UserInfo4Reply `json:"user_info"`
	// 被回复用户的信息
	RepliedUserInfo UserInfo4Reply `json:"replied_user_info"`
	ApprovalStatus  bool           `json:"approval_status"`
	// 回复时间
	CreatedAt int64 `json:"created_at"`
	// 修改时间
	UpdatedAt int64 `json:"updated_at"`
}

type UserInfo4Reply UserInfo

type UserInfo4Comment UserInfo

type UserInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Ip      string `json:"ip"`
	Picture string `json:"picture"`
	Website string `json:"website"`
}

type AdminUserInfoVO UserInfoVO

type UserInfoVO struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Ip      string `json:"ip"`
	Website string `json:"website"`
}

type IdVO struct {
	Id string `json:"id"`
}
