import { getPostList, type PostListResponse } from "@/src/api/posts";
import {
  DEFAULT_COMMON_CONFIG,
  DEFAULT_WEBSITE_OWNER_CONFIG,
  getCommonConfig,
} from "@/src/api/config";
import { DEFAULT_WEBSITE_STATS, getWebsiteStats } from "@/src/api/stats";
import type { Metadata } from "next";
import SearchPageClient from "../../SearchPageClient";
import { resolvePublicUrl } from "@/src/utils/publicUrl";

const DEFAULT_POST_LIST: PostListResponse = {
  PageNo: 1,
  PageSize: 10,
  totalPages: 0,
  totalCount: 0,
  list: [],
};

async function settleWithFallback<T>(promise: Promise<T>, fallback: T) {
  try {
    return { data: await promise, failed: false };
  } catch {
    return { data: fallback, failed: true };
  }
}

export async function generateMetadata({
  searchParams,
}: {
  searchParams: Promise<{ keyword?: string }>;
}): Promise<Metadata> {
  const resolvedSearchParams = await searchParams;
  const keyword = resolvedSearchParams?.keyword || "";
  const config = await getCommonConfig().catch(() => DEFAULT_COMMON_CONFIG);
  const siteTitle = config.seo_meta.title || config.website_meta.website_name;
  const title = keyword ? `搜索：${keyword} - ${siteTitle}` : `搜索文章 - ${siteTitle}`;
  const description = keyword
    ? `搜索与“${keyword}”相关的全部文章。`
    : "搜索本站全部文章。";

  return {
    title,
    description,
    openGraph: {
      title,
      description,
      url:
        process.env.BASE_HOST +
        `/search` +
        (keyword ? `?keyword=${encodeURIComponent(keyword)}` : ""),
      images: config.seo_meta.og_image
        ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }]
        : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function SearchPageWithPagination({
  params,
  searchParams,
}: {
  params: Promise<{ page: string }>;
  searchParams: Promise<{ filter?: string; pageSize?: string; keyword?: string }>;
}) {
  const { page } = await params;
  const resolvedSearchParams = await searchParams;
  const field =
    (resolvedSearchParams?.filter as "latest" | "oldest" | "likes") ||
    "latest";
  const pageNumber = Number(page || 1);
  const pageSize = Number(resolvedSearchParams?.pageSize || 10);
  const keyword = resolvedSearchParams?.keyword || "";

  const [posts, config, stats] = await Promise.all([
    settleWithFallback(
      getPostList({
        pageNo: pageNumber,
        pageSize,
        sortField: field === "likes" ? "like_count" : "created_at",
        sortOrder: field === "oldest" ? "ASC" : "DESC",
        keyword,
      }),
      { ...DEFAULT_POST_LIST, PageNo: pageNumber, PageSize: pageSize }
    ),
    settleWithFallback(getCommonConfig(), DEFAULT_COMMON_CONFIG),
    settleWithFallback(getWebsiteStats(), DEFAULT_WEBSITE_STATS),
  ]);

  return (
    <SearchPageClient
      keyword={keyword}
      field={field}
      page={pageNumber}
      pageSize={pageSize}
      list={posts.data.list}
      total={posts.data.totalCount}
      hasError={posts.failed}
      siteOwner={{
        name: config.failed
          ? DEFAULT_WEBSITE_OWNER_CONFIG.website_owner
          : config.data.website_meta.website_owner,
        avatar: config.failed ? "" : config.data.website_meta.website_owner_avatar,
        bio: config.failed
          ? DEFAULT_WEBSITE_OWNER_CONFIG.website_owner_profile
          : config.data.website_meta.website_owner_profile,
        stats: stats.data,
        hasError: config.failed || stats.failed,
      }}
    />
  );
}
