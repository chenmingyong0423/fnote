package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type IAIAssistantConfigService interface {
	GetAIAssistantConfig(ctx context.Context) (domain.AIAssistantConfig, error)
	UpdateAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error
	GetAIProviderKey(ctx context.Context, providerId string) (string, error)
	GetAIProvider(ctx context.Context, providerId string) (domain.AIProvider, error)
	TestAIProvider(ctx context.Context, provider domain.AIProvider) (domain.AIProviderTestResult, error)
}

var aiProviderIdPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)

var aiToolNames = map[string]bool{
	"search_posts": true, "read_post": true,
	"get_blog_overview": true, "get_author_profile": true,
}

func (s *WebsiteConfigService) GetAIAssistantConfig(ctx context.Context) (domain.AIAssistantConfig, error) {
	cfg, err := s.repo.GetAIAssistantConfig(ctx)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.AIAssistantConfig{
			AssistantName: "AI 助手小勇",
			SystemPrompt:  "你是博客 AI 助手。围绕博客公开内容回答问题，使用工具核实文章信息，并通过来源标识引用文章。资料不足时明确说明，不编造文章或博主信息。",
			EnabledTools:  []string{"search_posts", "read_post", "get_blog_overview", "get_author_profile"},
			Providers:     []domain.AIProvider{},
		}, nil
	}
	if err != nil {
		return domain.AIAssistantConfig{}, err
	}
	if cfg.Providers == nil {
		cfg.Providers = []domain.AIProvider{}
	}
	if cfg.EnabledTools == nil {
		cfg.EnabledTools = []string{}
	}
	return *cfg, nil
}

func (s *WebsiteConfigService) UpdateAIAssistantConfig(ctx context.Context, cfg domain.AIAssistantConfig) error {
	// Serialize read/merge/save so blank credentials always retain the latest saved key.
	s.aiConfigMutex.Lock()
	defer s.aiConfigMutex.Unlock()

	if err := normalizeAndValidateAIConfig(&cfg); err != nil {
		return err
	}
	previous, err := s.GetAIAssistantConfig(ctx)
	if err != nil {
		return err
	}
	keys := make(map[string]string, len(previous.Providers))
	for _, provider := range previous.Providers {
		keys[provider.Id] = provider.EncryptedApiKey
	}
	for i := range cfg.Providers {
		provider := &cfg.Providers[i]
		// Ciphertext is server-owned, regardless of the caller's domain value.
		provider.EncryptedApiKey = keys[provider.Id]
		if provider.ApiKey != "" {
			provider.EncryptedApiKey, err = s.aiKeyCipher.Encrypt(provider.Id, provider.ApiKey)
			if err != nil {
				return err
			}
		}
		provider.ApiKey = ""
		if cfg.Enabled && provider.Id == cfg.ActiveProviderId {
			if provider.EncryptedApiKey == "" {
				return invalidAIConfig("启用的模型服务需要配置 API Key")
			}
			if _, err = s.aiKeyCipher.Decrypt(provider.Id, provider.EncryptedApiKey); err != nil {
				return err
			}
		}
	}
	return s.repo.UpsertAIAssistantConfig(ctx, cfg)
}

func (s *WebsiteConfigService) GetAIProviderKey(ctx context.Context, providerId string) (string, error) {
	provider, err := s.GetAIProvider(ctx, providerId)
	return provider.ApiKey, err
}

// Return connection settings and their decrypted credential from the same read.
// Configuration edits must never pair a new key with an old provider address.
func (s *WebsiteConfigService) GetAIProvider(ctx context.Context, providerId string) (domain.AIProvider, error) {
	cfg, err := s.GetAIAssistantConfig(ctx)
	if err != nil {
		return domain.AIProvider{}, err
	}
	for _, provider := range cfg.Providers {
		if provider.Id == providerId && provider.EncryptedApiKey != "" {
			provider.ApiKey, err = s.aiKeyCipher.Decrypt(provider.Id, provider.EncryptedApiKey)
			if err != nil {
				return domain.AIProvider{}, err
			}
			return provider, nil
		}
	}
	return domain.AIProvider{}, invalidAIConfig("模型服务未配置 API Key")
}

func normalizeAndValidateAIConfig(cfg *domain.AIAssistantConfig) error {
	cfg.AssistantName = strings.TrimSpace(cfg.AssistantName)
	cfg.ActiveProviderId = strings.TrimSpace(cfg.ActiveProviderId)
	if cfg.AssistantName == "" || utf8.RuneCountInString(cfg.AssistantName) > 80 {
		return invalidAIConfig("助手名称需要为 1 至 80 个字符")
	}
	if utf8.RuneCountInString(cfg.SystemPrompt) > 32000 {
		return invalidAIConfig("提示词不能超过 32000 个字符")
	}
	if len(cfg.Providers) > 20 {
		return invalidAIConfig("最多配置 20 个模型服务")
	}
	if cfg.EnabledTools == nil {
		cfg.EnabledTools = []string{}
	}
	if cfg.Providers == nil {
		cfg.Providers = []domain.AIProvider{}
	}
	seenTools := make(map[string]bool)
	for _, name := range cfg.EnabledTools {
		if !aiToolNames[name] || seenTools[name] {
			return invalidAIConfig("工具名称不支持或重复")
		}
		seenTools[name] = true
	}
	seenProviders := make(map[string]bool)
	for i := range cfg.Providers {
		p := &cfg.Providers[i]
		if err := normalizeAndValidateAIProvider(p); err != nil {
			return err
		}
		if !aiProviderIdPattern.MatchString(p.Id) || seenProviders[p.Id] {
			return invalidAIConfig("模型服务 ID 不合法或重复")
		}
		seenProviders[p.Id] = true
		if p.Name == "" || utf8.RuneCountInString(p.Name) > 80 {
			return invalidAIConfig("模型服务名称需要为 1 至 80 个字符")
		}
	}
	if cfg.ActiveProviderId != "" && !seenProviders[cfg.ActiveProviderId] {
		return invalidAIConfig("启用的模型服务不存在")
	}
	if cfg.Enabled && cfg.ActiveProviderId == "" {
		return invalidAIConfig("请先选择启用的模型服务")
	}
	return nil
}

// Connection tests use the same connection validation as saved providers, but
// do not require a display name or an ID for a new, unsaved credential.
func normalizeAndValidateAIProvider(p *domain.AIProvider) error {
	p.Id = strings.TrimSpace(p.Id)
	p.Name = strings.TrimSpace(p.Name)
	p.BaseUrl = strings.TrimRight(strings.TrimSpace(p.BaseUrl), "/")
	p.Model = strings.TrimSpace(p.Model)
	p.ApiKey = strings.TrimSpace(p.ApiKey)
	if p.Id != "" && !aiProviderIdPattern.MatchString(p.Id) {
		return invalidAIConfig("模型服务 ID 不合法")
	}
	switch p.ApiType {
	case "openai_responses", "openai_chat_completions":
	default:
		return invalidAIConfig("不支持的模型接口协议")
	}
	u, err := url.Parse(p.BaseUrl)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Hostname() == "" || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || len(p.BaseUrl) > 2048 {
		return invalidAIConfig("服务地址需要为 HTTP(S) 地址，不能包含账号、查询参数或片段")
	}
	if p.Model == "" || utf8.RuneCountInString(p.Model) > 200 {
		return invalidAIConfig("模型名称需要为 1 至 200 个字符")
	}
	if len(p.ApiKey) > 8192 || strings.ContainsAny(p.ApiKey, "\r\n") {
		return invalidAIConfig("API Key 格式不合法")
	}
	if p.MaxOutputTokens != nil && (*p.MaxOutputTokens < 16 || *p.MaxOutputTokens > 32768) {
		return invalidAIConfig("输出 token 上限需要在 16 至 32768 之间")
	}
	if p.ApiType == "openai_responses" {
		if p.ThinkingMode != "" {
			return invalidAIConfig("GPT 请使用推理强度参数，不能设置 DeepSeek 思考模式")
		}
		switch p.ReasoningEffort {
		case "", "none", "minimal", "low", "medium", "high", "xhigh", "max":
		default:
			return invalidAIConfig("不支持的 GPT 推理强度")
		}
	} else {
		switch p.ThinkingMode {
		case "", "enabled", "disabled":
		default:
			return invalidAIConfig("不支持的 DeepSeek 思考模式")
		}
		switch p.ReasoningEffort {
		case "", "none", "low", "high", "max":
		default:
			return invalidAIConfig("不支持的 DeepSeek 推理强度")
		}
		if (p.ThinkingMode == "disabled" && p.ReasoningEffort != "" && p.ReasoningEffort != "none") || (p.ThinkingMode == "enabled" && p.ReasoningEffort == "none") {
			return invalidAIConfig("DeepSeek 思考模式与推理强度冲突")
		}
	}
	return nil
}

func invalidAIConfig(message string) error {
	return fmt.Errorf("%w：%s", domain.ErrInvalidAIConfig, message)
}
