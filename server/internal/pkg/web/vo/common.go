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

type PageVO[T any] struct {
	// 当前页
	PageNo int64 `json:"pageNo"`
	// 每页数量
	PageSize int64 `json:"pageSize"`
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
