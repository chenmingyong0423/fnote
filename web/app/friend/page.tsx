import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import { getFriendSummary, getFriends } from "@/src/api/friend";
import FriendPageClient from "@/src/components/FriendPageClient";
import type { Metadata } from "next";
import { resolvePublicUrl } from "@/src/utils/publicUrl";

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
    title: `友链 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: config.seo_meta.description || config.website_meta.website_name,
    openGraph: {
      title: `友链 - ${
        config.seo_meta.og_title || config.website_meta.website_name
      }`,
      description: config.seo_meta.description,
      url: process.env.BASE_HOST + "/friend",
      images: config.seo_meta.og_image
        ? [{ url: resolvePublicUrl(config.seo_meta.og_image) }]
        : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function FriendPage() {
  const [friends, summary] = await Promise.all([
    settleWithFallback(getFriends(), []),
    settleWithFallback(getFriendSummary(), ""),
  ]);

  return (
    <FriendPageClient
      friends={friends.data}
      summary={summary.data}
      hasFriendError={friends.failed}
      hasSummaryError={summary.failed}
    />
  );
}
