"use client";
import React from "react";
import { message } from "antd";
import { addReply, AddReplyBody } from "@/src/api/comments";
import { BaseCommentForm } from "./BaseCommentForm";

interface ReplyFormProps {
  postId: string;
  commentId: string;
  replyToId?: string;
  onSuccess: () => void;
  onCancel?: () => void;
}

export const ReplyForm: React.FC<ReplyFormProps> = ({ postId, commentId, replyToId, onSuccess, onCancel }) => {
  const [preview, setPreview] = React.useState(false);
  const [formKey, setFormKey] = React.useState(0);
  const handleFinish = async (values: AddReplyBody) => {
    try {
      await addReply(commentId, { ...values, postId, replyToId });
      message.success("回复成功，待审核后显示");
      onSuccess();
    } catch (e: any) {
      message.error(e.message || "回复失败");
    }
  };
  return (
    <BaseCommentForm
      key={formKey}
      onFinish={handleFinish}
      submitText="提交回复"
      contentLabel="回复内容"
      contentPlaceholder="请输入回复内容..."
      showPreview
      preview={preview}
      onPreviewToggle={setPreview}
      showClear
      onClear={() => {
        setFormKey(k => k + 1);
        setPreview(false);
      }}
    />
  );
};
