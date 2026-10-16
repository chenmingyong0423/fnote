//go:build wireinject

package ai_assistant

import (
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/repository"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/repository/dao"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/web"
	"github.com/chenmingyong0423/fnote/server/internal/post"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/chenmingyong0423/go-mongox/v2"
	"github.com/google/wire"
)

func InitAIAssistantModule(db *mongox.Database, cfg *website_config.Module, posts *post.Module) (*Module, error) {
	panic(wire.Build(initModelHTTPClient, dao.NewAIAssistantDao, repository.NewAIAssistantRepository,
		service.NewToolService, service.NewModelService, service.NewAIAssistantService, web.NewAIAssistantHandler,
		wire.FieldsOf(new(*website_config.Module), "Svc"), wire.FieldsOf(new(*post.Module), "Svc"),
		wire.Bind(new(dao.IAIAssistantDao), new(*dao.AIAssistantDao)),
		wire.Bind(new(repository.IAIAssistantRepository), new(*repository.AIAssistantRepository)), wire.Struct(new(Module), "*")))
}
