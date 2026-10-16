package web_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/ioc"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/jwtutil"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/web"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type aiHTTPRepository struct {
	repository.IWebsiteConfigRepository
	config *domain.AIAssistantConfig
}

func (r *aiHTTPRepository) GetAIAssistantConfig(context.Context) (*domain.AIAssistantConfig, error) {
	if r.config == nil {
		return nil, mongo.ErrNoDocuments
	}
	data, _ := bson.Marshal(r.config)
	var cfg domain.AIAssistantConfig
	if err := bson.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (r *aiHTTPRepository) UpsertAIAssistantConfig(_ context.Context, cfg domain.AIAssistantConfig) error {
	r.config = &cfg
	return nil
}

func newAIConfigHTTP(t *testing.T) (*gin.Engine, *aiHTTPRepository, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	repo := &aiHTTPRepository{}
	key := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{42}, 32))
	svc := service.NewWebsiteConfigService(repo, service.NewAIKeyCipher(key), service.NewAIProviderTester(http.DefaultClient))
	router := gin.New()
	router.Use(ioc.JwtParseMiddleware(func() bool { return true }))
	web.NewWebsiteConfigHandler(svc).RegisterGinRoutes(router)
	token, _, err := jwtutil.GenerateJwt()
	if err != nil {
		t.Fatal(err)
	}
	return router, repo, token
}

func aiConfigRequest(router http.Handler, method, path, token string, body []byte) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", token)
	}
	response := httptest.NewRecorder()
	router.ServeHTTP(response, req)
	return response
}

func aiConfigPayload(key string) []byte {
	data, _ := json.Marshal(map[string]any{
		"enabled": true, "assistant_name": "博客助手", "system_prompt": "只引用公开文章",
		"enabled_tools": []string{"search_posts"}, "active_provider_id": "p1",
		"providers": []map[string]any{{
			"id": "p1", "name": "主要服务", "api_type": "openai_responses",
			"base_url": "https://example.com/v1", "model": "configured-model", "api_key": key,
		}},
	})
	return data
}

func TestAIConfigHTTPAuthAndCredentialRoundTrip(t *testing.T) {
	router, repo, token := newAIConfigHTTP(t)
	path := "/admin-api/configs/ai-assistant"
	for _, method := range []string{http.MethodGet, http.MethodPut} {
		response := aiConfigRequest(router, method, path, "", aiConfigPayload("test-secret"))
		if response.Code != http.StatusUnauthorized {
			t.Fatal("AI configuration must require the existing admin JWT")
		}
	}
	response := aiConfigRequest(router, http.MethodGet, path, token, nil)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"providers":[]`) {
		t.Fatal("initial read must return a usable default configuration")
	}
	response = aiConfigRequest(router, http.MethodPut, path, token, aiConfigPayload("test-secret"))
	if response.Code != http.StatusOK || repo.config == nil {
		t.Fatal("authenticated save must persist configuration")
	}
	encrypted := repo.config.Providers[0].EncryptedApiKey
	response = aiConfigRequest(router, http.MethodGet, path, token, nil)
	for _, forbidden := range []string{"test-secret", encrypted, `"encrypted_api_key"`, `"api_key":`} {
		if strings.Contains(response.Body.String(), forbidden) {
			t.Fatal("admin reads must not disclose plaintext or ciphertext")
		}
	}
	if !strings.Contains(response.Body.String(), `"has_api_key":true`) {
		t.Fatal("admin reads should indicate that a credential is configured")
	}
	response = aiConfigRequest(router, http.MethodPut, path, token, aiConfigPayload(""))
	if response.Code != http.StatusOK || repo.config.Providers[0].EncryptedApiKey != encrypted {
		t.Fatal("saving an unchanged credential must preserve it")
	}
	response = aiConfigRequest(router, http.MethodGet, "/configs/ai-assistant", "", nil)
	if response.Code != http.StatusNotFound {
		t.Fatal("there must be no public AI configuration endpoint")
	}
}

func TestAIConfigHTTPRejectsMalformedAndOversizedRequests(t *testing.T) {
	router, repo, token := newAIConfigHTTP(t)
	for _, body := range [][]byte{
		[]byte(`{"assistant_name":"missing enabled","providers":[]}`),
		[]byte(`{"enabled":false,"assistant_name":"missing providers"}`),
		[]byte(`{"enabled":false,"assistant_name":"invalid","providers":{},"api_key":"never-echo-this"}`),
		[]byte(`{"enabled":false,"assistant_name":"oversized","providers":[],"system_prompt":"` + strings.Repeat("x", 300000) + `"}`),
	} {
		response := aiConfigRequest(router, http.MethodPut, "/admin-api/configs/ai-assistant", token, body)
		if response.Code != http.StatusBadRequest || repo.config != nil {
			t.Fatal("malformed or oversized configuration must fail before persistence")
		}
		if strings.Contains(response.Body.String(), "never-echo-this") {
			t.Fatal("binding errors must not echo submitted credentials")
		}
	}
}

func TestAIProviderHTTPTestRequiresAdminAndDoesNotSave(t *testing.T) {
	router, repo, token := newAIConfigHTTP(t)
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" || r.Header.Get("Authorization") != "Bearer ephemeral-key" {
			t.Error("handler must send the entered connection parameters to Eino")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"resp-test","object":"response","status":"completed","output":[{"type":"message","role":"assistant","content":[{"type":"output_text","text":"OK"}]}]}`))
	}))
	defer upstream.Close()
	body, _ := json.Marshal(map[string]any{
		"api_type": "openai_responses", "base_url": upstream.URL + "/v1", "model": "test-model", "api_key": "ephemeral-key",
	})
	path := "/admin-api/configs/ai-assistant/test-provider"
	if response := aiConfigRequest(router, http.MethodPost, path, "", body); response.Code != http.StatusUnauthorized {
		t.Fatal("testing a provider must require admin authentication")
	}
	response := aiConfigRequest(router, http.MethodPost, path, token, body)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), `"success":true`) || !strings.Contains(response.Body.String(), `"reply":"OK"`) || repo.config != nil {
		t.Fatal("test endpoint must return a chat result without saving configuration")
	}
	if strings.Contains(response.Body.String(), "ephemeral-key") {
		t.Fatal("test endpoint must not disclose credentials")
	}
	for _, invalid := range []string{
		`{"api_type":"anthropic_messages","base_url":"https://example.com","model":"unused","api_key":"never-echo"}`,
		`{"api_type":"openai_responses","base_url":"https://example.com","model":"unused"}`,
		`{"api_type":{ "api_key":"never-echo"}}`,
	} {
		response = aiConfigRequest(router, http.MethodPost, path, token, []byte(invalid))
		if response.Code != http.StatusBadRequest || strings.Contains(response.Body.String(), "never-echo") {
			t.Fatal("invalid test parameters must be rejected without echoing credentials")
		}
	}
}
