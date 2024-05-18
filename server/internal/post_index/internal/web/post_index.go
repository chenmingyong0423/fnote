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
	"github.com/chenmingyong0423/fnote/server/internal/post_index/internal/service"
	"github.com/gin-gonic/gin"
)

func NewPostIndexHandler(serv service.IPostIndexService) *PostIndexHandler {
	return &PostIndexHandler{
		serv: serv,
	}
}

type PostIndexHandler struct {
	serv service.IPostIndexService
}

func (h *PostIndexHandler) RegisterGinRoutes(engine *gin.Engine) {
	engine.POST("/post-index/baidu/push", apiwrap.WrapWithBody(h.BaiduPostIndex))
	adminGroup := engine.Group("/admin-api")
	adminGroup.POST("/post-index/sitemap", apiwrap.Wrap(h.GenerateSitemap))
}

func (h *PostIndexHandler) BaiduPostIndex(ctx *gin.Context, req PostIndexRequest) (*apiwrap.ResponseBody[BaiduPushVO], error) {
	baiduResponse, err := h.serv.PushUrls2Baidu(ctx, req.Urls)
	if err != nil {
		return nil, err
	}
	if baiduResponse == nil {
		return apiwrap.SuccessResponseWithData(BaiduPushVO{}), nil
	}
	return apiwrap.SuccessResponseWithData(BaiduPushVO{
		Remain:      baiduResponse.Remain,
		Success:     baiduResponse.Success,
		NotSameSite: baiduResponse.NotSameSite,
		NotValid:    baiduResponse.NotValid,
		Err:         baiduResponse.Err,
		Message:     baiduResponse.Message,
	}), nil
}

func (h *PostIndexHandler) GenerateSitemap(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.GenerateSitemap(ctx)
}
