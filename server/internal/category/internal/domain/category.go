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

type Category struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	ShowInNav   bool   `json:"show_in_nav"`
	Sort        int64  `json:"sort"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	PostCount   int64  `json:"post_count"`
}

type CategoryWithCount struct {
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
	Count       int64  `json:"count"`
}
