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

package request

type UpdateWebsiteConfigReq struct {
	WebsiteName  string `json:"website_name" binding:"required"`
	LiveTime     int64  `json:"live_time" binding:"required"`
	OwnerName    string `json:"owner_name" binding:"required"`
	OwnerProfile string `json:"owner_profile" binding:"required"`
	OwnerPicture string `json:"owner_picture"`
	Icon         string `json:"icon"`
}

type UpdateOwnerConfigReq struct {
	Name    string `json:"name" binding:"required"`
	Profile string `json:"profile" binding:"required"`
	Picture string `json:"picture" binding:"required"`
}

type UpdateSeoMetaConfigReq struct {
	Title                 string `json:"title" binding:"required"`
	Description           string `json:"description" binding:"required"`
	OgTitle               string `json:"og_title" binding:"required"`
	Keywords              string `json:"keywords" binding:"required"`
	Author                string `json:"author" binding:"required"`
	Robots                string `json:"robots" binding:"required"`
	OgImage               string `json:"og_image"`
	BaiduSiteVerification string `json:"baidu_site_verification"`
}

type UpdateCommentConfigReq struct {
	EnableComment *bool `json:"enable_comment" binding:"required"`
}

type UpdateFriendConfigReq struct {
	EnableFriendCommit *bool `json:"enable_friend_commit" binding:"required"`
}

type UpdateEmailConfigReq struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UpdateNoticeConfigReq struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateNoticeConfigEnabledReq struct {
	Enabled *bool `json:"enabled" binding:"required"`
}

type UpdateFPCConfigCountReq struct {
	Count int64 `json:"count" binding:"required"`
}

type AddRecordInWebsiteConfig struct {
	Record string `json:"record" binding:"required"`
}

type AddPayInfoRequest struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
}

type SocialInfoReq struct {
	SocialName  string `json:"social_name"`
	SocialValue string `json:"social_value"`
	CssClass    string `json:"css_class"`
	IsLink      bool   `json:"is_link"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type InitRequest struct {
	WebsiteName         string      `json:"website_name" binding:"required"`
	WebsiteIcon         string      `json:"website_icon" binding:"required"`
	WebsiteOwner        string      `json:"website_owner" binding:"required"`
	WebsiteOwnerProfile string      `json:"website_owner_profile" binding:"required"`
	WebsiteOwnerAvatar  string      `json:"website_owner_avatar" binding:"required"`
	WebsiteOwnerEmail   string      `json:"website_owner_email" binding:"required"`
	EmailServer         EmailServer `json:"email_server" binding:"required"`
	Admin               Admin       `json:"admin" binding:"required"`
}

type EmailServer struct {
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type Admin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
