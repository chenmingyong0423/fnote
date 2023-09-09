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
	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
)

type IConfigService interface {
	GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfigVO, error)
	GetSwitchStatusByTyp(ctx context.Context, typ string) (*domain.SwitchConfig, error)
	IncreaseWebsiteViews(ctx context.Context) error
	GetEmailConfig(ctx context.Context) (*domain.EmailConfig, error)
}

func NewConfigService(repo repository.IConfigRepository) *ConfigService {
	return &ConfigService{
		repo: repo,
	}
}

var _ IConfigService = (*ConfigService)(nil)

type ConfigService struct {
	repo repository.IConfigRepository
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

func (s *ConfigService) GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfigVO, error) {
	props, err := s.repo.FindByTyp(ctx, typ)
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.FindByTyp failed")
	}
	webMasterConfig := new(domain.WebMasterConfig)
	err = s.anyToStruct(props, webMasterConfig)
	if err != nil {
		return nil, err
	}
	return &domain.WebMasterConfigVO{Name: webMasterConfig.Name, PostCount: webMasterConfig.PostCount, ColumnCount: webMasterConfig.ColumnCount, WebsiteViews: webMasterConfig.WebsiteViews, WebsiteLiveTime: webMasterConfig.WebsiteLiveTime, Profile: webMasterConfig.Profile, Picture: webMasterConfig.Picture, WebsiteIcon: webMasterConfig.WebsiteIcon, Domain: webMasterConfig.Domain}, nil
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
