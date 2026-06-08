import { getCommonConfig } from "@/src/api/config";
import { getFriendSummary, getFriends } from "@/src/api/friend";
import FriendPageClient from "@/src/components/FriendPageClient";
import type { Metadata } from "next";

export async function generateMetadata(): Promise<Metadata> {
  const config = await getCommonConfig().catch(() => null);
  if (!config) return {};
  return {
    title: `友链 - ${config.seo_meta.title || config.website_meta.website_name}`,
    description: config.seo_meta.description || config.website_meta.website_name,
    openGraph: {
      title: `友链 - ${config.seo_meta.og_title || config.website_meta.website_name}`,
      description: config.seo_meta.description,
      url: process.env.BASE_HOST + "/friends",
      images: config.seo_meta.og_image ? [{ url: process.env.NEXT_PUBLIC_SERVER_HOST + config.seo_meta.og_image }] : undefined,
      siteName: config.website_meta.website_name,
      type: "website",
    },
  };
}

export default async function FriendPage() {
  const [friends, summary] = await Promise.all([
    getFriends().catch(() => []),
    getFriendSummary().catch(() => ""),
  ]);
  return <FriendPageClient friends={friends} summary={summary} />;
}
