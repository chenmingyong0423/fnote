package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/google/uuid"
)

var errStopped = errors.New("visitor stopped generation")

type activeGeneration struct {
	conversationID, requestID string
	cancel                    context.CancelCauseFunc
}
type requestWindow struct {
	count int
	until time.Time
}

type AIAssistantService struct {
	repo    repository.IAIAssistantRepository
	configs website_config.Service
	model   *ModelService
	mu      sync.Mutex
	active  map[string]*activeGeneration
	rates   map[string]requestWindow
}

func NewAIAssistantService(repo repository.IAIAssistantRepository, cfg website_config.Service, model *ModelService) *AIAssistantService {
	return &AIAssistantService{repo: repo, configs: cfg, model: model, active: map[string]*activeGeneration{}, rates: map[string]requestWindow{}}
}

func (s *AIAssistantService) Identify(ctx context.Context, id, ip string) error {
	return s.repo.Identify(ctx, id, ip)
}

func (s *AIAssistantService) List(ctx context.Context, visitor string) ([]domain.Conversation, website_config.AIAssistantConfig, error) {
	cfg, err := s.configs.GetAIAssistantConfig(ctx)
	if err != nil {
		return nil, cfg, err
	}
	list, err := s.repo.List(ctx, visitor)
	if err != nil {
		return nil, cfg, err
	}
	for i := range list {
		if list[i].Generation != nil && list[i].Generation.Status == "running" {
			c, err := s.Get(ctx, visitor, list[i].ID)
			if err != nil {
				return nil, cfg, err
			}
			list[i].Generation = c.Generation
		}
	}
	return list, cfg, nil
}

func (s *AIAssistantService) Get(ctx context.Context, visitor, id string) (*domain.Conversation, error) {
	c, err := s.repo.Get(ctx, visitor, id)
	if err != nil {
		return nil, err
	}
	if c.Generation != nil && c.Generation.Status == "running" {
		s.mu.Lock()
		a := s.active[visitor]
		active := a != nil && a.conversationID == id
		s.mu.Unlock()
		if !active {
			for i := range c.Messages {
				if c.Messages[i].ID == c.Generation.MessageID {
					c.Messages[i].Status = "interrupted"
					c.Messages[i].Error = "回答已中断，可继续提问"
					if err := s.repo.SaveAnswer(ctx, visitor, id, c.Messages[i], true); err != nil {
						return nil, err
					}
					c.Generation.Status = "interrupted"
					break
				}
			}
		}
	}
	return c, nil
}

func (s *AIAssistantService) Create(ctx context.Context, visitor, title string) (*domain.Conversation, error) {
	cfg, err := s.configs.GetAIAssistantConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, domain.ErrDisabled
	}
	count, err := s.repo.Count(ctx, visitor)
	if err != nil {
		return nil, err
	}
	if count >= 100 {
		return nil, &domain.Error{429, "对话数量达到上限，请继续已有对话"}
	}
	title = limitText(strings.TrimSpace(title), 40)
	if title == "" {
		title = "新对话"
	}
	now := time.Now()
	c := domain.Conversation{ID: uuid.NewString(), VisitorID: visitor, Title: title, Messages: []domain.Message{}, CreatedAt: now, UpdatedAt: now}
	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &c, nil
}

type Session struct {
	service                            *AIAssistantService
	conversationID, visitor, requestID string
	history                            []domain.Message
	answer                             domain.Message
	cfg                                website_config.AIAssistantConfig
	ctx                                context.Context
	cancel                             context.CancelFunc
	active                             *activeGeneration
	replay                             bool
}

func (s *AIAssistantService) Start(ctx context.Context, visitor, ip, id, requestID, content string) (*Session, error) {
	content = strings.TrimSpace(content)
	if _, err := uuid.Parse(requestID); err != nil {
		return nil, &domain.Error{400, "请求 ID 不正确"}
	}
	if content == "" || utf8.RuneCountInString(content) > 8000 {
		return nil, &domain.Error{400, "问题需要为 1 至 8000 字"}
	}
	c, err := s.Get(ctx, visitor, id)
	if err != nil {
		return nil, err
	}
	for i, m := range c.Messages {
		if m.Role == "user" && m.RequestID == requestID {
			if m.Content != content {
				return nil, &domain.Error{409, "请求 ID 已用于其他问题，请刷新后重试"}
			}
			if i+1 >= len(c.Messages) || c.Messages[i+1].Status == "typing" {
				return nil, domain.ErrBusy
			}
			return &Session{conversationID: id, requestID: requestID, answer: c.Messages[i+1], replay: true}, nil
		}
	}
	if len(c.Messages) >= 100 {
		return nil, &domain.Error{400, "此对话已达长度上限，请新建对话"}
	}
	cfg, err := s.configs.GetAIAssistantConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, domain.ErrDisabled
	}
	modelCtx, timeout := context.WithTimeout(ctx, 90*time.Second)
	modelCtx, cancel := context.WithCancelCause(modelCtx)
	a := &activeGeneration{conversationID: id, requestID: requestID, cancel: cancel}
	if err := s.reserve(visitor, ip, a); err != nil {
		timeout()
		cancel(err)
		return nil, err
	}
	now := time.Now()
	user := domain.Message{ID: requestID, RequestID: requestID, Role: "user", Content: content, Status: "complete", CreatedAt: now}
	answer := domain.Message{ID: uuid.NewString(), RequestID: requestID, Role: "assistant", Status: "typing", CreatedAt: now}
	// Only user messages carry the idempotency key in the Mongo query.
	answer.RequestID = ""
	if err := s.repo.Begin(ctx, visitor, id, user, answer); err != nil {
		s.release(visitor, a)
		timeout()
		cancel(err)
		return nil, err
	}
	return &Session{service: s, conversationID: id, visitor: visitor, requestID: requestID, history: append(c.Messages, user), answer: answer, cfg: cfg,
		ctx: modelCtx, cancel: func() { cancel(context.Canceled); timeout() }, active: a}, nil
}

func (s *AIAssistantService) reserve(visitor, ip string, a *activeGeneration) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active[visitor] != nil {
		return domain.ErrBusy
	}
	if len(s.active) >= 8 {
		return &domain.Error{429, "助手正忙，请稍后重试"}
	}
	now := time.Now()
	window := s.rates[ip]
	if !now.Before(window.until) {
		window = requestWindow{until: now.Add(time.Hour)}
	}
	if window.count >= 60 {
		return &domain.Error{429, "提问过于频繁，请稍后再试"}
	}
	if len(s.rates) >= 4096 {
		for key, w := range s.rates {
			if !now.Before(w.until) {
				delete(s.rates, key)
			}
		}
		if len(s.rates) >= 4096 {
			return &domain.Error{429, "助手正忙，请稍后重试"}
		}
	}
	window.count++
	s.rates[ip] = window
	s.active[visitor] = a
	return nil
}
func (s *AIAssistantService) release(visitor string, a *activeGeneration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.active[visitor] == a {
		delete(s.active, visitor)
	}
}

func (s *AIAssistantService) Stop(ctx context.Context, visitor, id, requestID string) error {
	if _, err := s.repo.Get(ctx, visitor, id); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if a := s.active[visitor]; a != nil && a.conversationID == id && a.requestID == requestID {
		a.cancel(errStopped)
	}
	return nil
}

func (session *Session) Run(emit func(domain.Event) error) {
	meta := domain.Event{Type: "meta", ConversationID: session.conversationID, RequestID: session.requestID, MessageID: session.answer.ID}
	if session.replay {
		if emit(meta) != nil {
			return
		}
		if session.answer.Content != "" {
			if emit(domain.Event{Type: "delta", Delta: session.answer.Content}) != nil {
				return
			}
		}
		if len(session.answer.References) > 0 {
			if emit(domain.Event{Type: "references", References: session.answer.References}) != nil {
				return
			}
		}
		_ = emit(domain.Event{Type: "done", Status: session.answer.Status, Message: session.answer.Error})
		return
	}
	s := session.service
	var runErr error
	defer func() {
		if recover() != nil {
			runErr = errModel
			slog.Error("AI generation interrupted by internal failure")
		}
		switch {
		case errors.Is(context.Cause(session.ctx), errStopped):
			session.answer.Status = "stopped"
		case errors.Is(session.ctx.Err(), context.DeadlineExceeded):
			session.answer.Status = "failed"
			session.answer.Error = "回答超时，已保存生成内容，可继续提问"
		case session.ctx.Err() != nil:
			session.answer.Status = "interrupted"
			session.answer.Error = "连接中断，已保存生成内容"
		case runErr != nil:
			session.answer.Status = "failed"
			session.answer.Error = runErr.Error()
		default:
			session.answer.Status = "complete"
		}
		// Finalization outlives browser cancellation, with a short bounded retry.
		persistCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var saveErr error
		for attempt := 0; attempt < 3; attempt++ {
			saveErr = s.repo.SaveAnswer(persistCtx, session.visitor, session.conversationID, session.answer, true)
			if saveErr == nil || persistCtx.Err() != nil {
				break
			}
		}
		if saveErr != nil {
			session.answer.Status = "failed"
			session.answer.Error = "回答保存失败，请刷新历史确认后再试"
			slog.Error("AI answer persistence failed")
		}
		s.release(session.visitor, session.active)
		session.cancel()
		_ = emit(domain.Event{Type: "done", Status: session.answer.Status, Message: session.answer.Error})
	}()
	if err := emit(meta); err != nil {
		runErr = errModel
		return
	}
	lastSave := time.Now()
	runErr = s.model.Run(session.ctx, session.cfg, session.history, func(event domain.Event) error {
		if session.ctx.Err() != nil {
			return session.ctx.Err()
		}
		if event.Type == "delta" {
			session.answer.Content += event.Delta
		}
		if event.Type == "references" {
			for _, ref := range event.References {
				found := false
				for _, saved := range session.answer.References {
					if saved.Href == ref.Href {
						found = true
						break
					}
				}
				if !found && len(session.answer.References) < 32 {
					session.answer.References = append(session.answer.References, ref)
				}
			}
			event.References = session.answer.References
		}
		if time.Since(lastSave) >= 500*time.Millisecond {
			if err := s.repo.SaveAnswer(session.ctx, session.visitor, session.conversationID, session.answer, false); err != nil {
				return errors.New("回答暂时无法保存，请稍后重试")
			}
			lastSave = time.Now()
		}
		return emit(event)
	})
}
