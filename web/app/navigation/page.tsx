import { getCategories } from "@/src/api/category";
import { getTags } from "@/src/api/tags";
import type { Metadata } from "next";
import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import NavigationContent from "@/src/components/NavigationContent";

async function settleWithFallback<T>(promise: Promise<T>, fallback: T) {
  try {
    return { data: await promise, failed: false };
  } catch {
    return { data: fallback, failed: true };
  }
}

export async function generateMetadata(): Promise<Metadata> {
  const config = await getCommonConfig().catch(() => DEFAULT_COMMON_CONFIG);

  return {
    title: `全部分类与标签 - ${
      config.seo_meta.title || config.website_meta.website_name
    }`,
    description: "浏览本站文章分类与标签。",
    openGraph: {
      title: `全部分类与标签 - ${
        config.seo_meta.og_title || config.website_meta.website_name
      }`,
      description: "浏览本站文章分类与标签。",
      url: process.env.BASE_HOST + "/navigation",
      images: config.seo_meta.og_image
        ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }]
        : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function NavigationPage() {
  const [categories, tags] = await Promise.all([
    settleWithFallback(getCategories(), []),
    settleWithFallback(getTags(), []),
  ]);

  return (
    <NavigationContent
      categories={categories.data}
      tags={tags.data}
      hasCategoryError={categories.failed}
      hasTagError={tags.failed}
    />
  );
}
