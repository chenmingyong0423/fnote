import { cache } from "react";
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import ArticleList from "./ArticleList";
import { getPostList } from "@/src/api/posts";
import { getCategoryNameStringByRoute } from "@/src/api/category";
import { getTagNameByRoute } from "@/src/api/tags";
import { getCommonConfig, getWebsiteOwnerConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { resolvePublicUrl } from "@/src/utils/publicUrl";
import { getSiteUrl } from "@/src/utils/seo";
import { ensurePageExists, resolvePagination, type ListSearchParams } from "@/src/utils/pagination";

type Kind = "categories" | "tags";

const loadArchive = cache(async (kind: Kind, slug: string, rawPage: string | undefined, query: ListSearchParams) => {
  const pagination = resolvePagination(`/${kind}/${encodeURIComponent(slug)}`, rawPage, query);
  const name = kind === "categories"
    ? await getCategoryNameStringByRoute(slug)
    : await getTagNameByRoute(slug);
  if (!name) notFound();
  const posts = await getPostList({
    pageNo: pagination.page,
    pageSize: pagination.pageSize,
    sortField: pagination.field === "likes" ? "like_count" : "created_at",
    sortOrder: pagination.field === "oldest" ? "ASC" : "DESC",
    ...(kind === "categories" ? { categories: [name] } : { tags: [name] }),
  });
  ensurePageExists(pagination.page, posts.totalPages);
  return { ...pagination, name, posts, label: kind === "categories" ? "分类" : "标签" };
});

export async function archiveMetadata(kind: Kind, slug: string, rawPage: string | undefined, query: ListSearchParams): Promise<Metadata> {
  const archive = await loadArchive(kind, slug, rawPage, query);
  const config = await getCommonConfig();
  const pageTitle = `${archive.name} - ${archive.label}文章${archive.page > 1 ? ` - 第 ${archive.page} 页` : ""}`;
  const title = `${pageTitle} - ${config.seo_meta.title || config.website_meta.website_name}`;
  const description = `浏览${archive.name}${archive.label}下的文章${archive.page > 1 ? `，第 ${archive.page} 页` : ""}。`;
  return {
    title,
    description,
    alternates: { canonical: archive.pathname },
    ...(archive.posts.totalCount === 0 ? { robots: { index: false, follow: true } } : {}),
    openGraph: {
      title: `${pageTitle} - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description,
      url: new URL(archive.pathname, getSiteUrl()).toString(),
      images: config.seo_meta.og_image ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export async function renderArchive(kind: Kind, slug: string, rawPage: string | undefined, query: ListSearchParams) {
  const archive = await loadArchive(kind, slug, rawPage, query);
  const [owner, stats] = await Promise.all([getWebsiteOwnerConfig(), getWebsiteStats()]);
  return <ArticleList
    list={archive.posts.list}
    total={archive.posts.totalCount}
    pageHeading={`${archive.name}${archive.label}文章${archive.page > 1 ? ` - 第 ${archive.page} 页` : ""}`}
    siteOwner={{ name: owner.website_owner, avatar: owner.website_owner_avatar, bio: owner.website_owner_profile, socialInfo: owner.social_info_list, stats }}
    hideSiteOwnerOnMobile
    field={archive.field}
    currentPage={archive.page}
    pageSize={archive.pageSize}
  />;
}
