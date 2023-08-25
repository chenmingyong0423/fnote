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

package api

import (
	"net/http"
)

type PageVO[T any] struct {
	Page
	// 总页数
	TotalPages int64 `json:"totalPages"`
	// 总数量
	TotalCount int64 `json:"totalCount"`
	List       []T   `json:"list"`
}

func (p *PageVO[T]) SetTotalCountAndCalculateTotalPages(totalCount int64) {
	if p.PageSize == 0 {
		p.TotalPages = 0
	} else {
		p.TotalPages = (totalCount + p.PageSize - 1) / p.PageSize
	}
	p.TotalCount = totalCount
}

type PageRequest struct {
	Page
	// 排序字段
	Sort *Sorting `form:"sort,omitempty"`
	// 搜索内容
	Search *string `form:"search,omitempty"`
}

func (p *PageRequest) ValidateAndSetDefault() {
	if p.PageNo <= 0 {
		p.PageNo = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

type Sorting struct {
	Filed string
	Order string
}

type Page struct {
	// 当前页
	PageNo int64 `form:"pageNo" binding:"required"`
	// 每页数量
	PageSize int64 `form:"pageSize" binding:"required"`
}

type ListVO[T any] struct {
	List []T `json:"list"`
}

type ResponseBody[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func ErrResponse(message string) ResponseBody[any] {
	return ResponseBody[any]{
		Code:    http.StatusInternalServerError,
		Message: http.StatusText(http.StatusInternalServerError),
		Data:    nil,
	}
}

func SuccessResponse[T any](data T) ResponseBody[T] {
	return ResponseBody[T]{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
		Data:    data,
	}
}
