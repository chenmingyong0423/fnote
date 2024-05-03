// Copyright 2024 chenmingyong0423

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
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post_visit/internal/service"
	"github.com/gin-gonic/gin"
)

func NewPostVisitHandler(serv service.IPostVisitService) *PostVisitHandler {
	return &PostVisitHandler{
		serv: serv,
	}
}

type PostVisitHandler struct {
	serv service.IPostVisitService
}

func (h *PostVisitHandler) RegisterGinRoutes(engine *gin.Engine) {
	routerGroup := engine.Group("/logs")
	routerGroup.POST("/post-visit", apiwrap.WrapWithBody(h.CollectPostVisit))
}

func (h *PostVisitHandler) CollectPostVisit(ctx *gin.Context, req PostVisitRequest) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.SavePostVisit(ctx, domain.PostVisit{
		PostId:    req.PostId,
		Ip:        ctx.ClientIP(),
		UserAgent: ctx.GetHeader("User-Agent"),
		Origin:    ctx.GetHeader("Origin"),
		Referer:   ctx.GetHeader("Referer"),
		StayTime:  req.StayTime,
		VisitAt:   req.VisitAt,
	})
}
