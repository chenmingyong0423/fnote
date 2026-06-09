import ArticleList from "@/src/components/ArticleList";
import { getPostList } from "@/src/api/posts";
import { getCategoryNameStringByRoute } from "@/src/api/category";
import { getCommonConfig } from "@/src/api/config";
import { getWebsiteStats } from "@/src/api/stats";
import { notFound } from "next/navigation";
import type { Metadata } from "next";

export async function generateMetadata({ params }: { params: Promise<{ category: string }> }): Promise<Metadata> {
  const { category } = await params;
  const categoryName = await getCategoryNameStringByRoute(category);
  const config = await getCommonConfig();
  return {
    title: `${categoryName} - 分类文章 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: `浏览${categoryName}分类下的全部文章。`,
    openGraph: {
      title: `${categoryName} - 分类文章 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: `浏览${categoryName}分类下的全部文章。`,
      url: process.env.BASE_HOST + `/categories/${category}`,
      images: config.seo_meta.og_image ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function CategoryPage({ params, searchParams }: { params: Promise<{ category: string }>; searchParams: Promise<{ filter?: string, pageSize?: string, page?: string }> }) {
  const { category } = await params;
  const resolvedSearchParams = await searchParams;
  const field = (resolvedSearchParams?.filter as "latest" | "oldest" | "likes") || "latest";
  const pageSize = Number(resolvedSearchParams?.pageSize || 10);
  const page = Number(resolvedSearchParams?.page || 1);

  const categoryName = await getCategoryNameStringByRoute(category);
  if (!categoryName) return notFound();
  const posts = await getPostList({
    pageNo: page,
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
      currentPage={page}
      pageSize={pageSize}
    />
  );
}
