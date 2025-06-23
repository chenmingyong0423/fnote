import LatestComments from "@/src/components/LatestComments";
import FeaturedCarousel from "../src/components/FeaturedCarousel";
import LatestPosts from "../src/components/LatestPosts";
import SiteOwnerCard from "../src/components/SiteOwnerCard";
import { getLatestPosts } from "../src/api/posts";
import { getCarouselList } from "../src/api/carousel";

export default async function Home() {
  // 动态获取最新文章
  const latestArticles = await getLatestPosts();
  // 动态获取轮播图
  const carouselItems = await getCarouselList();

  const latestComments = [
    {
      id: 1,
      user: "Alice",
      avatar: "/logo.png",
      content: "文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！文章写得很棒，受益匪浅！",
      article: { title: "深入理解 Next.js 13 App Router深入理解 Next.js 13 App Router深入理解 Next.js 13 App Router深入理解 Next.js 13 App Router深入理解 Next.js 13 App Router", link: "/articles/1" },
    },
    {
      id: 2,
      user: "Bob",
      avatar: "/logo.png",
      content: "期待更多 TypeScript 相关内容。",
      article: { title: "TypeScript 最佳实践", link: "/articles/2" },
    },
  ];

  return (
    <div className="w-full max-w-7xl mx-auto py-4 grid grid-cols-1 md:grid-cols-12 gap-8">
      {/* 左侧主内容区 8/12 */}
      <div className="md:col-span-8 flex flex-col gap-8 min-w-0">
        {/* Banner 区 */}
        <section>
          <FeaturedCarousel items={carouselItems} />
        </section>
        {/* 内容卡片区 2列 */}
        <section>
          <LatestPosts articles={latestArticles} />
        </section>
      </div>
      {/* 右侧信息区 4/12 */}
      <div className="md:col-span-4 flex flex-col gap-8 min-w-0">
        {/* 用户卡片 */}
       <SiteOwnerCard />
        {/* 最新评论卡片 */}
        <LatestComments comments={latestComments} />
      </div>
    </div>
  );
}
