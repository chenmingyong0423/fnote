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
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"

	"github.com/chenmingyong0423/fnote/backend/ineternal/config/repository"
	"github.com/pkg/errors"
)

type IConfigService interface {
	GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfigVO, error)
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

func (s *ConfigService) GetWebmasterInfo(ctx context.Context, typ string) (*domain.WebMasterConfigVO, error) {
	webMasterConfig, err := s.repo.FindByTyp(ctx, typ)
	if err != nil {
		return nil, errors.WithMessage(err, "s.repo.FindByTyp failed")
	}
	return &domain.WebMasterConfigVO{Name: webMasterConfig.Name, PostCount: webMasterConfig.PostCount, ColumnCount: webMasterConfig.ColumnCount, WebsiteViews: webMasterConfig.WebsiteViews, WebsiteLiveTime: webMasterConfig.WebsiteLiveTime, Profile: webMasterConfig.Profile, Picture: webMasterConfig.Picture, WebsiteIcon: webMasterConfig.WebsiteIcon, Domain: webMasterConfig.Domain}, nil
}
