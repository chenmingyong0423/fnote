package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/website_config/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type aiConfigRepositoryStub struct {
	repository.IWebsiteConfigRepository
	stored *domain.AIAssistantConfig
	writes int
	err    error
}

func (r *aiConfigRepositoryStub) GetAIAssistantConfig(context.Context) (*domain.AIAssistantConfig, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.stored == nil {
		return nil, mongo.ErrNoDocuments
	}
	return cloneAIConfig(*r.stored), nil
}

func (r *aiConfigRepositoryStub) UpsertAIAssistantConfig(_ context.Context, cfg domain.AIAssistantConfig) error {
	for _, p := range cfg.Providers {
		if p.ApiKey != "" {
			panic("plaintext reached the repository")
		}
	}
	r.stored = cloneAIConfig(cfg)
	r.writes++
	return nil
}

func cloneAIConfig(cfg domain.AIAssistantConfig) *domain.AIAssistantConfig {
	data, err := bson.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	var result domain.AIAssistantConfig
	if err = bson.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return &result
}

func validAIConfig() domain.AIAssistantConfig {
	return domain.AIAssistantConfig{
		Enabled: true, AssistantName: "博客助手", ActiveProviderId: "p1",
		EnabledTools: []string{"search_posts", "read_post"},
		Providers: []domain.AIProvider{{
			Id: "p1", Name: "主要服务", ApiType: "openai_responses",
			BaseUrl: "https://example.com/v1/", Model: "configured-model", ApiKey: "test-secret-1",
		}},
	}
}

func newAIConfigServiceForTest(repo repository.IWebsiteConfigRepository, cipher *AIKeyCipher) *WebsiteConfigService {
	return NewWebsiteConfigService(repo, cipher, NewAIProviderTester(http.DefaultClient))
}

func TestAIConfigDefaultsAndKeyLifecycle(t *testing.T) {
	ctx := context.Background()
	repo := &aiConfigRepositoryStub{}
	svc := newAIConfigServiceForTest(repo, testAIKeyCipher())
	defaults, err := svc.GetAIAssistantConfig(ctx)
	if err != nil || defaults.Enabled || defaults.Providers == nil || len(defaults.EnabledTools) != 4 || repo.writes != 0 {
		t.Fatal("missing config must return disabled defaults without writing")
	}
	if err = svc.UpdateAIAssistantConfig(ctx, validAIConfig()); err != nil {
		t.Fatal(err)
	}
	if repo.stored.Providers[0].BaseUrl != "https://example.com/v1" || repo.stored.Providers[0].EncryptedApiKey == "test-secret-1" {
		t.Fatal("config must normalize the URL and encrypt credentials")
	}
	originalCiphertext := repo.stored.Providers[0].EncryptedApiKey
	cfg, _ := svc.GetAIAssistantConfig(ctx)
	cfg.AssistantName = "新名称"
	// A forged ciphertext must not replace the server-owned credential.
	cfg.Providers[0].EncryptedApiKey = "untrusted-ciphertext"
	if err = svc.UpdateAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	if repo.stored.Providers[0].EncryptedApiKey != originalCiphertext {
		t.Fatal("blank key must preserve the saved credential")
	}
	key, err := newAIConfigServiceForTest(repo, testAIKeyCipher()).GetAIProviderKey(ctx, "p1")
	if err != nil || key != "test-secret-1" {
		t.Fatal("saved key must remain usable after service recreation")
	}
	cfg, _ = svc.GetAIAssistantConfig(ctx)
	cfg.Providers[0].ApiKey = "test-secret-2"
	if err = svc.UpdateAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	if key, err = svc.GetAIProviderKey(ctx, "p1"); err != nil || key != "test-secret-2" {
		t.Fatal("new key must replace the saved credential")
	}
	resolved, err := svc.GetAIProvider(ctx, "p1")
	if err != nil || resolved.ApiKey != "test-secret-2" || resolved.BaseUrl != "https://example.com/v1" || resolved.Model != "configured-model" {
		t.Fatal("a resolved provider must pair connection settings with their own saved key")
	}
	cfg, _ = svc.GetAIAssistantConfig(ctx)
	cfg.Enabled = false
	cfg.ActiveProviderId = ""
	cfg.Providers = []domain.AIProvider{}
	if err = svc.UpdateAIAssistantConfig(ctx, cfg); err != nil || len(repo.stored.Providers) != 0 {
		t.Fatal("removing providers must remove their credentials")
	}
}

func TestAIConfigSwitchesProvidersAndAllowsKeylessDrafts(t *testing.T) {
	ctx := context.Background()
	repo := &aiConfigRepositoryStub{}
	svc := newAIConfigServiceForTest(repo, testAIKeyCipher())
	cfg := validAIConfig()
	p2 := cfg.Providers[0]
	p2.Id, p2.Name, p2.ApiType, p2.ApiKey = "p2", "备用服务", "openai_chat_completions", "test-secret-2"
	cfg.Providers = append(cfg.Providers, p2)
	cfg.ActiveProviderId = "p2"
	if err := svc.UpdateAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal(err)
	}
	if key, err := svc.GetAIProviderKey(ctx, "p2"); err != nil || key != "test-secret-2" {
		t.Fatal("active provider must retain its own credential")
	}
	cfg.Enabled = false
	cfg.ActiveProviderId = ""
	cfg.Providers = cfg.Providers[:1]
	cfg.Providers[0].Id, cfg.Providers[0].ApiKey = "draft", ""
	if err := svc.UpdateAIAssistantConfig(ctx, cfg); err != nil {
		t.Fatal("an inactive provider may be saved without a key")
	}
	cfg.Enabled, cfg.ActiveProviderId = true, "draft"
	if err := svc.UpdateAIAssistantConfig(ctx, cfg); !errors.Is(err, domain.ErrInvalidAIConfig) {
		t.Fatal("keyless provider must not be enabled")
	}
}

func TestAIConfigRejectsInvalidChangesWithoutWriting(t *testing.T) {
	cases := map[string]func(*domain.AIAssistantConfig){
		"Claude disabled":   func(c *domain.AIAssistantConfig) { c.Providers[0].ApiType = "anthropic_messages" },
		"GPT thinking mode": func(c *domain.AIAssistantConfig) { c.Providers[0].ThinkingMode = "enabled" },
		"DeepSeek effort": func(c *domain.AIAssistantConfig) {
			c.Providers[0].ApiType = "openai_chat_completions"
			c.Providers[0].ReasoningEffort = "medium"
		},
		"DeepSeek conflicting reasoning": func(c *domain.AIAssistantConfig) {
			c.Providers[0].ApiType = "openai_chat_completions"
			c.Providers[0].ThinkingMode = "disabled"
			c.Providers[0].ReasoningEffort = "high"
		},
		"empty name":              func(c *domain.AIAssistantConfig) { c.AssistantName = " " },
		"duplicate provider":      func(c *domain.AIAssistantConfig) { c.Providers = append(c.Providers, c.Providers[0]) },
		"invalid id":              func(c *domain.AIAssistantConfig) { c.Providers[0].Id = "$invalid" },
		"unknown protocol":        func(c *domain.AIAssistantConfig) { c.Providers[0].ApiType = "unknown" },
		"unknown tool":            func(c *domain.AIAssistantConfig) { c.EnabledTools = []string{"run_script"} },
		"duplicate tool":          func(c *domain.AIAssistantConfig) { c.EnabledTools = []string{"read_post", "read_post"} },
		"credentials in URL":      func(c *domain.AIAssistantConfig) { c.Providers[0].BaseUrl = "https://user:secret@example.com" },
		"key in URL":              func(c *domain.AIAssistantConfig) { c.Providers[0].BaseUrl = "https://example.com?key=secret" },
		"invalid URL":             func(c *domain.AIAssistantConfig) { c.Providers[0].BaseUrl = "file:///etc/file" },
		"empty model":             func(c *domain.AIAssistantConfig) { c.Providers[0].Model = " " },
		"missing active provider": func(c *domain.AIAssistantConfig) { c.ActiveProviderId = "missing" },
		"no active provider":      func(c *domain.AIAssistantConfig) { c.ActiveProviderId = "" },
		"invalid key":             func(c *domain.AIAssistantConfig) { c.Providers[0].ApiKey = "secret\nheader" },
	}
	for name, change := range cases {
		t.Run(name, func(t *testing.T) {
			repo := &aiConfigRepositoryStub{}
			cfg := validAIConfig()
			change(&cfg)
			err := newAIConfigServiceForTest(repo, testAIKeyCipher()).UpdateAIAssistantConfig(context.Background(), cfg)
			if !errors.Is(err, domain.ErrInvalidAIConfig) || repo.writes != 0 {
				t.Fatalf("invalid config must fail without writing, got %v", err)
			}
		})
	}
}

func TestAIConfigFailsSafelyWhenDeploymentKeyIsMissing(t *testing.T) {
	repo := &aiConfigRepositoryStub{}
	svc := newAIConfigServiceForTest(repo, NewAIKeyCipher(""))
	if err := svc.UpdateAIAssistantConfig(context.Background(), validAIConfig()); !errors.Is(err, ErrAIKeyEncryption) || repo.writes != 0 {
		t.Fatal("credentials must not be persisted if deployment encryption is unavailable")
	}
	repo.err = errors.New("database unavailable")
	if _, err := svc.GetAIAssistantConfig(context.Background()); !errors.Is(err, repo.err) {
		t.Fatal("database failures must not be treated as empty defaults")
	}
}

func TestAIConfigConcurrentBlankSavesPreserveReplacementKey(t *testing.T) {
	ctx := context.Background()
	repo := &aiConfigRepositoryStub{}
	svc := newAIConfigServiceForTest(repo, testAIKeyCipher())
	if err := svc.UpdateAIAssistantConfig(ctx, validAIConfig()); err != nil {
		t.Fatal(err)
	}
	var workers sync.WaitGroup
	for i := 0; i < 10; i++ {
		workers.Add(1)
		go func(replace bool) {
			defer workers.Done()
			cfg := validAIConfig()
			cfg.Providers[0].ApiKey = ""
			if replace {
				cfg.Providers[0].ApiKey = "replacement-key"
			}
			if err := svc.UpdateAIAssistantConfig(ctx, cfg); err != nil {
				t.Error(err)
			}
		}(i == 5)
	}
	workers.Wait()
	if key, err := svc.GetAIProviderKey(ctx, "p1"); err != nil || key != "replacement-key" {
		t.Fatal("blank saves must not restore stale credentials over a replacement")
	}
}
