package service

import (
	"context"
	"encoding/json"
	"io"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/chenmingyong0423/fnote/server/internal/ai_assistant/internal/domain"
	"github.com/chenmingyong0423/fnote/server/internal/post"
	"github.com/chenmingyong0423/fnote/server/internal/website_config"
	"github.com/cloudwego/eino/schema"
)

type ToolService struct {
	posts   post.Service
	configs website_config.Service
}

func NewToolService(posts post.Service, configs website_config.Service) *ToolService {
	return &ToolService{posts, configs}
}

func toolDefinitions(enabled []string) []*schema.ToolInfo {
	definitions := map[string]*schema.ToolInfo{
		"search_posts":       {Name: "search_posts", Desc: "按标题关键词搜索公开文章。返回文章 ID、摘要及可引用的链接。", ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{"query": {Type: schema.String, Desc: "标题关键词，最多 100 字", Required: true}})},
		"read_post":          {Name: "read_post", Desc: "根据公开文章 ID 读取正文。内容过长时返回截取的正文。", ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{"post_id": {Type: schema.String, Desc: "搜索结果中的文章 ID", Required: true}})},
		"get_blog_overview":  {Name: "get_blog_overview", Desc: "获取博客公开名称、简介、公开文章数量和最近文章。", ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{})},
		"get_author_profile": {Name: "get_author_profile", Desc: "获取博主公开姓名和个人介绍。", ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{})},
	}
	result := []*schema.ToolInfo{}
	for _, name := range enabled {
		if tool := definitions[name]; tool != nil {
			result = append(result, tool)
		}
	}
	return result
}

func (s *ToolService) Execute(ctx context.Context, name, args string) (string, []domain.Reference) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fail := func(message string) (string, []domain.Reference) {
		b, _ := json.Marshal(map[string]string{"error": message})
		return string(b), nil
	}
	if len(args) > 8192 {
		return fail("工具参数过长")
	}
	var input struct {
		Query  string `json:"query"`
		PostID string `json:"post_id"`
	}
	if strings.TrimSpace(args) == "" {
		args = "{}"
	}
	decoder := json.NewDecoder(strings.NewReader(args))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&input) != nil || decoder.Decode(new(any)) != io.EOF {
		return fail("工具参数格式不正确")
	}
	var data any
	refs := []domain.Reference{}
	switch name {
	case "search_posts", "get_blog_overview":
		query := strings.TrimSpace(input.Query)
		if utf8.RuneCountInString(query) > 100 || (name == "search_posts" && query == "") {
			return fail("请提供 1 至 100 字的搜索关键词")
		}
		if name == "get_blog_overview" {
			query = ""
		}
		posts, total, err := s.posts.SearchPublicPosts(ctx, query, 6)
		if err != nil {
			return fail("暂时无法搜索公开文章")
		}
		articles := []map[string]any{}
		for _, p := range posts {
			if p == nil || !p.IsDisplayed {
				continue
			}
			ref := articleReference(p)
			refs = append(refs, ref)
			articles = append(articles, map[string]any{"id": p.Id, "title": limitText(p.Title, 200), "summary": limitText(p.Summary, 400), "url": ref.Href})
		}
		data = map[string]any{"articles": articles, "total": total}
		if name == "get_blog_overview" {
			cfg, err := s.configs.GetWebSiteConfig(ctx)
			if err != nil {
				return fail("暂时无法获取博客概况")
			}
			data = map[string]any{"name": cfg.WebsiteName, "description": limitText(cfg.WebsiteOwnerProfile, 1500), "public_post_count": total, "latest_posts": articles}
		}
	case "read_post":
		if len(input.PostID) == 0 || len(input.PostID) > 200 {
			return fail("文章 ID 不正确")
		}
		p, err := s.posts.GetPublicPostById(ctx, input.PostID)
		if err != nil || p == nil || !p.IsDisplayed {
			return fail("文章不存在或未公开")
		}
		ref := articleReference(p)
		refs = append(refs, ref)
		data = map[string]any{"id": p.Id, "title": limitText(p.Title, 200), "url": ref.Href, "content": limitText(p.Content, 15000), "truncated": utf8.RuneCountInString(p.Content) > 15000}
	case "get_author_profile":
		cfg, err := s.configs.GetWebSiteConfig(ctx)
		if err != nil {
			return fail("暂时无法获取博主公开信息")
		}
		data = map[string]string{"name": cfg.WebsiteOwner, "profile": limitText(cfg.WebsiteOwnerProfile, 1500)}
	default:
		return fail("工具不可用")
	}
	if ctx.Err() != nil {
		return fail("工具调用超时")
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fail("工具结果暂不可用")
	}
	return string(b), refs
}

func articleReference(p *post.Post) domain.Reference {
	return domain.Reference{Title: limitText(p.Title, 200), Href: "/posts/" + url.PathEscape(p.Id)}
}
func limitText(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n]) + "…"
	}
	return s
}
