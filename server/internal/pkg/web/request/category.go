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

package request

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Route       string `json:"route" binding:"required"`
	Description string `json:"description"`
	ShowInNav   bool   `json:"show_in_nav"`
	Enabled     bool   `json:"enabled"`
}

type CategoryEnabledRequest struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

type CategoryNavRequest struct {
	ShowInNav *bool `json:"show_in_nav" binding:"required"`
}

type UpdateCategoryRequest struct {
	Description string `json:"description" binding:"required"`
}
