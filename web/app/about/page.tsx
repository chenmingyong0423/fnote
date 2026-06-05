import PostDetail from "@/src/components/PostDetail";
import { getPostDetail } from "@/src/api/posts";
import { notFound } from "next/navigation";
import type {Metadata} from "next";
import {getCommonConfig} from "@/src/api/config";

export async function generateMetadata(): Promise<Metadata> {
  const res = await getPostDetail("about-me");
  if (res.code !== 0 || !res.data) return {};
  const post = res.data;
  const config = await getCommonConfig();
  return {
    title: `${post.title} - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: post.meta_description || post.summary,
    keywords: post.meta_keywords || config.seo_meta.keywords,
    openGraph: {
      title: `${post.title} - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: post.meta_description || post.summary,
      url: process.env.BASE_HOST + `/posts/about-me`,
      images: post.cover_img ? [{ url: post.cover_img }] : undefined,
      siteName: config.website_meta.website_name,
      type: "article",
    },
  };
}

export default async function AboutPage() {
  const res = await getPostDetail("about-me");
  if (res.code !== 0 || !res.data) return notFound();
  return <PostDetail post={res.data} />;
}