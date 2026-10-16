package web

import (
	"errors"
	"net/http"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"
	"github.com/gin-gonic/gin"
)

func (h *WebsiteConfigHandler) AdminGetAIAssistantConfig(ctx *gin.Context) (*apiwrap.ResponseBody[AIAssistantConfigVO], error) {
	cfg, err := h.serv.GetAIAssistantConfig(ctx.Request.Context())
	if err != nil {
		return nil, err
	}
	providers := make([]AIProviderVO, 0, len(cfg.Providers))
	for _, p := range cfg.Providers {
		providers = append(providers, AIProviderVO{
			Id: p.Id, Name: p.Name, ApiType: p.ApiType,
			BaseUrl: p.BaseUrl, Model: p.Model, HasApiKey: p.EncryptedApiKey != "",
			MaxOutputTokens: p.MaxOutputTokens, ReasoningEffort: p.ReasoningEffort, ThinkingMode: p.ThinkingMode,
		})
	}
	return apiwrap.SuccessResponseWithData(AIAssistantConfigVO{
		Enabled: cfg.Enabled, AssistantName: cfg.AssistantName,
		SystemPrompt: cfg.SystemPrompt, EnabledTools: cfg.EnabledTools,
		ActiveProviderId: cfg.ActiveProviderId, Providers: providers,
	}), nil
}

func (h *WebsiteConfigHandler) AdminUpdateAIAssistantConfig(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 256<<10)
	var req UpdateAIAssistantConfigReq
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Providers == nil {
		// Binding errors may contain submitted values. Never echo credentials.
		apiwrap.ErrorHandler(ctx, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "AI 配置请求格式不正确，请提交完整配置"))
		return
	}
	providers := make([]domain.AIProvider, 0, len(req.Providers))
	for _, p := range req.Providers {
		providers = append(providers, toAIProvider(p))
	}
	err := h.serv.UpdateAIAssistantConfig(ctx.Request.Context(), domain.AIAssistantConfig{
		Enabled: *req.Enabled, AssistantName: req.AssistantName,
		SystemPrompt: req.SystemPrompt, EnabledTools: req.EnabledTools,
		ActiveProviderId: req.ActiveProviderId, Providers: providers,
	})
	if err != nil {
		apiwrap.ErrorHandler(ctx, aiConfigHTTPError(err))
		return
	}
	ctx.JSON(http.StatusOK, apiwrap.SuccessResponse())
}

func (h *WebsiteConfigHandler) AdminTestAIProvider(ctx *gin.Context) {
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, 32<<10)
	var req UpdateAIProviderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		apiwrap.ErrorHandler(ctx, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "模型测试请求格式不正确"))
		return
	}
	result, err := h.serv.TestAIProvider(ctx.Request.Context(), toAIProvider(req))
	if err != nil {
		apiwrap.ErrorHandler(ctx, aiConfigHTTPError(err))
		return
	}
	ctx.JSON(http.StatusOK, apiwrap.SuccessResponseWithData(AIProviderTestResultVO{
		Success: result.Success, Message: result.Message, Reply: result.Reply, DurationMs: result.DurationMs,
	}))
}

func toAIProvider(req UpdateAIProviderReq) domain.AIProvider {
	return domain.AIProvider{
		Id: req.Id, Name: req.Name, ApiType: req.ApiType, BaseUrl: req.BaseUrl,
		Model: req.Model, ApiKey: req.ApiKey,
		MaxOutputTokens: req.MaxOutputTokens, ReasoningEffort: req.ReasoningEffort, ThinkingMode: req.ThinkingMode,
	}
}

func aiConfigHTTPError(err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidAIConfig):
		return apiwrap.NewErrorResponseBody(http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrAIKeyEncryption):
		return apiwrap.NewErrorResponseBody(http.StatusServiceUnavailable, err.Error())
	default:
		return err
	}
}
