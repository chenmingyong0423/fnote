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

package repository

import (
	"context"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/repository/dao"
	"github.com/pkg/errors"
)

type IWebsiteConfigRepository interface {
	FindByTyp(ctx context.Context, typ string) (any, error)
	Increase(ctx context.Context, field string) error
	FindConfigByTypes(ctx context.Context, types ...string) ([]domain.Config, error)
	Decrease(ctx context.Context, field string) error
	UpdateWebSiteConfig(ctx context.Context, webSiteConfig domain.WebSiteConfig) error
	UpdateOwnerConfig(ctx context.Context, ownerConfig domain.OwnerConfig) error
	UpdateSeoMetaConfig(ctx context.Context, cfg *domain.SeoMetaConfig) error
	UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error
}

func NewWebsiteConfigRepository(dao dao.IWebsiteConfigDao) *WebsiteConfigRepository {
	return &WebsiteConfigRepository{
		dao: dao,
	}
}

var _ IWebsiteConfigRepository = (*WebsiteConfigRepository)(nil)

type WebsiteConfigRepository struct {
	dao dao.IWebsiteConfigDao
}

func (r *WebsiteConfigRepository) UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "comment"),
		update.BsonBuilder().SetSimple("props.enable_comment", commentConfig.EnableComment).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateSeoMetaConfig(ctx context.Context, cfg *domain.SeoMetaConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "seo meta"),
		update.Set(map[string]any{
			"props.title":                   cfg.Title,
			"props.description":             cfg.Description,
			"props.og_title":                cfg.OgTitle,
			"props.og_image":                cfg.OgImage,
			"props.baidu_site_verification": cfg.BaiduSiteVerification,
			"props.keywords":                cfg.Keywords,
			"props.author":                  cfg.Author,
			"props.robots":                  cfg.Robots,
			"update_time":                   time.Now().Unix(),
		}),
	)
}

func (r *WebsiteConfigRepository) UpdateOwnerConfig(ctx context.Context, ownerConfig domain.OwnerConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "owner"),
		update.BsonBuilder().SetSimple("props.name", ownerConfig.Name).SetSimple("props.profile", ownerConfig.Profile).SetSimple("props.picture", ownerConfig.Picture).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateWebSiteConfig(ctx context.Context, webSiteConfig domain.WebSiteConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().SetSimple("props.name", webSiteConfig.Name).SetSimple("props.icon", webSiteConfig.Icon).SetSimple("props.live_time", webSiteConfig.LiveTime).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) Decrease(ctx context.Context, field string) error {
	return r.dao.Decrease(ctx, field)
}

func (r *WebsiteConfigRepository) FindConfigByTypes(ctx context.Context, types ...string) ([]domain.Config, error) {
	configs, err := r.dao.GetByTypes(ctx, types...)
	if err != nil {
		return nil, err
	}
	return r.toConfigs(configs), nil
}

func (r *WebsiteConfigRepository) Increase(ctx context.Context, field string) error {
	return r.dao.Increase(ctx, field)
}

func (r *WebsiteConfigRepository) FindByTyp(ctx context.Context, typ string) (any, error) {
	config, err := r.dao.FindByTyp(ctx, typ)
	if err != nil {
		return nil, errors.WithMessage(err, "r.dao.FindByTyp failed")
	}
	return config.Props, nil
}

func (r *WebsiteConfigRepository) toConfigs(configs []*dao.WebsiteConfig) []domain.Config {
	result := make([]domain.Config, len(configs))
	for i, config := range configs {
		result[i] = domain.Config{Typ: config.Typ, Props: config.Props}
	}
	return result
}
