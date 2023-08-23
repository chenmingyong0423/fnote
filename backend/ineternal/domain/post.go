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

type PostVO struct {
	BasePost
}

type Post struct {
	BasePost
}

type BasePost struct {
	Sug        string   `json:"sug"`
	Author     string   `json:"author"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	CoverImg   string   `json:"cover_img"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	LikeCount  int      `json:"likeCount"`
	Comments   int      `json:"comments"`
	Visits     int      `json:"visit"`
	Priority   int      `json:"priority"`
	CreateTime int64    `json:"createTime"`
}
