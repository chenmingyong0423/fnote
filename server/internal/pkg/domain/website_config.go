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

type Config struct {
	Id         string
	Props      any
	Typ        string
	CreateTime int64
	UpdateTime int64
}

type IndexConfig struct {
	WebSiteConfig    WebSiteConfig
	NoticeConfig     NoticeConfig
	SocialInfoConfig SocialInfoConfig
	PayInfoConfig    []PayInfoConfigElem
	SeoMetaConfig    SeoMetaConfig
}

type PayInfoConfigElem struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
}

type PayInfoConfig struct {
	List []PayInfoConfigElem `bson:"list"`
}

// WebSiteConfig 站点信息
type WebSiteConfig struct {
	// 站点名称
	WebsiteName string `bson:"website_name"`
	// 站点图标
	Icon string `bson:"icon"`
	// 网站运行时间
	LiveTime int64 `bson:"live_time"`
	// 备案信息
	Records []string `bson:"records"`
	// 站长名称
	OwnerName string `bson:"owner_name"`
	// 站长简介
	OwnerProfile string `bson:"owner_profile"`
	// 站长照片
	OwnerPicture string `bson:"owner_picture"`
}

// NoticeConfig 公告配置
type NoticeConfig struct {
	Title       string `bson:"title"`
	Content     string `bson:"content"`
	Enabled     bool   `bson:"enabled"`
	PublishTime int64  `bson:"publish_time"`
}

type SwitchConfig struct {
	Enable bool `bson:"enable" json:"enable"`
}

type CommentConfig struct {
	EnableComment bool `bson:"enable_comment" json:"enable_comment"`
}

type FriendConfig struct {
	EnableFriendCommit bool `bson:"enable_friend_commit" json:"enable_friend_commit"`
}

type EmailConfig struct {
	Host     string `bson:"host"`
	Port     int    `bson:"port"`
	Username string `bson:"username"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type SocialInfoConfig struct {
	SocialInfoList []SocialInfo `bson:"social_info_list" json:"social_info_list"`
}

type SocialInfo struct {
	Id          []byte `bson:"id"`
	SocialName  string `bson:"social_name" json:"social_name"`
	SocialValue string `bson:"social_value" json:"social_value"`
	CssClass    string `bson:"css_class" json:"css_class"`
	IsLink      bool   `bson:"is_link" json:"is_link"`
}

type FrontPostCountConfig struct {
	Count int64 `bson:"count"`
}

type SeoMetaConfig struct {
	Title                 string `bson:"title"`
	Description           string `bson:"description"`
	OgTitle               string `bson:"og_title"`
	OgImage               string `bson:"og_image"`
	BaiduSiteVerification string `bson:"baidu_site_verification"`
	Keywords              string `bson:"keywords"`
	Author                string `bson:"author"`
	Robots                string `bson:"robots"`
}
