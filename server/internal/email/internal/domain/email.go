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

type Email struct {
	// STMP 地址
	Host string
	Port int
	// 邮箱账号
	Username string
	// 密码
	Password string
	// 发件人
	Name string
	// 收件人
	To []string
	// 标题
	Subject string
	// 内容
	Body string
	// 内容类型
	ContentType string
}
