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

package vo

type DataAnalysis struct {
	// 文章总数
	PostCount uint `json:"post_count"`
	// 分类总数
	CategoryCount uint `json:"category_count"`
	// 评论总数
	CommentCount int64 `json:"comment_count"`
	// 点赞数
	LikeCount int64 `json:"like_count"`
	// 标签总数
	TagCount uint `json:"tag_count"`
	// 今日访问量
	TodayViewCount int64 `json:"today_view_count"`
	// 总访问量
	TotalViewCount int64 `json:"total_view_count"`
	// 今日用户访问量
	TodayUserVisitCount int64 `json:"today_user_visit_count"`
}
