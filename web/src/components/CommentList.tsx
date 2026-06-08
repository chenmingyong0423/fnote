"use client";
import React from "react";
import { Button, Avatar } from "antd";
import { CommentItem } from "@/src/api/comments";
import { ReplyForm } from "./ReplyForm";
import { MarkdownPreview } from "./MarkdownPreview";

interface CommentListProps {
  comments: CommentItem[];
  loading: boolean;
  onReplyAction: (commentId: string, replyToId?: string, replyToName?: string) => void;
  postId: string;
}

export const CommentList: React.FC<CommentListProps & {
  replying?: { commentId: string; replyToId?: string; replyToName?: string } | null;
  onReplyFormFinish?: () => void;
}> = ({ comments, loading, onReplyAction, replying, onReplyFormFinish, postId }) => {
  if (loading) {
    return <div className="py-6 text-center text-gray-400">加载中...</div>;
  }

  if (comments.length === 0) {
    return <div className="py-6 text-center text-gray-400">暂无评论</div>;
  }

  return (
    <div className="divide-y divide-gray-100 dark:divide-gray-800">
      {comments.map((item) => (
        <div key={item.id} className="flex gap-3 py-4">
          <Avatar src={item.picture} />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center justify-between gap-2">
              <span>
                {item.username}
                <span className="text-xs text-gray-400 ml-2">{new Date(item.comment_time * 1000).toLocaleString()}</span>
              </span>
              <div className="flex items-center gap-1">
                {(!replying || replying.commentId !== item.id || replying.replyToId) && (
                  <Button size="small" type="link" onClick={() => onReplyAction(item.id)}>回复</Button>
                )}
                {replying && replying.commentId === item.id && !replying.replyToId && (
                  <Button size="small" type="link" danger onClick={onReplyFormFinish || (() => {})}>取消回复</Button>
                )}
              </div>
            </div>
            <MarkdownPreview content={item.content} className="p-4" />
            {replying && replying.commentId === item.id && !replying.replyToId && (
              <div className="w-full mt-2 p-3 bg-gray-50 dark:bg-[#232426] rounded border border-gray-200 dark:border-gray-700">
                <ReplyForm
                  postId={postId}
                  commentId={item.id}
                  onSuccessAction={onReplyFormFinish || (() => {})}
                  onCancel={onReplyFormFinish || (() => {})}
                />
              </div>
            )}
            {item.replies && item.replies.length > 0 && (
              <div className="mt-3 divide-y divide-gray-100 dark:divide-gray-800">
                {item.replies.map((reply) => (
                  <div key={reply.id} className="flex gap-3 py-3">
                    <Avatar src={reply.picture} />
                    <div className="min-w-0 flex-1">
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <span>
                          {reply.name}
                          <span className="text-xs text-gray-400 ml-2">{new Date(reply.reply_time * 1000).toLocaleString()}</span>
                          {reply.reply_to && <span className="text-xs text-blue-400 ml-2">@{reply.reply_to}</span>}
                        </span>
                        <Button size="small" type="link" onClick={() => onReplyAction(item.id, reply.id, reply.name)}>回复</Button>
                      </div>
                      <div className="mt-1">{reply.content}</div>
                      {replying && replying.commentId === item.id && replying.replyToId === reply.id && (
                        <div className="w-full mt-2 p-3 bg-gray-50 dark:bg-[#232426] rounded border border-gray-200 dark:border-gray-700">
                          <ReplyForm
                            postId={postId}
                            commentId={item.id}
                            replyToId={reply.id}
                            onSuccessAction={onReplyFormFinish || (() => {})}
                            onCancel={onReplyFormFinish || (() => {})}
                          />
                        </div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      ))}
    </div>
  );
};
