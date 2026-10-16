package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/aimodel"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/cloudwego/eino/schema"
)

type AIProviderTester struct {
	client   *http.Client
	capacity chan struct{}
}

func NewAIProviderTester(client *http.Client) *AIProviderTester {
	boundedClient := *client
	boundedClient.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &AIProviderTester{client: &boundedClient, capacity: make(chan struct{}, 2)}
}

func (s *WebsiteConfigService) TestAIProvider(ctx context.Context, provider domain.AIProvider) (domain.AIProviderTestResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := normalizeAndValidateAIProvider(&provider); err != nil {
		return domain.AIProviderTestResult{}, err
	}
	if provider.ApiKey == "" {
		if provider.Id == "" {
			return domain.AIProviderTestResult{}, invalidAIConfig("请填写 API Key，或使用已保存的模型服务")
		}
		key, err := s.GetAIProviderKey(ctx, provider.Id)
		if err != nil {
			return domain.AIProviderTestResult{}, err
		}
		provider.ApiKey = key
	}
	return s.aiProviderTester.Test(ctx, provider), nil
}

func (t *AIProviderTester) Test(ctx context.Context, provider domain.AIProvider) (result domain.AIProviderTestResult) {
	started := time.Now()
	defer func() { result.DurationMs = time.Since(started).Milliseconds() }()
	select {
	case t.capacity <- struct{}{}:
		defer func() { <-t.capacity }()
	default:
		result.Message = "已有连接测试正在进行，请稍后重试"
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	transport := t.client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	observed := &aiTestTransport{base: transport}
	client := *t.client
	client.Transport = observed
	maxTokens := provider.MaxOutputTokens
	if maxTokens == nil {
		defaultBudget := 1024
		maxTokens = &defaultBudget
	}
	chat, err := aimodel.New(ctx, aimodel.Config{
		APIType: provider.ApiType, BaseURL: provider.BaseUrl, Model: provider.Model, APIKey: provider.ApiKey,
		MaxOutputTokens: maxTokens, ReasoningEffort: provider.ReasoningEffort, ThinkingMode: provider.ThinkingMode,
	}, &client)
	if err != nil {
		result.Message = "模型初始化失败，请检查接口协议和参数"
		return
	}
	message, err := chat.Generate(ctx, []*schema.AgenticMessage{schema.UserAgenticMessage("请只回复 OK。")})
	if err != nil {
		var networkError net.Error
		switch {
		case errors.Is(ctx.Err(), context.Canceled):
			result.Message = "测试已取消"
		case errors.Is(ctx.Err(), context.DeadlineExceeded) || (errors.As(err, &networkError) && networkError.Timeout()):
			result.Message = "测试超时，请检查服务地址和模型响应速度"
		case observed.status.Load() >= 300:
			result.Message = aiProviderHTTPError(int(observed.status.Load()))
		case observed.status.Load() >= 200:
			result.Message = "模型服务回复格式不正确或内容过大，请检查接口协议"
		default:
			result.Message = "无法连接模型服务，请检查服务地址、网络或 HTTPS 证书"
		}
		// SDK errors can include response bodies and credentials; never return them.
		return
	}
	var texts []string
	if message != nil && message.Role == schema.AgenticRoleTypeAssistant {
		for _, block := range message.ContentBlocks {
			if block != nil && block.AssistantGenText != nil {
				texts = append(texts, block.AssistantGenText.Text)
			}
		}
		if meta := message.ResponseMeta; meta != nil && meta.OpenAIExtension != nil {
			if status := meta.OpenAIExtension.Status; status != "" && status != "completed" {
				texts = nil
			}
		}
	}
	reply := strings.TrimSpace(strings.Join(texts, "\n"))
	if reply == "" {
		result.Message = "模型未返回有效的对话文本，请检查模型、推理强度或输出 token 上限"
		return
	}
	reply = strings.ReplaceAll(reply, provider.ApiKey, "[密钥已隐藏]")
	runes := []rune(reply)
	if len(runes) > 500 {
		reply = string(runes[:500]) + "…"
	}
	result.Success, result.Message, result.Reply = true, "对话测试成功", reply
	return
}

// Bound the SDK response and retain only the HTTP status for safe diagnostics.
// This transport is used only by the non-streaming connection test.
type aiTestTransport struct {
	base   http.RoundTripper
	status atomic.Int32
}

func (t *aiTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	response, err := t.base.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	t.status.Store(int32(response.StatusCode))
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > 1<<20 {
		return nil, errors.New("model test response exceeds limit")
	}
	response.Body = io.NopCloser(bytes.NewReader(data))
	return response, nil
}

func aiProviderHTTPError(status int) string {
	message := "模型服务拒绝请求，请检查接口协议和模型参数"
	switch {
	case status >= 300 && status < 400:
		message = "服务返回重定向，请填写最终的 API 服务地址"
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		message = "认证失败，请检查 API Key 和模型访问权限"
	case status == http.StatusNotFound:
		message = "接口或模型不存在，请检查服务地址、接口协议和模型名称"
	case status == http.StatusPaymentRequired:
		message = "模型服务余额不足，请检查账户额度"
	case status == http.StatusTooManyRequests:
		message = "模型服务限流或额度不足，请稍后重试并检查账户额度"
	case status >= 500:
		message = "模型服务暂时不可用，请稍后重试"
	}
	return fmt.Sprintf("%s（HTTP %d）", message, status)
}
