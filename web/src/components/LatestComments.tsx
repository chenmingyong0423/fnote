"use client";

import { BookOutlined } from "@ant-design/icons";
import { Avatar, Card, Empty } from "antd";
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
    <Card
      title="最新评论"
      className="[&_.ant-card-body]:!p-4 md:[&_.ant-card-body]:!p-6"
    >
      {comments.length === 0 ? (
        <Empty description={hasError ? "网站数据暂时异常" : "暂无评论"} />
      ) : (
        <ul className="divide-y divide-gray-100 dark:divide-gray-800">
          {comments.map((item) => (
            <li key={item.id} className="flex gap-2 py-3 md:gap-3">
              <Avatar src={item.avatar} />
              <div className="min-w-0 flex-1">
                <div className="flex flex-col">
                  <span className="break-all font-medium text-sm text-gray-800 dark:text-gray-200">
                    {item.user}
                  </span>
                  <span className="text-xs text-gray-400 mt-1 dark:text-gray-400">
                    {new Date(
                      item.created_at ? item.created_at * 1000 : Date.now()
                    ).toLocaleString()}
                  </span>
                </div>
                <div className="mt-1 flex flex-col gap-1">
                  <div className="truncate text-sm text-gray-700 dark:text-gray-200">
                    {item.content}
                  </div>
                  <a
                    href={item.article.link}
                    className="mt-1 flex min-w-0 items-center gap-1 truncate text-xs text-blue-600 hover:underline"
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
    </Card>
  );
}
