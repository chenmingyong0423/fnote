package web_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/post"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// This fixture exercises real Gin handlers, ownership, service orchestration,
// Eino serialization/stream parsing, tool calls and terminal persistence.
type memoryStore struct {
	mu            sync.Mutex
	conversations map[string]domain.Conversation
	visitors      map[string]string
}

func clone(c domain.Conversation) domain.Conversation {
	b, _ := bson.Marshal(c)
	var result domain.Conversation
	_ = bson.Unmarshal(b, &result)
	return result
}
func (r *memoryStore) Identify(_ context.Context, id, ip string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.visitors[id] = ip
	return nil
}
func (r *memoryStore) Count(_ context.Context, v string) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var n int64
	for _, c := range r.conversations {
		if c.VisitorID == v {
			n++
		}
	}
	return n, nil
}
func (r *memoryStore) Create(_ context.Context, c domain.Conversation) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conversations[c.ID] = clone(c)
	return nil
}
func (r *memoryStore) List(_ context.Context, v string) ([]domain.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := []domain.Conversation{}
	for _, c := range r.conversations {
		if c.VisitorID == v {
			copy := clone(c)
			copy.Messages = nil
			out = append(out, copy)
		}
	}
	return out, nil
}
func (r *memoryStore) Get(_ context.Context, v, id string) (*domain.Conversation, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.conversations[id]
	if !ok || c.VisitorID != v {
		return nil, domain.ErrNotFound
	}
	copy := clone(c)
	return &copy, nil
}
func (r *memoryStore) Begin(_ context.Context, v, id string, user, answer domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.conversations[id]
	if !ok || c.VisitorID != v {
		return domain.ErrNotFound
	}
	if c.Generation != nil && c.Generation.Status == "running" {
		return domain.ErrBusy
	}
	for _, m := range c.Messages {
		if m.RequestID == user.RequestID {
			return domain.ErrBusy
		}
	}
	c.Messages = append(c.Messages, user, answer)
	c.Generation = &domain.Generation{RequestID: user.RequestID, MessageID: answer.ID, Status: "running"}
	r.conversations[id] = clone(c)
	return nil
}
func (r *memoryStore) SaveAnswer(ctx context.Context, v, id string, answer domain.Message, final bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	c, ok := r.conversations[id]
	if !ok || c.VisitorID != v {
		return domain.ErrNotFound
	}
	if c.Generation == nil || c.Generation.MessageID != answer.ID || c.Generation.Status != "running" {
		return domain.ErrBusy
	}
	for i := range c.Messages {
		if c.Messages[i].ID == answer.ID {
			c.Messages[i] = answer
		}
	}
	if final {
		c.Generation.Status = answer.Status
	}
	r.conversations[id] = clone(c)
	return nil
}

type configService struct {
	website_config.Service
	cfg website_config.AIAssistantConfig
}

func (c *configService) GetAIAssistantConfig(context.Context) (website_config.AIAssistantConfig, error) {
	return c.cfg, nil
}
func (c *configService) GetAIProviderKey(context.Context, string) (string, error) {
	return "provider-private-key", nil
}
func (c *configService) GetAIProvider(_ context.Context, id string) (website_config.AIProvider, error) {
	for _, p := range c.cfg.Providers {
		if p.Id == id {
			p.ApiKey = "provider-private-key"
			return p, nil
		}
	}
	return website_config.AIProvider{}, fmt.Errorf("provider not found")
}
func (c *configService) GetWebSiteConfig(context.Context) (*website_config.WebsiteConfig, error) {
	return &website_config.WebsiteConfig{WebsiteName: "公开博客", WebsiteOwner: "博主", WebsiteOwnerProfile: "公开简介"}, nil
}

type postService struct {
	post.Service
	private bool
}

func (p *postService) SearchPublicPosts(context.Context, string, int64) ([]*post.Post, int64, error) {
	return []*post.Post{{PrimaryPost: post.PrimaryPost{Id: "public-post", Title: "Go 入门", Summary: "公开摘要"}, ExtraPost: post.ExtraPost{IsDisplayed: true}}}, 1, nil
}
func (p *postService) GetPublicPostById(context.Context, string) (*post.Post, error) {
	return &post.Post{PrimaryPost: post.PrimaryPost{Id: "public-post", Title: "Go 入门"}, ExtraPost: post.ExtraPost{IsDisplayed: !p.private, Content: "文章正文"}}, nil
}

type fixture struct {
	server  *httptest.Server
	client  *http.Client
	store   *memoryStore
	configs *configService
	posts   *postService
}

func newFixture(t *testing.T, provider http.HandlerFunc) *fixture {
	t.Helper()
	gin.SetMode(gin.TestMode)
	upstream := httptest.NewServer(provider)
	t.Cleanup(upstream.Close)
	store := &memoryStore{conversations: map[string]domain.Conversation{}, visitors: map[string]string{}}
	configs := &configService{cfg: website_config.AIAssistantConfig{Enabled: true, AssistantName: "测试助手", SystemPrompt: "只回答博客问题", ActiveProviderId: "provider", EnabledTools: []string{"search_posts", "read_post", "get_blog_overview", "get_author_profile"}, Providers: []website_config.AIProvider{{Id: "provider", ApiType: "openai_chat_completions", BaseUrl: upstream.URL + "/v1", Model: "test-model", EncryptedApiKey: "never-expose-ciphertext"}}}}
	posts := &postService{}
	model := service.NewModelService(configs, service.NewToolService(posts, configs), upstream.Client())
	svc := service.NewAIAssistantService(store, configs, model)
	router := gin.New()
	web.NewAIAssistantHandler(svc).RegisterGinRoutes(router)
	server := httptest.NewServer(router)
	t.Cleanup(server.Close)
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar, Timeout: 10 * time.Second}
	return &fixture{server: server, client: client, store: store, configs: configs, posts: posts}
}
func (f *fixture) request(t *testing.T, method, path string, body any) *http.Response {
	t.Helper()
	var reader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reader = bytes.NewReader(b)
	}
	req, _ := http.NewRequest(method, f.server.URL+path, reader)
	req.Header.Set("Content-Type", "application/json")
	resp, err := f.client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}
func (f *fixture) create(t *testing.T) string {
	t.Helper()
	resp := f.request(t, "POST", "/ai/conversations", map[string]string{"title": "测试对话"})
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		t.Fatalf("create status %d", resp.StatusCode)
	}
	var out struct {
		Data domain.Conversation `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out.Data.ID
}
func (f *fixture) get(t *testing.T, id string) domain.Conversation {
	t.Helper()
	resp := f.request(t, "GET", "/ai/conversations/"+id, nil)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("get status %d", resp.StatusCode)
	}
	var out struct {
		Data domain.Conversation `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return out.Data
}
func readEvents(t *testing.T, resp *http.Response) []domain.Event {
	t.Helper()
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("stream status %d: %s", resp.StatusCode, b)
	}
	if !strings.Contains(resp.Header.Get("Cache-Control"), "no-transform") || resp.Header.Get("X-Accel-Buffering") != "no" {
		t.Fatal("SSE must disable proxy compression and buffering")
	}
	events := []domain.Event{}
	scanner := bufio.NewScanner(resp.Body)
	scanner.Buffer(make([]byte, 4096), 256<<10)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			var event domain.Event
			if err := json.Unmarshal([]byte(line[6:]), &event); err != nil {
				t.Fatal(err)
			}
			events = append(events, event)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
	return events
}
func writeText(w http.ResponseWriter, text string, finish bool) {
	b, _ := json.Marshal(map[string]any{"id": "chat-test", "object": "chat.completion.chunk", "choices": []any{map[string]any{"index": 0, "delta": map[string]string{"role": "assistant", "content": text}, "finish_reason": nil}}})
	_, _ = fmt.Fprintf(w, "data: %s\n\n", b)
	if finish {
		_, _ = fmt.Fprint(w, "data: {\"id\":\"chat-test\",\"object\":\"chat.completion.chunk\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}]}\n\ndata: [DONE]\n\n")
	}
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func TestChatStreamsToolsHistoryIdempotencyAndVisitorOwnership(t *testing.T) {
	var calls atomic.Int32
	f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer provider-private-key" || r.URL.Path != "/v1/chat/completions" {
			t.Error("incorrect Eino authentication or endpoint")
		}
		var body struct {
			Messages []map[string]any `json:"messages"`
			Stream   bool             `json:"stream"`
			Tools    []any            `json:"tools"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		if !body.Stream || len(body.Tools) != 4 {
			t.Error("real streaming and configured tools are required")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		if calls.Add(1) == 1 {
			_, _ = fmt.Fprint(w, "data: {\"id\":\"chat-test\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"tool_calls\":[{\"index\":0,\"id\":\"call-search\",\"type\":\"function\",\"function\":{\"name\":\"search_posts\",\"arguments\":\"{\\\"query\\\":\\\"Go\\\"}\"}}]},\"finish_reason\":null}]}\n\ndata: {\"id\":\"chat-test\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"tool_calls\"}]}\n\ndata: [DONE]\n\n")
		} else {
			if calls.Load() == 2 {
				found := false
				for _, m := range body.Messages {
					if m["role"] == "tool" {
						found = true
					}
				}
				if !found {
					t.Error("tool results must return through Eino's tool message protocol")
				}
			}
			writeText(w, "请阅读 ", false)
			writeText(w, "[Go 入门](/posts/public-post)", true)
		}
	})
	resp := f.request(t, "GET", "/ai/conversations", nil)
	body, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if strings.Contains(string(body), "private") || strings.Contains(string(body), "ciphertext") || strings.Contains(string(body), "system_prompt") {
		t.Fatal("bootstrap must not expose internal model settings")
	}
	if len(resp.Cookies()) != 1 || !resp.Cookies()[0].HttpOnly || resp.Cookies()[0].SameSite != http.SameSiteLaxMode {
		t.Fatal("visitor credentials must be isolated in an HttpOnly SameSite cookie")
	}
	id := f.create(t)
	requestID := uuid.NewString()
	payload := map[string]string{"request_id": requestID, "content": "推荐 Go 文章"}
	events := readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", payload))
	if len(events) < 5 || events[0].Type != "meta" || events[len(events)-1].Status != "complete" || calls.Load() != 2 {
		t.Fatalf("unexpected Eino run: %+v calls=%d", events, calls.Load())
	}
	c := f.get(t, id)
	if len(c.Messages) != 2 || c.Messages[1].Content != "请阅读 [Go 入门](/posts/public-post)" || len(c.Messages[1].References) != 1 || c.Messages[1].References[0].Href != "/posts/public-post" {
		t.Fatalf("history must retain text and verified references: %+v", c)
	}
	readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", payload))
	if calls.Load() != 2 || len(f.get(t, id).Messages) != 2 {
		t.Fatal("duplicate request IDs must replay without another model call")
	}
	payload["content"] = "different"
	resp = f.request(t, "POST", "/ai/conversations/"+id+"/messages", payload)
	_ = resp.Body.Close()
	if resp.StatusCode != 409 {
		t.Fatal("same request ID cannot change its question")
	}
	payload["request_id"] = uuid.NewString()
	readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", payload))
	if len(f.get(t, id).Messages) != 4 {
		t.Fatal("a conversation must support continued multi-turn chat")
	}
	other := &http.Client{Timeout: time.Second}
	resp, err := other.Get(f.server.URL + "/ai/conversations/" + id)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal("a second visitor must not read another visitor's history")
	}
	f.store.mu.Lock()
	defer f.store.mu.Unlock()
	if len(f.store.visitors) != 2 {
		t.Fatal("visitors must have distinct persistent identifiers")
	}
	for visitor, ip := range f.store.visitors {
		if len(visitor) != 64 || ip == "" {
			t.Fatal("visitor identity and IP must be recorded")
		}
	}
}

func TestChatStopAndDisconnectPersistPartialAnswers(t *testing.T) {
	for _, disconnect := range []bool{false, true} {
		t.Run(fmt.Sprint("disconnect=", disconnect), func(t *testing.T) {
			started := make(chan struct{}, 1)
			f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				writeText(w, "已生成的部分回答", false)
				started <- struct{}{}
				<-r.Context().Done()
			})
			id := f.create(t)
			requestID := uuid.NewString()
			payload := map[string]string{"request_id": requestID, "content": "请继续"}
			resp := f.request(t, "POST", "/ai/conversations/"+id+"/messages", payload)
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("provider did not start")
			}
			scanner := bufio.NewScanner(resp.Body)
			seen := false
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "已生成的部分回答") {
					seen = true
					break
				}
			}
			if !seen {
				t.Fatal("text must reach the browser before completion")
			}
			conflict := f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "并发提问"})
			_ = conflict.Body.Close()
			if conflict.StatusCode != 409 {
				t.Fatal("concurrent generation must be rejected")
			}
			if disconnect {
				_ = resp.Body.Close()
			} else {
				stop := f.request(t, "POST", "/ai/conversations/"+id+"/stop", map[string]string{"request_id": requestID})
				_ = stop.Body.Close()
				_, _ = io.Copy(io.Discard, resp.Body)
				_ = resp.Body.Close()
			}
			deadline := time.Now().Add(3 * time.Second)
			for {
				c := f.get(t, id)
				if c.Generation.Status != "running" {
					want := "stopped"
					if disconnect {
						want = "interrupted"
					}
					if c.Messages[1].Content != "已生成的部分回答" || c.Messages[1].Status != want {
						t.Fatalf("partial answer must survive cancellation: %+v", c.Messages[1])
					}
					break
				}
				if time.Now().After(deadline) {
					t.Fatal("generation lock was not released")
				}
				time.Sleep(20 * time.Millisecond)
			}
		})
	}
}

func TestChatProviderErrorsAreSafeAndDisabledChatCannotStart(t *testing.T) {
	f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = fmt.Fprint(w, `{"error":{"message":"provider-private-key internal-details"}}`)
	})
	id := f.create(t)
	events := readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "你好"}))
	last := events[len(events)-1]
	if last.Status != "failed" || strings.Contains(last.Message, "private") || strings.Contains(last.Message, "internal") {
		t.Fatal("upstream errors and keys must not leak")
	}
	if f.get(t, id).Generation.Status != "failed" {
		t.Fatal("model errors must release persistent generation state")
	}
	f.configs.cfg.Enabled = false
	resp := f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "你好"})
	_ = resp.Body.Close()
	if resp.StatusCode != 503 {
		t.Fatal("disabled assistants must not generate")
	}
}

func writeTool(w http.ResponseWriter, name string) {
	b, _ := json.Marshal(map[string]any{"id": "chat-test", "choices": []any{map[string]any{"index": 0, "delta": map[string]any{"role": "assistant", "tool_calls": []any{map[string]any{"index": 0, "id": "tool-call", "type": "function", "function": map[string]string{"name": name, "arguments": `{"post_id":"public-post"}`}}}}, "finish_reason": "tool_calls"}}})
	_, _ = fmt.Fprintf(w, "data: %s\n\ndata: [DONE]\n\n", b)
}

func TestChatRejectsDisabledToolsAndLimitsToolLoops(t *testing.T) {
	for _, test := range []struct {
		name     string
		disabled bool
		maxCalls int32
	}{{"disabled", true, 1}, {"loop", false, 4}} {
		t.Run(test.name, func(t *testing.T) {
			var calls atomic.Int32
			f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				w.Header().Set("Content-Type", "text/event-stream")
				writeTool(w, "read_post")
			})
			if test.disabled {
				f.configs.cfg.EnabledTools = []string{}
			}
			id := f.create(t)
			events := readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "读取文章"}))
			if events[len(events)-1].Status != "failed" || calls.Load() != test.maxCalls {
				t.Fatal("tool permissions and execution loops must be bounded", events, calls.Load())
			}
		})
	}
}

func TestChatNeverReturnsNonPublicArticleContent(t *testing.T) {
	var calls atomic.Int32
	f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		if calls.Add(1) == 1 {
			writeTool(w, "read_post")
			return
		}
		b, _ := io.ReadAll(r.Body)
		if strings.Contains(string(b), "文章正文") {
			t.Error("private article content reached the model")
		}
		writeText(w, "文章未公开，无法读取。", true)
	})
	f.posts.private = true
	id := f.create(t)
	readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "读取文章"}))
	if len(f.get(t, id).Messages[1].References) != 0 {
		t.Fatal("private articles must never become reference links")
	}
}

func TestChatGPTResponsesStreamingThroughEino(t *testing.T) {
	f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/responses" {
			t.Error("GPT must use Responses")
		}
		w.Header().Set("Content-Type", "text/event-stream")
		frames := []string{
			`{"type":"response.created","response":{"id":"resp-1","status":"in_progress","output":[]}}`,
			`{"type":"response.output_item.added","output_index":0,"item":{"id":"msg-1","type":"message","role":"assistant","status":"in_progress","content":[]}}`,
			`{"type":"response.content_part.added","item_id":"msg-1","output_index":0,"content_index":0,"part":{"type":"output_text","text":"","annotations":[]}}`,
			`{"type":"response.output_text.delta","item_id":"msg-1","output_index":0,"content_index":0,"delta":"GPT 流式回复"}`,
			`{"type":"response.output_text.done","item_id":"msg-1","output_index":0,"content_index":0,"text":"GPT 流式回复"}`,
			`{"type":"response.output_item.done","output_index":0,"item":{"id":"msg-1","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"GPT 流式回复","annotations":[]}]}}`,
			`{"type":"response.completed","response":{"id":"resp-1","status":"completed","output":[{"id":"msg-1","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"GPT 流式回复","annotations":[]}]}]}}`,
		}
		for _, frame := range frames {
			var e struct {
				Type string `json:"type"`
			}
			_ = json.Unmarshal([]byte(frame), &e)
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", e.Type, frame)
			w.(http.Flusher).Flush()
		}
	})
	f.configs.cfg.Providers[0].ApiType = "openai_responses"
	id := f.create(t)
	events := readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "你好"}))
	c := f.get(t, id)
	if c.Messages[1].Content != "GPT 流式回复" || events[len(events)-1].Status != "complete" {
		t.Fatalf("Responses text must stream once and persist: %+v %+v", c.Messages[1], events)
	}
}

func TestChatGPTResponsesToolResultsReturnThroughEino(t *testing.T) {
	var calls atomic.Int32
	f := newFixture(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "text/event-stream")
		frames := []string{}
		if calls.Add(1) == 1 {
			frames = []string{
				`{"type":"response.created","response":{"id":"resp-1","status":"in_progress","output":[]}}`,
				`{"type":"response.output_item.added","output_index":0,"item":{"id":"fc-1","type":"function_call","call_id":"read-call","name":"read_post","arguments":"","status":"in_progress"}}`,
				`{"type":"response.function_call_arguments.delta","item_id":"fc-1","output_index":0,"delta":"{\"post_id\":\"public-post\"}"}`,
				`{"type":"response.output_item.done","output_index":0,"item":{"id":"fc-1","type":"function_call","call_id":"read-call","name":"read_post","arguments":"{\"post_id\":\"public-post\"}","status":"completed"}}`,
				`{"type":"response.completed","response":{"id":"resp-1","status":"completed","output":[]}}`,
			}
		} else {
			input, _ := json.Marshal(body["input"])
			if !strings.Contains(string(input), "function_call_output") || !strings.Contains(string(input), "文章正文") {
				t.Error("GPT must receive Eino-serialized function results")
			}
			frames = []string{
				`{"type":"response.created","response":{"id":"resp-2","status":"in_progress","output":[]}}`,
				`{"type":"response.output_item.added","output_index":0,"item":{"id":"msg-2","type":"message","role":"assistant","status":"in_progress","content":[]}}`,
				`{"type":"response.content_part.added","item_id":"msg-2","output_index":0,"content_index":0,"part":{"type":"output_text","text":"","annotations":[]}}`,
				`{"type":"response.output_text.delta","item_id":"msg-2","output_index":0,"content_index":0,"delta":"已阅读公开文章"}`,
				`{"type":"response.output_item.done","output_index":0,"item":{"id":"msg-2","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"已阅读公开文章","annotations":[]}]}}`,
				`{"type":"response.completed","response":{"id":"resp-2","status":"completed","output":[]}}`,
			}
		}
		for _, frame := range frames {
			var event struct {
				Type string `json:"type"`
			}
			_ = json.Unmarshal([]byte(frame), &event)
			_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, frame)
			w.(http.Flusher).Flush()
		}
	})
	f.configs.cfg.Providers[0].ApiType = "openai_responses"
	f.configs.cfg.EnabledTools = []string{"read_post"}
	id := f.create(t)
	events := readEvents(t, f.request(t, "POST", "/ai/conversations/"+id+"/messages", map[string]string{"request_id": uuid.NewString(), "content": "读取文章"}))
	c := f.get(t, id)
	if calls.Load() != 2 || events[len(events)-1].Status != "complete" || c.Messages[1].Content != "已阅读公开文章" || len(c.Messages[1].References) != 1 {
		t.Fatalf("GPT tool loop must stream and retain verified references: %+v %+v", c.Messages[1], events)
	}
}
