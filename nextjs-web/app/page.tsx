import LatestComments from "@/src/components/LatestComments";
import FeaturedCarousel from "../src/components/FeaturedCarousel";
import LatestArticles from "../src/components/LatestArticles";
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
    <div className="w-full max-w-7xl mx-auto px-4 py-8 grid grid-cols-1 md:grid-cols-12 gap-8">
      {/* 左侧主内容区 8/12 */}
      <div className="md:col-span-8 flex flex-col gap-8 min-w-0">
        {/* Banner 区 */}
        <section>
          <FeaturedCarousel items={carouselItems} />
        </section>
        {/* 内容卡片区 2列 */}
        <section>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            {latestArticles.map((item) => (
              <div key={item.sug} className="bg-white rounded-xl shadow p-6 flex flex-col h-full">
                <div className="text-lg font-bold mb-2">{item.title}</div>
                <div className="text-gray-500 text-sm mb-4 flex-1">{item.summary}</div>
                <a href={`/articles/${item.sug}`} className="inline-block bg-blue-600 text-white rounded px-4 py-1 text-sm font-medium mb-4 self-start hover:bg-blue-700 transition">阅读全文</a>
                <div className="flex items-center justify-between text-xs text-gray-400 mt-auto pt-2 border-t border-gray-100">
                  <span>作者：{item.author}</span>
                  <span className="flex items-center gap-2">
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6 6 0 10-12 0v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" /></svg> {item.visit_count}</span>
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M14 9l-1-1m0 0l-1 1m1-1v6m0 0l-1 1m1-1l1 1" /></svg> {item.like_count}</span>
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M17 8h2a2 2 0 012 2v10a2 2 0 01-2 2H5a2 2 0 01-2-2V10a2 2 0 012-2h2m10 0V6a4 4 0 00-8 0v2m8 0H7" /></svg> {item.comment_count}</span>
                  </span>
                </div>
              </div>
            ))}
          </div>
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
