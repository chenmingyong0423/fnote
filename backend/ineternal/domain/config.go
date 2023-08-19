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

type WebMasterConfig struct {
	Name            string
	ArticleCount    uint
	ColumnCount     uint
	WebsiteViews    uint
	WebsiteLiveTime string
	Profile         string
	Picture         string
	WebsiteIcon     string
}

type WebMasterConfigVO struct {
	Name            string `json:"name"`
	ArticleCount    uint   `json:"article_count"`
	ColumnCount     uint   `json:"column_count"`
	WebsiteViews    uint   `json:"website_views"`
	WebsiteLiveTime string `json:"website_live_time"`
	Profile         string `json:"profile"`
	Picture         string `json:"picture"`
	WebsiteIcon     string `json:"website_icon"`
}
