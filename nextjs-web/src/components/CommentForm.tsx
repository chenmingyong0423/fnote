"use client";
import React from "react";
import { Button, Input, message, Avatar, Form } from "antd";
import md5 from "blueimp-md5";
import { addComment, AddCommentBody } from "@/src/api/comments";
import { MarkdownPreview } from "@/src/components/MarkdownPreview";
import { BaseCommentForm } from "./BaseCommentForm";

interface CommentFormProps {
  postId: string;
  onSuccess: () => void;
}

export const CommentForm: React.FC<CommentFormProps> = ({ postId, onSuccess }) => {
  const [preview, setPreview] = React.useState(false);
  const handleFinish = async (values: AddCommentBody) => {
    try {
      await addComment({ ...values, postId });
      message.success("评论成功，待审核后显示");
      onSuccess();
    } catch (e: any) {
      message.error(e.message || "评论失败");
    }
  };
  return (
    <BaseCommentForm
      onFinish={handleFinish}
      submitText="发表评论"
      contentLabel="评论内容"
      contentPlaceholder="请输入评论内容..."
      showPreview
      preview={preview}
      onPreviewToggle={setPreview}
      showClear
      onClear={() => {
        // 通过 ref 或 form 实例清空内容
        document.querySelector('textarea[name="content"]')?.dispatchEvent(new Event('input', { bubbles: true }));
      }}
    />
  );
}
