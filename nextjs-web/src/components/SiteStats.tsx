"use client";
import React from "react";
import { useConfigStore } from "../store/config";

export default function SiteStats() {
  const stats = useConfigStore((s) => s.config?.stats);
  if (!stats) return null;
  return (
    <div className="w-full border-t border-b border-gray-200 py-3 grid grid-cols-3 gap-y-2 text-xs text-gray-600 divide-x divide-gray-200">
      <div className="col-span-1 flex flex-col items-center">
        <span>文章</span>
        <span className="font-bold">{stats.post_count}</span>
      </div>
      <div className="col-span-1 flex flex-col items-center">
        <span>分类</span>
        <span className="font-bold">{stats.category_count}</span>
      </div>
      <div className="col-span-1 flex flex-col items-center">
        <span>标签</span>
        <span className="font-bold">{stats.tag_count}</span>
      </div>
      <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
        <span>评论</span>
        <span className="font-bold">{stats.comment_count}</span>
      </div>
      <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
        <span>点赞</span>
        <span className="font-bold">{stats.like_count}</span>
      </div>
      <div className="col-span-1 flex flex-col items-center pt-2 border-t border-gray-100">
        <span>浏览</span>
        <span className="font-bold">{stats.website_view_count}</span>
      </div>
    </div>
  );
}
