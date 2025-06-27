"use client";
import React from "react";
import { Button, Input, message, Avatar, Form } from "antd";
import md5 from "blueimp-md5";
import { addComment, AddCommentBody } from "@/src/api/comments";

interface CommentFormProps {
  postId: string;
  onSuccess: () => void;
}

export const CommentForm: React.FC<CommentFormProps> = ({ postId, onSuccess }) => {
  const [form] = Form.useForm();
  // 获取头像
  const getGravatar = (email: string) => {
    if (!email) return undefined;
    return `https://1.gravatar.com/avatar/${md5(email.trim().toLowerCase())}?s=48&d=identicon`;
  };
  const handleFinish = async (values: AddCommentBody) => {
    try {
      await addComment({ ...values, postId });
      message.success("评论成功，待审核后显示");
      form.resetFields();
      onSuccess();
    } catch (e: any) {
      message.error(e.message || "评论失败");
    }
  };
  return (
    <Form form={form} layout="vertical" onFinish={handleFinish} className="mb-8">
      <Form.Item name="content" rules={[{ required: true, message: '请输入评论内容' }]}> 
        <Input.TextArea rows={4} placeholder="请输入评论内容..." /> 
      </Form.Item>
      <Form.Item noStyle shouldUpdate>
        {() => (
          <div className="w-full flex items-center gap-3 mb-4">
            <Form.Item>
            <Avatar src={getGravatar(form.getFieldValue('email') || '')} size={40} />
            </Form.Item>
            <Form.Item name="username" rules={[{ required: true, message: '请输入昵称' }]} className="mb-0 flex-1">
              <Input placeholder="昵称" />
            </Form.Item>
            <Form.Item name="email" rules={[{ required: true, message: '请输入邮箱' }]} className="mb-0 flex-1">
              <Input placeholder="邮箱" />
            </Form.Item>
            <Form.Item name="website" className="mb-0 flex-1">
              <Input placeholder="个人网站（可选）" />
            </Form.Item>
          </div>
        )}
      </Form.Item>
      <Form.Item> <Button type="primary" htmlType="submit">发表评论</Button> </Form.Item>
    </Form>
  );
}
