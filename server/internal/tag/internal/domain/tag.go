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

type Tag struct {
	Id        string
	Name      string
	Route     string
	PostCount int64
	Enabled   bool
	CreatedAt int64
	UpdatedAt int64
}

type TagWithCount struct {
	Name  string
	Route string
	Count int64
}

type TagEvent struct {
	TagId string `json:"tag_id"`
	Type  string `json:"type"`
}

type PageDTO struct {
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

func (p *PageDTO) OrderConvertToInt() int {
	switch p.Order {
	case "ASC":
		return 1
	case "DESC":
		return -1
	default:
		return -1
	}
}
