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
	"os"
	"strings"

	"github.com/chenmingyong0423/fnote/server/internal/category"
	"github.com/chenmingyong0423/fnote/server/internal/tag"

	jsoniter "github.com/json-iterator/go"

	"github.com/chenmingyong0423/fnote/server/internal/file"
	"github.com/chenmingyong0423/fnote/server/internal/post"

	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
)

type IPostIndexService interface {
	PushUrls2Baidu(ctx context.Context, urls string) (*domain.BaiduResponse, error)
	GenerateSitemap(ctx context.Context) error
	GetSitemap(ctx context.Context) (string, bool, error)
	GetRobotsTxt(ctx context.Context) (string, bool, error)
	GenerateRobotsTxt(ctx context.Context) (string, error)
	SaveRobotsTxt(ctx context.Context, content string) error
}

var _ IPostIndexService = (*PostIndexService)(nil)

func NewPostIndexService(baiduServ *BaiduService, cfgServ website_config.Service, postServ post.Service, fileServ file.Service, categoryServ category.Service, tagServ tag.Service) *PostIndexService {
	return &PostIndexService{
		baiduServ:    baiduServ,
		cfgServ:      cfgServ,
		postServ:     postServ,
		fileServ:     fileServ,
		categoryServ: categoryServ,
		tagServ:      tagServ,
	}
}

type PostIndexService struct {
	baiduServ    *BaiduService
	cfgServ      website_config.Service
	postServ     post.Service
	fileServ     file.Service
	categoryServ category.Service
	tagServ      tag.Service
}

func (s *PostIndexService) GenerateSitemap(ctx context.Context) error {
	posts, err := s.postServ.FindDisplayedPosts(ctx)
	if err != nil {
		return err
	}
	postBytes, err := jsoniter.Marshal(posts)
	if err != nil {
		return err
	}
	categories, err := s.categoryServ.FindEnabledCategories(ctx)
	if err != nil {
		return err
	}
	categoryBytes, err := jsoniter.Marshal(categories)
	if err != nil {
		return err
	}
	tags, err := s.tagServ.FindEnabledTags(ctx)
	if err != nil {
		return err

	}
	tagBytes, err := jsoniter.Marshal(tags)
	if err != nil {
		return err
	}
	return s.fileServ.GenerateSitemap(ctx, postBytes, categoryBytes, tagBytes)
}

func (s *PostIndexService) GetSitemap(ctx context.Context) (string, bool, error) {
	return s.fileServ.GetSitemap(ctx)
}

func (s *PostIndexService) GetRobotsTxt(ctx context.Context) (string, bool, error) {
	return s.fileServ.GetRobotsTxt(ctx)
}

func (s *PostIndexService) GenerateRobotsTxt(ctx context.Context) (string, error) {
	websiteBaseHost := strings.TrimRight(os.Getenv("WEBSITE_BASE_HOST"), "/")
	if websiteBaseHost == "" {
		websiteBaseHost = "http://localhost:3000"
	}
	content := "User-agent: *\nAllow: /\nDisallow: /admin\n\nSitemap: " + websiteBaseHost + "/sitemap.xml\n"
	return content, s.fileServ.SaveRobotsTxt(ctx, content)
}

func (s *PostIndexService) SaveRobotsTxt(ctx context.Context, content string) error {
	return s.fileServ.SaveRobotsTxt(ctx, content)
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
