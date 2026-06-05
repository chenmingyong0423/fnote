import PostDetail from "@/src/components/PostDetail";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";
import { getCommonConfig } from "@/src/api/config";
import type { Metadata } from "next";

export async function generateMetadata({ params }: { params: Promise<{ id: string }> }): Promise<Metadata> {
  const { id } = await params;
  const res = await getPostDetail(id);
  if (res.code !== 0 || !res.data) return {};
  const post = res.data;
  const config = await getCommonConfig();
  return {
    title: `${post.title} - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: post.meta_description || post.summary,
    keywords: post.meta_keywords,
    openGraph: {
      title: `${post.title} - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: post.meta_description || post.summary,
      url: process.env.BASE_HOST + `/posts/${id}`,
      images: post.cover_img ? [{ url: post.cover_img }] : (config.seo_meta.og_image ? [{ url: process.env.SERVER_HOST + config.seo_meta.og_image }] : undefined),
      siteName: config.website_meta.website_name,
      type: "article",
    },
  };
}

type Params = Promise<{
    id: string;
  }>

export default async function PostDetailPage({ params }: { params: Params }) {
  const { id } = await params
  const res = await getPostDetail(id);
  if (res.code !== 0 || !res.data) return notFound();
  return <PostDetail post={res.data} />;
}
