"use client";
import React, { forwardRef, useImperativeHandle } from "react";
import { Button, Input, Avatar, Form } from "antd";
import md5 from "blueimp-md5";
import { MarkdownPreview } from "@/src/components/MarkdownPreview";

export interface BaseCommentFormProps {
  initialValues?: Record<string, any>;
  onFinish: (values: any) => void;
  loading?: boolean;
  submitText?: string;
  contentLabel?: string;
  contentPlaceholder?: string;
  showPreview?: boolean;
  onPreviewToggle?: (v: boolean) => void;
  preview?: boolean;
  showClear?: boolean;
  onClear?: () => void;
}

export const BaseCommentForm = forwardRef<any, BaseCommentFormProps>(
  (
    {
      initialValues = {},
      onFinish,
      loading = false,
      submitText = "发表评论",
      contentLabel = "评论内容",
      contentPlaceholder = "请输入内容...",
      showPreview = false,
      onPreviewToggle,
      preview = false,
      showClear = false,
      onClear,
    },
    ref
  ) => {
    const [form] = Form.useForm();
    React.useEffect(() => {
      form.setFieldsValue(initialValues);
    }, [initialValues]);
    useImperativeHandle(ref, () => ({
      resetFields: () => form.resetFields(),
    }));
    const getGravatar = (email: string) => {
      if (!email) return undefined;
      return `https://1.gravatar.com/avatar/${md5(
        email.trim().toLowerCase()
      )}?s=48&d=identicon`;
    };
    return (
      <Form form={form} layout="vertical" onFinish={onFinish} className="mb-8">
        <Form.Item
          name="content"
          label={contentLabel}
          rules={[{ required: true, message: "请输入内容" }]}
        >
          {preview ? (
            <div className="border border-gray-200 rounded p-3 mb-2 bg-gray-50 dark:bg-[#232426] text-sm">
              <MarkdownPreview content={form.getFieldValue("content") || ""} className="p-4" />
            </div>
          ) : (
            <Input.TextArea rows={4} placeholder={contentPlaceholder} />
          )}
        </Form.Item>
        <Form.Item noStyle shouldUpdate>
          {() => (
            <div className="w-full flex items-center gap-3 mb-4">
              <Form.Item>
                <Avatar
                  src={getGravatar(form.getFieldValue("email") || "")}
                  size={40}
                />
              </Form.Item>
              <Form.Item
                name="username"
                rules={[{ required: true, message: "请输入昵称" }]}
                className="mb-0 flex-1"
              >
                <Input placeholder="昵称" maxLength={10} showCount />
              </Form.Item>
              <Form.Item
                name="email"
                rules={[
                  { required: true, message: "请输入邮箱" },
                  { type: "email", message: "邮箱格式不正确" },
                ]}
                className="mb-0 flex-1"
              >
                <Input placeholder="邮箱" />
              </Form.Item>
              <Form.Item
                name="website"
                className="mb-0 flex-1"
                rules={[
                  {
                    validator: (_, value) => {
                      if (!value) return Promise.resolve();
                      if (/^https:\/\//.test(value)) return Promise.resolve();
                      return Promise.reject("个人网站需以 https:// 开头");
                    },
                  },
                ]}
              >
                <Input placeholder="个人网站，以 https:// 开头（可选）" />
              </Form.Item>
            </div>
          )}
        </Form.Item>
        <div className="flex gap-2 items-center justify-center">
          <Form.Item noStyle>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
            >
              {submitText}
            </Button>
          </Form.Item>
          {showPreview && onPreviewToggle && (
            <Button
              type="default"
              onClick={() => onPreviewToggle(!preview)}
            >
              {preview ? "关闭预览" : "预览"}
            </Button>
          )}
          {showClear && onClear && (
            <Button type="default" onClick={onClear}>
              清空
            </Button>
          )}
        </div>
      </Form>
    );
  }
);
