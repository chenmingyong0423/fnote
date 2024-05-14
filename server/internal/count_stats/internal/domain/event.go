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

type TagEvent struct {
	TagId string `json:"tag_id"`
	Type  string `json:"type"`
}

type CommentEvent struct {
	PostId    string   `json:"post_id"`
	CommentId string   `json:"comment_id"`
	RepliesId []string `json:"replies_id"`
	Count     int      `json:"count"`
	Type      string   `json:"type"`
}

type CategoryEvent struct {
	CategoryId string `json:"category_id"`
}
