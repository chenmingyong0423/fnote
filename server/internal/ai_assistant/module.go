package ai_assistant

import (
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/web"
)

type Handler = web.AIAssistantHandler
type Service = service.AIAssistantService
type Module struct {
	Svc *Service
	Hdl *Handler
}
