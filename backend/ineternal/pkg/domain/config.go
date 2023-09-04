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

// WebMasterConfig 站长信息
type WebMasterConfig struct {
	Name            string `bson:"name" json:"name"`
	PostCount       uint   `bson:"postCount" json:"postCount"`
	ColumnCount     uint   `bson:"columnCount" json:"columnCount"`
	WebsiteViews    uint   `bson:"websiteViews" json:"websiteViews"`
	WebsiteLiveTime string `bson:"websiteLiveTime" json:"websiteLiveTime"`
	Profile         string `bson:"profile" json:"profile"`
	Picture         string `bson:"picture" json:"picture"`
	WebsiteIcon     string `bson:"websiteIcon" json:"websiteIcon"`
	Domain          string `bson:"domain" json:"domain"`
}

type WebMasterConfigVO struct {
	Name            string `json:"name"`
	PostCount       uint   `json:"post_count"`
	ColumnCount     uint   `json:"column_count"`
	WebsiteViews    uint   `json:"website_views"`
	WebsiteLiveTime string `json:"website_live_time"`
	Profile         string `json:"profile"`
	Picture         string `json:"picture"`
	WebsiteIcon     string `json:"website_icon"`
	Domain          string `json:"domain"`
}
