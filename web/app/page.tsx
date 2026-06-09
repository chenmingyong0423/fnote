import LatestComments from "@/src/components/LatestComments";
import FeaturedCarousel from "../src/components/FeaturedCarousel";
import LatestPosts from "../src/components/LatestPosts";
import SiteOwnerCard from "../src/components/SiteOwnerCard";
import { getLatestPosts } from "@/src/api/posts";
import { getCarouselList } from "@/src/api/carousel";
import { getLatestComments } from "@/src/api/comments";
import {
  DEFAULT_WEBSITE_OWNER_CONFIG,
  getWebsiteOwnerConfig,
} from "@/src/api/config";
import { DEFAULT_WEBSITE_STATS, getWebsiteStats } from "@/src/api/stats";

async function settleWithFallback<T>(promise: Promise<T>, fallback: T) {
  try {
    return { data: await promise, failed: false };
  } catch {
    return { data: fallback, failed: true };
  }
}

export default async function Home() {
  const [latestArticles, carouselItems, latestComments, config, stats] =
    await Promise.all([
      settleWithFallback(getLatestPosts(), []),
      settleWithFallback(getCarouselList(), []),
      settleWithFallback(getLatestComments(), []),
      settleWithFallback(getWebsiteOwnerConfig(), DEFAULT_WEBSITE_OWNER_CONFIG),
      settleWithFallback(getWebsiteStats(), DEFAULT_WEBSITE_STATS),
    ]);

  return (
    <div className="w-full max-w-7xl mx-auto px-4 md:px-0">
      <div className="flex flex-col md:grid md:grid-cols-12 gap-6 md:gap-8">
        <div className="w-full md:col-span-8 flex flex-col gap-6 md:gap-8 min-w-0">
          <section>
            <FeaturedCarousel
              items={carouselItems.data}
              hasError={carouselItems.failed}
            />
          </section>
          <section>
            <LatestPosts
              articles={latestArticles.data}
              hasError={latestArticles.failed}
            />
          </section>
        </div>
        <div className="w-full md:col-span-4 flex flex-col gap-6 md:gap-8 min-w-0">
          <SiteOwnerCard
            name={config.data.website_owner}
            avatar={config.data.website_owner_avatar}
            bio={config.data.website_owner_profile}
            stats={stats.data}
            hasError={config.failed || stats.failed}
          />
          <LatestComments
            comments={latestComments.data}
            hasError={latestComments.failed}
          />
        </div>
      </div>
    </div>
  );
}
