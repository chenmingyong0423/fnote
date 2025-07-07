import ArticleList from "@/src/components/ArticleList";
import { getPostList } from "@/src/api/posts";
import { getTagNameByRoute } from "@/src/api/tags";
import { getCommonConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { notFound } from "next/navigation";
import type { Metadata } from "next";

export async function generateMetadata({ params }: { params: { tag: string; page: string } }): Promise<Metadata> {
  const tagName = await getTagNameByRoute(params.tag);
  const config = await getCommonConfig();
  return {
    title: `${tagName} - 标签文章 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: `浏览${tagName}标签下的全部文章。`,
    openGraph: {
      title: `${tagName} - 标签文章 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: `浏览${tagName}标签下的全部文章。`,
      url: process.env.BASE_HOST + `/tags/${params.tag}/page/${params.page}`,
      images: config.seo_meta.og_image ? [{ url: process.env.SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function TagPageWithPagination({ params, searchParams }: { params: { tag: string; page: string }, searchParams: { filter?: string, pageSize?: string } }) {
  const field = (searchParams?.filter as "latest" | "oldest" | "likes") || "latest";
  const pageNumber = Number(params.page || 1);
  const pageSize = Number(searchParams?.pageSize || 10);

  const tagName = await getTagNameByRoute(params.tag);
  if (!tagName) return notFound();
  const posts = await getPostList({
    pageNo: pageNumber,
    pageSize,
    sortField: field === "likes" ? "like_count" : "created_at",
    sortOrder: field === "oldest" ? "ASC" : "DESC",
    tags: [tagName],
  });
  const config = await getCommonConfig();
  const stats = await getWebsiteStats();

  return (
    <ArticleList
      list={posts.list}
      total={posts.totalCount}
      siteOwner={{
        name: config.website_meta.website_owner,
        avatar: config.website_meta.website_owner_avatar,
        bio: config.website_meta.website_owner_profile,
        stats,
      }}
      field={field}
      currentPage={pageNumber}
      pageSize={pageSize}
    />
  );
}
