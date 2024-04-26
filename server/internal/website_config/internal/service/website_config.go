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

package service

import (
	"context"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository"

	"github.com/gin-gonic/gin"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type IWebsiteConfigService interface {
	GetWebSiteConfig(ctx context.Context) (*domain.WebsiteConfig, error)
	GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error)
	GetIndexConfig(ctx context.Context) (*domain.IndexConfig, error)
	GetFrontPostCount(ctx context.Context) (*domain.FrontPostCountConfig, error)
	GetSeoMetaConfig(ctx context.Context) (*domain.SeoMetaConfig, error)
	UpdateSeoMetaConfig(ctx context.Context, seoCfg *domain.SeoMetaConfig) error
	GetCommentConfig(ctx context.Context) (domain.CommentConfig, error)
	UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error
	GetFriendConfig(ctx context.Context) (domain.FriendConfig, error)
	UpdateFriendConfig(ctx context.Context, friendConfig domain.FriendConfig) error
	UpdateEmailConfig(ctx context.Context, emailCfg *domain.EmailConfig) error
	GetNoticeConfig(ctx context.Context) (domain.NoticeConfig, error)
	UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error
	UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error
	GetFrontPostCountConfig(ctx context.Context) (domain.FrontPostCountConfig, error)
	UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error
	AddRecordInWebsiteConfig(ctx context.Context, record string) error
	DeleteRecordInWebsiteConfig(ctx context.Context, record string) error
	GetPayConfig(ctx context.Context) (domain.PayInfoConfig, error)
	AddPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error
	DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error
	GetSocialConfig(ctx context.Context) (domain.SocialInfoConfig, error)
	AddSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error
	UpdateSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error
	DeleteSocialInfo(ctx context.Context, id []byte) error
	GetAdminConfig(ctx context.Context) (*domain.AdminConfig, error)
	InitializeWebsite(ctx context.Context, adminConfig domain.AdminConfig, webSiteConfig domain.WebsiteConfig, emailConfig domain.EmailConfig) error
	UpdateWebsiteConfig(ctx context.Context, websiteConfig domain.WebsiteConfig, now time.Time) error
	GetTPSVConfig(ctx context.Context) (*domain.TPSVConfig, error)
	AddTPSVConfig(ctx context.Context, tpsv domain.TPSV) error
	DeleteTPSVConfigByKey(ctx context.Context, key string) error
	GetBaiduPushConfig(ctx context.Context) (*domain.Baidu, error)
	UpdatePushConfigByKey(ctx context.Context, key string, updates map[string]any) error
}

var _ IWebsiteConfigService = (*WebsiteConfigService)(nil)

func NewWebsiteConfigService(repo repository.IWebsiteConfigRepository) *WebsiteConfigService {
	return &WebsiteConfigService{
		repo: repo,
	}
}

type WebsiteConfigService struct {
	repo repository.IWebsiteConfigRepository
}

func (s *WebsiteConfigService) UpdatePushConfigByKey(ctx context.Context, key string, updates map[string]any) error {
	return s.repo.UpdatePushConfigByKey(ctx, key, updates)
}

func (s *WebsiteConfigService) GetBaiduPushConfig(ctx context.Context) (*domain.Baidu, error) {
	return s.repo.GetBaiduPushConfig(ctx)
}

func (s *WebsiteConfigService) DeleteTPSVConfigByKey(ctx context.Context, key string) error {
	return s.repo.DeleteTPSVConfigByKey(ctx, key)
}

func (s *WebsiteConfigService) AddTPSVConfig(ctx context.Context, tpsv domain.TPSV) error {
	return s.repo.AddTPSVConfig(ctx, tpsv)
}

func (s *WebsiteConfigService) GetTPSVConfig(ctx context.Context) (*domain.TPSVConfig, error) {
	return s.repo.GetTPSVConfig(ctx)
}

func (s *WebsiteConfigService) InitializeWebsite(ctx context.Context, adminConfig domain.AdminConfig, webSiteConfig domain.WebsiteConfig, emailConfig domain.EmailConfig) error {
	now := time.Now()
	err := s.UpdateAdminConfig(ctx, adminConfig, now)
	if err != nil {
		return err
	}

	err = s.repo.UpdateEmailConfig(ctx, &emailConfig, now)
	if err != nil {
		return err
	}

	err = s.UpdateWebsiteConfig(ctx, webSiteConfig, now)
	if err != nil {
		return err
	}
	return nil
}

func (s *WebsiteConfigService) GetAdminConfig(ctx context.Context) (*domain.AdminConfig, error) {
	cfg := &domain.AdminConfig{}
	err := s.getConfigAndConvertTo(ctx, "admin", cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) DeleteSocialInfo(ctx context.Context, id []byte) error {
	return s.repo.DeleteSocialInfo(ctx, id)
}

func (s *WebsiteConfigService) UpdateSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error {
	return s.repo.UpdateSocialInfo(ctx, socialInfo)
}

func (s *WebsiteConfigService) AddSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error {
	return s.repo.AddSocialInfo(ctx, socialInfo)
}

func (s *WebsiteConfigService) GetSocialConfig(ctx context.Context) (domain.SocialInfoConfig, error) {
	cfg := domain.SocialInfoConfig{}
	err := s.getConfigAndConvertTo(ctx, "social", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return s.repo.DeletePayInfo(ctx, payInfoConfigElem)
}

func (s *WebsiteConfigService) AddPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return s.repo.PushPayInfo(ctx, payInfoConfigElem)
}

func (s *WebsiteConfigService) GetPayConfig(ctx context.Context) (domain.PayInfoConfig, error) {
	cfg := domain.PayInfoConfig{List: make([]domain.PayInfoConfigElem, 0)}
	err := s.getConfigAndConvertTo(ctx, "pay", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) DeleteRecordInWebsiteConfig(ctx context.Context, record string) error {
	return s.repo.DeleteRecordInWebsiteConfig(ctx, record)
}

func (s *WebsiteConfigService) AddRecordInWebsiteConfig(ctx context.Context, record string) error {
	return s.repo.AddRecordInWebsiteConfig(ctx, record)
}

func (s *WebsiteConfigService) UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error {
	return s.repo.UpdateFrontPostCountConfig(ctx, cfg)
}

func (s *WebsiteConfigService) GetFrontPostCountConfig(ctx context.Context) (domain.FrontPostCountConfig, error) {
	cfg := domain.FrontPostCountConfig{}
	err := s.getConfigAndConvertTo(ctx, "front-post-count", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error {
	return s.repo.UpdateNoticeConfigEnabled(ctx, enabled)
}

func (s *WebsiteConfigService) UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error {
	return s.repo.UpdateNoticeConfig(ctx, noticeCfg)
}

func (s *WebsiteConfigService) GetNoticeConfig(ctx context.Context) (domain.NoticeConfig, error) {
	cfg := domain.NoticeConfig{}
	err := s.getConfigAndConvertTo(ctx, "notice", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) UpdateEmailConfig(ctx context.Context, emailCfg *domain.EmailConfig) error {
	return s.repo.UpdateEmailConfig(ctx, emailCfg, time.Now())
}

func (s *WebsiteConfigService) UpdateFriendConfig(ctx context.Context, friendConfig domain.FriendConfig) error {
	return s.repo.UpdateFriendConfig(ctx, friendConfig)
}

func (s *WebsiteConfigService) GetFriendConfig(ctx context.Context) (domain.FriendConfig, error) {
	cfg := domain.FriendConfig{}
	err := s.getConfigAndConvertTo(ctx, "friend", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error {
	return s.repo.UpdateCommentConfig(ctx, commentConfig)
}

func (s *WebsiteConfigService) GetCommentConfig(ctx context.Context) (domain.CommentConfig, error) {
	cfg := domain.CommentConfig{}
	err := s.getConfigAndConvertTo(ctx, "comment", &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) UpdateSeoMetaConfig(ctx context.Context, seoCfg *domain.SeoMetaConfig) error {
	return s.repo.UpdateSeoMetaConfig(ctx, seoCfg)
}

func (s *WebsiteConfigService) GetSeoMetaConfig(ctx context.Context) (*domain.SeoMetaConfig, error) {
	cfg := &domain.SeoMetaConfig{}
	err := s.getConfigAndConvertTo(ctx, "seo meta", cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) GetFrontPostCount(ctx context.Context) (*domain.FrontPostCountConfig, error) {
	cfg := &domain.FrontPostCountConfig{}
	err := s.getConfigAndConvertTo(ctx, "front-post-count", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *WebsiteConfigService) GetIndexConfig(ctx context.Context) (*domain.IndexConfig, error) {
	configs, err := s.repo.FindConfigByTypes(ctx, "website", "owner", "notice", "social", "pay", "seo meta", "third party site verification")
	if err != nil {
		return nil, err
	}
	cfg := &domain.IndexConfig{}
	for _, config := range configs {
		if config.Typ == "website" {
			wsc := domain.WebsiteConfig{}
			err = s.anyToStruct(config.Props, &wsc)
			if err != nil {
				return nil, err
			}
			cfg.WebSiteConfig = wsc
		} else if config.Typ == "notice" {
			noticeCfg := domain.NoticeConfig{}
			err = s.anyToStruct(config.Props, &noticeCfg)
			if err != nil {
				return nil, err
			}
			if noticeCfg.Enabled {
				cfg.NoticeConfig = noticeCfg
			}
		} else if config.Typ == "social" {
			socialInfoConfig := domain.SocialInfoConfig{}
			err = s.anyToStruct(config.Props, &socialInfoConfig)
			if err != nil {
				return nil, err
			}
			cfg.SocialInfoConfig = socialInfoConfig
		} else if config.Typ == "pay" {
			payInfoConfig := domain.PayInfoConfig{List: make([]domain.PayInfoConfigElem, 0)}
			err = s.anyToStruct(config.Props, &payInfoConfig)
			if err != nil {
				return nil, err
			}
			cfg.PayInfoConfig = payInfoConfig.List
		} else if config.Typ == "seo meta" {
			seoMetaConfig := domain.SeoMetaConfig{}
			err = s.anyToStruct(config.Props, &seoMetaConfig)
			if err != nil {
				return nil, err
			}
			cfg.SeoMetaConfig = seoMetaConfig
		} else if config.Typ == "third party site verification" {
			tpsvConfig := domain.TPSVConfig{}
			err = s.anyToStruct(config.Props, &tpsvConfig)
			if err != nil {
				return nil, err
			}
			cfg.TPSVConfig = tpsvConfig.List
		}
	}
	return cfg, nil
}

func (s *WebsiteConfigService) GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error) {
	emailConfig := new(domain.EmailConfig)
	err := s.getConfigAndConvertTo(ctx, "email", emailConfig)
	if err != nil {
		return nil, err
	}
	return emailConfig, nil
}

func (s *WebsiteConfigService) getConfigAndConvertTo(ctx context.Context, typ string, config any) error {
	props, err := s.repo.FindByTyp(ctx, typ)
	if err != nil {
		return err
	}
	err = s.anyToStruct(props, config)
	if err != nil {
		return err
	}
	return nil
}

func (s *WebsiteConfigService) GetWebSiteConfig(ctx context.Context) (*domain.WebsiteConfig, error) {
	props, err := s.repo.FindByTyp(ctx, "website")
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.FindByTyp failed")
	}
	wsc := new(domain.WebsiteConfig)
	err = s.anyToStruct(props, wsc)
	return wsc, err
}

func (s *WebsiteConfigService) anyToStruct(props any, cfgInfos any) error {
	marshal, err := bson.Marshal(props)
	if err != nil {
		return errors.Wrapf(err, "bson.Marshal failed, val=%v", props)
	}
	err = bson.Unmarshal(marshal, cfgInfos)
	if err != nil {
		return errors.Wrapf(err, "bson.Unmarshal failed, data=%v, val=%v", marshal, cfgInfos)
	}
	return nil
}

func (s *WebsiteConfigService) UpdateAdminConfig(ctx context.Context, adminConfig domain.AdminConfig, now time.Time) error {
	return s.repo.UpdateAdminConfig(ctx, adminConfig, now)
}

func (s *WebsiteConfigService) UpdateWebsiteConfig(ctx context.Context, websiteConfig domain.WebsiteConfig, now time.Time) error {
	return s.repo.UpdateWebSiteConfig(ctx, websiteConfig, now)
}
