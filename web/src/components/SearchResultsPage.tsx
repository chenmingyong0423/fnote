import { cache } from "react";
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { getPostList } from "@/src/api/posts";
import { DEFAULT_COMMON_CONFIG, DEFAULT_WEBSITE_OWNER_CONFIG, getCommonConfig, getWebsiteOwnerConfig } from "@/src/api/config";
import { DEFAULT_WEBSITE_STATS, getWebsiteStats } from "@/src/api/stats";
import { ensurePageExists, resolvePagination, type ListSearchParams } from "@/src/utils/pagination";
import { buildPageMetadata } from "@/src/utils/seo";
import SearchPageClient from "@/app/search/SearchPageClient";

const loadSearch = cache(async (rawPage: string | undefined, query: ListSearchParams) => {
  const pagination = resolvePagination("/search", rawPage, query);
  if (Array.isArray(query.keyword)) notFound();
  const keyword = query.keyword || "";
  const posts = await getPostList({
    pageNo: pagination.page,
    pageSize: pagination.pageSize,
    sortField: pagination.field === "likes" ? "like_count" : "created_at",
    sortOrder: pagination.field === "oldest" ? "ASC" : "DESC",
    keyword,
  });
  ensurePageExists(pagination.page, posts.totalPages);
  return { ...pagination, keyword, posts };
});

export async function searchMetadata(rawPage: string | undefined, query: ListSearchParams): Promise<Metadata> {
  const search = await loadSearch(rawPage, query);
  const config = await getCommonConfig().catch(() => DEFAULT_COMMON_CONFIG);
  const pageTitle = `${search.keyword ? `搜索：${search.keyword}` : "搜索文章"}${search.page > 1 ? ` - 第 ${search.page} 页` : ""}`;
  const description = search.keyword ? `搜索与“${search.keyword}”相关的文章。` : "搜索本站全部文章。";
  return buildPageMetadata(config, {
    title: pageTitle,
    description,
    noindex: true,
    pathname: search.pathname,
  });
}

export async function renderSearch(rawPage: string | undefined, query: ListSearchParams) {
  const search = await loadSearch(rawPage, query);
  const [owner, stats] = await Promise.all([
    getWebsiteOwnerConfig().then(data => ({ data, failed: false })).catch(() => ({ data: DEFAULT_WEBSITE_OWNER_CONFIG, failed: true })),
    getWebsiteStats().then(data => ({ data, failed: false })).catch(() => ({ data: DEFAULT_WEBSITE_STATS, failed: true })),
  ]);
  return <SearchPageClient
    keyword={search.keyword}
    field={search.field}
    page={search.page}
    pageSize={search.pageSize}
    list={search.posts.list}
    total={search.posts.totalCount}
    siteOwner={{ name: owner.data.website_owner, avatar: owner.data.website_owner_avatar, bio: owner.data.website_owner_profile, socialInfo: owner.data.social_info_list, stats: stats.data, hasError: owner.failed || stats.failed }}
  />;
}
