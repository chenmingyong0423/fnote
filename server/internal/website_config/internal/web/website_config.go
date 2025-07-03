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
	adminGroup := engine.Group("/admin-api/configs")

	// 获取首页的配置信息
	routerGroup.GET("/index", apiwrap.Wrap(h.GetIndexConfig))
	routerGroup.GET("/common", apiwrap.Wrap(h.GetCommonConfig))

	// 轮播图
	routerGroup.GET("/index/carousel", apiwrap.Wrap(h.GetCarouselConfig))
	adminGroup.POST("/carousel", apiwrap.WrapWithBody(h.AddCarouselConfig))
	adminGroup.PUT("/carousel/:id", apiwrap.WrapWithBody(h.UpdateCarouselElem))
	adminGroup.DELETE("/carousel/:id", apiwrap.Wrap(h.DeleteCarouselElem))
	adminGroup.GET("/carousel", apiwrap.Wrap(h.AdminGetCarouselConfig))
	adminGroup.PUT("/carousel/:id/show", apiwrap.WrapWithBody(h.AdminUpdateCarouselShowStatus))

	// 初始化相关
	routerGroup.GET("/check-initialization", apiwrap.Wrap(h.GetInitStatus))
	adminGroup.GET("/check-initialization", apiwrap.Wrap(h.GetInitStatus))
	adminGroup.POST("/initialization", apiwrap.WrapWithBody(h.InitializeWebsite))

	// website
	routerGroup.GET("/website", apiwrap.Wrap(h.GetWebsiteConfig))
	adminGroup.GET("/website", apiwrap.Wrap(h.AdminGetWebsiteConfig))
	adminGroup.GET("/website/meta", apiwrap.Wrap(h.AdminGetWebsiteConfig4Meta))
	adminGroup.PUT("/website", apiwrap.WrapWithBody(h.AdminUpdateWebsiteConfig))
	adminGroup.POST("/website/records", apiwrap.WrapWithBody(h.AdminAddRecordInWebsiteConfig))
	adminGroup.DELETE("/website/records", apiwrap.Wrap(h.AdminDeleteRecordInWebsiteConfig))

	// seo
	adminGroup.GET("/seo", apiwrap.Wrap(h.AdminGetSeoConfig))
	adminGroup.PUT("/seo", apiwrap.WrapWithBody(h.AdminUpdateSeoConfig))

	// 评论
	adminGroup.GET("/comment", apiwrap.Wrap(h.AdminGetCommentConfig))
	adminGroup.PUT("/comment", apiwrap.WrapWithBody(h.AdminUpdateCommentConfig))
	// 友链
	adminGroup.GET("/friend", apiwrap.Wrap(h.AdminGetFriendConfig))
	adminGroup.PUT("/friend/switch", apiwrap.WrapWithBody(h.AdminUpdateSwitch4FriendConfig))
	adminGroup.PUT("/friend/introduction", apiwrap.WrapWithBody(h.AdminUpdateIntroduction4FriendConfig))
	// 邮件配置
	adminGroup.GET("/email", apiwrap.Wrap(h.AdminGetEmailConfig))
	adminGroup.PUT("/email", apiwrap.WrapWithBody(h.AdminUpdateEmailConfig))

	// 公告配置
	routerGroup.GET("/notice", apiwrap.Wrap(h.GetNoticeConfig))
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
	adminGroup.GET("/third-party-site-verification", apiwrap.Wrap(h.AdminGetTPSVConfig))
	adminGroup.POST("/third-party-site-verification", apiwrap.WrapWithBody(h.AdminAddTPSVConfig))
	adminGroup.DELETE("/third-party-site-verification/:key", apiwrap.Wrap(h.AdminDeleteTPSVConfig))
	adminGroup.GET("/post-index/:key", apiwrap.Wrap(h.AdminGetPushConfigByKey))
	adminGroup.PUT("/post-index/:key", apiwrap.WrapWithBody(h.AdminUpdatePushConfigByKey))

	engine.POST("/admin-api/login", apiwrap.WrapWithBody(h.AdminLogin))

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
		TPSVVO:             h.toTPSVVO(config.TPSVConfig),
	}), nil
}

func (h *WebsiteConfigHandler) toWebsiteConfigVO(webMasterCfg *domain.WebsiteConfig) WebsiteConfigVO {
	return WebsiteConfigVO{
		WebsiteName:         webMasterCfg.WebsiteName,
		WebsiteIcon:         webMasterCfg.WebsiteIcon,
		WebsiteOwner:        webMasterCfg.WebsiteOwner,
		WebsiteOwnerProfile: webMasterCfg.WebsiteOwnerProfile,
		WebsiteOwnerAvatar:  webMasterCfg.WebsiteOwnerAvatar,
		WebsiteRuntime:      gkit.GetValueOrDefault(webMasterCfg.WebsiteRuntime).Unix(),
		WebsiteRecords:      webMasterCfg.WebsiteRecords,
	}
}

func (h *WebsiteConfigHandler) toNoticeConfigVO(noticeCfg *domain.NoticeConfig) NoticeConfigVO {
	return NoticeConfigVO{Title: noticeCfg.Title, Content: noticeCfg.Content, PublishTime: noticeCfg.PublishTime.Unix()}
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
		WebsiteRuntime:      gkit.ToPtr(time.Unix(req.WebsiteRuntime, 0).Local()),
	}, time.Now().Local())
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
		Introduction:       config.Introduction,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateSwitch4FriendConfig(ctx *gin.Context, req UpdateFriendSwitchConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateSwitch4FriendConfig(ctx, gkit.GetValueOrDefault(req.EnableFriendCommit))
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
	record := ctx.Query("website_record")
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

func (h *WebsiteConfigHandler) AdminGetTPSVConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[TPSVVO]], error) {
	config, err := h.serv.GetTPSVConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toTPSVVO(config.List))), nil
}

func (h *WebsiteConfigHandler) toTPSVVO(list []domain.TPSV) []TPSVVO {
	result := make([]TPSVVO, len(list))
	for i, tpsv := range list {
		result[i] = TPSVVO{
			Key:         tpsv.Key,
			Value:       tpsv.Value,
			Description: tpsv.Description,
		}
	}
	return result
}

func (h *WebsiteConfigHandler) AdminAddTPSVConfig(ctx *gin.Context, req TPSVRequest) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddTPSVConfig(ctx, domain.TPSV{
		Key:         req.Key,
		Value:       req.Value,
		Description: req.Description,
	})
}

func (h *WebsiteConfigHandler) AdminDeleteTPSVConfig(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.DeleteTPSVConfigByKey(ctx, ctx.Param("key"))
}

func (h *WebsiteConfigHandler) AdminGetPushConfigByKey(ctx *gin.Context) (*apiwrap.ResponseBody[BaiduPushConfigVO], error) {
	key := ctx.Param("key")
	var (
		cfg *domain.Baidu
		err error
	)
	if key == "baidu" {
		cfg, err = h.serv.GetBaiduPushConfig(ctx)
		if err != nil {
			return nil, err
		}
	}
	return apiwrap.SuccessResponseWithData(BaiduPushConfigVO{
		Site:  cfg.Site,
		Token: cfg.Token,
	}), nil
}

func (h *WebsiteConfigHandler) AdminUpdatePushConfigByKey(ctx *gin.Context, req map[string]any) (*apiwrap.ResponseBody[any], error) {
	if len(req) == 0 {
		return nil, apiwrap.NewErrorResponseBody(400, "request body is nil.")
	}
	return apiwrap.SuccessResponse(), h.serv.UpdatePushConfigByKey(ctx, ctx.Param("key"), req)
}

func (h *WebsiteConfigHandler) AdminGetWebsiteConfig4Meta(ctx *gin.Context) (*apiwrap.ResponseBody[WebsiteConfigMetaVO], error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(WebsiteConfigMetaVO{
		WebsiteName: config.WebsiteName,
		WebsiteIcon: config.WebsiteIcon,
	}), nil
}

func (h *WebsiteConfigHandler) GetCarouselConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[CarouselVO]], error) {
	config, err := h.serv.GetCarouselConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toCarouselVO(config.List, func(elem domain.CarouselElem) bool {
		return elem.Show
	}))), nil
}

func (h *WebsiteConfigHandler) toCarouselVO(list []domain.CarouselElem, filterFunc func(elem domain.CarouselElem) bool) []CarouselVO {
	result := make([]CarouselVO, 0, len(list))
	for _, elem := range list {
		if filterFunc(elem) {
			result = append(result, CarouselVO{
				Id:        elem.Id,
				Title:     elem.Title,
				Summary:   elem.Summary,
				CoverImg:  elem.CoverImg,
				Show:      elem.Show,
				Color:     elem.Color,
				CreatedAt: elem.CreatedAt.Unix(),
				UpdatedAt: elem.UpdatedAt.Unix(),
			})
		}
	}
	return result
}

func (h *WebsiteConfigHandler) AddCarouselConfig(ctx *gin.Context, req CarouselRequest) (*apiwrap.ResponseBody[any], error) {
	exist, err := h.serv.IsCarouselElemExist(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, apiwrap.NewErrorResponseBody(409, "carousel element already exists")
	}
	now := time.Now().Local()
	return apiwrap.SuccessResponse(), h.serv.AddCarouselConfig(ctx, domain.CarouselElem{
		Id:        req.Id,
		Title:     req.Title,
		Summary:   req.Summary,
		CoverImg:  req.CoverImg,
		Show:      req.Show,
		Color:     req.Color,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (h *WebsiteConfigHandler) AdminGetCarouselConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[CarouselVO]], error) {
	config, err := h.serv.GetCarouselConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toCarouselVO(config.List, func(_ domain.CarouselElem) bool {
		return true
	}))), nil
}

func (h *WebsiteConfigHandler) AdminUpdateCarouselShowStatus(ctx *gin.Context, req CarouselShowRequest) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateCarouselShowStatus(ctx, ctx.Param("id"), req.Show)
}

func (h *WebsiteConfigHandler) UpdateCarouselElem(ctx *gin.Context, req CarouselRequest) (*apiwrap.ResponseBody[any], error) {
	now := time.Now().Local()
	return apiwrap.SuccessResponse(), h.serv.UpdateCarouselElem(ctx, domain.CarouselElem{
		Id:        ctx.Param("id"),
		Title:     req.Title,
		Summary:   req.Summary,
		CoverImg:  req.CoverImg,
		Show:      req.Show,
		Color:     req.Color,
		CreatedAt: now,
		UpdatedAt: now,
	})
}

func (h *WebsiteConfigHandler) DeleteCarouselElem(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.DeleteCarouselElem(ctx, ctx.Param("id"))
}

func (h *WebsiteConfigHandler) AdminUpdateIntroduction4FriendConfig(ctx *gin.Context, req UpdateFriendIntroConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateIntroduction4FriendConfig(ctx, gkit.GetValueOrDefault(req.Introduction))
}

func (h *WebsiteConfigHandler) GetWebsiteConfig(ctx *gin.Context) (*apiwrap.ResponseBody[WebsiteConfigVO], error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toWebsiteConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) GetCommonConfig(ctx *gin.Context) (*apiwrap.ResponseBody[CommonConfigVO], error) {
	config, err := h.serv.GetCommonConfig(ctx)
	if err != nil {
		return nil, err
	}

	return apiwrap.SuccessResponseWithData(CommonConfigVO{
		WebsiteMeta: h.toMetaConfigVO(&config.WebSiteConfig),
		SeoMeta:     h.toSeoMetaConfigVO(&config.SeoMetaConfig),
		TPSVVO:      h.toTPSVVO(config.TPSVConfig),
		Records:     config.WebSiteConfig.WebsiteRecords,
	}), nil

}

func (h *WebsiteConfigHandler) toMetaConfigVO(wc *domain.WebsiteConfig) WebsiteConfigMetaVO {
	return WebsiteConfigMetaVO{
		WebsiteName:  wc.WebsiteName,
		WebsiteIcon:  wc.WebsiteIcon,
		WebsiteOwner: wc.WebsiteOwner,
	}
}

func (h *WebsiteConfigHandler) GetNoticeConfig(ctx *gin.Context) (*apiwrap.ResponseBody[NoticeConfigVO], error) {
	config, err := h.serv.GetNoticeConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toNoticeConfigVO(&config)), nil
}
