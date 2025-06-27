"use client";
import React, { useEffect, useState } from "react";
import { getCommentsByPostId, addComment, addReply, CommentItem, AddCommentBody, AddReplyBody } from "@/src/api/comments";
import { Button, Input, message, Avatar, List, Form, Modal, Card } from "antd";
import { CommentOutlined } from "@ant-design/icons";
import md5 from "blueimp-md5";
import { CommentForm } from "./CommentForm";
import { CommentList } from "./CommentList";

interface CommentsProps {
  postId: string;
}

export const Comments: React.FC<CommentsProps> = ({ postId }) => {
  const [comments, setComments] = useState<CommentItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [replying, setReplying] = useState<string | null>(null);
  const [replyModal, setReplyModal] = useState<{ open: boolean; commentId: string; replyToId?: string; replyToName?: string }>({ open: false, commentId: '', replyToId: undefined, replyToName: undefined });

  const fetchComments = async () => {
    setLoading(true);
    try {
      const res = await getCommentsByPostId(postId);
      setComments(res);
    } catch (e: any) {
      message.error(e.message || "获取评论失败");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComments();
    // eslint-disable-next-line
  }, [postId]);

  const handleReply = (commentId: string, replyToId?: string, replyToName?: string) => {
    setReplyModal({ open: true, commentId, replyToId, replyToName });
  };

  const handleReplyFinish = async (values: AddReplyBody) => {
    try {
      await addReply(replyModal.commentId, { ...values, postId, replyToId: replyModal.replyToId });
      message.success("回复成功，待审核后显示");
      setReplyModal({ open: false, commentId: '', replyToId: undefined, replyToName: undefined });
      fetchComments();
    } catch (e: any) {
      message.error(e.message || "回复失败");
    }
  };

  return (
    <Card
      className="mt-10 !rounded-xl shadow-sm"
      style={{padding: '12px'}}
      id="comments"
      title={
        <div className="flex items-center gap-2">
          <CommentOutlined className="text-lg text-blue-500" />
          <span className="text-base font-semibold">评论</span>
        </div>
      }
    >
      <CommentForm postId={postId} onSuccess={fetchComments} />
      <CommentList comments={comments} loading={loading} onReply={handleReply} />
      <Modal
        open={replyModal.open}
        title={replyModal.replyToName ? `回复 @${replyModal.replyToName}` : '回复评论'}
        onCancel={() => setReplyModal({ open: false, commentId: '', replyToId: undefined, replyToName: undefined })}
        footer={null}
        destroyOnHidden
      >
        <Form layout="vertical" onFinish={handleReplyFinish}>
          <Form.Item name="username" label="昵称" rules={[{ required: true, message: '请输入昵称' }]}> <Input /> </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ required: true, message: '请输入邮箱' }]}> <Input /> </Form.Item>
          <Form.Item name="website" label="个人网站"> <Input /> </Form.Item>
          <Form.Item name="content" label="回复内容" rules={[{ required: true, message: '请输入回复内容' }]}> <Input.TextArea rows={3} /> </Form.Item>
          <Form.Item> <Button type="primary" htmlType="submit">提交回复</Button> </Form.Item>
        </Form>
      </Modal>
    </Card>
  );
};
