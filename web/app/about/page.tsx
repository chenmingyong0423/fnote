import PostDetail from "@/src/components/PostDetail";
import { getPostDetailOrNull, type PostDetail as PostDetailType } from "@/src/api/posts";
import type { Metadata } from "next";
import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import { resolvePublicUrl } from "@/src/utils/publicUrl";

async function getAboutPost() {
  try {
    return {
      data: await getPostDetailOrNull("about-me"),
      failed: false,
    };
  } catch {
    return {
      data: null as PostDetailType | null,
      failed: true,
    };
  }
}

export async function generateMetadata(): Promise<Metadata> {
  const [post, config] = await Promise.all([
    getAboutPost(),
    getCommonConfig().catch(() => DEFAULT_COMMON_CONFIG),
  ]);
  const siteTitle = config.seo_meta.title || config.website_meta.website_name;

  if (!post.data) {
    const description = post.failed
      ? "网站数据暂时异常"
      : "关于页面暂未配置";
    return {
      title: `关于 - ${siteTitle}`,
      description,
      openGraph: {
        title: `关于 - ${config.seo_meta.og_title || siteTitle}`,
        description,
        url: process.env.BASE_HOST + "/about",
        images: config.seo_meta.og_image
          ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }]
          : undefined,
        siteName: config.website_meta.website_name,
        type: "website",
      },
    };
  }

  return {
    title: `${post.data.title} - ${siteTitle}`,
    description: post.data.meta_description || post.data.summary,
    keywords: post.data.meta_keywords || config.seo_meta.keywords,
    openGraph: {
      title: `${post.data.title} - ${config.seo_meta.og_title || siteTitle}`,
      description: post.data.meta_description || post.data.summary,
      url: process.env.BASE_HOST + "/about",
      images: post.data.cover_img ? [{ url: resolvePublicUrl(post.data.cover_img) }] : undefined,
      siteName: config.website_meta.website_name,
      type: "article",
    },
  };
}

function AboutState({ hasError }: { hasError: boolean }) {
  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      <section className="rounded-xl border border-dashed border-gray-200 bg-white p-5 text-center shadow-sm dark:border-gray-700 dark:bg-[#141414] md:p-8">
        <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-gray-100 text-xl dark:bg-[#232426]">
          i
        </div>
        <h1 className="text-xl font-bold text-gray-900 dark:text-gray-100 md:text-2xl">
          {hasError ? "网站数据暂时异常" : "关于页面未配置"}
        </h1>
        <p className="mt-3 break-words text-sm text-gray-500 dark:text-gray-400">
          {hasError ? (
            "关于页面内容暂时无法加载，请稍后再试。"
          ) : (
            <>
              请在后台创建 ID 为{" "}
              <span className="font-mono text-gray-700 dark:text-gray-200">
                about-me
              </span>{" "}
              的文章后，关于页面会自动展示文章内容。
            </>
          )}
        </p>
      </section>
    </div>
  );
}

export default async function AboutPage() {
  const post = await getAboutPost();
  if (!post.data) return <AboutState hasError={post.failed} />;
  return <PostDetail post={post.data} />;
}
