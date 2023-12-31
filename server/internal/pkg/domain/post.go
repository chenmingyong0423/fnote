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

import "github.com/chenmingyong0423/fnote/backend/internal/pkg/api"

type PostsQueryCondition struct {
	Size int64
	Skip int64

	Keyword *string

	api.Sorting

	Categories []string
	Tags       []string
}

type PostRequest struct {
	api.PageRequest
	Categories []string `form:"categories"`
	Tags       []string `form:"tags"`
}

type DetailPostVO struct {
	PrimaryPost
	ExtraPost
	IsLiked bool `json:"is_liked"`
}

type Post struct {
	PrimaryPost
	ExtraPost
	IsCommentAllowed bool     `json:"is_comment_allowed"`
	Likes            []string `json:"-"`
}

type ExtraPost struct {
	Content         string `json:"content"`
	MetaDescription string `json:"meta_description"`
	MetaKeywords    string `json:"meta_keywords"`
	WordCount       int    `json:"word_count"`
	UpdateTime      int64  `json:"update_time"`
}

type PrimaryPost struct {
	Sug          string   `json:"sug"`
	Author       string   `json:"author"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	CoverImg     string   `json:"cover_img"`
	Categories   []string `json:"category"`
	Tags         []string `json:"tags"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	VisitCount   int      `json:"visit_count"`
	StickyWeight int      `json:"sticky_weight"`
	CreateTime   int64    `json:"create_time"`
}

type PostStatus uint

const (
	// PostStatusDraft 草稿
	PostStatusDraft PostStatus = iota
	// PostStatusPunished 已发布
	PostStatusPunished
	// PostStatusDeleted 已删除
	PostStatusDeleted
)
