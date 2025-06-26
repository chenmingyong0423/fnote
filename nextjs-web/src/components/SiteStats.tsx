"use client";
import React from "react";

export default function SiteStats({ stats }: { stats?: any }) {
  if (!stats) return null;
  return (
    <div className="w-full rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-[#141414] p-4 text-xs text-gray-600 dark:text-gray-400">
      <div className="grid grid-cols-3 gap-y-2">
        <div className="flex flex-col items-center">
          <span>文章</span>
          <span className="font-bold">{stats.post_count}</span>
        </div>
        <div className="flex flex-col items-center">
          <span>分类</span>
          <span className="font-bold">{stats.category_count}</span>
        </div>
        <div className="flex flex-col items-center">
          <span>标签</span>
          <span className="font-bold">{stats.tag_count}</span>
        </div>
      </div>
      <div className="my-3 border-t border-gray-100 dark:border-gray-700" />
      <div className="grid grid-cols-3 gap-y-2">
        <div className="flex flex-col items-center">
          <span>评论</span>
          <span className="font-bold">{stats.comment_count}</span>
        </div>
        <div className="flex flex-col items-center">
          <span>点赞</span>
          <span className="font-bold">{stats.like_count}</span>
        </div>
        <div className="flex flex-col items-center">
          <span>浏览</span>
          <span className="font-bold">{stats.website_view_count}</span>
        </div>
      </div>
    </div>
  );
}