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

package handler

import (
	"github.com/chenmingyong0423/fnote/backend/internal/config/service"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/gin-gonic/gin"
)

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
}

type WebsiteConfigVO struct {
	// 站点名称
	Name string `json:"name"`
	// 站点图标
	Icon string `json:"icon"`
	// 文章数量
	PostCount uint `json:"post_count"`
	// 分类数量
	CategoryCount uint `json:"category_count"`
	// 访问量
	ViewCount uint `json:"view_count"`
	// 网站运行时间
	LiveTime int64 `json:"live_time"`
	// 域名
	Domain string `json:"domain"`
	// 备案信息
	Records []string `json:"records"`
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

type SeoMetaConfigVO struct {
	Title                 string `json:"title"`
	Description           string `json:"description"`
	OgTitle               string `json:"ogTitle"`
	OgImage               string `json:"ogImage"`
	TwitterCard           string `json:"twitterCard"`
	BaiduSiteVerification string `json:"baidu-site-verification"`
	Keywords              string `json:"keywords"`
	Author                string `json:"author"`
	Robots                string `json:"robots"`
}

func NewConfigHandler(serv service.IConfigService) *ConfigHandler {
	return &ConfigHandler{
		serv: serv,
	}
}

type ConfigHandler struct {
	serv service.IConfigService
}

func (h *ConfigHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/configs")
	// 获取首页的配置信息
	routerGroup.GET("/index", api.Wrap(h.GetIndexConfig))

	adminGroup := engine.Group("/admin/configs")
	adminGroup.GET("/website", api.Wrap(h.AdminGetWebsiteConfig))
}

func (h *ConfigHandler) GetIndexConfig(ctx *gin.Context) (*IndexConfigVO, error) {
	config, err := h.serv.GetIndexConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &IndexConfigVO{
		WebsiteConfig:      *h.toWebsiteConfigVO(&config.WebSiteConfig),
		OwnerConfig:        *h.toOwnerConfigVO(&config.OwnerConfig),
		NoticeConfigVO:     *h.toNoticeConfigVO(&config.NoticeConfig),
		SocialInfoConfigVO: *h.toSocialInfoConfigVO(&config.SocialInfoConfig),
		PayInfoConfigVO:    h.toPayInfoConfigVO(config.PayInfoConfig),
		SeoMetaConfigVO:    *h.toSeoMetaConfigVO(&config.SeoMetaConfig),
	}, nil
}

func (h *ConfigHandler) toWebsiteConfigVO(webMasterCfg *domain.WebSiteConfig) *WebsiteConfigVO {
	return &WebsiteConfigVO{
		Name:          webMasterCfg.Name,
		Icon:          webMasterCfg.Icon,
		PostCount:     webMasterCfg.PostCount,
		CategoryCount: webMasterCfg.CategoryCount,
		ViewCount:     webMasterCfg.ViewCount,
		LiveTime:      webMasterCfg.LiveTime,
		Domain:        webMasterCfg.Domain,
		Records:       webMasterCfg.Records,
	}
}

func (h *ConfigHandler) toNoticeConfigVO(noticeCfg *domain.NoticeConfig) *NoticeConfigVO {
	return &NoticeConfigVO{Title: noticeCfg.Title, Content: noticeCfg.Content, PublishTime: noticeCfg.PublishTime}
}

func (h *ConfigHandler) toSocialInfoConfigVO(socialINfoConfig *domain.SocialInfoConfig) *SocialInfoConfigVO {
	socialInfoVOList := make([]SocialInfoVO, len(socialINfoConfig.SocialInfoList))
	for i, socialInfo := range socialINfoConfig.SocialInfoList {
		socialInfoVOList[i] = SocialInfoVO{SocialName: socialInfo.SocialName, SocialValue: socialInfo.SocialValue, CssClass: socialInfo.CssClass, IsLink: socialInfo.IsLink}
	}
	return &SocialInfoConfigVO{SocialInfoList: socialInfoVOList}
}

func (h *ConfigHandler) toPayInfoConfigVO(config []domain.PayInfoConfigElem) []PayInfoConfigVO {
	result := make([]PayInfoConfigVO, len(config))
	for i, payInfoConfig := range config {
		result[i] = PayInfoConfigVO{
			Name:  payInfoConfig.Name,
			Image: payInfoConfig.Image,
		}
	}
	return result
}

func (h *ConfigHandler) toSeoMetaConfigVO(config *domain.SeoMetaConfig) *SeoMetaConfigVO {
	return &SeoMetaConfigVO{
		Title:                 config.Title,
		Description:           config.Description,
		OgTitle:               config.OgTitle,
		OgImage:               config.OgImage,
		TwitterCard:           config.TwitterCard,
		BaiduSiteVerification: config.BaiduSiteVerification,
		Keywords:              config.Keywords,
		Author:                config.Author,
		Robots:                config.Robots,
	}
}

func (h *ConfigHandler) toOwnerConfigVO(ownerConfig *domain.OwnerConfig) *OwnerConfigVO {
	return &OwnerConfigVO{
		Name:    ownerConfig.Name,
		Profile: ownerConfig.Profile,
		Picture: ownerConfig.Picture,
	}
}

func (h *ConfigHandler) AdminGetWebsiteConfig(ctx *gin.Context) (WebsiteConfigVO, error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return WebsiteConfigVO{}, err
	}
	return *h.toWebsiteConfigVO(config), nil
}
