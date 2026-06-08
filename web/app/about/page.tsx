import PostDetail from "@/src/components/PostDetail";
import { getPostDetailOrNull } from "@/src/api/posts";
import type {Metadata} from "next";
import {getCommonConfig} from "@/src/api/config";

export async function generateMetadata(): Promise<Metadata> {
  const post = await getPostDetailOrNull("about-me");
  const config = await getCommonConfig().catch(() => null);
  const siteTitle = config?.seo_meta.title || config?.website_meta.website_name || "网站";
  if (!post) {
    return {
      title: `关于 - ${siteTitle}`,
      description: "关于页面暂未配置",
      openGraph: {
        title: `关于 - ${config?.seo_meta.og_title || siteTitle}`,
        description: "关于页面暂未配置",
        url: process.env.BASE_HOST + "/about",
        images: config?.seo_meta.og_image ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }] : undefined,
        siteName: config?.website_meta.website_name,
        type: "website",
      },
    };
  }
  return {
    title: `${post.title} - ${siteTitle}`,
    description: post.meta_description || post.summary,
    keywords: post.meta_keywords || config?.seo_meta.keywords,
    openGraph: {
      title: `${post.title} - ${config?.seo_meta.og_title || siteTitle}`,
      description: post.meta_description || post.summary,
      url: process.env.BASE_HOST + "/about",
      images: post.cover_img ? [{ url: post.cover_img }] : undefined,
      siteName: config?.website_meta.website_name,
      type: "article",
    },
  };
}

function AboutNotConfigured() {
  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      <section className="rounded-xl border border-dashed border-gray-200 bg-white p-8 text-center shadow-sm dark:border-gray-700 dark:bg-[#141414]">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 text-xl dark:bg-[#232426]">
          i
        </div>
        <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">关于页面未配置</h1>
        <p className="mt-3 text-sm text-gray-500 dark:text-gray-400">
          请在后台创建 ID 为 <span className="font-mono text-gray-700 dark:text-gray-200">about-me</span> 的文章后，关于页面会自动展示文章内容。
        </p>
      </section>
    </div>
  );
}

export default async function AboutPage() {
  const post = await getPostDetailOrNull("about-me");
  if (!post) return <AboutNotConfigured />;
  return <PostDetail post={post} />;
}
