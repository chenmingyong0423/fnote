// Copyright 2023 chenmingyong0423

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package web

import (
	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/message_template/internal/service"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/gin-gonic/gin"
)

func NewMessageTemplateHandler(serv service.IMessageTemplateService) *MessageTemplateHandler {
	return &MessageTemplateHandler{
		serv: serv,
	}
}

type MessageTemplateHandler struct {
	serv service.IMessageTemplateService
}

func (h *MessageTemplateHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/admin-api/message-templates")
	group.GET("", apiwrap.Wrap(h.FindAll))
	group.PUT("/:id", apiwrap.WrapWithBody(h.Update))
	group.PUT("/:id/active", apiwrap.WrapWithBody(h.UpdateActive))
}

func (h *MessageTemplateHandler) FindAll(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[MessageTemplateVO]], error) {
	templates, err := h.serv.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]MessageTemplateVO, 0, len(templates))
	for _, messageTemplate := range templates {
		result = append(result, toVO(messageTemplate))
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(result)), nil
}

func (h *MessageTemplateHandler) Update(ctx *gin.Context, req UpdateMessageTemplateRequest) (*apiwrap.ResponseBody[any], error) {
	if err := h.serv.Update(ctx, ctx.Param("id"), req.Title, req.Content); err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *MessageTemplateHandler) UpdateActive(ctx *gin.Context, req UpdateMessageTemplateActiveRequest) (*apiwrap.ResponseBody[any], error) {
	if err := h.serv.UpdateActive(ctx, ctx.Param("id"), *req.Active); err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponse(), nil
}

func toVO(messageTemplate domain.MessageTemplate) MessageTemplateVO {
	return MessageTemplateVO{
		Id:            messageTemplate.Id,
		Name:          string(messageTemplate.Name),
		Title:         messageTemplate.Title,
		Content:       messageTemplate.Content,
		Active:        messageTemplate.Active,
		RecipientType: uint(messageTemplate.RecipientType),
		Variables:     domain.RequiredVariables(messageTemplate.Name),
		CreatedAt:     messageTemplate.CreatedAt,
		UpdatedAt:     messageTemplate.UpdatedAt,
	}
}
