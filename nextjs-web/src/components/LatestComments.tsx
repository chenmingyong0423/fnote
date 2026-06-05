"use client";
import { List, Avatar, Card } from "antd";
import { BookOutlined } from "@ant-design/icons";
import React from "react";

export interface LatestComment {
  id: number;
  user: string;
  avatar: string;
  content: string;
  article: { title: string; link: string };
  created_at: number;
}

export default function LatestComments({ comments }: { comments: LatestComment[] }) {
  return (
    <Card title="最新评论">
      <List
        itemLayout="horizontal"
        dataSource={comments}
        renderItem={(item) => (
          <List.Item>
            <List.Item.Meta
              avatar={<Avatar src={item.avatar} />}
              title={
                <div className="flex flex-col">
                  <span className="font-medium text-sm text-gray-800 dark:text-gray-200">{item.user}</span>
                  <span className="text-xs text-gray-400 mt-1 dark:text-gray-400">{new Date(item.created_at ? item.created_at * 1000 : Date.now()).toLocaleString()}</span>
                </div>
              }
              description={
                <div className="flex flex-col gap-1">
                  <div className="text-gray-700 text-sm truncate dark:text-gray-200">{item.content}</div>
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
    </Card>
  );
}
