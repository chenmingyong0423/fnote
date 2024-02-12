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

package vo

// IndexConfigVO 首页信息
type IndexConfigVO struct {
	WebsiteConfig      WebsiteConfigVO    `json:"website_config"`
	OwnerConfig        OwnerConfigVO      `json:"owner_config"`
	NoticeConfigVO     NoticeConfigVO     `json:"notice_config"`
	SocialInfoConfigVO SocialInfoConfigVO `json:"social_info_config"`
	PayInfoConfigVO    []PayInfoConfigVO  `json:"pay_info_config"`
	SeoMetaConfigVO    SeoMetaConfigVO    `json:"seo_meta_config"`
}

type OwnerConfigVO struct {
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Picture string `json:"picture"`
}

type PayInfoConfigVO struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type NoticeConfigVO struct {
	Title       string `json:"title" `
	Content     string `json:"content"`
	PublishTime int64  `json:"publish_time"`
	Enabled     bool   `json:"enabled"`
}

type WebsiteConfigVO struct {
	// 站点名称
	WebsiteName string `json:"website_name"`
	// 站点图标
	Icon string `json:"icon"`
	// 网站运行时间
	LiveTime int64 `json:"live_time"`
	// 备案信息
	Records []string `json:"records"`
	// 站长名称
	OwnerName string `json:"owner_name"`
	// 站长简介
	OwnerProfile string `json:"owner_profile"`
	// 站长照片
	OwnerPicture string `json:"owner_picture"`
}

type SocialInfoConfigVO struct {
	SocialInfoList []SocialInfoVO `json:"social_info_list"`
}

type SocialInfoVO struct {
	SocialName  string `json:"social_name"`
	SocialValue string `json:"social_value"`
	CssClass    string `json:"css_class"`
	IsLink      bool   `json:"is_link"`
}

type AdminSocialInfoVO struct {
	Id string `json:"id"`
	SocialInfoVO
}

type SeoMetaConfigVO struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	OgTitle               string `json:"og_title"`
	OgImage               string `json:"og_image"`
	BaiduSiteVerification string `json:"baidu_site_verification"`
	Keywords              string `json:"keywords"`
	Author                string `json:"author"`
	Robots                string `json:"robots"`
}

type CommentConfigVO struct {
	EnableComment bool `json:"enable_comment"`
}

type FriendConfigVO struct {
	EnableFriendCommit bool `json:"enable_friend_commit"`
}

type EmailConfigVO struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type FrontPostCountConfigVO struct {
	Count int64 `json:"count"`
}

type LoginVO struct {
	AdminInfo  AdminInfoVO `json:"admin_info"`
	Expiration int64       `json:"expiration"`
	Token      string      `json:"token"`
}

type AdminInfoVO struct {
	Username string `json:"username"`
	Picture  string `json:"picture"`
}
