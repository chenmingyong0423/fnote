import { getCategories } from "@/src/api/category";
import { getTags } from "@/src/api/tags";
import type { Metadata } from "next";
import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import NavigationContent from "@/src/components/NavigationContent";
import { buildPageMetadata } from "@/src/utils/seo";

async function settleWithFallback<T>(promise: Promise<T>, fallback: T) {
  try {
    return { data: await promise, failed: false };
  } catch {
    return { data: fallback, failed: true };
  }
}

export async function generateMetadata(): Promise<Metadata> {
  const config = await getCommonConfig().catch(() => DEFAULT_COMMON_CONFIG);

  return buildPageMetadata(config, {
    title: "全部分类与标签",
    description: "浏览本站文章分类与标签。",
    pathname: "/navigation",
  });
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
