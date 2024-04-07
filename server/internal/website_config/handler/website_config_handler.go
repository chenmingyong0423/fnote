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
	"encoding/hex"

	"github.com/chenmingyong0423/fnote/server/internal/global"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/jwtutil"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

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
	routerGroup.GET("/init_status", api.Wrap(h.GetInitStatus))

	adminGroup := engine.Group("/admin/configs")
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

func (h *WebsiteConfigHandler) GetIndexConfig(ctx *gin.Context) (*vo.IndexConfigVO, error) {
	config, err := h.serv.GetIndexConfig(ctx)
	if err != nil {
		return nil, err
	}
	return &vo.IndexConfigVO{
		WebsiteConfig:      *h.toWebsiteConfigVO(&config.WebSiteConfig),
		NoticeConfigVO:     *h.toNoticeConfigVO(&config.NoticeConfig),
		SocialInfoConfigVO: *h.toSocialInfoConfigVO(&config.SocialInfoConfig),
		PayInfoConfigVO:    h.toPayInfoConfigVO(config.PayInfoConfig),
		SeoMetaConfigVO:    *h.toSeoMetaConfigVO(&config.SeoMetaConfig),
	}, nil
}

func (h *WebsiteConfigHandler) toWebsiteConfigVO(webMasterCfg *domain.WebSiteConfig) *vo.WebsiteConfigVO {
	return &vo.WebsiteConfigVO{
		WebsiteName:  webMasterCfg.WebsiteName,
		Icon:         webMasterCfg.Icon,
		LiveTime:     webMasterCfg.LiveTime,
		Records:      webMasterCfg.Records,
		OwnerName:    webMasterCfg.OwnerName,
		OwnerProfile: webMasterCfg.OwnerProfile,
		OwnerPicture: webMasterCfg.OwnerPicture,
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

func (h *WebsiteConfigHandler) AdminGetWebsiteConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.WebsiteConfigVO], error) {
	config, err := h.serv.GetWebSiteConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(*h.toWebsiteConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateWebsiteConfig(ctx *gin.Context, req request.UpdateWebsiteConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateWebSiteConfig(ctx, domain.WebSiteConfig{
		WebsiteName:  req.WebsiteName,
		Icon:         req.Icon,
		LiveTime:     req.LiveTime,
		OwnerName:    req.OwnerName,
		OwnerProfile: req.OwnerProfile,
		OwnerPicture: req.OwnerPicture,
	})
}

func (h *WebsiteConfigHandler) AdminGetSeoConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.SeoMetaConfigVO], error) {
	config, err := h.serv.GetSeoMetaConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(*h.toSeoMetaConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateSeoConfig(ctx *gin.Context, req request.UpdateSeoMetaConfigReq) (*apiwrap.ResponseBody[any], error) {
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

func (h *WebsiteConfigHandler) AdminGetCommentConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.CommentConfigVO], error) {
	config, err := h.serv.GetCommentConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toCommentConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toCommentConfigVO(config domain.CommentConfig) vo.CommentConfigVO {
	return vo.CommentConfigVO{
		EnableComment: config.EnableComment,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateCommentConfig(ctx *gin.Context, req request.UpdateCommentConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateCommentConfig(ctx, domain.CommentConfig{
		EnableComment: gkit.GetValueOrDefault(req.EnableComment),
	})
}

func (h *WebsiteConfigHandler) AdminGetFriendConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.FriendConfigVO], error) {
	config, err := h.serv.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toFriendConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toFriendConfigVO(config domain.FriendConfig) vo.FriendConfigVO {
	return vo.FriendConfigVO{
		EnableFriendCommit: config.EnableFriendCommit,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFriendConfig(ctx *gin.Context, req request.UpdateFriendConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateFriendConfig(ctx, domain.FriendConfig{
		EnableFriendCommit: gkit.GetValueOrDefault(req.EnableFriendCommit),
	})
}

func (h *WebsiteConfigHandler) AdminGetEmailConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.EmailConfigVO], error) {
	config, err := h.serv.GetEmailConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toEmailConfigVO(config)), nil
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

func (h *WebsiteConfigHandler) AdminUpdateEmailConfig(ctx *gin.Context, req request.UpdateEmailConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateEmailConfig(ctx, &domain.EmailConfig{
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
}

func (h *WebsiteConfigHandler) AdminGetNoticeConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.NoticeConfigVO], error) {
	config, err := h.serv.GetNoticeConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(*h.toNoticeConfigVO(&config)), nil
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeConfig(ctx *gin.Context, req request.UpdateNoticeConfigReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateNoticeConfig(ctx, &domain.NoticeConfig{
		Title:   req.Title,
		Content: req.Content,
	})
}

func (h *WebsiteConfigHandler) AdminUpdateNoticeEnabled(ctx *gin.Context, req request.UpdateNoticeConfigEnabledReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateNoticeConfigEnabled(ctx, gkit.GetValueOrDefault(req.Enabled))
}

func (h *WebsiteConfigHandler) AdminGetFPCConfig(ctx *gin.Context) (*apiwrap.ResponseBody[vo.FrontPostCountConfigVO], error) {
	config, err := h.serv.GetFrontPostCountConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(h.toFrontPostCountConfigVO(config)), nil
}

func (h *WebsiteConfigHandler) toFrontPostCountConfigVO(config domain.FrontPostCountConfig) vo.FrontPostCountConfigVO {
	return vo.FrontPostCountConfigVO{
		Count: config.Count,
	}
}

func (h *WebsiteConfigHandler) AdminUpdateFPCConfig(ctx *gin.Context, req request.UpdateFPCConfigCountReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdateFrontPostCountConfig(ctx, domain.FrontPostCountConfig{
		Count: req.Count,
	})
}

func (h *WebsiteConfigHandler) AdminAddRecordInWebsiteConfig(ctx *gin.Context, req request.AddRecordInWebsiteConfig) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddRecordInWebsiteConfig(ctx, req.Record)
}

func (h *WebsiteConfigHandler) AdminDeleteRecordInWebsiteConfig(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	record := ctx.Query("record")
	if record == "" {
		return nil, errors.New("record is empty")
	}
	return apiwrap.SuccessResponse(), h.serv.DeleteRecordInWebsiteConfig(ctx, record)
}

func (h *WebsiteConfigHandler) AdminGetPayConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.PayInfoConfigVO]], error) {
	config, err := h.serv.GetPayConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toPayInfoConfigVO(config.List))), nil
}

func (h *WebsiteConfigHandler) AdminAddPayInfo(ctx *gin.Context, req request.AddPayInfoRequest) (*apiwrap.ResponseBody[any], error) {
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

func (h *WebsiteConfigHandler) AdminGetSocialConfig(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[vo.AdminSocialInfoVO]], error) {
	socialConfig, err := h.serv.GetSocialConfig(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.toAdminSocialInfoVO(socialConfig.SocialInfoList))), nil
}

func (h *WebsiteConfigHandler) toAdminSocialInfoVO(list []domain.SocialInfo) []vo.AdminSocialInfoVO {
	result := make([]vo.AdminSocialInfoVO, len(list))
	for i, socialInfo := range list {
		result[i] = vo.AdminSocialInfoVO{
			Id: hex.EncodeToString(socialInfo.Id),
			SocialInfoVO: vo.SocialInfoVO{
				SocialName:  socialInfo.SocialName,
				SocialValue: socialInfo.SocialValue,
				CssClass:    socialInfo.CssClass,
				IsLink:      socialInfo.IsLink,
			},
		}
	}
	return result
}

func (h *WebsiteConfigHandler) AdminAddSocialConfig(ctx *gin.Context, req request.SocialInfoReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.AddSocialInfo(ctx, domain.SocialInfo{
		SocialName:  req.SocialName,
		SocialValue: req.SocialValue,
		CssClass:    req.CssClass,
		IsLink:      req.IsLink,
	})
}

func (h *WebsiteConfigHandler) AdminUpdateSocialConfig(ctx *gin.Context, req request.SocialInfoReq) (*apiwrap.ResponseBody[any], error) {
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

func (h *WebsiteConfigHandler) AdminLogin(ctx *gin.Context, req request.LoginRequest) (*apiwrap.ResponseBody[any], error) {
	adminConfig, err := h.serv.GetAdminConfig(ctx)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if adminConfig == nil || adminConfig.Username != req.Username || adminConfig.Password != req.Password {
		return apiwrap.NewResponseBody[any](100001, "username or password is wrong.", nil), nil
	}
	jwt, exp, err := jwtutil.GenerateJwt()
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData[any](vo.LoginVO{Token: jwt, Expiration: exp}), nil
}

func (h *WebsiteConfigHandler) GetInitStatus(_ *gin.Context) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponseWithData[any](map[string]bool{
		"initStatus": global.IsWebsiteInitialized(),
	}), nil
}
