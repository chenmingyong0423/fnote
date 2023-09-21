// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,0
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package domain

type Friend struct {
	Id          string
	Name        string
	Url         string
	Logo        string
	Description string
	Email       string
	Status      FriendStatus
	Priority    int
}

type FriendVO struct {
	Name        string       `json:"name"`
	Url         string       `json:"url"`
	Logo        string       `json:"logo"`
	Description string       `json:"description"`
	Status      FriendStatus `json:"status"`
	Priority    int          `json:"priority"`
}

// 0 hiding，1 pending，2 showing

type FriendStatus uint

const (
	// FriendStatusPending 审核中
	FriendStatusPending FriendStatus = iota
	// FriendStatusShowing 展示
	FriendStatusShowing
	// FriendStatusHiding 隐藏
	FriendStatusHiding
)
