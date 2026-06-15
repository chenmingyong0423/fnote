"use client";

import { BookOutlined } from "@ant-design/icons";
import { Avatar, Empty } from "antd";
import React from "react";

export interface LatestComment {
  id: number;
  user: string;
  avatar: string;
  content: string;
  article: { title: string; link: string };
  created_at: number;
}

export default function LatestComments({
  comments,
  hasError = false,
}: {
  comments: LatestComment[];
  hasError?: boolean;
}) {
  return (
    <section className="rounded-lg border border-gray-200/80 bg-white shadow-sm dark:border-[#303030] dark:bg-[#141414]">
      <div className="flex items-center justify-between border-b border-gray-100 px-5 py-4 dark:border-[#303030]">
        <h2 className="text-base font-bold text-gray-950 dark:text-gray-100">
          最新评论
        </h2>
      </div>
      {comments.length === 0 ? (
        <div className="p-6">
          <Empty description={hasError ? "网站数据暂时异常" : "暂无评论"} />
        </div>
      ) : (
        <ul className="px-5 py-2">
          {comments.map((item) => (
            <li
              key={item.id}
              className="flex gap-3 border-b border-gray-100 py-4 last:border-b-0 dark:border-[#252525]"
            >
              <Avatar src={item.avatar} size={36} />
              <div className="min-w-0 flex-1">
                <div className="flex items-baseline justify-between gap-2">
                  <span className="break-all text-sm font-medium text-gray-800 dark:text-gray-200">
                    {item.user}
                  </span>
                  <span className="shrink-0 text-xs text-gray-400 dark:text-gray-500">
                    {new Date(
                      item.created_at ? item.created_at * 1000 : Date.now()
                    ).toLocaleDateString()}
                  </span>
                </div>
                <div className="mt-2 flex flex-col gap-1.5">
                  <div className="line-clamp-2 text-sm leading-6 text-gray-600 dark:text-gray-300">
                    {item.content}
                  </div>
                  <a
                    href={item.article.link}
                    className="flex min-w-0 items-center gap-1.5 truncate text-xs text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
                  >
                    <BookOutlined />
                    <span className="truncate">{item.article.title}</span>
                  </a>
                </div>
              </div>
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
