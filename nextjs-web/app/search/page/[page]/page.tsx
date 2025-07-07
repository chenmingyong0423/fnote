import { getPostList } from "@/src/api/posts";
import { getCommonConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import type { Metadata } from "next";
import SearchPageClient from "../../SearchPageClient";

export async function generateMetadata({ searchParams }: { searchParams: { keyword?: string } }): Promise<Metadata> {
  const keyword = searchParams?.keyword || "";
  const config = await getCommonConfig();
  return {
    title: keyword ? `搜索：${keyword} - ${config.seo_meta.title || config.website_meta.website_name}` : `搜索文章 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: keyword ? `搜索与“${keyword}”相关的全部文章。` : "搜索本站全部文章。",
    openGraph: {
      title: keyword ? `搜索：${keyword} - ${config.seo_meta.og_title || config.website_meta.website_name}` : `搜索文章 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: keyword ? `搜索与“${keyword}”相关的全部文章。` : "搜索本站全部文章。",
      url: process.env.BASE_HOST + `/search` + (keyword ? `?keyword=${encodeURIComponent(keyword)}` : ""),
      images: config.seo_meta.og_image ? [{ url: process.env.SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function SearchPageWithPagination({ searchParams }: { searchParams: { filter?: string, pageSize?: string, page?: string, keyword?: string } }) {
  const field = (searchParams?.filter as "latest" | "oldest" | "likes") || "latest";
  const pageNumber = Number(searchParams?.page || 1);
  const pageSize = Number(searchParams?.pageSize || 10);
  const keyword = searchParams?.keyword || "";

  const posts = await getPostList({
    pageNo: pageNumber,
    pageSize,
    sortField: field === "likes" ? "like_count" : "created_at",
    sortOrder: field === "oldest" ? "ASC" : "DESC",
    keyword,
  });
  const config = await getCommonConfig();
  const stats = await getWebsiteStats();

  return (
    <SearchPageClient
      keyword={keyword}
      field={field}
      page={pageNumber}
      pageSize={pageSize}
      list={posts.list}
      total={posts.totalCount}
      siteOwner={{
        name: config.website_meta.website_owner,
        avatar: config.website_meta.website_owner_avatar,
        bio: config.website_meta.website_owner_profile,
        stats,
      }}
    />
  );
}
