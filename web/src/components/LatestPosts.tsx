"use client";

import { Card, Empty } from "antd";
import Image from "next/image";
import React from "react";
import type { LatestPostVO } from "../api/posts";
import { EyeOutlined, LikeOutlined, MessageOutlined } from "@ant-design/icons";
import { formatDate } from "../utils/date";

function applyInteractiveCardMotion(event: React.PointerEvent<HTMLElement>) {
  if (event.pointerType === "touch") return;

  const card = event.currentTarget;
  const rect = card.getBoundingClientRect();
  const x = event.clientX - rect.left;
  const y = event.clientY - rect.top;
  const rotateY = ((x / rect.width) - 0.5) * 6;
  const rotateX = ((0.5 - y / rect.height) * 6);

  card.style.transform = `perspective(900px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) translateY(-4px)`;
  card.style.setProperty("--card-glow-x", `${(x / rect.width) * 100}%`);
  card.style.setProperty("--card-glow-y", `${(y / rect.height) * 100}%`);
}

function resetInteractiveCardMotion(event: React.PointerEvent<HTMLElement>) {
  const card = event.currentTarget;
  card.style.transform = "";
  card.style.removeProperty("--card-glow-x");
  card.style.removeProperty("--card-glow-y");
}

function PostTags({ item }: { item: LatestPostVO }) {
  return (
    <div className="flex flex-wrap gap-1.5">
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
  );
}

function PostMeta({ item }: { item: LatestPostVO }) {
  return (
    <div className="flex flex-col gap-2 text-xs text-gray-400 dark:text-gray-500 sm:flex-row sm:items-center sm:justify-between">
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
  );
}

function PostCard({ item }: { item: LatestPostVO }) {
  return (
    <a
      href={`/posts/${item.sug}`}
      target="_blank"
      rel="noopener noreferrer"
      className="block h-full"
      tabIndex={0}
    >
      <article
        onPointerMove={applyInteractiveCardMotion}
        onPointerLeave={resetInteractiveCardMotion}
        className="glass-card group relative flex h-full transform-gpu flex-col overflow-hidden rounded-lg transition-[transform,border-color,box-shadow] duration-200 hover:border-blue-200 hover:shadow-md hover:shadow-gray-200/80 dark:hover:border-blue-900/60 dark:hover:shadow-black/20"
      >
        <div
          aria-hidden="true"
          className="pointer-events-none absolute inset-0 z-10 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
          style={{
            background:
              "radial-gradient(circle at var(--card-glow-x, 50%) var(--card-glow-y, 30%), rgba(59, 130, 246, 0.16), transparent 42%)",
          }}
        />
        <div className="relative aspect-[16/9] w-full overflow-hidden bg-gradient-to-br from-sky-50 to-blue-100 dark:from-[#18202a] dark:to-[#111827]">
          <Image
            src={item.cover_img}
            alt={item.title}
            fill
            sizes="(max-width: 768px) 100vw, 33vw"
            className="object-contain transition-transform duration-500 group-hover:scale-[1.02]"
          />
        </div>

        <div className="flex flex-1 flex-col p-4">
          <PostTags item={item} />
          <div className="mt-3 line-clamp-2 text-base font-bold leading-snug text-gray-950 transition-colors group-hover:text-blue-600 dark:text-gray-100 dark:group-hover:text-blue-300 md:text-lg">
            {item.title}
          </div>
          <div className="line-clamp-2 pt-2 text-sm leading-6 text-gray-500 dark:text-gray-400">
            {item.summary}
          </div>
          <div className="mt-auto pt-4">
            <PostMeta item={item} />
          </div>
        </div>
      </article>
    </a>
  );
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

  const sortedArticles = articles
    .slice()
    .sort((a, b) => (b.sticky_weight || 0) - (a.sticky_weight || 0));

  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-xl font-bold text-gray-900 drop-shadow-sm dark:text-gray-100">
          最新文章
        </h2>
        <span className="ml-4 h-px flex-1 bg-white/50 dark:bg-white/15" />
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 md:gap-6">
        {sortedArticles.map((item) => (
          <PostCard key={item.sug} item={item} />
        ))}
      </div>
    </section>
  );
}
