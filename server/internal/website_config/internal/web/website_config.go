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

package web

import (
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"

	"github.com/chenmingyong0423/fnote/server/internal/global"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/jwtutil"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

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
	initMutex sync.Mutex
	serv      service.IWebsiteConfigService
}

func (h *WebsiteConfigHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/configs")
	// 获取首页的配置信息
	routerGroup.GET("/index", apiwrap.Wrap(h.GetIndexConfig))
	routerGroup.GET("/check-initialization", apiwrap.Wrap(h.GetInitStatus))

	adminGroup := engine.Group("/admin/configs")
	adminGroup.GET("/check-initialization", apiwrap.Wrap(h.GetInitStatus))
	adminGroup.POST("/initialization", apiwrap.WrapWithBody(h.InitializeWebsite))
	adminGroup.GET("/website", apiwrap.Wrap(h.AdminGetWebsiteConfig))
	adminGroup.PUT("/website", apiwrap.WrapWithBody(h.AdminUpdateWebsiteConfig))
	adminGroup.POST("/website/records", apiwrap.WrapWithBody(h.AdminAddRecordInWebsiteConfig))
	adminGroup.DELETE("/website/records", apiwrap.Wrap(h.AdminDeleteRecordInWebsiteConfig))
	adminGroup.GET("/seo", apiwrap.Wrap(h.AdminGetSeoConfig))
	adminGroup.PUT("/seo", apiwrap.WrapWithBody(h.AdminUpdateSeoConfig))
	adminGroup.GET("/comment", apiwrap.Wrap(h.AdminGetCommentConfig))
	adminGroup.PUT("/comment", apiwrap.WrapWithBody(h.AdminUpdateCommentConfig))
	adminGroup.GET("/friend", apiwrap.Wrap(h.AdminGetFriendConfig))
	adminGroup.PUT("/friend", apiwrap.WrapWithBody(h.AdminUpdateFriendConfig))
	adminGroup.GET("/email", apiwrap.Wrap(h.AdminGetEmailConfig))
	adminGroup.PUT("/email", apiwrap.WrapWithBody(h.AdminUpdateEmailConfig))

	adminGroup.GET("/notice", apiwrap.Wrap(h.AdminGetNoticeConfig))
	adminGroup.PUT("/notice", apiwrap.WrapWithBody(h.AdminUpdateNoticeConfig))
	adminGroup.PUT("/notice/enabled", apiwrap.WrapWithBody(h.AdminUpdateNoticeEnabled))

	adminGroup.GET("/front-post-count", apiwrap.Wrap(h.AdminGetFPCConfig))
	adminGroup.PUT("/front-post-count", apiwrap.WrapWithBody(h.AdminUpdateFPCConfig))
	adminGroup.GET("/pay", apiwrap.Wrap(h.AdminGetPayConfig))
	adminGroup.POST("/pay", apiwrap.WrapWithBody(h.AdminAddPayInfo))
	adminGroup.DELETE("/pay/:name", apiwrap.Wrap(h.AdminDeletePayInfo))
	adminGroup.GET("/social", apiwrap.Wrap(h.AdminGetSocialConfig))
	adminGroup.POST("/social", apiwrap.WrapWithBody(h.AdminAddSocialConfig))
	adminGroup.PUT("/social/:id", apiwrap.WrapWithBody(h.AdminUpdateSocialConfig))
	adminGroup.DELETE("/social/:id", apiwrap.Wrap(h.AdminDeleteSocialConfig))

	engine.POST("/admin/login", apiwrap.WrapWithBody(h.AdminLogin))

}

func (h *WebsiteConfigHandler) GetIndexConfig(ctx *gin.Context) (*apiwrap.ResponseBody[IndexConfigVO], error) {
	config, err := h.serv.GetIndexConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(IndexConfigVO{
		WebsiteConfig:      h.toWebsiteConfigVO(&config.WebSiteConfig),
		NoticeConfigVO:     h.toNoticeConfigVO(&config.NoticeConfig),
		SocialInfoConfigVO: h.toSocialInfoConfigVO(&config.SocialInfoConfig),
		PayInfoConfigVO:    h.toPayInfoConfigVO(config.PayInfoConfig),
		SeoMetaConfigVO:    h.toSeoMetaConfigVO(&config.SeoMetaConfig),
	}), nil
}

func (h *WebsiteConfigHandler) toWebsiteConfigVO(webMasterCfg *domain.WebsiteConfig) WebsiteConfigVO {
	return WebsiteConfigVO{
		WebsiteName:         webMasterCfg.WebsiteName,
		WebsiteIcon:         webMasterCfg.WebsiteIcon,
		WebsiteOwner:        webMasterCfg.WebsiteOwner,
		WebsiteOwnerProfile: webMasterCfg.WebsiteOwnerProfile,
		WebsiteOwnerAvatar:  webMasterCfg.WebsiteOwnerAvatar,
		WebsiteOwnerEmail:   webMasterCfg.WebsiteOwnerEmail,
		WebsiteRuntime:      webMasterCfg.WebsiteRuntime.Unix(),
		WebsiteRecords:      webMasterCfg.WebsiteRecords,
	}
}

func (h *WebsiteConfigHandler) toNoticeConfigVO(noticeCfg *domain.NoticeConfig) NoticeConfigVO {
	return NoticeConfigVO{Title: noticeCfg.Title, Content: noticeCfg.Content, Enabled: noticeCfg.Enabled, PublishTime: noticeCfg.PublishTime.Unix()}
}

func (h *WebsiteConfigHandler) toSocialInfoConfigVO(socialINfoConfig *domain.SocialInfoConfig) SocialInfoConfigVO {
	socialInfoVOList := make([]SocialInfoVO, len(socialINfoConfig.SocialInfoList))
	for i, socialInfo := range socialINfoConfig.SocialInfoList {
		socialInfoVOList[i] = SocialInfoVO{SocialName: socialInfo.SocialName, SocialValue: socialInfo.SocialValue, CssClass: socialInfo.CssClass, IsLink: socialInfo.IsLink}
	}
	return SocialInfoConfigVO{SocialInfoList: socialInfoVOList}
}

func (h *WebsiteConfigHandler) toPayInfoConfigVO(config []domain.PayInfoConfigElem) []PayInfoConfigVO {
	result := make([]PayInfoConfigVO, len(config))
	for i, payInfoConfig := range config {
		result[i] = PayInfoConfigVO{
			Name:  payInfoConfig.Name,
			Image: payInfoConfig.Image,
		}
	}
	return result
}

func (h *WebsiteConfigHandler) toSeoMetaConfigVO(config *domain.SeoMetaConfig) SeoMetaConfigVO {
	return SeoMetaConfigVO{
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

func (h *WebsiteConfigHandler) AdminGetWebsiteConfig(ctx *gin.Context) (*apiwrap.ResponseBody[WebsiteConfigVO], error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toWebsiteConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateWebsiteConfig(ctx *gin.Context, req UpdateWebsiteConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateWebsiteConfig(ctx, domain.WebsiteConfig{
		WebsiteName:         req.WebsiteName,
		WebsiteIcon:         req.WebsiteIcon,
		WebsiteOwner:        req.WebsiteOwner,
		WebsiteOwnerProfile: req.WebsiteOwnerProfile,
		WebsiteOwnerAvatar:  req.WebsiteOwnerAvatar,
		WebsiteOwnerEmail:   req.WebsiteOwnerEmail,
		WebsiteRuntime:      time.Unix(req.WebsiteRuntime, 0).Local(),
		WebsiteRecords:      nil,
	}, time.Now())
}

func (h *WebsiteConfigHandler) AdminGetSeoConfig(ctx *gin.Context) (*apiwrap.ResponseBody[SeoMetaConfigVO], error) {
	config, err := h.serv.GetSeoMetaConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toSeoMetaConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateSeoConfig(ctx *gin.Context, req UpdateSeoMetaConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateSeoMetaConfig(ctx, &domain.SeoMetaConfig{
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

func (h *WebsiteConfigHandler) AdminGetCommentConfig(ctx *gin.Context) (*apiwrap.ResponseBody[CommentConfigVO], error) {
	config, err := h.serv.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toCommentConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toCommentConfigVO(config domain.CommentConfig) CommentConfigVO {
	return CommentConfigVO{
		EnableComment: config.EnableComment,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateCommentConfig(ctx *gin.Context, req UpdateCommentConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateCommentConfig(ctx, domain.CommentConfig{
		EnableComment: gkit.GetValueOrDefault(req.EnableComment),
	})
}

func (h *WebsiteConfigHandler) AdminGetFriendConfig(ctx *gin.Context) (*apiwrap.ResponseBody[FriendConfigVO], error) {
	config, err := h.serv.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toFriendConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toFriendConfigVO(config domain.FriendConfig) FriendConfigVO {
	return FriendConfigVO{
		EnableFriendCommit: config.EnableFriendCommit,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFriendConfig(ctx *gin.Context, req UpdateFriendConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateFriendConfig(ctx, domain.FriendConfig{
		EnableFriendCommit: gkit.GetValueOrDefault(req.EnableFriendCommit),
	})
}

func (h *WebsiteConfigHandler) AdminGetEmailConfig(ctx *gin.Context) (*apiwrap.ResponseBody[EmailConfigVO], error) {
	config, err := h.serv.GetEmailConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toEmailConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toEmailConfigVO(config *domain.EmailConfig) EmailConfigVO {
	return EmailConfigVO{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Email:    config.Email,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateEmailConfig(ctx *gin.Context, req UpdateEmailConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateEmailConfig(ctx, &domain.EmailConfig{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
}

func (h *WebsiteConfigHandler) AdminGetNoticeConfig(ctx *gin.Context) (*apiwrap.ResponseBody[NoticeConfigVO], error) {
	config, err := h.serv.GetNoticeConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toNoticeConfigVO(&config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeConfig(ctx *gin.Context, req UpdateNoticeConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateNoticeConfig(ctx, &domain.NoticeConfig{
		Title:   req.Title,
		Content: req.Content,
	})
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeEnabled(ctx *gin.Context, req UpdateNoticeConfigEnabledReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateNoticeConfigEnabled(ctx, gkit.GetValueOrDefault(req.Enabled))
}

func (h *WebsiteConfigHandler) AdminGetFPCConfig(ctx *gin.Context) (*apiwrap.ResponseBody[FrontPostCountConfigVO], error) {
	config, err := h.serv.GetFrontPostCountConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toFrontPostCountConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toFrontPostCountConfigVO(config domain.FrontPostCountConfig) FrontPostCountConfigVO {
	return FrontPostCountConfigVO{
		Count: config.Count,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFPCConfig(ctx *gin.Context, req UpdateFPCConfigCountReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateFrontPostCountConfig(ctx, domain.FrontPostCountConfig{
		Count: req.Count,
	})
}

func (h *WebsiteConfigHandler) AdminAddRecordInWebsiteConfig(ctx *gin.Context, req AddRecordInWebsiteConfig) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddRecordInWebsiteConfig(ctx, req.Record)
}

func (h *WebsiteConfigHandler) AdminDeleteRecordInWebsiteConfig(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	record := ctx.Query("record")
	if record == "" {
		return nil, errors.New("record is empty")
	}
	return apiwrap.SuccessResponse(), h.serv.DeleteRecordInWebsiteConfig(ctx, record)
}

func (h *WebsiteConfigHandler) AdminGetPayConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[PayInfoConfigVO]], error) {
	config, err := h.serv.GetPayConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toPayInfoConfigVO(config.List))), nil
}

func (h *WebsiteConfigHandler) AdminAddPayInfo(ctx *gin.Context, req AddPayInfoRequest) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddPayInfo(ctx, domain.PayInfoConfigElem{
		Name:  req.Name,
		Image: req.Image,
	})
}

func (h *WebsiteConfigHandler) AdminDeletePayInfo(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	name := ctx.Param("name")
	if name == "" {
		return nil, errors.New("name is empty")
	}
	image := ctx.Query("image")
	if image == "" {
		return nil, errors.New("image is empty")
	}
	return apiwrap.SuccessResponse(), h.serv.DeletePayInfo(ctx, domain.PayInfoConfigElem{
		Name:  name,
		Image: image,
	})
}

func (h *WebsiteConfigHandler) AdminGetSocialConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[AdminSocialInfoVO]], error) {
	socialConfig, err := h.serv.GetSocialConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toAdminSocialInfoVO(socialConfig.SocialInfoList))), nil
}

func (h *WebsiteConfigHandler) toAdminSocialInfoVO(list []domain.SocialInfo) []AdminSocialInfoVO {
	result := make([]AdminSocialInfoVO, len(list))
	for i, socialInfo := range list {
		result[i] = AdminSocialInfoVO{
			Id: hex.EncodeToString(socialInfo.Id),
			SocialInfoVO: SocialInfoVO{
				SocialName:  socialInfo.SocialName,
				SocialValue: socialInfo.SocialValue,
				CssClass:    socialInfo.CssClass,
				IsLink:      socialInfo.IsLink,
			},
		}
	}
	return result
}

func (h *WebsiteConfigHandler) AdminAddSocialConfig(ctx *gin.Context, req SocialInfoReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddSocialInfo(ctx, domain.SocialInfo{
		SocialName:  req.SocialName,
		SocialValue: req.SocialValue,
		CssClass:    req.CssClass,
		IsLink:      req.IsLink,
	})
}

func (h *WebsiteConfigHandler) AdminUpdateSocialConfig(ctx *gin.Context, req SocialInfoReq) (*apiwrap.ResponseBody[any], error) {
	sid := ctx.Param("id")
	id, err := hex.DecodeString(sid)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), h.serv.UpdateSocialInfo(ctx, domain.SocialInfo{
		Id:          id,
		SocialName:  req.SocialName,
		SocialValue: req.SocialValue,
		CssClass:    req.CssClass,
		IsLink:      req.IsLink,
	})
}

func (h *WebsiteConfigHandler) AdminDeleteSocialConfig(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	sid := ctx.Param("id")
	id, err := hex.DecodeString(sid)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), h.serv.DeleteSocialInfo(ctx, id)
}

func (h *WebsiteConfigHandler) AdminLogin(ctx *gin.Context, req LoginRequest) (*apiwrap.ResponseBody[LoginVO], error) {
	adminConfig, err := h.serv.GetAdminConfig(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if adminConfig == nil || adminConfig.Username != req.Username || adminConfig.Password != req.Password {
		return nil, *apiwrap.NewResponseBody[any](40101, "username or password is incorrect", nil)
	}
	jwt, exp, err := jwtutil.GenerateJwt()
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(LoginVO{Token: jwt, Expiration: exp}), nil
}

func (h *WebsiteConfigHandler) GetInitStatus(_ *gin.Context) (*apiwrap.ResponseBody[map[string]bool], error) {
	return apiwrap.SuccessResponseWithData(map[string]bool{
		"initStatus": global.IsWebsiteInitialized(),
	}), nil
}

func (h *WebsiteConfigHandler) InitializeWebsite(ctx *gin.Context, req InitRequest) (*apiwrap.ResponseBody[any], error) {
	h.initMutex.Lock()
	defer h.initMutex.Unlock()
	if global.IsWebsiteInitialized() {
		return nil, apiwrap.NewErrorResponseBody(http.StatusForbidden, "website has been initialized")
	}
	err := h.serv.InitializeWebsite(ctx, domain.AdminConfig{
		Username: req.Admin.Username,
		Password: req.Admin.Password,
	}, domain.WebsiteConfig{
		WebsiteName:         req.WebsiteName,
		WebsiteIcon:         req.WebsiteIcon,
		WebsiteOwner:        req.WebsiteOwner,
		WebsiteOwnerProfile: req.WebsiteOwnerProfile,
		WebsiteOwnerAvatar:  req.WebsiteOwnerAvatar,
		WebsiteOwnerEmail:   req.WebsiteOwnerEmail,
		WebsiteInit:         gkit.ToPtr(true),
	}, domain.EmailConfig{
		Host:     req.EmailServer.Host,
		Port:     req.EmailServer.Port,
		Username: req.EmailServer.Username,
		Password: req.EmailServer.Password,
		Email:    req.EmailServer.Email,
	})
	if err == nil {
		global.Config.IsWebsiteInitialized = true
	}
	return apiwrap.SuccessResponse(), err
}
