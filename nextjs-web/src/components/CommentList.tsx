"use client";
import React from "react";
import { Button, Avatar, List } from "antd";
import { CommentItem } from "@/src/api/comments";
import { ReplyForm } from "./ReplyForm";
import { MarkdownPreview } from "./MarkdownPreview";

interface CommentListProps {
  comments: CommentItem[];
  loading: boolean;
  onReply: (commentId: string, replyToId?: string, replyToName?: string) => void;
  postId: string;
}

export const CommentList: React.FC<CommentListProps & {
  replying?: { commentId: string; replyToId?: string; replyToName?: string } | null;
  onReplyFormFinish?: () => void;
}> = ({ comments, loading, onReply, replying, onReplyFormFinish, postId }) => (
  <List
    loading={loading}
    dataSource={comments}
    locale={{ emptyText: '暂无评论' }}
    renderItem={item => (
      <List.Item
        key={item.id}
        actions={[
          (!replying || replying.commentId !== item.id || replying.replyToId) && (
            <Button size="small" type="link" onClick={() => onReply(item.id)}>回复</Button>
          ),
          replying && replying.commentId === item.id && !replying.replyToId && (
            <Button size="small" type="link" danger onClick={onReplyFormFinish || (() => {})} key="cancel-reply">取消回复</Button>
          )
        ].filter(Boolean)}
      >
        <List.Item.Meta
          avatar={<Avatar src={item.picture} />}
          title={<span>{item.username} <span className="text-xs text-gray-400 ml-2">{new Date(item.comment_time * 1000).toLocaleString()}</span></span>}
          description={
            <>
              <span><MarkdownPreview content={item.content} /></span>
              {replying && replying.commentId === item.id && !replying.replyToId && (
                <div className="w-full mt-2 p-3 bg-gray-50 dark:bg-[#232426] rounded border border-gray-200 dark:border-gray-700">
                  <ReplyForm
                    postId={postId}
                    commentId={item.id}
                    onSuccess={onReplyFormFinish || (() => {})}
                    onCancel={onReplyFormFinish || (() => {})}
                  />
                </div>
              )}
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
                        description={
                          <>
                            <span>{reply.content}</span>
                            {replying && replying.commentId === item.id && replying.replyToId === reply.id && (
                              <div className="w-full mt-2 p-3 bg-gray-50 dark:bg-[#232426] rounded border border-gray-200 dark:border-gray-700">
                                <ReplyForm
                                  postId={postId}
                                  commentId={item.id}
                                  replyToId={reply.id}
                                  onSuccess={onReplyFormFinish || (() => {})}
                                  onCancel={onReplyFormFinish || (() => {})}
                                />
                              </div>
                            )}
                          </>
                        }
                      />
                    </List.Item>
                  )}
                />
              )}
            </>
          }
        />
      </List.Item>
    )}
  />
);
