"use client";
import React from "react";
import {message} from "antd";
import {addComment, AddCommentBody} from "@/src/api/comments";
import {BaseCommentForm} from "./BaseCommentForm";
import type { FormInstance } from 'antd';

interface CommentFormProps {
    postId: string;
    onSuccessAction: () => void;
}

export const CommentForm: React.FC<CommentFormProps> = ({postId, onSuccessAction}) => {
    const [preview, setPreview] = React.useState(false);
    const formRef = React.useRef<FormInstance>(null);
    const handleFinish = async (values: AddCommentBody) => {
        try {
            await addComment({...values, postId});
            message.success("评论成功，待审核后显示");
            formRef.current?.resetFields();
            setPreview(false);
            onSuccessAction();
        } catch (e: unknown) {
            message.error((e instanceof Error ? e.message : String(e)) || "评论失败");
        }
    };
    return (
        <BaseCommentForm
            ref={formRef}
            onFinish={handleFinish}
            submitText="发表评论"
            contentLabel="评论内容"
            contentPlaceholder="请输入评论内容..."
            showPreview
            preview={preview}
            onPreviewToggle={setPreview}
            showClear
            onClear={() => {
                formRef.current?.resetFields();
                setPreview(false);
            }}
        />
    );
}
