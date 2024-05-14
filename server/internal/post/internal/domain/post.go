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
	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
)

type Page struct {
	Size    int64
	Skip    int64
	Keyword string
	Field   string
	Order   string

	CategoryFilter []string
	TagFilter      []string
}

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
	Likes []string `json:"-"`
}

type ExtraPost struct {
	Content          string `json:"content"`
	MetaDescription  string `json:"meta_description"`
	MetaKeywords     string `json:"meta_keywords"`
	WordCount        int    `json:"word_count"`
	UpdateTime       int64  `json:"update_time"`
	IsDisplayed      bool   `json:"is_displayed"`
	IsCommentAllowed bool   `json:"is_comment_allowed"`
}

type PrimaryPost struct {
	Id           string          `json:"_id"`
	Author       string          `json:"author"`
	Title        string          `json:"title"`
	Summary      string          `json:"summary"`
	CoverImg     string          `json:"cover_img"`
	Categories   []Category4Post `json:"category"`
	Tags         []Tag4Post      `json:"tags"`
	LikeCount    int             `json:"like_count"`
	CommentCount int             `json:"comment_count"`
	VisitCount   int             `json:"visit_count"`
	StickyWeight int             `json:"sticky_weight"`
	CreateTime   int64           `json:"create_time"`
}

type Category4Post struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tag4Post struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PostEvent struct {
	PostId            string   `json:"post_id"`
	AddedCategoryId   []string `json:"added_category_id,omitempty"`
	DeletedCategoryId []string `json:"deleted_category_id,omitempty"`
	AddedTagId        []string `json:"added_tag_id,omitempty"`
	DeletedTagId      []string `json:"deleted_tag_id,omitempty"`
	NewFileId         string   `json:"new_file_id,omitempty"`
	OldFileId         string   `json:"old_file_id,omitempty"`
	Type              string   `json:"type"`
}

type LikePostEvent struct {
	PostId string `json:"post_id"`
}
