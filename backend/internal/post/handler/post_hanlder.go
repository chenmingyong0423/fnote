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
	"net/http"
	"slices"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chenmingyong0423/fnote/backend/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/backend/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/backend/internal/post/service"
	"github.com/gin-gonic/gin"
)

type SummaryPostVO struct {
	Sug          string   `json:"sug"`
	Author       string   `json:"author"`
	Title        string   `json:"title"`
	Summary      string   `json:"summary"`
	CoverImg     string   `json:"cover_img"`
	Categories   []string `json:"categories"`
	Tags         []string `json:"tags"`
	LikeCount    int      `json:"like_count"`
	CommentCount int      `json:"comment_count"`
	VisitCount   int      `json:"visit_count"`
	Priority     int      `json:"priority"`
	CreateTime   int64    `json:"create_time"`
}

func NewPostHandler(serv service.IPostService) *PostHandler {
	return &PostHandler{
		serv: serv,
	}
}

type PostHandler struct {
	serv service.IPostService
}

func (h *PostHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/posts")
	group.GET("/latest", api.Wrap(h.GetLatestPosts))
	group.GET("", api.WrapWithBody(h.GetPosts))
	group.GET("/:id", api.Wrap(h.GetPostBySug))
	group.POST("/:id/likes", api.Wrap(h.AddLike))
	group.DELETE("/:id/likes", api.Wrap(h.DeleteLike))
}

func (h *PostHandler) GetLatestPosts(ctx *gin.Context) (listVO api.ListVO[*SummaryPostVO], err error) {
	posts, err := h.serv.GetLatestPosts(ctx)
	if err != nil {
		return
	}
	listVO.List = h.postsToPostVOs(posts)
	return
}

func (h *PostHandler) postsToPostVOs(posts []*domain.Post) []*SummaryPostVO {
	postVOs := make([]*SummaryPostVO, 0, len(posts))
	for _, post := range posts {
		postVOs = append(postVOs, &SummaryPostVO{
			Sug:          post.PrimaryPost.Sug,
			Author:       post.PrimaryPost.Author,
			Title:        post.PrimaryPost.Title,
			Summary:      post.PrimaryPost.Summary,
			CoverImg:     post.PrimaryPost.CoverImg,
			Categories:   post.PrimaryPost.Categories,
			Tags:         post.PrimaryPost.Tags,
			LikeCount:    post.PrimaryPost.LikeCount,
			CommentCount: post.PrimaryPost.CommentCount,
			VisitCount:   post.PrimaryPost.VisitCount,
			Priority:     post.PrimaryPost.Priority,
			CreateTime:   post.PrimaryPost.CreateTime,
		})
	}
	return postVOs
}

func (h *PostHandler) GetPosts(ctx *gin.Context, req *domain.PostRequest) (pageVO api.PageVO[*SummaryPostVO], err error) {
	req.ValidateAndSetDefault()
	posts, cnt, err := h.serv.GetPosts(ctx, req)
	if err != nil {
		return
	}
	pageVO.Page = req.Page
	pageVO.List = h.postsToPostVOs(posts)
	pageVO.SetTotalCountAndCalculateTotalPages(cnt)
	return
}

func (h *PostHandler) GetPostBySug(ctx *gin.Context) (vo *domain.DetailPostVO, err error) {
	sug := ctx.Param("id")
	post, err := h.serv.GetPunishedPostById(ctx, sug)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, api.NewErrorResponseBody(http.StatusBadRequest, "The postId does not exist.")
		}
		return
	}
	vo = new(domain.DetailPostVO)
	vo.PrimaryPost, vo.ExtraPost, vo.IsLiked = post.PrimaryPost, post.ExtraPost, slices.Contains(post.Likes, ctx.ClientIP())
	return
}

func (h *PostHandler) AddLike(ctx *gin.Context) (r any, err error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	sug := ctx.Param("id")
	return r, h.serv.AddLike(ctx, sug, ip)
}

func (h *PostHandler) DeleteLike(ctx *gin.Context) (r any, err error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, api.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	sug := ctx.Param("id")
	return r, h.serv.DeleteLike(ctx, sug, ip)
}
