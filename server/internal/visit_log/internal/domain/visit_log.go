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

type VisitLog struct{}

type WebsiteVisitEvent struct {
	Url       string `json:"url"`
	Ip        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Origin    string `json:"origin"`
	Referer   string `json:"referer"`
}

type VisitHistory struct {
	Url       string
	Ip        string
	UserAgent string
	Origin    string
	Type      string
	Referer   string
}

type TendencyData struct {
	Timestamp int64
	ViewCount int64
}
