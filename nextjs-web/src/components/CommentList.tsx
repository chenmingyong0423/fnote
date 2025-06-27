"use client";
import React from "react";
import { Button, Avatar, List } from "antd";
import { CommentItem } from "@/src/api/comments";

interface CommentListProps {
  comments: CommentItem[];
  loading: boolean;
  onReply: (commentId: string, replyToId?: string, replyToName?: string) => void;
}

export const CommentList: React.FC<CommentListProps> = ({ comments, loading, onReply }) => (
  <List
    loading={loading}
    dataSource={comments}
    locale={{ emptyText: '暂无评论' }}
    renderItem={item => (
      <List.Item
        key={item.id}
        actions={[
          <Button size="small" type="link" onClick={() => onReply(item.id)}>回复</Button>
        ]}
      >
        <List.Item.Meta
          avatar={<Avatar src={item.picture} />}
          title={<span>{item.username} <span className="text-xs text-gray-400 ml-2">{new Date(item.comment_time * 1000).toLocaleString()}</span></span>}
          description={<span>{item.content}</span>}
        />
        {item.replies && item.replies.length > 0 && (
          <List
            size="small"
            dataSource={item.replies}
            renderItem={reply => (
              <List.Item
                key={reply.id}
                actions={[
                  <Button size="small" type="link" onClick={() => onReply(item.id, reply.id, reply.name)}>回复</Button>
                ]}
              >
                <List.Item.Meta
                  avatar={<Avatar src={reply.picture} />}
                  title={<span>{reply.name} <span className="text-xs text-gray-400 ml-2">{new Date(reply.reply_time * 1000).toLocaleString()}</span> {reply.reply_to && <span className="text-xs text-blue-400 ml-2">@{reply.reply_to}</span>}</span>}
                  description={<span>{reply.content}</span>}
                />
              </List.Item>
            )}
          />
        )}
      </List.Item>
    )}
  />
);
