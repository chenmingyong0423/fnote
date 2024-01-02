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

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/dto"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"

	configServ "github.com/chenmingyong0423/fnote/server/internal/website_config/service"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/api"
	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post/service"
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
	StickyWeight int      `json:"sticky_weight"`
	CreateTime   int64    `json:"create_time"`
}

func NewPostHandler(serv service.IPostService, cfgService configServ.IWebsiteConfigService) *PostHandler {
	return &PostHandler{
		serv:       serv,
		cfgService: cfgService,
	}
}

type PostHandler struct {
	serv       service.IPostService
	cfgService configServ.IWebsiteConfigService
}

func (h *PostHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/posts")
	group.GET("/latest", api.Wrap(h.GetLatestPosts))
	group.GET("", api.WrapWithBody(h.GetPosts))
	group.GET("/:id", api.Wrap(h.GetPostBySug))
	group.POST("/:id/likes", api.Wrap(h.AddLike))
	group.DELETE("/:id/likes", api.Wrap(h.DeleteLike))

	adminGroup := engine.Group("/admin/posts")
	adminGroup.GET("", api.WrapWithBody(h.AdminGetPosts))
}

func (h *PostHandler) GetLatestPosts(ctx *gin.Context) (listVO api.ListVO[*SummaryPostVO], err error) {
	countCfg, err := h.cfgService.GetFrontPostCount(ctx)
	if err != nil {
		return
	}
	posts, err := h.serv.GetLatestPosts(ctx, countCfg.Count)
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
			Sug:          post.PrimaryPost.Id,
			Author:       post.PrimaryPost.Author,
			Title:        post.PrimaryPost.Title,
			Summary:      post.PrimaryPost.Summary,
			CoverImg:     post.PrimaryPost.CoverImg,
			Categories:   post.PrimaryPost.Categories,
			Tags:         post.PrimaryPost.Tags,
			LikeCount:    post.PrimaryPost.LikeCount,
			CommentCount: post.PrimaryPost.CommentCount,
			VisitCount:   post.PrimaryPost.VisitCount,
			StickyWeight: post.PrimaryPost.StickyWeight,
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

func (h *PostHandler) AdminGetPosts(ctx *gin.Context, req request.PageRequest) (pageVO vo.PageVO[vo.AdminPostVO], err error) {
	posts, total, err := h.serv.AdminGetPosts(ctx, dto.PageDTO{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		Field:    req.Field,
		Order:    req.Order,
		Keyword:  req.Keyword,
	})
	if err != nil {
		return vo.PageVO[vo.AdminPostVO]{}, err
	}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.postsToAdminPost(posts)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return
}

func (h *PostHandler) postsToAdminPost(posts []*domain.Post) []vo.AdminPostVO {
	adminPostVOs := make([]vo.AdminPostVO, len(posts))
	for i, post := range posts {
		adminPostVOs[i] = vo.AdminPostVO{
			Id:         post.Id,
			CoverImg:   post.CoverImg,
			Title:      post.Title,
			Summary:    post.Summary,
			Categories: post.Categories,
			Tags:       post.Tags,
			CreateTime: post.CreateTime,
			UpdateTime: post.UpdateTime,
		}
	}
	return adminPostVOs
}
