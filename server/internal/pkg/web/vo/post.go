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

// {
//    key: '1',
//    coverImg: 'John Brown',
//    title: 32,
//    summary: 'New York No. 1 Lake Park',
//    categories: ['nice', 'developer'],
//    tags: ['nice', 'developer'],
//    createTime: 121123123312132,
//    updateTime: 121123123312132,
//  },

type AdminPostVO struct {
	Id         string   `json:"_id"`
	CoverImg   string   `json:"coverImg"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Categories []string `json:"categories"`
	Tags       []string `json:"tags"`
	CreateTime int64    `json:"createTime"`
	UpdateTime int64    `json:"updateTime"`
}
