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
	Name            string `bson:"name" json:"name"`
	ArticleCount    uint   `bson:"articleCount" json:"articleCount"`
	ColumnCount     uint   `bson:"columnCount" json:"columnCount"`
	WebsiteViews    uint   `bson:"websiteViews" json:"websiteViews"`
	WebsiteLiveTime string `bson:"websiteLiveTime" json:"websiteLiveTime"`
	Profile         string `bson:"profile" json:"profile"`
	Picture         string `bson:"picture" json:"picture"`
	WebsiteIcon     string `bson:"websiteIcon" json:"websiteIcon"`
}
