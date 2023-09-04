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

package handler

import (
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/ineternal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/ineternal/post/service"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log/slog"
	"net/http"
)

func NewPostHandler(engine *gin.Engine, serv service.IPostService) *PostHandler {
	h := &PostHandler{
		serv: serv,
	}

	engine.GET("/home/posts", h.GetHomePosts)
	engine.GET("/posts", h.GetPosts)
	engine.GET("/post/:sug", h.GetPostBySug)
	engine.POST("/post/:sug/likes", h.AddLike)
	engine.DELETE("/post/:sug/likes", h.DeleteLike)

	return h
}

type PostHandler struct {
	serv service.IPostService
}

func (h *PostHandler) GetHomePosts(ctx *gin.Context) {
	listVO, err := h.serv.GetHomePosts(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponseWithData[api.ListVO[*domain.SummaryPostVO]](listVO))
}

func (h *PostHandler) GetPosts(ctx *gin.Context) {
	pageRequest := &domain.PostRequest{}
	err := ctx.ShouldBindQuery(pageRequest)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrResponse)
		return
	}
	pageRequest.ValidateAndSetDefault()
	pageVO, err := h.serv.GetPosts(ctx, pageRequest)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponseWithData[*api.PageVO[*domain.SummaryPostVO]](pageVO))
}

func (h *PostHandler) GetPostBySug(ctx *gin.Context) {
	sug := ctx.Param("sug")
	detailPostVO, err := h.serv.GetPostBySug(ctx, sug)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponseWithData(detailPostVO))
}

func (h *PostHandler) AddLike(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if ip == "" {
		slog.ErrorContext(ctx, "post", errors.New("Fails to like without ip."))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrResponse)
		return
	}
	sug := ctx.Param("sug")
	err := h.serv.AddLike(ctx, sug, ip)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse())
}

func (h *PostHandler) DeleteLike(ctx *gin.Context) {
	ip := ctx.ClientIP()
	if ip == "" {
		slog.ErrorContext(ctx, "post", errors.New("Fails to unlike without ip."))
		ctx.AbortWithStatusJSON(http.StatusBadRequest, api.ErrResponse)
		return
	}
	sug := ctx.Param("sug")
	err := h.serv.DeleteLike(ctx, sug, ip)
	if err != nil {
		slog.ErrorContext(ctx, "post", err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, api.ErrResponse)
		return
	}
	ctx.JSON(http.StatusOK, api.SuccessResponse())
}
