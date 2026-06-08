package service

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
)

type configHealthRule struct {
	key           string
	label         string
	level         string
	description   string
	href          string
	missingFields []string
}

type completionField struct {
	label string
	value string
}

func (s *WebsiteConfigService) GetConfigHealth(ctx context.Context) (*domain.ConfigHealth, error) {
	websiteConfig, err := s.GetWebSiteConfig(ctx)
	if err != nil {
		return nil, err
	}
	emailConfig, err := s.GetEmailConfig(ctx)
	if err != nil {
		return nil, err
	}
	seoConfig, err := s.GetSeoMetaConfig(ctx)
	if err != nil {
		return nil, err
	}
	noticeConfig, err := s.GetNoticeConfig(ctx)
	if err != nil {
		return nil, err
	}
	friendConfig, err := s.GetFriendConfig(ctx)
	if err != nil {
		return nil, err
	}
	frontPostCountConfig, err := s.GetFrontPostCountConfig(ctx)
	if err != nil {
		return nil, err
	}
	payConfig, err := s.GetPayConfig(ctx)
	if err != nil {
		return nil, err
	}
	socialConfig, err := s.GetSocialConfig(ctx)
	if err != nil {
		return nil, err
	}
	tpsvConfig, err := s.GetTPSVConfig(ctx)
	if err != nil {
		return nil, err
	}
	baiduPushConfig, err := s.GetBaiduPushConfig(ctx)
	if err != nil {
		return nil, err
	}
	carouselConfig, err := s.GetCarouselConfig(ctx)
	if err != nil {
		return nil, err
	}
	states, err := s.repo.GetConfigCheckStates(ctx)
	if err != nil {
		return nil, err
	}

	stateByKey := make(map[string]domain.ConfigCheckState, len(states))
	for _, state := range states {
		stateByKey[state.CheckKey] = state
	}

	rules := []configHealthRule{
		{
			key:         "website.basic",
			label:       "站点基础信息",
			level:       domain.ConfigHealthLevelRequired,
			description: "站点名称、图标和站长资料会展示在前台页面。",
			href:        "/home/setting?tab=basic",
			missingFields: missingStringFields([]completionField{
				{label: "站点名称", value: websiteConfig.WebsiteName},
				{label: "站点图标", value: websiteConfig.WebsiteIcon},
				{label: "站长昵称", value: websiteConfig.WebsiteOwner},
				{label: "站长简介", value: websiteConfig.WebsiteOwnerProfile},
				{label: "站长头像", value: websiteConfig.WebsiteOwnerAvatar},
			}),
		},
		{
			key:           "email.smtp",
			label:         "邮件配置",
			level:         domain.ConfigHealthLevelRecommended,
			description:   "用于评论、友链、审核等通知邮件。",
			href:          "/home/setting?tab=email",
			missingFields: missingEmailFields(emailConfig),
		},
		{
			key:         "seo.meta",
			label:       "SEO 配置",
			level:       domain.ConfigHealthLevelRecommended,
			description: "用于搜索引擎摘要和分享元信息。",
			href:        "/home/setting?tab=seo",
			missingFields: missingStringFields([]completionField{
				{label: "SEO 标题", value: seoConfig.Title},
				{label: "SEO 描述", value: seoConfig.Description},
				{label: "关键词", value: seoConfig.Keywords},
				{label: "作者", value: seoConfig.Author},
				{label: "robots", value: seoConfig.Robots},
			}),
		},
		{
			key:           "notice.content",
			label:         "公告配置",
			level:         domain.ConfigHealthLevelRecommended,
			description:   "开启公告时需要配置标题和内容。",
			href:          "/home/setting?tab=notice",
			missingFields: missingNoticeFields(noticeConfig),
		},
		{
			key:           "friend.introduction",
			label:         "友链说明",
			level:         domain.ConfigHealthLevelRecommended,
			description:   "开启友链提交时建议配置申请说明。",
			href:          "/home/setting?tab=friend",
			missingFields: missingFriendFields(friendConfig),
		},
		{
			key:           "post.front_count",
			label:         "首页文章数量",
			level:         domain.ConfigHealthLevelRecommended,
			description:   "控制首页列表展示数量。",
			href:          "/home/setting?tab=front-post-count",
			missingFields: missingPositiveInt("展示数量", frontPostCountConfig.Count),
		},
		{
			key:           "pay.qrcode",
			label:         "支付二维码",
			level:         domain.ConfigHealthLevelOptional,
			description:   "用于前台展示打赏二维码。",
			href:          "/home/setting?tab=pay",
			missingFields: missingList("支付二维码", len(payConfig.List)),
		},
		{
			key:           "social.links",
			label:         "社交信息",
			level:         domain.ConfigHealthLevelOptional,
			description:   "用于前台展示站长社交链接。",
			href:          "/home/setting?tab=social",
			missingFields: missingList("社交信息", len(socialConfig.SocialInfoList)),
		},
		{
			key:           "site.verification",
			label:         "站点验证",
			level:         domain.ConfigHealthLevelOptional,
			description:   "用于第三方站长平台验证站点归属。",
			href:          "/home/setting?tab=verification",
			missingFields: missingList("站点验证码", len(tpsvConfig.List)),
		},
		{
			key:         "post.baidu_push",
			label:       "文章推送",
			level:       domain.ConfigHealthLevelOptional,
			description: "用于向搜索引擎推送文章索引。",
			href:        "/home/setting?tab=push",
			missingFields: missingStringFields([]completionField{
				{label: "Baidu site", value: baiduPushConfig.Site},
				{label: "Baidu token", value: baiduPushConfig.Token},
			}),
		},
		{
			key:           "carousel.items",
			label:         "轮播图",
			level:         domain.ConfigHealthLevelOptional,
			description:   "用于首页轮播展示。",
			href:          "/home/setting?tab=carousel",
			missingFields: missingList("轮播图", len(carouselConfig.List)),
		},
	}

	items := make([]domain.ConfigHealthItem, 0, len(rules))
	now := time.Now().Local()
	for _, rule := range rules {
		item := domain.ConfigHealthItem{
			Key:           rule.key,
			Label:         rule.label,
			Level:         rule.level,
			Configured:    len(rule.missingFields) == 0,
			MissingFields: rule.missingFields,
			Description:   rule.description,
			Href:          rule.href,
		}
		item.Status = configHealthStatus(item.Configured, stateByKey[rule.key], now)
		if state, ok := stateByKey[rule.key]; ok {
			item.SnoozedUntil = state.SnoozedUntil
			item.IgnoredReason = state.IgnoredReason
		}
		items = append(items, item)
	}

	required := summarizeConfigHealthGroup(items, domain.ConfigHealthLevelRequired)
	recommended := summarizeConfigHealthGroup(items, domain.ConfigHealthLevelRecommended)
	optional := summarizeConfigHealthGroup(items, domain.ConfigHealthLevelOptional)

	return &domain.ConfigHealth{
		Score:       configHealthScore(required, recommended),
		Required:    required,
		Recommended: recommended,
		Optional:    optional,
		Items:       items,
	}, nil
}

func (s *WebsiteConfigService) UpdateConfigCheckState(ctx context.Context, state domain.ConfigCheckState) error {
	if state.CheckKey == "" {
		return errors.New("check_key is empty")
	}
	switch state.Status {
	case domain.ConfigCheckStateActive:
		state.IgnoredReason = ""
		state.SnoozedUntil = nil
	case domain.ConfigCheckStateIgnored:
		state.SnoozedUntil = nil
	case domain.ConfigCheckStateSnoozed:
		if state.SnoozedUntil == nil {
			return errors.New("snoozed_until is empty")
		}
	default:
		return errors.Errorf("unsupported config check state: %s", state.Status)
	}
	return s.repo.UpsertConfigCheckState(ctx, state)
}

func configHealthStatus(configured bool, state domain.ConfigCheckState, now time.Time) string {
	if configured {
		return domain.ConfigHealthStatusOK
	}
	if state.Status == domain.ConfigCheckStateIgnored {
		return domain.ConfigHealthStatusIgnored
	}
	if state.Status == domain.ConfigCheckStateSnoozed && state.SnoozedUntil != nil && state.SnoozedUntil.After(now) {
		return domain.ConfigHealthStatusSnoozed
	}
	return domain.ConfigHealthStatusMissing
}

func summarizeConfigHealthGroup(items []domain.ConfigHealthItem, level string) domain.ConfigHealthGroup {
	group := domain.ConfigHealthGroup{}
	for _, item := range items {
		if item.Level != level {
			continue
		}
		group.Total++
		if item.Configured {
			group.Done++
		}
	}
	group.Completed = group.Done == group.Total
	return group
}

func configHealthScore(required, recommended domain.ConfigHealthGroup) int {
	total := required.Total + recommended.Total
	if total == 0 {
		return 100
	}
	return (required.Done + recommended.Done) * 100 / total
}

func missingStringFields(fields []completionField) []string {
	missingFields := make([]string, 0, len(fields))
	for _, field := range fields {
		if strings.TrimSpace(field.value) == "" {
			missingFields = append(missingFields, field.label)
		}
	}
	return missingFields
}

func missingEmailFields(config *domain.EmailConfig) []string {
	missingFields := missingStringFields([]completionField{
		{label: "host", value: config.Host},
		{label: "username", value: config.Username},
		{label: "password", value: config.Password},
		{label: "email", value: config.Email},
	})
	if config.Port <= 0 {
		missingFields = append(missingFields, "port")
	}
	return missingFields
}

func missingNoticeFields(config domain.NoticeConfig) []string {
	if !config.Enabled {
		return nil
	}
	return missingStringFields([]completionField{
		{label: "公告标题", value: config.Title},
		{label: "公告内容", value: config.Content},
	})
}

func missingFriendFields(config domain.FriendConfig) []string {
	if !config.EnableFriendCommit {
		return nil
	}
	return missingStringFields([]completionField{
		{label: "友链说明", value: config.Introduction},
	})
}

func missingPositiveInt(label string, value int64) []string {
	if value > 0 {
		return nil
	}
	return []string{label}
}

func missingList(label string, count int) []string {
	if count > 0 {
		return nil
	}
	return []string{label}
}
