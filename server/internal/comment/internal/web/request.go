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

type CommentRequest struct {
	PostId   string `json:"postId" binding:"required"`
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,validateEmailFormat"`
	Website  string `json:"website"`
	Content  string `json:"content" binding:"required,max=200"`
}

type PageRequest struct {
	// 当前页
	PageNo int64 `form:"pageNo" binding:"required"`
	// 每页数量
	PageSize int64 `form:"pageSize" binding:"required"`
	// 排序字段
	Field string `form:"sortField,omitempty"`
	// 排序规则
	Order string `form:"sortOrder,omitempty"`
	// 搜索内容
	Keyword string `form:"keyword,omitempty"`
}

type BatchApprovedCommentRequest struct {
	CommentIds []string `json:"comment_ids" binding:"required,lt=0"`
	// key 为 commentId, value 为 reply_ids
	Replies map[string][]string `json:"replies"`
}
