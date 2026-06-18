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
	"net/http"

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
	adminGroup.GET("/post-index/sitemap", apiwrap.Wrap(h.GetSitemap))
	adminGroup.POST("/post-index/sitemap", apiwrap.Wrap(h.GenerateSitemap))
	adminGroup.GET("/post-index/robots", apiwrap.Wrap(h.GetRobotsTxt))
	adminGroup.PUT("/post-index/robots", apiwrap.WrapWithBody(h.SaveRobotsTxt))
	adminGroup.POST("/post-index/robots/generate", apiwrap.Wrap(h.GenerateRobotsTxt))
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

func (h *PostIndexHandler) GetSitemap(ctx *gin.Context) (*apiwrap.ResponseBody[SitemapVO], error) {
	content, exists, err := h.serv.GetSitemap(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(SitemapVO{Content: content, Exists: exists}), nil
}

func (h *PostIndexHandler) GetRobotsTxt(ctx *gin.Context) (*apiwrap.ResponseBody[RobotsTxtVO], error) {
	content, exists, err := h.serv.GetRobotsTxt(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(RobotsTxtVO{Content: content, Exists: exists}), nil
}

func (h *PostIndexHandler) GenerateRobotsTxt(ctx *gin.Context) (*apiwrap.ResponseBody[RobotsTxtVO], error) {
	content, err := h.serv.GenerateRobotsTxt(ctx)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(RobotsTxtVO{Content: content, Exists: true}), nil
}

func (h *PostIndexHandler) SaveRobotsTxt(ctx *gin.Context, req RobotsTxtRequest) (*apiwrap.ResponseBody[any], error) {
	if len(req.Content) > 100*1024 {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "robots.txt content is too large")
	}
	return apiwrap.SuccessResponse(), h.serv.SaveRobotsTxt(ctx, req.Content)
}
