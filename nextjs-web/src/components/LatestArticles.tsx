"use client";
import { Card } from "antd";
import Meta from "antd/es/card/Meta";
import Image from "next/image";
import React from "react";

export interface LatestArticle {
  id: number;
  title: string;
  summary: string;
  cover: string;
  link: string;
  views: number;
  likes: number;
  comments: number;
  date: string;
}

export default function LatestArticles({ articles }: { articles: LatestArticle[] }) {
  return (
    <section>
      <h2 className="text-xl font-bold mb-4">最新发布</h2>
      <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
        {articles.map((item) => (
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
  );
}
