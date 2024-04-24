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
	postServ "github.com/chenmingyong0423/fnote/server/internal/post/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft"
	"github.com/gin-gonic/gin"
)

func NewAggregatePostHandler(postServ postServ.IPostService, postDraftServ post_draft.Service) *AggregatePostHandler {
	return &AggregatePostHandler{
		postServ:      postServ,
		postDraftServ: postDraftServ,
	}
}

type AggregatePostHandler struct {
	postServ      postServ.IPostService
	postDraftServ post_draft.Service
}

func (h *AggregatePostHandler) RegisterGinRoutes(engine *gin.Engine) {
	adminGroup := engine.Group("/admin-api")
	adminGroup.GET("/post-draft/:id", apiwrap.Wrap(h.GetPostDraftById))
}

func (h *AggregatePostHandler) GetPostDraftById(ctx *gin.Context) (*apiwrap.ResponseBody[*PostDraftVO], error) {
	//postDraft, err := h.serv.GetPostDraftById(ctx, ctx.Param("id"))
	//if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
	//	return nil, err
	//}
	//if postDraft == nil {
	//	return nil, apiwrap.NewErrorResponseBody(404, "post draft not found")
	//}
	//return apiwrap.SuccessResponseWithData(h.toVO(postDraft)), nil
	return nil, nil
}
