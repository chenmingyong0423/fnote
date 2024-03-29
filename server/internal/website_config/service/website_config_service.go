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

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/repository"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type IWebsiteConfigService interface {
	GetWebSiteConfig(ctx context.Context) (*domain.WebSiteConfig, error)
	GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error)
	GetIndexConfig(ctx context.Context) (*domain.IndexConfig, error)
	GetFrontPostCount(ctx context.Context) (*domain.FrontPostCountConfig, error)
	UpdateWebSiteConfig(ctx context.Context, webSiteConfig domain.WebSiteConfig) error
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
	return s.repo.UpdateEmailConfig(ctx, emailCfg)
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

func (s *WebsiteConfigService) UpdateWebSiteConfig(ctx context.Context, webSiteConfig domain.WebSiteConfig) error {
	return s.repo.UpdateWebSiteConfig(ctx, webSiteConfig)
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
	configs, err := s.repo.FindConfigByTypes(ctx, "website", "owner", "notice", "social", "pay", "seo meta")
	if err != nil {
		return nil, err
	}
	cfg := &domain.IndexConfig{}
	for _, config := range configs {
		if config.Typ == "website" {
			wsc := domain.WebSiteConfig{}
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

func (s *WebsiteConfigService) GetWebSiteConfig(ctx context.Context) (*domain.WebSiteConfig, error) {
	props, err := s.repo.FindByTyp(ctx, "website")
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.FindByTyp failed")
	}
	wsc := new(domain.WebSiteConfig)
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
