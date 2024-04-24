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
	"time"

	"github.com/chenmingyong0423/fnote/server/internal/pkg/domain"
	apiwrap "github.com/chenmingyong0423/fnote/server/internal/pkg/web/wrap"
	postServ "github.com/chenmingyong0423/fnote/server/internal/post/service"
	"github.com/chenmingyong0423/fnote/server/internal/post_draft"
	"github.com/chenmingyong0423/gkit/slice"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
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
	var (
		postDraft *post_draft.PostDraft
		post      *domain.Post
		err       error
	)
	id := ctx.Param("id")
	postDraft, err = h.postDraftServ.GetPostDraftById(ctx, id)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}
	if postDraft == nil {
		// 查询文章是否存在
		post, err = h.postServ.AdminGetPostById(ctx, id)
		if err != nil {
			return nil, err
		}
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		if post == nil {
			return nil, apiwrap.NewErrorResponseBody(404, "post draft not found")
		}
		// 保存草稿
		createdAt := time.Now().Local().Unix()
		_, err = h.postDraftServ.SavePostDraft(ctx, post_draft.PostDraft{
			Id:       post.Id,
			Author:   post.Author,
			Title:    post.Title,
			Summary:  post.Summary,
			CoverImg: post.CoverImg,
			Categories: slice.Map(post.Categories, func(idx int, c domain.Category4Post) post_draft.Category4PostDraft {
				return post_draft.Category4PostDraft{
					Id:   c.Id,
					Name: c.Name,
				}
			}),
			Tags: slice.Map(post.Tags, func(idx int, t domain.Tag4Post) post_draft.Tag4PostDraft {
				return post_draft.Tag4PostDraft{
					Id:   t.Id,
					Name: t.Name,
				}
			}),
			StickyWeight:     post.StickyWeight,
			Content:          post.Content,
			MetaDescription:  post.MetaDescription,
			MetaKeywords:     post.MetaKeywords,
			WordCount:        post.WordCount,
			IsDisplayed:      post.IsDisplayed,
			IsCommentAllowed: post.IsCommentAllowed,
			CreatedAt:        createdAt,
		})
		if err != nil {
			return nil, err
		}
		draftVO := h.postToPostDraftVO(post)
		draftVO.CreatedAt = createdAt
		return apiwrap.SuccessResponseWithData(draftVO), nil
	}
	return apiwrap.SuccessResponseWithData(h.postDraftToPostDraftVO(postDraft)), nil
}

func (h *AggregatePostHandler) postDraftToPostDraftVO(postDraft *post_draft.PostDraft) *PostDraftVO {
	return &PostDraftVO{
		Id:       postDraft.Id,
		Author:   postDraft.Author,
		Title:    postDraft.Title,
		Summary:  postDraft.Summary,
		Content:  postDraft.Content,
		CoverImg: postDraft.CoverImg,
		Categories: slice.Map(postDraft.Categories, func(idx int, c post_draft.Category4PostDraft) Category4PostDraft {
			return Category4PostDraft{
				Id:   c.Id,
				Name: c.Name,
			}
		}),
		Tags: slice.Map(postDraft.Tags, func(idx int, c post_draft.Tag4PostDraft) Tag4PostDraft {
			return Tag4PostDraft{
				Id:   c.Id,
				Name: c.Name,
			}
		}),
		StickyWeight:     postDraft.StickyWeight,
		IsDisplayed:      postDraft.IsDisplayed,
		MetaDescription:  postDraft.MetaDescription,
		MetaKeywords:     postDraft.MetaKeywords,
		WordCount:        postDraft.WordCount,
		IsCommentAllowed: postDraft.IsCommentAllowed,
		CreatedAt:        postDraft.CreatedAt,
	}
}

func (h *AggregatePostHandler) postToPostDraftVO(post *domain.Post) *PostDraftVO {
	return &PostDraftVO{
		Id:       post.Id,
		Author:   post.Author,
		Title:    post.Title,
		Summary:  post.Summary,
		Content:  post.Content,
		CoverImg: post.CoverImg,
		Categories: slice.Map(post.Categories, func(idx int, c domain.Category4Post) Category4PostDraft {
			return Category4PostDraft{
				Id:   c.Id,
				Name: c.Name,
			}
		}),
		Tags: slice.Map(post.Tags, func(idx int, c domain.Tag4Post) Tag4PostDraft {
			return Tag4PostDraft{
				Id:   c.Id,
				Name: c.Name,
			}
		}),
		StickyWeight:     post.StickyWeight,
		IsDisplayed:      post.IsDisplayed,
		MetaDescription:  post.MetaDescription,
		MetaKeywords:     post.MetaKeywords,
		WordCount:        post.WordCount,
		IsCommentAllowed: post.IsCommentAllowed,
	}
}
