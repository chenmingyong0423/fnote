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

package domain

type LatestComment struct {
	PostInfo
	Name       string
	Content    string
	Email      string
	CreateTime int64
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
	UserInfo   UserInfo
	Status     CommentStatus
	CreateTime int64
}

func (c *Comment) IsApproved() bool {
	return c.Status == CommentStatusApproved
}
func (c *Comment) IsDisapproved() bool {
	return c.Status == CommentStatusDisapproved
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
	Status          CommentStatus
	CreateTime      int64
}

type CommentReplyWithPostInfo struct {
	CommentReply
	PostInfo PostInfo
}

func (c *CommentReply) IsApproved() bool {
	return c.Status == CommentStatusApproved
}

func (c *CommentReply) IsDisapproved() bool {
	return c.Status == CommentStatusDisapproved
}

type UserInfo4Reply UserInfo

type PostInfo struct {
	// 文章 ID
	PostId string
	// 文章标题字段
	PostTitle string
	// 文章链接
	PostUrl string
}

type UserInfo struct {
	Name    string
	Email   string
	Ip      string
	Website string
}

type CommentStatus uint

const (
	// CommentStatusPending 审核中
	CommentStatusPending CommentStatus = iota
	// CommentStatusApproved 审核通过
	CommentStatusApproved
	// CommentStatusHidden 隐藏
	CommentStatusHidden
	// CommentStatusDisapproved 审核不通过
	CommentStatusDisapproved
)

type AdminComment struct {
	Id string `json:"_id"`
	// 评论的内容
	PostInfo   PostInfo `json:"post_info"`
	Content    string   `json:"content"`
	UserInfo   UserInfo `json:"user_info"`
	Fid        string   `json:"fid"`
	Type       int      `json:"type"`
	Status     int      `json:"status"`
	CreateTime int64    `json:"create_time"`
}
