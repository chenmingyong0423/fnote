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
	WebMasterConfig    WebMasterConfigVO  `json:"web_master_config"`
	NoticeConfigVO     NoticeConfigVO     `json:"notice_config"`
	SocialInfoConfigVO SocialInfoConfigVO `json:"social_info_config"`
	PayInfoConfigVO    []PayInfoConfigVO  `json:"pay_info_config"`
	SeoMetaConfigVO    SeoMetaConfigVO    `json:"seo_meta_config"`
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

type WebMasterConfigVO struct {
	Name            string   `json:"name"`
	PostCount       uint     `json:"post_count"`
	CategoryCount   uint     `json:"category_count"`
	WebsiteViews    uint     `json:"website_views"`
	WebsiteLiveTime int64    `json:"website_live_time"`
	Profile         string   `json:"profile"`
	Picture         string   `json:"picture"`
	WebsiteIcon     string   `json:"website_icon"`
	Domain          string   `json:"domain"`
	Records         []string `json:"records"`
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
}

func (h *ConfigHandler) GetWebmasterInfo(ctx *gin.Context) (*WebMasterConfigVO, error) {
	webMasterConfig, err := h.serv.GetWebmasterInfo(ctx, "webmaster")
	if err != nil {
		return nil, err
	}
	return h.toWebMasterConfigVO(webMasterConfig), nil
}

func (h *ConfigHandler) GetIndexConfig(ctx *gin.Context) (*IndexConfigVO, error) {
	config, err := h.serv.GetIndexConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &IndexConfigVO{
		WebMasterConfig:    *h.toWebMasterConfigVO(&config.WebMasterConfig),
		NoticeConfigVO:     *h.toNoticeConfigVO(&config.NoticeConfig),
		SocialInfoConfigVO: *h.toSocialInfoConfigVO(&config.SocialInfoConfig),
		PayInfoConfigVO:    h.toPayInfoConfigVO(config.PayInfoConfig),
		SeoMetaConfigVO:    *h.toSeoMetaConfigVO(&config.SeoMetaConfig),
	}, nil
}

func (h *ConfigHandler) toWebMasterConfigVO(webMasterCfg *domain.WebMasterConfig) *WebMasterConfigVO {
	return &WebMasterConfigVO{Name: webMasterCfg.Name, PostCount: webMasterCfg.PostCount, CategoryCount: webMasterCfg.CategoryCount, WebsiteViews: webMasterCfg.WebsiteViews, WebsiteLiveTime: webMasterCfg.WebsiteLiveTime, Profile: webMasterCfg.Profile, Picture: webMasterCfg.Picture, WebsiteIcon: webMasterCfg.WebsiteIcon, Domain: webMasterCfg.Domain, Records: webMasterCfg.Records}
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
