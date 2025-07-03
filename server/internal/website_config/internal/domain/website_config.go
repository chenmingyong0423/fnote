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

import "time"

type Config struct {
	Id        string
	Props     any
	Typ       string
	CreatedAt int64
	UpdatedAt int64
}

type IndexConfig struct {
	WebSiteConfig    WebsiteConfig
	NoticeConfig     NoticeConfig
	SocialInfoConfig SocialInfoConfig
	PayInfoConfig    []PayInfoConfigElem
	SeoMetaConfig    SeoMetaConfig
	TPSVConfig       []TPSV
}

type CommonConfig struct {
	WebSiteConfig WebsiteConfig
	SeoMetaConfig SeoMetaConfig
	TPSVConfig    []TPSV
}

type PayInfoConfigElem struct {
	Name  string `bson:"name"`
	Image string `bson:"image"`
}

type PayInfoConfig struct {
	List []PayInfoConfigElem `bson:"list"`
}

type CarouselConfig struct {
	List []CarouselElem `bson:"list"`
}

type CarouselElem struct {
	Id        string    `bson:"id"`
	Title     string    `bson:"title"`
	Summary   string    `bson:"summary"`
	CoverImg  string    `bson:"cover_img"`
	Show      bool      `bson:"show"`
	Color     string    `bson:"color"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type WebsiteConfig struct {
	WebsiteName         string     `bson:"website_name"`
	WebsiteIcon         string     `bson:"website_icon"`
	WebsiteOwner        string     `bson:"website_owner"`
	WebsiteOwnerProfile string     `bson:"website_owner_profile"`
	WebsiteOwnerAvatar  string     `bson:"website_owner_avatar"`
	WebsiteRuntime      *time.Time `bson:"website_runtime,omitempty"`
	WebsiteRecords      []string   `bson:"website_records,omitempty"`
	WebsiteInit         *bool      `bson:"website_init,omitempty"`
}

// NoticeConfig 公告配置
type NoticeConfig struct {
	Title       string    `bson:"title"`
	Content     string    `bson:"content"`
	Enabled     bool      `bson:"enabled"`
	PublishTime time.Time `bson:"publish_time"`
}

type SwitchConfig struct {
	Enable bool `bson:"enable" json:"enable"`
}

type CommentConfig struct {
	EnableComment bool `bson:"enable_comment" json:"enable_comment"`
}

type FriendConfig struct {
	EnableFriendCommit bool   `bson:"enable_friend_commit,omitempty" json:"enable_friend_commit"`
	Introduction       string `bson:"introduction,omitempty" json:"introduction"`
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

type AdminConfig struct {
	Username string `bson:"username"`
	Password string `bson:"password"`
}

type TokenInfo struct {
	Expiration int64
	Token      string
}

type TPSVConfig struct {
	List []TPSV `bson:"list"`
}
type TPSV struct {
	Key         string `bson:"key"`
	Value       string `bson:"value"`
	Description string `bson:"description"`
}

type BaiduPushConfig struct {
	Baidu Baidu `bson:"baidu"`
}

type Baidu struct {
	Site  string `bson:"site"`
	Token string `bson:"token"`
}
