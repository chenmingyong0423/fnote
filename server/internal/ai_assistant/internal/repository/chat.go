package repository

import (
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/repository/dao"
)

// Persistence has no cross-module business queries; blog data comes from Services.
type IAIAssistantRepository interface{ dao.IAIAssistantDao }
type AIAssistantRepository struct{ dao.IAIAssistantDao }

func NewAIAssistantRepository(d dao.IAIAssistantDao) *AIAssistantRepository {
	return &AIAssistantRepository{d}
}
