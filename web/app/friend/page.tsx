import { DEFAULT_COMMON_CONFIG, getCommonConfig } from "@/src/api/config";
import { getFriendSummary, getFriends } from "@/src/api/friend";
import FriendPageClient from "@/src/components/FriendPageClient";
import type { Metadata } from "next";
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
    title: "友链",
    description: `浏览${config.website_meta.website_name}的友情链接及友链申请说明。`,
    pathname: "/friend",
  });
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
