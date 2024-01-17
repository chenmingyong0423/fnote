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

package vo

type AdminPostVO struct {
	Id         string          `json:"id"`
	CoverImg   string          `json:"cover_img"`
	Title      string          `json:"title"`
	Summary    string          `json:"summary"`
	Categories []Category4Post `json:"categories"`
	Tags       []Tag4Post      `json:"tags"`
	CreateTime int64           `json:"create_time"`
	UpdateTime int64           `json:"update_time"`
}

type Category4Post struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tag4Post struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type PostDetailVO struct {
	Id               string          `json:"id"`
	Author           string          `json:"author"`
	Title            string          `json:"title"`
	Summary          string          `json:"summary"`
	Content          string          `json:"content"`
	CoverImg         string          `json:"cover_img"`
	Categories       []Category4Post `json:"categories"`
	Tags             []Tag4Post      `json:"tags"`
	IsDisplayed      bool            `json:"is_displayed"`
	StickyWeight     int             `json:"sticky_weight"`
	MetaDescription  string          `json:"meta_description"`
	MetaKeywords     string          `json:"meta_keywords"`
	IsCommentAllowed bool            `json:"is_comment_allowed"`
}
