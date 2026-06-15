"use client";

import React from "react";

export interface SiteStatsProps {
  post_count: number;
  category_count: number;
  tag_count: number;
  comment_count: number;
  like_count: number;
  website_view_count: number;
}

export default function SiteStats({ stats }: { stats?: SiteStatsProps }) {
  if (!stats) return null;

  const items = [
    { label: "文章", value: stats.post_count },
    { label: "分类", value: stats.category_count },
    { label: "标签", value: stats.tag_count },
    { label: "评论", value: stats.comment_count },
    { label: "点赞", value: stats.like_count },
    { label: "浏览", value: stats.website_view_count },
  ];

  return (
    <div className="mt-1 w-full rounded-lg border border-white/45 bg-white/35 p-2.5 text-xs text-gray-500 backdrop-blur dark:border-white/10 dark:bg-white/5 dark:text-gray-400 md:p-3">
      <div className="grid grid-cols-3 gap-2">
        {items.map((item) => (
          <div
            key={item.label}
            className="rounded-md px-1.5 py-2 transition-colors hover:bg-white/55 dark:hover:bg-white/10"
          >
            <div>{item.label}</div>
            <div className="mt-1 truncate font-semibold text-gray-900 dark:text-gray-100">
              {item.value}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
