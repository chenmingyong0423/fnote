// Copyright 2024 chenmingyong0423

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

	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
)

type IPostIndexService interface {
	PushUrls2Baidu(ctx context.Context, urls string) (*domain.BaiduResponse, error)
}

var _ IPostIndexService = (*PostIndexService)(nil)

func NewPostIndexService(baiduServ *BaiduService, cfgServ website_config.Service) *PostIndexService {
	return &PostIndexService{
		baiduServ: baiduServ,
		cfgServ:   cfgServ,
	}
}

type PostIndexService struct {
	baiduServ *BaiduService
	cfgServ   website_config.Service
}

func (s *PostIndexService) PushUrls2Baidu(ctx context.Context, urls string) (*domain.BaiduResponse, error) {
	// 查询百度推送配置
	bdCfg, err := s.cfgServ.GetBaiduPushConfig(ctx)
	if err != nil {
		return nil, err
	}
	if bdCfg.Token == "" {
		return nil, nil
	}
	return s.baiduServ.Push(ctx, bdCfg.Site, bdCfg.Token, urls)
}
