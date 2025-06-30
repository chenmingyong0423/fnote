"use client";
import { Card, Tag } from "antd";
import Meta from "antd/es/card/Meta";
import Image from "next/image";
import React from "react";
import type { LatestPostVO } from "../api/posts";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";

export default function LatestArticles({ articles }: { articles: LatestPostVO[] }) {
  return (
    <section>
      <h2 className="text-xl font-bold mb-4">最新发布</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
        {articles
          .slice()
          .sort((a, b) => (b.sticky_weight || 0) - (a.sticky_weight || 0))
          .map((item) => (
            <a
              key={item.sug}
              href={`/posts/${item.sug}`}
              className="block h-full"
              tabIndex={0}
            >
              <Card
                hoverable
                className="h-full flex flex-col relative"
                styles={{ body: { padding: 12, display: 'flex', flexDirection: 'column', height: '100%' } }}
              >
                {/* 标签区：左上角绝对定位，使用 antd Tag */}
                <div className="absolute top-4 left-4 flex flex-wrap gap-2 z-10">
                  {item.sticky_weight > 0 && (
                    <Tag color="success">置顶</Tag>
                  )}
                  {item.categories?.map((cat) => (
                    <Tag key={cat} color="blue">{cat}</Tag>
                  ))}
                  {item.tags?.map((tag) => (
                    <Tag key={tag} color="orange">#{tag}</Tag>
                  ))}
                </div>
                {/* 封面图片 */}
                <div className="h-40 w-full relative rounded-t-lg overflow-hidden bg-gray-50">
                  <Image src={item.cover_img} alt={item.title} fill sizes="60" className="object-cover" />
                </div>
                {/* 标题 */}
                <div className="pt-4 text-lg font-bold">
                  {item.title}
                </div>
                {/* 摘要 */}
                <div className="pt-2 pb-4 text-gray-500 h-13 text-sm line-clamp-2">{item.summary}</div>
                {/* 数据区 */}
                <div className="flex items-center justify-between pb-2 text-xs text-gray-400 mt-auto">
                  <div className="flex items-center gap-4">
                    {/* 浏览量 */}
                    <span className="flex items-center gap-1"><EyeOutlined /> {item.visit_count}</span>
                    {/* 点赞数 */}
                    <span className="flex items-center gap-1"><LikeOutlined /> {item.like_count}</span>
                    {/* 评论数 */}
                    <span className="flex items-center gap-1"><MessageOutlined /> {item.comment_count}</span>
                  </div>
                  <div className="text-right whitespace-nowrap">{new Date(item.created_at * 1000).toLocaleDateString()}</div>
                </div>
              </Card>
            </a>
          ))}
      </div>
    </section>
  );
}
