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
	"encoding/hex"
	"time"

	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/go-mongox/bsonx"

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
	UpdateFriendConfig(ctx context.Context, friendConfig domain.FriendConfig) error
	UpdateEmailConfig(ctx context.Context, emailConfig *domain.EmailConfig) error
	UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error
	UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error
	UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error
	AddRecordInWebsiteConfig(ctx context.Context, record string) error
	DeleteRecordInWebsiteConfig(ctx context.Context, record string) error
	PushPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error
	DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error
	AddSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error
	UpdateSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error
	DeleteSocialInfo(ctx context.Context, id []byte) error
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

func (r *WebsiteConfigRepository) DeleteSocialInfo(ctx context.Context, id []byte) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "social"),
		update.BsonBuilder().Pull(bsonx.M("props.social_info_list", bsonx.M("id", id))).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.BsonBuilder().Eq("typ", "social").ElemMatch("props.social_info_list", bsonx.M("id", socialInfo.Id)).Build(),
		update.Set(map[string]any{
			"props.social_info_list.$.social_name":  socialInfo.SocialName,
			"props.social_info_list.$.social_value": socialInfo.SocialValue,
			"props.social_info_list.$.css_class":    socialInfo.CssClass,
			"props.social_info_list.$.is_link":      socialInfo.IsLink,
		}),
	)
}

func (r *WebsiteConfigRepository) AddSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error {
	id, err := hex.DecodeString(uuidx.RearrangeUUID4())
	if err != nil {
		return err
	}
	socialInfo.Id = id
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "social"),
		update.BsonBuilder().Push(bsonx.M("props.social_info_list", socialInfo)).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Pull(bsonx.M("props.list", payInfoConfigElem)).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) PushPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Push(bsonx.M("props.list", payInfoConfigElem)).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeleteRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Pull(bsonx.M("props.records", record)).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) AddRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Push(bsonx.M("props.records", record)).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "front-post-count"),
		update.BsonBuilder().SetSimple("props.count", cfg.Count).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().SetSimple("props.enabled", enabled).SetSimple("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().SetSimple("props.title", noticeCfg.Title).SetSimple("props.content", noticeCfg.Content).SetSimple("props.publish_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateEmailConfig(ctx context.Context, emailConfig *domain.EmailConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "email"),
		update.Set(map[string]any{
			"props.host":     emailConfig.Host,
			"props.port":     emailConfig.Port,
			"props.username": emailConfig.Username,
			"props.password": emailConfig.Password,
			"props.email":    emailConfig.Email,
			"update_time":    time.Now().Unix(),
		}),
	)
}

func (r *WebsiteConfigRepository) UpdateFriendConfig(ctx context.Context, friendConfig domain.FriendConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "friend"),
		update.BsonBuilder().SetSimple("props.enable_friend_commit", friendConfig.EnableFriendCommit).SetSimple("update_time", time.Now().Unix()).Build(),
	)
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
