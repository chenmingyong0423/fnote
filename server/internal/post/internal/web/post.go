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
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/chenmingyong0423/go-eventbus"

	"github.com/chenmingyong0423/fnote/server/internal/post/internal/domain"

	"github.com/chenmingyong0423/fnote/server/internal/post/internal/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_like"

	"github.com/chenmingyong0423/fnote/server/internal/website_config"

	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"

	"github.com/chenmingyong0423/gkit/slice"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/vo"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/web/request"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"

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

func NewPostHandler(serv service.IPostService, cfgService website_config.Service, postLikeServ post_like.Service, eventBus *eventbus.EventBus) *PostHandler {
	return &PostHandler{
		serv:         serv,
		cfgService:   cfgService,
		postLikeServ: postLikeServ,
		eventBus:     eventBus,
	}
}

type PostHandler struct {
	serv         service.IPostService
	cfgService   website_config.Service
	postLikeServ post_like.Service
	ipMap        sync.Map
	eventBus     *eventbus.EventBus
}

func (h *PostHandler) RegisterGinRoutes(engine *gin.Engine) {
	group := engine.Group("/posts")
	group.GET("/latest", apiwrap.Wrap(h.GetLatestPosts))
	group.GET("", apiwrap.WrapWithBody(h.GetPosts))
	group.GET("/:id", apiwrap.Wrap(h.GetPostBySug))
	group.POST("/:id/likes", apiwrap.Wrap(h.AddLike))

	adminGroup := engine.Group("/admin-api/posts")
	adminGroup.GET("", apiwrap.WrapWithBody(h.AdminGetPosts))
	adminGroup.GET("/:id", apiwrap.Wrap(h.AdminGetPostById))
	adminGroup.POST("", apiwrap.WrapWithBody(h.AddPost))
	adminGroup.DELETE("/:id", apiwrap.Wrap(h.DeletePost))
	adminGroup.PUT("/:id/display", apiwrap.WrapWithBody(h.UpdatePostIsDisplayed))
	adminGroup.PUT("/:id/comment-allowed", apiwrap.WrapWithBody(h.UpdatePostIsCommentAllowed))
}

func (h *PostHandler) GetLatestPosts(ctx *gin.Context) (*apiwrap.ResponseBody[apiwrap.ListVO[*SummaryPostVO]], error) {
	countCfg, err := h.cfgService.GetFrontPostCount(ctx)
	if err != nil {
		return nil, err
	}
	posts, err := h.serv.GetLatestPosts(ctx, countCfg.Count)
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(apiwrap.NewListVO(h.postsToPostVOs(posts))), nil
}

func (h *PostHandler) postsToPostVOs(posts []*domain.Post) []*SummaryPostVO {
	postVOs := make([]*SummaryPostVO, 0, len(posts))
	for _, post := range posts {
		categories := slice.Map[domain.Category4Post, string](post.PrimaryPost.Categories, func(_ int, c domain.Category4Post) string {
			return c.Name
		})
		tags := slice.Map[domain.Tag4Post, string](post.PrimaryPost.Tags, func(_ int, t domain.Tag4Post) string {
			return t.Name
		})
		postVOs = append(postVOs, &SummaryPostVO{
			Sug:          post.PrimaryPost.Id,
			Author:       post.PrimaryPost.Author,
			Title:        post.PrimaryPost.Title,
			Summary:      post.PrimaryPost.Summary,
			CoverImg:     post.PrimaryPost.CoverImg,
			Categories:   categories,
			Tags:         tags,
			LikeCount:    post.PrimaryPost.LikeCount,
			CommentCount: post.PrimaryPost.CommentCount,
			VisitCount:   post.PrimaryPost.VisitCount,
			StickyWeight: post.PrimaryPost.StickyWeight,
			CreateTime:   post.PrimaryPost.CreatedAt,
		})
	}
	return postVOs
}

func (h *PostHandler) GetPosts(ctx *gin.Context, req *domain.PostRequest) (*apiwrap.ResponseBody[apiwrap.PageVO[*SummaryPostVO]], error) {
	req.ValidateAndSetDefault()
	posts, cnt, err := h.serv.GetPosts(ctx, req)
	if err != nil {
		return nil, err
	}
	pageVO := apiwrap.PageVO[*SummaryPostVO]{
		Page: apiwrap.Page{
			PageNo:   req.Page.PageNo,
			PageSize: req.Page.PageSize,
		},
		List: h.postsToPostVOs(posts),
	}
	pageVO.SetTotalCountAndCalculateTotalPages(cnt)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *PostHandler) GetPostBySug(ctx *gin.Context) (*apiwrap.ResponseBody[domain.DetailPostVO], error) {
	sug := ctx.Param("id")
	post, err := h.serv.GetPunishedPostById(ctx, sug)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "The postId does not exist.")
		}
		return nil, err
	}
	// 查询点赞状态
	liked, err := h.postLikeServ.GetLikeStatus(ctx, post.PrimaryPost.Id, ctx.ClientIP())
	if err != nil {
		return nil, err
	}
	return apiwrap.SuccessResponseWithData(domain.DetailPostVO{
		PrimaryPost: post.PrimaryPost,
		ExtraPost:   post.ExtraPost,
		IsLiked:     liked,
	}), nil
}

func (h *PostHandler) AddLike(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	ip := ctx.ClientIP()
	if ip == "" {
		return nil, apiwrap.NewErrorResponseBody(http.StatusBadRequest, "Ip is empty.")
	}
	postId := ctx.Param("id")
	key := fmt.Sprintf("%s:%s", postId, ip)
	_, isExist := h.ipMap.LoadOrStore(key, struct{}{})
	if !isExist {
		defer h.ipMap.Delete(key)
		var likePostEvent = domain.LikePostEvent{PostId: postId}
		marshal, err := json.Marshal(likePostEvent)
		if err != nil {
			return nil, err
		}

		id, err := h.postLikeServ.Add(ctx, post_like.PostLike{
			PostId:    postId,
			Ip:        ip,
			UserAgent: ctx.GetHeader("User-Agent"),
		})
		if err != nil {
			// 已点过赞
			if mongo.IsDuplicateKeyError(err) {
				slog.WarnContext(ctx, "post like", fmt.Sprintf("%+v", err), nil)
				return apiwrap.SuccessResponse(), nil
			}
			return nil, err
		}
		// 文章点赞数自增
		err = h.serv.IncreasePostLikeCount(ctx, postId)
		if err != nil {
			err = h.postLikeServ.DeleteById(ctx, id)
			if err != nil {
				return nil, err
			}
		}

		h.eventBus.Publish("post-like", eventbus.Event{Payload: marshal})
	}
	return apiwrap.SuccessResponse(), nil
}

func (h *PostHandler) AdminGetPosts(ctx *gin.Context, req PageRequest) (*apiwrap.ResponseBody[vo.PageVO[vo.AdminPostVO]], error) {
	posts, total, err := h.serv.AdminGetPosts(ctx, domain.Page{
		Size:           req.PageSize,
		Skip:           (req.PageNo - 1) * req.PageSize,
		Field:          req.Field,
		Order:          req.Order,
		Keyword:        req.Keyword,
		CategoryFilter: req.CategoryFilter,
		TagFilter:      req.TagFilter,
	})
	if err != nil {
		return nil, err
	}
	pageVO := vo.PageVO[vo.AdminPostVO]{}
	pageVO.PageNo = req.PageNo
	pageVO.PageSize = req.PageSize
	pageVO.List = h.postsToAdminPost(posts)
	pageVO.SetTotalCountAndCalculateTotalPages(total)
	return apiwrap.SuccessResponseWithData(pageVO), nil
}

func (h *PostHandler) postsToAdminPost(posts []*domain.Post) []vo.AdminPostVO {
	adminPostVOs := make([]vo.AdminPostVO, len(posts))
	for i, post := range posts {
		categories := slice.Map[domain.Category4Post, vo.Category4Post](post.PrimaryPost.Categories, func(_ int, c domain.Category4Post) vo.Category4Post {
			return vo.Category4Post{
				Id:   c.Id,
				Name: c.Name,
			}
		})
		tags := slice.Map[domain.Tag4Post, vo.Tag4Post](post.PrimaryPost.Tags, func(_ int, t domain.Tag4Post) vo.Tag4Post {
			return vo.Tag4Post{
				Id:   t.Id,
				Name: t.Name,
			}
		})
		adminPostVOs[i] = vo.AdminPostVO{
			Id:               post.Id,
			CoverImg:         post.CoverImg,
			Title:            post.Title,
			Summary:          post.Summary,
			Categories:       categories,
			Tags:             tags,
			IsDisplayed:      post.IsDisplayed,
			IsCommentAllowed: post.IsCommentAllowed,
			CreateTime:       post.CreatedAt,
			UpdateTime:       post.UpdateTime,
		}
	}
	return adminPostVOs
}

func (h *PostHandler) AddPost(ctx *gin.Context, req request.PostReq) (*apiwrap.ResponseBody[any], error) {
	existPost, err := h.serv.AdminGetPostById(ctx, req.Id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if existPost != nil {
		return nil, apiwrap.NewErrorResponseBody(http.StatusConflict, "The postId already exists.")
	}
	categories := slice.Map[request.Category4Post, domain.Category4Post](req.Categories, func(_ int, c request.Category4Post) domain.Category4Post {
		return domain.Category4Post{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	tags := slice.Map[request.Tag4Post, domain.Tag4Post](req.Tags, func(_ int, t request.Tag4Post) domain.Tag4Post {
		return domain.Tag4Post{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	return apiwrap.SuccessResponse(), h.serv.AddPost(ctx, &domain.Post{
		PrimaryPost: domain.PrimaryPost{
			Id:           req.Id,
			Author:       req.Author,
			Title:        req.Title,
			Summary:      req.Summary,
			CoverImg:     req.CoverImg,
			Categories:   categories,
			Tags:         tags,
			StickyWeight: req.StickyWeight,
		},
		ExtraPost: domain.ExtraPost{
			Content:          req.Content,
			MetaDescription:  req.MetaDescription,
			MetaKeywords:     req.MetaKeywords,
			WordCount:        req.WordCount,
			IsDisplayed:      req.IsDisplayed,
			IsCommentAllowed: req.IsCommentAllowed,
		},
	})
}

func (h *PostHandler) DeletePost(ctx *gin.Context) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.DeletePost(ctx, ctx.Param("id"))
}

func (h *PostHandler) AdminGetPostById(ctx *gin.Context) (*apiwrap.ResponseBody[vo.PostDetailVO], error) {
	post, err := h.serv.AdminGetPostById(ctx, ctx.Param("id"))
	if err != nil {
		return nil, err
	}
	categories := slice.Map[domain.Category4Post, vo.Category4Post](post.PrimaryPost.Categories, func(_ int, c domain.Category4Post) vo.Category4Post {
		return vo.Category4Post{
			Id:   c.Id,
			Name: c.Name,
		}
	})
	tags := slice.Map[domain.Tag4Post, vo.Tag4Post](post.PrimaryPost.Tags, func(_ int, t domain.Tag4Post) vo.Tag4Post {
		return vo.Tag4Post{
			Id:   t.Id,
			Name: t.Name,
		}
	})
	return apiwrap.SuccessResponseWithData(vo.PostDetailVO{
		Id:               post.Id,
		Author:           post.Author,
		Title:            post.Title,
		Summary:          post.Summary,
		Content:          post.Content,
		CoverImg:         post.CoverImg,
		Categories:       categories,
		Tags:             tags,
		IsDisplayed:      post.IsDisplayed,
		StickyWeight:     post.StickyWeight,
		MetaDescription:  post.MetaDescription,
		MetaKeywords:     post.MetaKeywords,
		IsCommentAllowed: post.IsCommentAllowed,
	}), nil
}

func (h *PostHandler) UpdatePostIsDisplayed(ctx *gin.Context, req request.PostDisplayReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdatePostIsDisplayed(ctx, ctx.Param("id"), req.IsDisplayed)
}

func (h *PostHandler) UpdatePostIsCommentAllowed(ctx *gin.Context, req request.PostCommentAllowedReq) (*apiwrap.ResponseBody[any], error) {
	return apiwrap.SuccessResponse(), h.serv.UpdatePostIsCommentAllowed(ctx, ctx.Param("id"), req.IsCommentAllowed)
}
