package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
)

const gptTestReply = `{"id":"resp-test","object":"response","created_at":1,"status":"completed","model":"test-model","output":[{"id":"msg-test","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"OK","annotations":[]}]}],"usage":{"input_tokens":2,"output_tokens":1,"total_tokens":3}}`
const deepSeekTestReply = `{"id":"chat-test","object":"chat.completion","created":1,"model":"test-model","choices":[{"index":0,"message":{"role":"assistant","content":"OK","reasoning_content":"private reasoning"},"finish_reason":"stop"}],"usage":{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}}`

func TestAIProviderUsesEinoAndProviderSpecificParameters(t *testing.T) {
	for _, apiType := range []string{"openai_responses", "openai_chat_completions"} {
		t.Run(apiType, func(t *testing.T) {
			var calls atomic.Int32
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				if r.Method != http.MethodPost || r.Header.Get("Authorization") != "Bearer unsaved-test-key" {
					t.Error("test must make an authenticated model POST")
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Error(err)
				}
				if body["model"] != "test-model" || body["temperature"] != nil || body["top_p"] != nil {
					t.Error("must use entered model without incompatible sampling parameters")
				}
				w.Header().Set("Content-Type", "application/json")
				if apiType == "openai_responses" {
					if r.URL.Path != "/v1/responses" || body["max_output_tokens"] != float64(512) || body["max_tokens"] != nil || body["store"] != false || body["input"] == nil {
						t.Error("GPT must use Responses and max_output_tokens")
					}
					reasoning, _ := body["reasoning"].(map[string]any)
					if reasoning["effort"] != "high" || body["thinking"] != nil {
						t.Error("incorrect GPT reasoning parameters")
					}
					_, _ = w.Write([]byte(gptTestReply))
				} else {
					if r.URL.Path != "/v1/chat/completions" || body["max_tokens"] != float64(512) || body["max_output_tokens"] != nil || body["max_completion_tokens"] != nil || body["messages"] == nil {
						t.Error("DeepSeek must use Chat Completions and max_tokens")
					}
					thinking, _ := body["thinking"].(map[string]any)
					if thinking["type"] != "enabled" || body["reasoning_effort"] != "high" || body["reasoning"] != nil {
						t.Error("incorrect DeepSeek reasoning parameters")
					}
					_, _ = w.Write([]byte(deepSeekTestReply))
				}
			}))
			defer upstream.Close()
			repo := &aiConfigRepositoryStub{err: errors.New("database unavailable")}
			svc := NewWebsiteConfigService(repo, NewAIKeyCipher(""), NewAIProviderTester(upstream.Client()))
			budget := 512
			provider := domain.AIProvider{ApiType: apiType, BaseUrl: upstream.URL + "/v1/", Model: "test-model", ApiKey: "unsaved-test-key", MaxOutputTokens: &budget, ReasoningEffort: "high"}
			if apiType == "openai_chat_completions" {
				provider.ThinkingMode = "enabled"
			}
			result, err := svc.TestAIProvider(context.Background(), provider)
			if err != nil || !result.Success || result.Reply != "OK" || calls.Load() != 1 || repo.writes != 0 {
				t.Fatalf("unsaved parameters must produce one Eino chat without persistence: result=%+v, err=%v, calls=%d", result, err, calls.Load())
			}
		})
	}
}

func TestAIProviderBlankKeyUsesSavedCredentialWithoutSavingChanges(t *testing.T) {
	var auth string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(gptTestReply))
	}))
	defer upstream.Close()
	repo := &aiConfigRepositoryStub{}
	svc := NewWebsiteConfigService(repo, testAIKeyCipher(), NewAIProviderTester(upstream.Client()))
	cfg := validAIConfig()
	if err := svc.UpdateAIAssistantConfig(context.Background(), cfg); err != nil {
		t.Fatal(err)
	}
	provider := cfg.Providers[0]
	provider.ApiKey, provider.BaseUrl = "", upstream.URL+"/v1"
	result, err := svc.TestAIProvider(context.Background(), provider)
	if err != nil || !result.Success || auth != "Bearer test-secret-1" || repo.writes != 1 || repo.stored.Providers[0].BaseUrl != "https://example.com/v1" {
		t.Fatal("test must reuse saved credentials without saving edited settings")
	}
	provider.ApiKey = "replacement-for-test-only"
	result, err = svc.TestAIProvider(context.Background(), provider)
	if err != nil || !result.Success || auth != "Bearer replacement-for-test-only" {
		t.Fatal("entered credentials must take precedence")
	}
	key, _ := svc.GetAIProviderKey(context.Background(), provider.Id)
	if key != "test-secret-1" {
		t.Fatal("testing a replacement key must not save it")
	}
}

func TestAIProviderFailureResponsesAreSafe(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
	}{
		{"authentication", 401, `{"error":{"message":"secret-credential"}}`},
		{"missing model", 404, `{"error":{"message":"secret-credential"}}`},
		{"rate limit", 429, `{"error":{"message":"secret-credential"}}`},
		{"provider failure", 503, `{"error":{"message":"secret-credential"}}`},
		{"invalid JSON", 200, `<html>secret-credential</html>`},
		{"empty response", 200, `{"output":[]}`},
		{"API error with 200", 200, `{"error":{"message":"secret-credential"}}`},
		{"oversized response", 200, strings.Repeat("x", (1<<20)+1)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer upstream.Close()
			result := NewAIProviderTester(upstream.Client()).Test(context.Background(), domain.AIProvider{ApiType: "openai_responses", BaseUrl: upstream.URL + "/v1", Model: "test-model", ApiKey: "secret-credential"})
			if result.Success || result.Message == "" || result.Reply != "" || strings.Contains(result.Message, "secret-credential") || calls.Load() != 1 {
				t.Fatal("failed tests must not retry, disclose credentials, or report success")
			}
		})
	}
}

func TestAIProviderDoesNotFollowRedirects(t *testing.T) {
	var redirectedCalls atomic.Int32
	destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { redirectedCalls.Add(1) }))
	defer destination.Close()
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, destination.URL, http.StatusTemporaryRedirect)
	}))
	defer upstream.Close()
	result := NewAIProviderTester(upstream.Client()).Test(context.Background(), domain.AIProvider{ApiType: "openai_responses", BaseUrl: upstream.URL, Model: "test-model", ApiKey: "test-key"})
	if result.Success || redirectedCalls.Load() != 0 || !strings.Contains(result.Message, "重定向") {
		t.Fatal("must not forward credentials on redirects")
	}
}

func TestAIProviderCancellationStopsUpstream(t *testing.T) {
	started, stopped := make(chan struct{}), make(chan struct{})
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body any
		_ = json.NewDecoder(r.Body).Decode(&body)
		close(started)
		<-r.Context().Done()
		close(stopped)
	}))
	defer upstream.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan domain.AIProviderTestResult, 1)
	go func() {
		result <- NewAIProviderTester(upstream.Client()).Test(ctx, domain.AIProvider{ApiType: "openai_responses", BaseUrl: upstream.URL, Model: "test-model", ApiKey: "test-key"})
	}()
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("request did not start")
	}
	cancel()
	select {
	case got := <-result:
		if got.Success || !strings.Contains(got.Message, "取消") {
			t.Fatal("must report cancellation")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("model did not stop")
	}
	select {
	case <-stopped:
	case <-time.After(3 * time.Second):
		t.Fatal("upstream was not canceled")
	}
}

func TestAIProviderRedactsCredentialEcho(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(strings.Replace(gptTestReply, `"text":"OK"`, `"text":"test-key"`, 1)))
	}))
	defer upstream.Close()
	result := NewAIProviderTester(upstream.Client()).Test(context.Background(), domain.AIProvider{
		ApiType: "openai_responses", BaseUrl: upstream.URL, Model: "test-model", ApiKey: "test-key",
	})
	if !result.Success || strings.Contains(result.Reply, "test-key") {
		t.Fatal("even successful replies must not expose an echoed API key")
	}
}
