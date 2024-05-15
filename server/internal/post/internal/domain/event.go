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

type CommentEvent struct {
	PostId    string   `json:"post_id"`
	CommentId string   `json:"comment_id"`
	RepliesId []string `json:"replies_id"`
	Count     int      `json:"count"`
	Type      string   `json:"type"`
}

type PostEvent struct {
	PostId            string   `json:"post_id"`
	AddedCategoryId   []string `json:"added_category_id,omitempty"`
	DeletedCategoryId []string `json:"deleted_category_id,omitempty"`
	AddedTagId        []string `json:"added_tag_id,omitempty"`
	DeletedTagId      []string `json:"deleted_tag_id,omitempty"`
	NewFileId         string   `json:"new_file_id,omitempty"`
	OldFileId         string   `json:"old_file_id,omitempty"`
	// 删除文章时，需要传入文章的评论数，用于更新网站的评论数
	CommentCount int    `json:"comment_count,omitempty"`
	Type         string `json:"type"`
}

type LikePostEvent struct {
	PostId string `json:"post_id"`
}
