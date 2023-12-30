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

	"github.com/chenmingyong0423/fnote/backend/internal/config/repository"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type IConfigService interface {
	GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfig, error)
	GetSwitchStatusByTyp(ctx context.Context, typ string) (*domain.SwitchConfig, error)
	IncreaseWebsiteViews(ctx context.Context) error
	GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error)
	GetIndexConfig(ctx context.Context) (*domain.IndexConfig, error)
	GetFrontPostCount(ctx context.Context) (*domain.FrontPostCount, error)
}

var _ IConfigService = (*ConfigService)(nil)

func NewConfigService(repo repository.IConfigRepository) *ConfigService {
	return &ConfigService{
		repo: repo,
	}
}

type ConfigService struct {
	repo repository.IConfigRepository
}

func (s *ConfigService) GetFrontPostCount(ctx context.Context) (*domain.FrontPostCount, error) {
	cfg := &domain.FrontPostCount{}
	err := s.getConfigAndConvertTo(ctx, "front-post-count", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *ConfigService) GetIndexConfig(ctx context.Context) (*domain.IndexConfig, error) {
	configs, err := s.repo.GetConfigByTypes(ctx, "webmaster", "notice", "social", "pay")
	if err != nil {
		return nil, err
	}
	cfg := &domain.IndexConfig{}
	for _, config := range configs {
		if config.Typ == "webmaster" {
			wmCfg := domain.WebMasterConfig{}
			err = s.anyToStruct(config.Props, &wmCfg)
			if err != nil {
				return nil, err
			}
			cfg.WebMasterConfig = wmCfg
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
		}
	}
	return cfg, nil
}

func (s *ConfigService) GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error) {
	emailConfig := new(domain.EmailConfig)
	err := s.getConfigAndConvertTo(ctx, "emailConfig", emailConfig)
	if err != nil {
		return nil, err
	}
	return emailConfig, nil
}

func (s *ConfigService) getConfigAndConvertTo(ctx context.Context, typ string, config any) error {
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

func (s *ConfigService) IncreaseWebsiteViews(ctx context.Context) error {
	return s.repo.Increase(ctx, "websiteViews")
}

func (s *ConfigService) GetSwitchStatusByTyp(ctx context.Context, typ string) (*domain.SwitchConfig, error) {
	switchConfig := new(domain.SwitchConfig)
	err := s.getConfigAndConvertTo(ctx, typ, switchConfig)
	if err != nil {
		return nil, err
	}
	return switchConfig, nil
}

func (s *ConfigService) GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfig, error) {
	props, err := s.repo.FindByTyp(ctx, typ)
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.FindByTyp failed")
	}
	webMasterConfig := new(domain.WebMasterConfig)
	err = s.anyToStruct(props, webMasterConfig)
	return webMasterConfig, err
}

func (s *ConfigService) anyToStruct(props any, cfgInfos any) error {
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
