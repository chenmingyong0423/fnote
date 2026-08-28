import ArticleList from "@/src/components/ArticleList";
import { getPostList } from "@/src/api/posts";
import { getTagNameByRoute } from "@/src/api/tags";
import { getCommonConfig, getWebsiteOwnerConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { notFound } from "next/navigation";
import type { Metadata } from "next";
import { resolvePublicUrl } from "@/src/utils/publicUrl";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ tag: string }>;
}): Promise<Metadata> {
  const { tag } = await params;
  const tagName = await getTagNameByRoute(tag);
  const config = await getCommonConfig();
  return {
    title: `${tagName} - 标签文章 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: `浏览${tagName}标签下的全部文章。`,
    openGraph: {
      title: `${tagName} - 标签文章 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: `浏览${tagName}标签下的全部文章。`,
      url: process.env.BASE_HOST + `/tags/${tag}`,
      images: config.seo_meta.og_image
        ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }]
        : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function TagPage({
  params,
  searchParams,
}: {
  params: Promise<{ tag: string }>;
  searchParams: Promise<{ filter?: string; pageSize?: string }>;
}) {
  const { tag } = await params;
  const resolvedSearchParams = await searchParams;
  const field =
    (resolvedSearchParams?.filter as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(resolvedSearchParams?.pageSize || 10);
  const page = 1;

  const tagName = await getTagNameByRoute(tag);
  if (!tagName) return notFound();
  const posts = await getPostList({
    pageNo: page,
    pageSize,
    sortField: field === "likes" ? "like_count" : "created_at",
    sortOrder: field === "oldest" ? "ASC" : "DESC",
    tags: [tagName],
  });
  const owner = await getWebsiteOwnerConfig();
  const stats = await getWebsiteStats();

  return (
    <ArticleList
      list={posts.list}
      total={posts.totalCount}
      pageHeading={`${tagName}标签文章`}
      siteOwner={{
        name: owner.website_owner,
        avatar: owner.website_owner_avatar,
        bio: owner.website_owner_profile,
        socialInfo: owner.social_info_list,
        stats,
      }}
      hideSiteOwnerOnMobile
      field={field}
      currentPage={page}
      pageSize={pageSize}
    />
  );
}
