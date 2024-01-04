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
	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/service"
	"github.com/chenmingyong0423/gkit"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func NewWebsiteConfigHandler(serv service.IWebsiteConfigService) *WebsiteConfigHandler {
	return &WebsiteConfigHandler{
		serv: serv,
	}
}

type WebsiteConfigHandler struct {
	serv service.IWebsiteConfigService
}

func (h *WebsiteConfigHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/configs")
	// 获取首页的配置信息
	routerGroup.GET("/index", api.Wrap(h.GetIndexConfig))

	adminGroup := engine.Group("/admin/configs")
	adminGroup.GET("/website", api.Wrap(h.AdminGetWebsiteConfig))
	adminGroup.PUT("/website", api.WrapWithBody(h.AdminUpdateWebsiteConfig))
	adminGroup.POST("/website/records", api.WrapWithBody(h.AdminAddRecordInWebsiteConfig))
	adminGroup.DELETE("/website/records", api.Wrap(h.AdminDeleteRecordInWebsiteConfig))
	adminGroup.GET("/owner", api.Wrap(h.AdminGetOwnerConfig))
	adminGroup.PUT("/owner", api.WrapWithBody(h.AdminUpdateOwnerConfig))
	adminGroup.GET("/seo", api.Wrap(h.AdminGetSeoConfig))
	adminGroup.PUT("/seo", api.WrapWithBody(h.AdminUpdateSeoConfig))
	adminGroup.GET("/comment", api.Wrap(h.AdminGetCommentConfig))
	adminGroup.PUT("/comment", api.WrapWithBody(h.AdminUpdateCommentConfig))
	adminGroup.GET("/friend", api.Wrap(h.AdminGetFriendConfig))
	adminGroup.PUT("/friend", api.WrapWithBody(h.AdminUpdateFriendConfig))
	adminGroup.GET("/email", api.Wrap(h.AdminGetEmailConfig))
	adminGroup.PUT("/email", api.WrapWithBody(h.AdminUpdateEmailConfig))

	adminGroup.GET("/notice", api.Wrap(h.AdminGetNoticeConfig))
	adminGroup.PUT("/notice", api.WrapWithBody(h.AdminUpdateNoticeConfig))
	adminGroup.PUT("/notice/enabled", api.WrapWithBody(h.AdminUpdateNoticeEnabled))

	adminGroup.GET("/front-post-count", api.Wrap(h.AdminGetFPCConfig))
	adminGroup.PUT("/front-post-count", api.WrapWithBody(h.AdminUpdateFPCConfig))
	adminGroup.GET("/pay", api.Wrap(h.AdminGetPayConfig))
	adminGroup.POST("/pay", api.WrapWithBody(h.AdminAddPayInfo))
	adminGroup.DELETE("/pay/:name", api.Wrap(h.AdminDeletePayInfo))
}

func (h *WebsiteConfigHandler) GetIndexConfig(ctx *gin.Context) (*vo.IndexConfigVO, error) {
	config, err := h.serv.GetIndexConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &vo.IndexConfigVO{
		WebsiteConfig:      *h.toWebsiteConfigVO(&config.WebSiteConfig),
		OwnerConfig:        *h.toOwnerConfigVO(&config.OwnerConfig),
		NoticeConfigVO:     *h.toNoticeConfigVO(&config.NoticeConfig),
		SocialInfoConfigVO: *h.toSocialInfoConfigVO(&config.SocialInfoConfig),
		PayInfoConfigVO:    h.toPayInfoConfigVO(config.PayInfoConfig),
		SeoMetaConfigVO:    *h.toSeoMetaConfigVO(&config.SeoMetaConfig),
	}, nil
}

func (h *WebsiteConfigHandler) toWebsiteConfigVO(webMasterCfg *domain.WebSiteConfig) *vo.WebsiteConfigVO {
	return &vo.WebsiteConfigVO{
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

func (h *WebsiteConfigHandler) toNoticeConfigVO(noticeCfg *domain.NoticeConfig) *vo.NoticeConfigVO {
	return &vo.NoticeConfigVO{Title: noticeCfg.Title, Content: noticeCfg.Content, Enabled: noticeCfg.Enabled, PublishTime: noticeCfg.PublishTime}
}

func (h *WebsiteConfigHandler) toSocialInfoConfigVO(socialINfoConfig *domain.SocialInfoConfig) *vo.SocialInfoConfigVO {
	socialInfoVOList := make([]vo.SocialInfoVO, len(socialINfoConfig.SocialInfoList))
	for i, socialInfo := range socialINfoConfig.SocialInfoList {
		socialInfoVOList[i] = vo.SocialInfoVO{SocialName: socialInfo.SocialName, SocialValue: socialInfo.SocialValue, CssClass: socialInfo.CssClass, IsLink: socialInfo.IsLink}
	}
	return &vo.SocialInfoConfigVO{SocialInfoList: socialInfoVOList}
}

func (h *WebsiteConfigHandler) toPayInfoConfigVO(config []domain.PayInfoConfigElem) []vo.PayInfoConfigVO {
	result := make([]vo.PayInfoConfigVO, len(config))
	for i, payInfoConfig := range config {
		result[i] = vo.PayInfoConfigVO{
			Name:  payInfoConfig.Name,
			Image: payInfoConfig.Image,
		}
	}
	return result
}

func (h *WebsiteConfigHandler) toSeoMetaConfigVO(config *domain.SeoMetaConfig) *vo.SeoMetaConfigVO {
	return &vo.SeoMetaConfigVO{
		Title:                 config.Title,
		Description:           config.Description,
		OgTitle:               config.OgTitle,
		OgImage:               config.OgImage,
		BaiduSiteVerification: config.BaiduSiteVerification,
		Keywords:              config.Keywords,
		Author:                config.Author,
		Robots:                config.Robots,
	}
}

func (h *WebsiteConfigHandler) toOwnerConfigVO(ownerConfig *domain.OwnerConfig) *vo.OwnerConfigVO {
	return &vo.OwnerConfigVO{
		Name:    ownerConfig.Name,
		Profile: ownerConfig.Profile,
		Picture: ownerConfig.Picture,
	}
}

func (h *WebsiteConfigHandler) AdminGetWebsiteConfig(ctx *gin.Context) (vo.WebsiteConfigVO, error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return vo.WebsiteConfigVO{}, err
	}
	return *h.toWebsiteConfigVO(config), nil
}

func (h *WebsiteConfigHandler) AdminUpdateWebsiteConfig(ctx *gin.Context, req request.UpdateWebsiteConfigReq) (any, error) {
	return nil, h.serv.UpdateWebSiteConfig(ctx, domain.WebSiteConfig{
		Name:     req.Name,
		Icon:     req.Icon,
		LiveTime: req.LiveTime,
	})
}

func (h *WebsiteConfigHandler) AdminGetOwnerConfig(ctx *gin.Context) (vo.OwnerConfigVO, error) {
	config, err := h.serv.GetOwnerConfig(ctx)
	if err != nil {
		return vo.OwnerConfigVO{}, err
	}
	return *h.toOwnerConfigVO(&config), nil
}

func (h *WebsiteConfigHandler) AdminUpdateOwnerConfig(ctx *gin.Context, req request.UpdateOwnerConfigReq) (any, error) {
	return nil, h.serv.UpdateOwnerConfig(ctx, domain.OwnerConfig{
		Name:    req.Name,
		Profile: req.Profile,
		Picture: req.Picture,
	})
}

func (h *WebsiteConfigHandler) AdminGetSeoConfig(ctx *gin.Context) (vo.SeoMetaConfigVO, error) {
	config, err := h.serv.GetSeoMetaConfig(ctx)
	if err != nil {
		return vo.SeoMetaConfigVO{}, err
	}
	return *h.toSeoMetaConfigVO(config), nil
}

func (h *WebsiteConfigHandler) AdminUpdateSeoConfig(ctx *gin.Context, req request.UpdateSeoMetaConfigReq) (any, error) {
	return nil, h.serv.UpdateSeoMetaConfig(ctx, &domain.SeoMetaConfig{
		Title:                 req.Title,
		Description:           req.Description,
		OgTitle:               req.OgTitle,
		OgImage:               req.OgImage,
		BaiduSiteVerification: req.BaiduSiteVerification,
		Keywords:              req.Keywords,
		Author:                req.Author,
		Robots:                req.Robots,
	})
}

func (h *WebsiteConfigHandler) AdminGetCommentConfig(ctx *gin.Context) (vo.CommentConfigVO, error) {
	config, err := h.serv.GetCommentConfig(ctx)
	if err != nil {
		return vo.CommentConfigVO{}, err
	}
	return h.toCommentConfigVO(config), nil
}

func (h *WebsiteConfigHandler) toCommentConfigVO(config domain.CommentConfig) vo.CommentConfigVO {
	return vo.CommentConfigVO{
		EnableComment: config.EnableComment,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateCommentConfig(ctx *gin.Context, req request.UpdateCommentConfigReq) (any, error) {
	return nil, h.serv.UpdateCommentConfig(ctx, domain.CommentConfig{
		EnableComment: gkit.GetValueOrDefault(req.EnableComment),
	})
}

func (h *WebsiteConfigHandler) AdminGetFriendConfig(ctx *gin.Context) (vo.FriendConfigVO, error) {
	config, err := h.serv.GetFriendConfig(ctx)
	if err != nil {
		return vo.FriendConfigVO{}, err
	}
	return h.toFriendConfigVO(config), nil
}

func (h *WebsiteConfigHandler) toFriendConfigVO(config domain.FriendConfig) vo.FriendConfigVO {
	return vo.FriendConfigVO{
		EnableFriendCommit: config.EnableFriendCommit,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFriendConfig(ctx *gin.Context, req request.UpdateFriendConfigReq) (any, error) {
	return nil, h.serv.UpdateFriendConfig(ctx, domain.FriendConfig{
		EnableFriendCommit: gkit.GetValueOrDefault(req.EnableFriendCommit),
	})
}

func (h *WebsiteConfigHandler) AdminGetEmailConfig(ctx *gin.Context) (vo.EmailConfigVO, error) {
	config, err := h.serv.GetEmailConfig(ctx)
	if err != nil {
		return vo.EmailConfigVO{}, err
	}
	return h.toEmailConfigVO(config), nil
}

func (h *WebsiteConfigHandler) toEmailConfigVO(config *domain.EmailConfig) vo.EmailConfigVO {
	return vo.EmailConfigVO{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Email:    config.Email,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateEmailConfig(ctx *gin.Context, req request.UpdateEmailConfigReq) (any, error) {
	return nil, h.serv.UpdateEmailConfig(ctx, &domain.EmailConfig{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
}

func (h *WebsiteConfigHandler) AdminGetNoticeConfig(ctx *gin.Context) (vo.NoticeConfigVO, error) {
	config, err := h.serv.GetNoticeConfig(ctx)
	if err != nil {
		return vo.NoticeConfigVO{}, err
	}
	return *h.toNoticeConfigVO(&config), nil
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeConfig(ctx *gin.Context, req request.UpdateNoticeConfigReq) (any, error) {
	return nil, h.serv.UpdateNoticeConfig(ctx, &domain.NoticeConfig{
		Title:   req.Title,
		Content: req.Content,
	})
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeEnabled(ctx *gin.Context, req request.UpdateNoticeConfigEnabledReq) (any, error) {
	return nil, h.serv.UpdateNoticeConfigEnabled(ctx, gkit.GetValueOrDefault(req.Enabled))
}

func (h *WebsiteConfigHandler) AdminGetFPCConfig(ctx *gin.Context) (vo.FrontPostCountConfigVO, error) {
	config, err := h.serv.GetFrontPostCountConfig(ctx)
	if err != nil {
		return vo.FrontPostCountConfigVO{}, err
	}
	return h.toFrontPostCountConfigVO(config), nil
}

func (h *WebsiteConfigHandler) toFrontPostCountConfigVO(config domain.FrontPostCountConfig) vo.FrontPostCountConfigVO {
	return vo.FrontPostCountConfigVO{
		Count: config.Count,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFPCConfig(ctx *gin.Context, req request.UpdateFPCConfigCountReq) (any, error) {
	return nil, h.serv.UpdateFrontPostCountConfig(ctx, domain.FrontPostCountConfig{
		Count: req.Count,
	})
}

func (h *WebsiteConfigHandler) AdminAddRecordInWebsiteConfig(ctx *gin.Context, req request.AddRecordInWebsiteConfig) (any, error) {
	return nil, h.serv.AddRecordInWebsiteConfig(ctx, req.Record)
}

func (h *WebsiteConfigHandler) AdminDeleteRecordInWebsiteConfig(ctx *gin.Context) (any, error) {
	record := ctx.Query("record")
	if record == "" {
		return nil, errors.New("record is empty")
	}
	return nil, h.serv.DeleteRecordInWebsiteConfig(ctx, record)
}

func (h *WebsiteConfigHandler) AdminGetPayConfig(ctx *gin.Context) (vo api.ListVO[vo.PayInfoConfigVO], err error) {
	config, err := h.serv.GetPayConfig(ctx)
	if err != nil {
		return
	}
	vo.List = h.toPayInfoConfigVO(config.List)
	return
}

func (h *WebsiteConfigHandler) AdminAddPayInfo(ctx *gin.Context, req request.AddPayInfoRequest) (any, error) {
	return nil, h.serv.AddPayInfo(ctx, domain.PayInfoConfigElem{
		Name:  req.Name,
		Image: req.Image,
	})
}

func (h *WebsiteConfigHandler) AdminDeletePayInfo(ctx *gin.Context) (any, error) {
	name := ctx.Param("name")
	if name == "" {
		return nil, errors.New("name is empty")
	}
	image := ctx.Query("image")
	if image == "" {
		return nil, errors.New("image is empty")
	}
	return nil, h.serv.DeletePayInfo(ctx, domain.PayInfoConfigElem{
		Name:  name,
		Image: image,
	})
}
