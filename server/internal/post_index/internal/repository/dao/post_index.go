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

package dao

import (
	"github.com/chenmingyong0423/go-mongox"
	"time"
)

type PostPushLogs struct {
	mongox.Model
	PostId string `bson:"post_id"`
	Url    string `bson:"url"`
	// 推送状态，false 表示失败，true 表示成功
	Status bool `bson:"status"`
	// 最后一次尝试推送的时间
	LastAttemptAt time.Time `bson:"last_attempt_date"`
	// 推送尝试次数
	AttemptCount int `bson:"try_count"`
	// 搜索引擎
	Engine       string        `bson:"engine"`
	ErrorDetails []ErrorDetail `bson:"error_details"`
}

type ErrorDetail struct {
	Response ErrorResponse `bson:"response"`
	// 失败的时间
	FailedAt time.Time `bson:"failed_at"`
}

type ErrorResponse struct {
	// http 状态码
	StatusCode int `bson:"status_code"`
	// 错误码
	Code int `bson:"code"`
	// 错误信息
	Message string `bson:"message"`
}
