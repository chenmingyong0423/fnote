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
	OwnerConfig      OwnerConfig
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

// OwnerConfig 站长信息
type OwnerConfig struct {
	Name    string `bson:"name" json:"name"`
	Profile string `bson:"profile" json:"profile"`
	Picture string `bson:"picture" json:"picture"`
}

// WebSiteConfig 站点信息
type WebSiteConfig struct {
	// 站点名称
	Name string `bson:"name"`
	// 站点图标
	Icon string `bson:"icon"`
	// 文章数量
	PostCount uint `bson:"post_count"`
	// 分类数量
	CategoryCount uint `bson:"category_count"`
	// 访问量
	ViewCount uint `bson:"view_count"`
	// 网站运行时间
	LiveTime int64 `bson:"live_time"`
	// 域名
	Domain string `bson:"domain"`
	// 备案信息
	Records []string `bson:"records"`
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
	Account  string `bson:"account"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
}

type SocialInfoConfig struct {
	SocialInfoList []SocialInfo `bson:"social_info_list" json:"social_info_list"`
}

type SocialInfo struct {
	SocialName  string `bson:"social_name" json:"social_name"`
	SocialValue string `bson:"social_value" json:"social_value"`
	CssClass    string `bson:"css_class" json:"css_class"`
	IsLink      bool   `bson:"is_link" json:"is_link"`
}

type FrontPostCount struct {
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
