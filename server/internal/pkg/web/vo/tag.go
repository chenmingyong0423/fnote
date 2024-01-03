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

type Tag struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Route       string `json:"route"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreateTime  int64  `json:"create_time"`
	UpdateTime  int64  `json:"update_time"`
}

type SelectTag struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}
