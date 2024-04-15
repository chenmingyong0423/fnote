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

package web

type BaiduPushVO struct {
	// 当天剩余的可推送url条数
	Remain int `json:"remain,omitempty"`
	// 成功推送的url条数
	Success int `json:"success,omitempty"`
	// 由于不是本站url而未处理的url列表
	NotSameSite []string `json:"not_same_site,omitempty"`
	// 不合法的url列表
	NotValid []string `json:"not_valid,omitempty"`
	// 错误码，与状态码相同
	Err int `json:"error,omitempty"`
	// 错误描述
	Message string `json:"message,omitempty"`
}
