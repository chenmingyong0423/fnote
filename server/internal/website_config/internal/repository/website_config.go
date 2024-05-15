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
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository/dao"
	"github.com/chenmingyong0423/gkit/uuidx"

	"github.com/gin-gonic/gin"

	"github.com/chenmingyong0423/go-mongox/bsonx"

	"github.com/chenmingyong0423/go-mongox/builder/query"
	"github.com/chenmingyong0423/go-mongox/builder/update"

	"github.com/pkg/errors"
)

type IWebsiteConfigRepository interface {
	FindByTyp(ctx context.Context, typ string) (any, error)
	Increase(ctx context.Context, field string) error
	FindConfigByTypes(ctx context.Context, types ...string) ([]domain.Config, error)
	Decrease(ctx context.Context, field string) error
	UpdateSeoMetaConfig(ctx context.Context, cfg *domain.SeoMetaConfig) error
	UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error
	UpdateSwitch4FriendConfig(ctx context.Context, enableFriendCommit bool) error
	UpdateEmailConfig(ctx context.Context, emailConfig *domain.EmailConfig, now time.Time) error
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
	UpdateAdminConfig(ctx context.Context, adminConfig domain.AdminConfig, now time.Time) error
	UpdateWebSiteConfig(ctx context.Context, websiteConfig domain.WebsiteConfig, now time.Time) error
	GetTPSVConfig(ctx context.Context) (*domain.TPSVConfig, error)
	AddTPSVConfig(ctx context.Context, tpsv domain.TPSV) error
	DeleteTPSVConfigByKey(ctx context.Context, key string) error
	GetBaiduPushConfig(ctx context.Context) (*domain.Baidu, error)
	UpdatePushConfigByKey(ctx context.Context, key string, updates map[string]any) error
	AddCarouselConfig(ctx context.Context, carouselElem domain.CarouselElem) error
	UpdateCarouselShowStatus(ctx context.Context, id string, show bool) error
	UpdateCarouselElem(ctx context.Context, carouselElem domain.CarouselElem) error
	DeleteCarouselElem(ctx context.Context, id string) error
	FindCarouselById(ctx context.Context, id string) (*domain.Config, error)
	UpdateIntroduction4FriendConfig(ctx context.Context, introduction string) error
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

func (r *WebsiteConfigRepository) UpdateIntroduction4FriendConfig(ctx context.Context, introduction string) error {
	return r.dao.UpdateByConditionAndUpdates(ctx, query.Eq("typ", "friend"), update.BsonBuilder().Set("props.introduction", introduction).Set("updated_at", time.Now().Local()).Build())
}

func (r *WebsiteConfigRepository) FindCarouselById(ctx context.Context, id string) (*domain.Config, error) {
	websiteConfig, err := r.dao.FindByFilter(ctx, query.BsonBuilder().Eq("typ", "carousel").Eq("props.list.id", id).Build())
	if err != nil {
		return nil, err
	}
	return r.toConfig(websiteConfig), nil
}

func (r *WebsiteConfigRepository) DeleteCarouselElem(ctx context.Context, id string) error {
	return r.dao.UpdateByConditionAndUpdates(ctx, query.Eq("typ", "carousel"),
		update.Pull("props.list", bsonx.M("id", id)))
}

func (r *WebsiteConfigRepository) UpdateCarouselElem(ctx context.Context, carouselElem domain.CarouselElem) error {
	return r.dao.UpdateByConditionAndUpdates(ctx, query.BsonBuilder().Eq("typ", "carousel").Eq("props.list.id", carouselElem.Id).Build(),
		update.BsonBuilder().Set("props.list.$", carouselElem).Set("updated_at", carouselElem.UpdatedAt).Build())
}

func (r *WebsiteConfigRepository) UpdateCarouselShowStatus(ctx context.Context, id string, show bool) error {
	return r.dao.UpdateByConditionAndUpdates(ctx, query.BsonBuilder().Eq("typ", "carousel").Eq("props.list.id", id).Build(),
		update.BsonBuilder().Set("props.list.$.show", show).Set("updated_at", time.Now().Local()).Build())
}

func (r *WebsiteConfigRepository) AddCarouselConfig(ctx context.Context, carouselElem domain.CarouselElem) error {
	return r.dao.PushCarouselConfig(ctx, carouselElem)
}

func (r *WebsiteConfigRepository) UpdatePushConfigByKey(ctx context.Context, key string, updates map[string]any) error {
	return r.dao.UpdatePostIndexProps(ctx, update.Set(fmt.Sprintf("props.%s", key), updates))
}

func (r *WebsiteConfigRepository) GetBaiduPushConfig(ctx context.Context) (*domain.Baidu, error) {
	cfg, err := r.dao.FindByTyp(ctx, "post index")
	if err != nil {
		return nil, err
	}
	baiduCfg := &domain.BaiduPushConfig{}
	err = r.anyToStruct(cfg.Props, baiduCfg)
	if err != nil {
		return nil, err
	}
	return &baiduCfg.Baidu, nil
}

func (r *WebsiteConfigRepository) DeleteTPSVConfigByKey(ctx context.Context, key string) error {
	return r.dao.DeleteTPSVConfigByKey(ctx, key)
}

func (r *WebsiteConfigRepository) AddTPSVConfig(ctx context.Context, tpsv domain.TPSV) error {
	return r.dao.AddTPSVConfig(ctx, tpsv)
}

func (r *WebsiteConfigRepository) GetTPSVConfig(ctx context.Context) (*domain.TPSVConfig, error) {
	cfg, err := r.dao.FindByTyp(ctx, "third party site verification")
	if err != nil {
		return nil, err
	}
	tpsvCfg := &domain.TPSVConfig{List: make([]domain.TPSV, 0)}
	err = r.anyToStruct(cfg.Props, tpsvCfg)
	if err != nil {
		return nil, err
	}
	return tpsvCfg, nil
}

func (r *WebsiteConfigRepository) anyToStruct(props any, cfgInfos any) error {
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

func (r *WebsiteConfigRepository) UpdateWebSiteConfig(ctx context.Context, websiteConfig domain.WebsiteConfig, now time.Time) error {
	b := update.BsonBuilder().
		Set("props.website_name", websiteConfig.WebsiteName).
		Set("props.website_icon", websiteConfig.WebsiteIcon).
		Set("props.website_owner", websiteConfig.WebsiteOwner).
		Set("props.website_owner_profile", websiteConfig.WebsiteOwnerProfile).
		Set("props.website_owner_avatar", websiteConfig.WebsiteOwnerAvatar).
		Set("updated_at", now)
	if websiteConfig.WebsiteInit != nil {
		b.Set("props.website_init", websiteConfig.WebsiteInit)
	}
	if websiteConfig.WebsiteRuntime != nil {
		b.Set("props.website_runtime", websiteConfig.WebsiteRuntime)
	}
	return r.dao.UpdateByConditionAndUpdates(ctx, query.Eq("typ", "website"), b.Build())
}

func (r *WebsiteConfigRepository) UpdateAdminConfig(ctx context.Context, adminConfig domain.AdminConfig, now time.Time) error {
	return r.dao.UpdatePropsByTyp(ctx, "admin", adminConfig, now)
}

func (r *WebsiteConfigRepository) DeleteSocialInfo(ctx context.Context, id []byte) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "social"),
		update.BsonBuilder().Pull("props.social_info_list", bsonx.M("id", id)).Set("updated_at", time.Now().Local()).Build(),
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
			Set("updated_at", time.Now().Local()).Build(),
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
		update.BsonBuilder().Push("props.social_info_list", socialInfo).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeletePayInfo(ctx context.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Pull("props.list", payInfoConfigElem).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) PushPayInfo(ctx *gin.Context, payInfoConfigElem domain.PayInfoConfigElem) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "pay"),
		update.BsonBuilder().Push("props.list", payInfoConfigElem).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) DeleteRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Pull("props.website_records", record).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) AddRecordInWebsiteConfig(ctx context.Context, record string) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "website"),
		update.BsonBuilder().Push("props.website_records", record).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateFrontPostCountConfig(ctx context.Context, cfg domain.FrontPostCountConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "front-post-count"),
		update.BsonBuilder().Set("props.count", cfg.Count).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfigEnabled(ctx context.Context, enabled bool) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().Set("props.enabled", enabled).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateNoticeConfig(ctx context.Context, noticeCfg *domain.NoticeConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "notice"),
		update.BsonBuilder().Set("props.title", noticeCfg.Title).Set("props.content", noticeCfg.Content).Set("props.publish_time", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateEmailConfig(ctx context.Context, emailConfig *domain.EmailConfig, now time.Time) error {
	return r.dao.UpdatePropsByTyp(ctx, "email", emailConfig, now)
}

func (r *WebsiteConfigRepository) UpdateSwitch4FriendConfig(ctx context.Context, enableFriendCommit bool) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "friend"),
		update.BsonBuilder().Set("props.enable_friend_commit", enableFriendCommit).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateCommentConfig(ctx context.Context, commentConfig domain.CommentConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "comment"),
		update.BsonBuilder().Set("props.enable_comment", commentConfig.EnableComment).Set("updated_at", time.Now().Local()).Build(),
	)
}

func (r *WebsiteConfigRepository) UpdateSeoMetaConfig(ctx context.Context, cfg *domain.SeoMetaConfig) error {
	return r.dao.UpdateByConditionAndUpdates(
		ctx,
		query.Eq("typ", "seo meta"),
		update.BsonBuilder().
			Set("props.description", cfg.Description).
			Set("props.title", cfg.Title).
			Set("props.og_title", cfg.OgTitle).
			Set("props.og_image", cfg.OgImage).
			Set("props.baidu_site_verification", cfg.BaiduSiteVerification).
			Set("props.keywords", cfg.Keywords).
			Set("props.author", cfg.Author).
			Set("props.robots", cfg.Robots).
			Set("updated_at", time.Now().Local()).Build(),
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

func (r *WebsiteConfigRepository) toConfig(config *dao.WebsiteConfig) *domain.Config {
	return &domain.Config{Typ: config.Typ, Props: config.Props}
}
