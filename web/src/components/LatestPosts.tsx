"use client";

import { Card, Empty } from "antd";
import Image from "next/image";
import React from "react";
import type { LatestPostVO } from "../api/posts";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";

function formatDate(timestamp: number) {
  return new Date(timestamp * 1000).toLocaleDateString();
}

export default function LatestArticles({
  articles,
  hasError = false,
}: {
  articles: LatestPostVO[];
  hasError?: boolean;
}) {
  const hasArticles = Array.isArray(articles) && articles.length > 0;

  if (!hasArticles) {
    return (
      <section>
        <h2 className="mb-4 text-xl font-bold">最新文章</h2>
        <Card className="flex items-center justify-center py-10">
          <Empty description={hasError ? "网站数据暂时异常" : "暂无文章"} />
        </Card>
      </section>
    );
  }

  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">
          最新文章
        </h2>
        <span className="ml-4 h-px flex-1 bg-gray-200/80 dark:bg-gray-800" />
      </div>
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6">
        {articles
          .slice()
          .sort((a, b) => (b.sticky_weight || 0) - (a.sticky_weight || 0))
          .map((item) => (
            <a
              key={item.sug}
              href={`/posts/${item.sug}`}
              target="_blank"
              rel="noopener noreferrer"
              className="block h-full"
              tabIndex={0}
            >
              <article className="group relative flex h-full flex-col overflow-hidden rounded-lg border border-gray-200/80 bg-white shadow-sm transition-all duration-200 hover:-translate-y-1 hover:border-blue-200 hover:shadow-md hover:shadow-gray-200/80 dark:border-[#303030] dark:bg-[#141414] dark:hover:border-blue-900/60 dark:hover:shadow-none">
                <div className="absolute left-3 right-3 top-3 z-10 flex flex-wrap gap-1.5">
                  {item.sticky_weight > 0 && (
                    <span className="rounded border border-emerald-200 bg-emerald-50/95 px-2 py-0.5 text-xs font-medium text-emerald-700 backdrop-blur dark:border-emerald-900 dark:bg-emerald-950/80 dark:text-emerald-300">
                      置顶
                    </span>
                  )}
                  {item.categories?.slice(0, 2).map((cat) => (
                    <span
                      key={cat}
                      className="rounded border border-sky-200 bg-sky-50/95 px-2 py-0.5 text-xs font-medium text-sky-700 backdrop-blur dark:border-sky-900 dark:bg-sky-950/80 dark:text-sky-300"
                    >
                      {cat}
                    </span>
                  ))}
                  {item.tags?.slice(0, 2).map((tag) => (
                    <span
                      key={tag}
                      className="rounded border border-amber-200 bg-amber-50/95 px-2 py-0.5 text-xs font-medium text-amber-700 backdrop-blur dark:border-amber-900 dark:bg-amber-950/80 dark:text-amber-300"
                    >
                      #{tag}
                    </span>
                  ))}
                </div>

                <div className="relative h-36 w-full overflow-hidden bg-gray-100 dark:bg-[#1f1f1f] sm:h-40">
                  <Image
                    src={item.cover_img}
                    alt={item.title}
                    fill
                    sizes="(max-width: 768px) 100vw, 33vw"
                    className="object-cover transition-transform duration-500 group-hover:scale-[1.035]"
                  />
                  <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-black/10 via-transparent to-black/5" />
                </div>

                <div className="flex flex-1 flex-col p-3.5 md:p-4">
                  <div className="line-clamp-2 text-base font-bold leading-snug text-gray-950 transition-colors group-hover:text-blue-600 dark:text-gray-100 dark:group-hover:text-blue-300 md:text-lg">
                    {item.title}
                  </div>
                  <div className="line-clamp-2 pt-2 text-sm leading-6 text-gray-500 dark:text-gray-400">
                    {item.summary}
                  </div>
                  <div className="mt-auto flex flex-col gap-2 pt-4 text-xs text-gray-400 dark:text-gray-500 sm:flex-row sm:items-center sm:justify-between">
                    <div className="flex flex-wrap items-center gap-3 md:gap-4">
                      <span className="flex items-center gap-1">
                        <EyeOutlined /> {item.visit_count}
                      </span>
                      <span className="flex items-center gap-1">
                        <LikeOutlined /> {item.like_count}
                      </span>
                      <span className="flex items-center gap-1">
                        <MessageOutlined /> {item.comment_count}
                      </span>
                    </div>
                    <div className="whitespace-nowrap sm:text-right">
                      {formatDate(item.created_at)}
                    </div>
                  </div>
                </div>
              </article>
            </a>
          ))}
      </div>
    </section>
  );
}
