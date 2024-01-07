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
		update.BsonBuilder().Pull("props.social_info_list", bsonx.M("id", id)).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateSocialInfo(ctx context.Context, socialInfo domain.SocialInfo) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.BsonBuilder().Eq("typ", "social").ElemMatch("props.social_info_list", bsonx.M("id", socialInfo.Id)).Build(),
		update.BsonBuilder().
			Set("props.social_info_list.$.social_name", socialInfo.SocialName).
			Set("props.social_info_list.$.social_value", socialInfo.SocialValue).
			Set("props.social_info_list.$.css_class", socialInfo.CssClass).
			Set("props.social_info_list.$.is_link", socialInfo.IsLink).
			Set("update_time", time.Now().Unix()).Build(),
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
		update.BsonBuilder().Push("props.social_info_list", socialInfo).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Pull("props.list", payInfoConfigElem).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) PushPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Push("props.list", payInfoConfigElem).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeleteRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Pull("props.records", record).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) AddRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Push("props.records", record).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "front-post-count"),
		update.BsonBuilder().Set("props.count", cfg.Count).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().Set("props.enabled", enabled).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().Set("props.title", noticeCfg.Title).Set("props.content", noticeCfg.Content).Set("props.publish_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateEmailConfig(ctx context.Context, emailConfig *domain.EmailConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "email"),
		update.BsonBuilder().
			Set("props.host", emailConfig.Host).
			Set("props.port", emailConfig.Port).
			Set("props.username", emailConfig.Username).
			Set("props.password", emailConfig.Password).
			Set("props.email", emailConfig.Email).
			Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateFriendConfig(ctx context.Context, friendConfig domain.FriendConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "friend"),
		update.BsonBuilder().Set("props.enable_friend_commit", friendConfig.EnableFriendCommit).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "comment"),
		update.BsonBuilder().Set("props.enable_comment", commentConfig.EnableComment).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateSeoMetaConfig(ctx context.Context, cfg *domain.SeoMetaConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "seo meta"),
		update.BsonBuilder().
			Set("props.title", cfg.Title).
			Set("props.description", cfg.Description).
			Set("props.og_title", cfg.OgTitle).
			Set("props.og_image", cfg.OgImage).
			Set("props.baidu_site_verification", cfg.BaiduSiteVerification).
			Set("props.keywords", cfg.Keywords).
			Set("props.author", cfg.Author).
			Set("props.robots", cfg.Robots).
			Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateOwnerConfig(ctx context.Context, ownerConfig domain.OwnerConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "owner"),
		update.BsonBuilder().Set("props.name", ownerConfig.Name).Set("props.profile", ownerConfig.Profile).Set("props.picture", ownerConfig.Picture).Set("update_time", time.Now().Unix()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateWebSiteConfig(ctx context.Context, webSiteConfig domain.WebSiteConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Set("props.name", webSiteConfig.Name).Set("props.icon", webSiteConfig.Icon).Set("props.live_time", webSiteConfig.LiveTime).Set("update_time", time.Now().Unix()).Build(),
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
