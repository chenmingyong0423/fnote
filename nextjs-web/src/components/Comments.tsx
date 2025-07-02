"use client";
import React, { useEffect, useState } from "react";
import { getCommentsByPostId, CommentItem } from "@/src/api/comments";
import { message, Divider, Card } from "antd";
import { CommentOutlined } from "@ant-design/icons";

import { CommentForm } from "./CommentForm";
import { CommentList } from "./CommentList";
import '@ant-design/v5-patch-for-react-19';

interface CommentsProps {
  postId: string;
}

export const Comments: React.FC<CommentsProps> = ({ postId }) => {
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [loading, setLoading] = useState(false);
  // 用于标记当前在哪条评论下方显示回复表单
  const [replying, setReplying] = useState<
    { commentId: string; replyToId?: string; replyToName?: string } | null
  >(null);

  const fetchComments = React.useCallback(async () => {
    setLoading(true);
    try {
      const res = await getCommentsByPostId(postId);
      setComments(res);
    } catch (e: unknown) {
      message.error((e instanceof Error ? e.message : String(e)) || "获取评论失败");
    } finally {
      setLoading(false);
    }
  }, [postId]);

  useEffect(() => {
    fetchComments().catch();
  }, [postId, fetchComments]);

  // 点击回复按钮时，设置当前回复的评论id
  const handleReply = (
    commentId: string,
    replyToId?: string,
    replyToName?: string
  ) => {
    setReplying({ commentId, replyToId, replyToName });
  };

  // 回复成功后，清空replying并刷新评论
  const handleReplyFormFinish = () => {
    setReplying(null);
    fetchComments().catch();
  };

  return (
    <Card
      className="mt-10 !rounded-xl shadow-sm"
      style={{ padding: "12px" }}
      id="comments"
      title={
        <div className="flex items-center gap-2">
          <CommentOutlined className="text-lg text-blue-500" />
          <span className="text-base font-semibold">评论</span>
        </div>
      }
    >
      <CommentForm postId={postId} onSuccessAction={fetchComments} />
      <Divider />
      <CommentList
        comments={comments}
        loading={loading}
        onReplyAction={handleReply}
        replying={replying}
        onReplyFormFinish={handleReplyFormFinish}
        postId={postId}
      />
    </Card>
  );
};
