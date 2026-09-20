import { cache } from "react";
import type { Metadata } from "next";
import { notFound } from "next/navigation";
import ArticleList from "./ArticleList";
import { getPostList } from "@/src/api/posts";
import { getCategoryNameStringByRoute } from "@/src/api/category";
import { getTagNameByRoute } from "@/src/api/tags";
import { getCommonConfig, getWebsiteOwnerConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { buildPageMetadata, buildCollectionJsonLd, serializeJsonLd } from "@/src/utils/seo";
import Breadcrumbs from "./Breadcrumbs";
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
  const description = `浏览${archive.name}${archive.label}下的文章${archive.page > 1 ? `，第 ${archive.page} 页` : ""}。`;
  return buildPageMetadata(config, {
    title: pageTitle,
    description,
    pathname: archive.pathname,
    noindex: archive.posts.totalCount === 0 || archive.field !== "latest" || archive.pageSize !== 10,
  });
}

export async function renderArchive(kind: Kind, slug: string, rawPage: string | undefined, query: ListSearchParams) {
  const archive = await loadArchive(kind, slug, rawPage, query);
  const [owner, stats] = await Promise.all([getWebsiteOwnerConfig(), getWebsiteStats()]);
  const heading = `${archive.name}${archive.label}文章${archive.page > 1 ? ` - 第 ${archive.page} 页` : ""}`;
  return <>
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      <Breadcrumbs items={[
        { name: "首页", pathname: "/" },
        { name: "全部分类与标签", pathname: "/navigation" },
        ...(archive.page > 1 ? [{ name: `${archive.name}${archive.label}文章`, pathname: `/${kind}/${encodeURIComponent(slug)}` }] : []),
        { name: heading, pathname: archive.pathname },
      ]} />
    </div>
    <script type="application/ld+json" dangerouslySetInnerHTML={{ __html: serializeJsonLd(buildCollectionJsonLd(heading, archive.pathname, archive.posts.list)) }} />
    <ArticleList
    list={archive.posts.list}
    total={archive.posts.totalCount}
    pageHeading={heading}
    siteOwner={{ name: owner.website_owner, avatar: owner.website_owner_avatar, bio: owner.website_owner_profile, socialInfo: owner.social_info_list, stats }}
    hideSiteOwnerOnMobile
    field={archive.field}
    currentPage={archive.page}
    pageSize={archive.pageSize}
  /></>;
}
