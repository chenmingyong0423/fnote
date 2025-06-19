"use client";
import React from "react";
import { Carousel, Card, Avatar } from "antd";
import Meta from "antd/es/card/Meta";
import Image from "next/image";
import dynamic from "next/dynamic";
import { useConfigStore } from "../store/config";

const LatestComments = dynamic(() => import("./LatestComments"), { ssr: false });

// mock 数据（可替换为接口数据）
const featuredArticles = [
  { id: 1, title: "深入理解 Next.js 13 App Router", cover: "/file.svg", description: "全新架构，拥抱未来的 React Web 开发方式。", link: "/articles/1" },
  { id: 2, title: "TypeScript 最佳实践", cover: "/window.svg", description: "让你的前端项目更健壮、更易维护。", link: "/articles/2" },
  { id: 3, title: "Tailwind CSS 高级用法", cover: "/globe.svg", description: "极致灵活的原子化 CSS 框架。", link: "/articles/3" },
];

const latestArticles = [
  { id: 4, title: "React 18 新特性全解读", summary: "Concurrent 模式、自动批处理等新特性详解。", cover: "/next.svg", link: "/articles/4", views: 1234, likes: 88, comments: 12, date: "2024-06-19" },
  { id: 5, title: "前端性能优化实战", summary: "从构建到运行时的全链路优化方案。", cover: "/vercel.svg", link: "/articles/5", views: 2345, likes: 99, comments: 22, date: "2024-06-18" },
  { id: 6, title: "Node.js 服务端开发入门", summary: "快速搭建高性能 API 服务。", cover: "/file.svg", link: "/articles/6", views: 3456, likes: 77, comments: 33, date: "2024-06-17" },
];

export default function HomeMain() {
  const config = useConfigStore((s) => s.config);
  const siteOwner = config
    ? {
        name: config.website_config.website_owner,
        avatar: config.website_config.website_owner_avatar,
        bio: config.website_config.website_owner_profile,
        links: [], // 可根据 social_info_config 或其他字段补充
      }
    : {
        name: "站长小明",
        avatar: "/logo.png",
        bio: "全栈开发者，热爱开源与分享。专注于前端、Node.js、云原生。",
        links: [
          { label: "GitHub", url: "https://github.com/owner" },
          { label: "博客", url: "/about" },
        ],
      };
  const latestComments = [
    { id: 1, user: "Alice", avatar: "/logo.png", content: "文章写得很棒，受益匪浅！", article: { title: "深入理解 Next.js 13 App Router", link: "/articles/1" } },
    { id: 2, user: "Bob", avatar: "/logo.png", content: "期待更多 TypeScript 相关内容。", article: { title: "TypeScript 最佳实践", link: "/articles/2" } },
  ];

  return (
    <div className="w-4/5 mx-auto p-8 flex flex-col md:flex-row gap-8">
      {/* 左侧主内容区 */}
      <div className="flex-[7] flex flex-col gap-8 min-w-0">
        {/* 轮播图 */}
        <section>
          <Carousel autoplay className="rounded-lg overflow-hidden shadow">
            {featuredArticles.map((item) => (
              <div key={item.id} className="relative h-56 sm:h-72 flex items-center justify-center bg-gray-100">
                <Image src={item.cover} alt={item.title} fill className="object-contain" />
                <div className="absolute bottom-0 left-0 right-0 bg-black/50 text-white p-4">
                  <a href={item.link} className="text-lg font-bold hover:underline">{item.title}</a>
                  <div className="text-sm mt-1">{item.description}</div>
                </div>
              </div>
            ))}
          </Carousel>
        </section>
        {/* 最新文章 */}
        <section>
          <h2 className="text-xl font-bold mb-4">最新发布</h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
            {latestArticles.map((item) => (
              <Card
                key={item.id}
                hoverable
                className="h-full flex flex-col p-0"
                styles={{ body: { padding: 0 } }}
              >
                {/* 封面图片 */}
                <div className="h-40 w-full relative rounded-t-lg overflow-hidden bg-gray-50">
                  <Image src={item.cover} alt={item.title} fill className="object-cover" />
                </div>
                {/* 标题 */}
                <div className="px-4 pt-4 text-lg font-bold">
                  <a href={item.link} className="hover:underline">{item.title}</a>
                </div>
                {/* 摘要 */}
                <div className="px-4 pt-2 pb-4 text-gray-500 text-sm min-h-[48px]">{item.summary}</div>
                {/* 数据区 */}
                <div className="flex items-center justify-between px-4 pb-4 text-xs text-gray-400">
                  <div className="flex items-center gap-4">
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6 6 0 10-12 0v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" /></svg> {item.views}</span>
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M14 9l-1-1m0 0l-1 1m1-1v6m0 0l-1 1m1-1l1 1" /></svg> {item.likes}</span>
                    <span className="flex items-center gap-1"><svg className="w-4 h-4" fill="none" stroke="currentColor" strokeWidth="2" viewBox="0 0 24 24"><path d="M17 8h2a2 2 0 012 2v10a2 2 0 01-2 2H5a2 2 0 01-2-2V10a2 2 0 012-2h2m10 0V6a4 4 0 00-8 0v2m8 0H7" /></svg> {item.comments}</span>
                  </div>
                  <div className="text-right whitespace-nowrap">{item.date}</div>
                </div>
              </Card>
            ))}
          </div>
        </section>
      </div>
      {/* 右侧侧边栏 */}
      <div className="flex-[3] flex flex-col gap-8 shrink-0 min-w-0">
        {/* 站长信息卡片重设计 */}
        <section>
          <Card className="text-center p-6">
            <div className="flex flex-col items-center gap-6">
              {/* 头像 */}
              <Avatar src={siteOwner.avatar} size={64} className="mb-3" />
              {/* 名字 */}
              <div className="font-bold text-lg mb-3">{siteOwner.name}</div>
              {/* 简介 */}
              <div className="text-gray-500 mb-5 text-sm">{siteOwner.bio}</div>
              {/* 指标区，两行三列布局 */}
              <div className="w-full border-t border-b border-gray-200 py-3 grid grid-cols-3 gap-y-2 text-xs text-gray-600 divide-x divide-gray-200">
                <div className="col-span-1 flex flex-col items-center">
                  <span>文章</span>
                  <span className="font-bold">66</span>
                </div>
                <div className="col-span-1 flex flex-col items-center">
                  <span>分类</span>
                  <span className="font-bold">8</span>
                </div>
                <div className="col-span-1 flex flex-col items-center">
                  <span>标签</span>
                  <span className="font-bold">20</span>
                </div>
                <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
                  <span>评论</span>
                  <span className="font-bold">123</span>
                </div>
                <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
                  <span>点赞</span>
                  <span className="font-bold">888</span>
                </div>
                <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
                  <span>浏览</span>
                  <span className="font-bold">9999</span>
                </div>
              </div>
            </div>
          </Card>
        </section>
        {/* 最新评论 */}
        <section>
          <h2 className="text-xl font-bold mb-4">最新评论</h2>
          <LatestComments comments={latestComments} />
        </section>
      </div>
    </div>
  );
}
