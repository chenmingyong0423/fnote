"use client";
import { List, Avatar } from "antd";
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
    <List
      itemLayout="horizontal"
      dataSource={comments}
      renderItem={(item) => (
        <List.Item>
          <List.Item.Meta
            avatar={<Avatar src={item.avatar} />}
            title={
              <span>
                {item.user} 评论于 <a href={item.article.link} className="hover:underline">{item.article.title}</a>
              </span>
            }
            description={item.content}
          />
        </List.Item>
      )}
    />
  );
}
