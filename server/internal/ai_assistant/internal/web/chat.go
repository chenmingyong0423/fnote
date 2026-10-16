package web

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/service"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/gin-gonic/gin"
)

const visitorCookie = "fnote_ai_visitor"

type AIAssistantHandler struct{ service *service.AIAssistantService }

func NewAIAssistantHandler(s *service.AIAssistantService) *AIAssistantHandler {
	return &AIAssistantHandler{s}
}

func (h *AIAssistantHandler) RegisterGinRoutes(e *gin.Engine) {
	g := e.Group("/ai", h.identify)
	g.GET("/conversations", h.list)
	g.POST("/conversations", h.create)
	g.GET("/conversations/:id", h.get)
	g.POST("/conversations/:id/messages", h.send)
	g.POST("/conversations/:id/stop", h.stop)
}

func (h *AIAssistantHandler) identify(c *gin.Context) {
	c.Header("Cache-Control", "no-store")
	if c.Request.Method != http.MethodGet {
		// Browser mutations use JSON over the same-origin web proxy. Reject cross-
		// site requests and simple form posts before touching visitor data.
		if c.GetHeader("Sec-Fetch-Site") == "cross-site" || !strings.HasPrefix(strings.ToLower(c.GetHeader("Content-Type")), "application/json") {
			c.AbortWithStatusJSON(403, gin.H{"code": 403, "message": "不允许跨站请求"})
			return
		}
	}
	token, _ := c.Cookie(visitorCookie)
	decoded, err := hex.DecodeString(token)
	if err != nil || len(decoded) != 32 {
		decoded = make([]byte, 32)
		if _, err = rand.Read(decoded); err != nil {
			replyError(c, err)
			c.Abort()
			return
		}
		token = hex.EncodeToString(decoded)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(visitorCookie, token, 365*24*60*60, "/", "", c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https", true)
	}
	hash := sha256.Sum256(decoded)
	visitor := hex.EncodeToString(hash[:])
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := h.service.Identify(ctx, visitor, c.ClientIP()); err != nil {
		replyError(c, err)
		c.Abort()
		return
	}
	c.Set("ai_visitor", visitor)
	c.Next()
}

func (h *AIAssistantHandler) list(c *gin.Context) {
	list, cfg, err := h.service.List(c.Request.Context(), c.GetString("ai_visitor"))
	if err != nil {
		replyError(c, err)
		return
	}
	c.JSON(200, apiwrap.SuccessResponseWithData(gin.H{"conversations": list, "assistant_name": cfg.AssistantName, "enabled": cfg.Enabled}))
}
func (h *AIAssistantHandler) get(c *gin.Context) {
	conversation, err := h.service.Get(c.Request.Context(), c.GetString("ai_visitor"), c.Param("id"))
	if err != nil {
		replyError(c, err)
		return
	}
	c.JSON(200, apiwrap.SuccessResponseWithData(conversation))
}
func (h *AIAssistantHandler) create(c *gin.Context) {
	var req struct {
		Title string `json:"title"`
	}
	if !bind(c, &req) {
		return
	}
	conversation, err := h.service.Create(c.Request.Context(), c.GetString("ai_visitor"), req.Title)
	if err != nil {
		replyError(c, err)
		return
	}
	c.JSON(201, apiwrap.SuccessResponseWithData(conversation))
}
func (h *AIAssistantHandler) stop(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id"`
	}
	if !bind(c, &req) {
		return
	}
	if err := h.service.Stop(c.Request.Context(), c.GetString("ai_visitor"), c.Param("id"), req.RequestID); err != nil {
		replyError(c, err)
		return
	}
	c.JSON(200, apiwrap.SuccessResponse())
}
func (h *AIAssistantHandler) send(c *gin.Context) {
	var req struct {
		RequestID string `json:"request_id"`
		Content   string `json:"content"`
	}
	if !bind(c, &req) {
		return
	}
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	session, err := h.service.Start(ctx, c.GetString("ai_visitor"), c.ClientIP(), c.Param("id"), req.RequestID, req.Content)
	if err != nil {
		replyError(c, err)
		return
	}
	events := make(chan domain.Event, 16)
	go func() {
		defer close(events)
		session.Run(func(event domain.Event) error {
			select {
			case events <- event:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		})
	}()
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("X-Accel-Buffering", "no")
	// no-transform prevents the web proxy's gzip middleware from buffering SSE.
	c.Header("Cache-Control", "no-cache, no-store, no-transform")
	c.Status(200)
	c.Writer.Flush()
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = http.NewResponseController(c.Writer).SetWriteDeadline(time.Now().Add(15 * time.Second))
			if _, err := fmt.Fprint(c.Writer, ": heartbeat\n\n"); err != nil {
				return
			}
			c.Writer.Flush()
		case event, ok := <-events:
			if !ok {
				return
			}
			data, _ := json.Marshal(event)
			_ = http.NewResponseController(c.Writer).SetWriteDeadline(time.Now().Add(15 * time.Second))
			if _, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event.Type, data); err != nil {
				return
			}
			c.Writer.Flush()
		}
	}
}

func bind(c *gin.Context, value any) bool {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 64<<10)
	if err := c.ShouldBindJSON(value); err != nil {
		c.JSON(400, gin.H{"code": 400, "message": "请求格式不正确"})
		return false
	}
	return true
}
func replyError(c *gin.Context, err error) {
	var known *domain.Error
	if errors.As(err, &known) {
		c.JSON(known.Status, gin.H{"code": known.Status, "message": known.Message})
		return
	}
	c.JSON(503, gin.H{"code": 503, "message": "对话服务暂不可用，请稍后再试"})
}
