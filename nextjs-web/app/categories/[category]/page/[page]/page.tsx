import ArticleList from "@/src/components/ArticleList";
import { getPostList } from "@/src/api/posts";
import { getCategoryNameStringByRoute } from "@/src/api/category";
import { getCommonConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { notFound } from "next/navigation";
import type { Metadata } from "next";

export async function generateMetadata({ params }: { params: { category: string; page: string } }): Promise<Metadata> {
  const categoryName = await getCategoryNameStringByRoute(params.category);
  const config = await getCommonConfig();
  return {
    title: `${categoryName} - 分类文章 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: `浏览${categoryName}分类下的全部文章。`,
    openGraph: {
      title: `${categoryName} - 分类文章 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: `浏览${categoryName}分类下的全部文章。`,
      url: process.env.BASE_HOST + `/categories/${params.category}/page/${params.page}`,
      images: config.seo_meta.og_image ? [{ url: process.env.SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function CategoryPageWithPagination({ params, searchParams }: { params: { category: string; page: string }, searchParams: { filter?: string, pageSize?: string } }) {
  const field = (searchParams?.filter as "latest" | "oldest" | "likes") || "latest";
  const pageNumber = Number(params.page || 1);
  const pageSize = Number(searchParams?.pageSize || 10);

  const categoryName = await getCategoryNameStringByRoute(params.category);
  if (!categoryName) return notFound();
  const posts = await getPostList({
    pageNo: pageNumber,
    pageSize,
    sortField: field === "likes" ? "like_count" : "created_at",
    sortOrder: field === "oldest" ? "ASC" : "DESC",
    categories: [categoryName],
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
