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
      content: "文章写得很棒，受益匪浅！",
      article: { title: "深入理解 Next.js 13 App Router", link: "/articles/1" },
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
        <section className="bg-white rounded-xl shadow p-6 flex flex-col items-center">
          <SiteOwnerCard />
        </section>
        {/* 最新评论卡片 */}
        <section className="bg-white rounded-xl shadow p-6 flex flex-col gap-4">
          <h2 className="text-xl font-bold mb-2 border-b border-gray-100 pb-2">最新评论</h2>
          <div className="divide-y divide-gray-100">
            {latestComments.map((c) => (
              <div key={c.id} className="flex items-start gap-3 py-3 first:pt-0 last:pb-0">
                <img src={c.avatar} alt={c.user} className="w-8 h-8 rounded-full" />
                <div className="flex-1">
                  <div className="font-medium text-sm text-gray-800">{c.user}</div>
                  <a href={c.article.link} className="text-blue-600 hover:underline text-xs font-medium">{c.article.title}</a>
                  <div className="text-gray-500 text-xs mt-1">{c.content}</div>
                </div>
              </div>
            ))}
          </div>
        </section>
      </div>
    </div>
  );
}
