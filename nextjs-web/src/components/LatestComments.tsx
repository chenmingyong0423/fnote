"use client";
import { List, Avatar } from "antd";
import { BookOutlined } from "@ant-design/icons";
import React from "react";

export interface LatestComment {
  id: number;
  user: string;
  avatar: string;
  content: string;
  article: { title: string; link: string };
}

export default function LatestComments({ comments }: { comments: LatestComment[] }) {
  return (
    <section className="bg-white rounded-xl shadow p-6 flex flex-col gap-4">
      <h2 className="text-xl font-bold mb-2 border-b border-gray-100 pb-2">最新评论</h2>
      <List
        itemLayout="horizontal"
        dataSource={comments}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              avatar={<Avatar src={item.avatar} />}
              title={
                <div className="flex flex-col">
                  <span className="font-medium text-sm text-gray-800">{item.user}</span>
                  <span className="text-xs text-gray-400 mt-1">{new Date(item.id ? item.id * 100000000 : Date.now()).toLocaleString()}</span>
                </div>
              }
              description={
                <div className="flex flex-col gap-1">
                  <div className="text-gray-700 text-sm truncate">{item.content}</div>
                  <a href={item.article.link} className="text-blue-600 hover:underline text-xs flex items-center gap-1 mt-1 truncate">
                    <BookOutlined />
                    <span className="truncate">{item.article.title}</span>
                  </a>
                </div>
              }
            />
          </List.Item>
        )}
      />
    </section>
  );
}
