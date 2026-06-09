"use client";

import React, { useMemo, useState } from "react";
import {
  Button,
  Card,
  Col,
  Empty,
  Form,
  Input,
  Row,
  Space,
  Typography,
  message,
} from "antd";
import { MarkdownPreview } from "@/src/components/MarkdownPreview";
import { createFriend, type FriendItem } from "@/src/api/friend";
import { getUserFriendlyError } from "@/src/utils/errorMessage";

const { Title, Paragraph } = Typography;

const sectionClassName =
  "p-6 rounded-md border border-gray-200 bg-white dark:border-[#303030] dark:bg-[#1b1c1f] dark:text-gray-300";
const titleClassName = "!text-gray-900 dark:!text-gray-100";
const secondaryTextClassName = "!text-gray-500 dark:!text-gray-400";
const inputClassName =
  "dark:!border-gray-700 dark:!bg-[#141414] dark:!text-gray-200 dark:placeholder:!text-gray-500";
const emptyClassName = "dark:[&_.ant-empty-description]:!text-gray-400";
const formClassName =
  "dark:[&_.ant-form-item-label>label]:!text-gray-300 dark:[&_.ant-input-data-count]:!text-gray-500";
const friendCardClassName =
  "h-full dark:!border-[#303030] dark:!bg-[#232426] dark:[&_.ant-card-head]:!border-[#303030] dark:[&_.ant-card-body]:!text-gray-300";

interface Props {
  friends: FriendItem[];
  summary: string;
  hasFriendError?: boolean;
  hasSummaryError?: boolean;
}

export const FriendPageClient: React.FC<Props> = ({
  friends,
  summary,
  hasFriendError = false,
  hasSummaryError = false,
}) => {
  const [submitting, setSubmitting] = useState(false);
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const friendCards = useMemo(() => {
    if (!friends || friends.length === 0) return null;
    return friends.map((item) => (
      <Col xs={24} sm={12} lg={8} key={item.url}>
        <Card
          hoverable
          className={friendCardClassName}
          title={<span className="dark:text-gray-100">{item.name}</span>}
          extra={
            <a
              href={item.url}
              target="_blank"
              rel="noreferrer"
              className="dark:text-blue-400"
            >
              访问
            </a>
          }
        >
          <Space align="start" direction="horizontal" size="middle">
            <div className="w-14 h-14 flex items-center justify-center rounded bg-gray-100 text-lg font-semibold text-gray-500 overflow-hidden dark:bg-[#141414] dark:text-gray-400 dark:border dark:border-gray-700">
              {item.logo ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={item.logo}
                  alt={item.name}
                  className="w-full h-full object-cover"
                />
              ) : (
                item.name.slice(0, 1)
              )}
            </div>
            <Paragraph className="!mb-0 dark:!text-gray-300">
              {item.description || "这个站长很神秘，还没有写描述。"}
            </Paragraph>
          </Space>
        </Card>
      </Col>
    ));
  }, [friends]);

  const onFinish = async (values: unknown) => {
    setSubmitting(true);
    try {
      const res = await createFriend(values as Parameters<typeof createFriend>[0]);
      if (res.code === 0) {
        messageApi.success("已提交，等待审核");
        form.resetFields();
      } else if (res.code === 400) {
        messageApi.error(res.message);
      } else if (res.code === 403) {
        messageApi.error("友链申请暂时关闭，无法提交");
      } else if (res.code === 429) {
        messageApi.error("请勿重复申请，该申请已通过审核或正在审核中");
      } else {
        messageApi.error(res.message || "申请失败，请稍后再试");
      }
    } catch (e: unknown) {
      messageApi.error(getUserFriendlyError(e));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="w-full max-w-7xl mx-auto md:px-4 py-6 md:py-8 flex flex-col gap-8 md:gap-12 bg-white p-4 md:p-6 rounded-xl shadow-sm dark:bg-[#141414] dark:text-gray-300 dark:border dark:border-[#303030]">
      {contextHolder}
      <section className={sectionClassName}>
        <Title level={2} className={titleClassName}>
          友链
        </Title>
        <Paragraph type="secondary" className={secondaryTextClassName}>
          欢迎互换友链，先看看申请须知和注意事项。
        </Paragraph>
        {summary ? (
          <MarkdownPreview content={summary} />
        ) : (
          <Empty
            description={hasSummaryError ? "网站数据暂时异常" : "暂无申请须知"}
            className={emptyClassName}
          />
        )}
      </section>

      <section className={sectionClassName}>
        <Title level={3} className={titleClassName}>
          友链列表
        </Title>
        {friends && friends.length > 0 ? (
          <Row gutter={[16, 16]} className="mt-4">
            {friendCards}
          </Row>
        ) : (
          <Empty
            description={hasFriendError ? "网站数据暂时异常" : "暂无友链"}
            className={`my-8 ${emptyClassName}`}
          />
        )}
      </section>

      <section className={sectionClassName}>
        <Title level={3} className={titleClassName}>
          申请友链
        </Title>
        <Paragraph type="secondary" className={`mb-4 ${secondaryTextClassName}`}>
          请填写以下信息，我们会尽快审核。
        </Paragraph>
        <Form
          form={form}
          layout="vertical"
          onFinish={onFinish}
          initialValues={{ url: "https://" }}
          className={formClassName}
        >
          <Form.Item
            label="名称"
            name="name"
            rules={[{ required: true, message: "请填写网站名称" }]}
          >
            <Input placeholder="例如：某某的博客" className={inputClassName} />
          </Form.Item>
          <Form.Item
            label="网址"
            name="url"
            rules={[
              { required: true, message: "请填写网址" },
              {
                type: "url",
                message: "请输入有效的 URL，例如 https://example.com",
              },
              {
                validator: (_, value) => {
                  if (!value || String(value).startsWith("https://")) {
                    return Promise.resolve();
                  }
                  return Promise.reject(
                    new Error("请使用 https:// 开头的地址")
                  );
                },
              },
            ]}
          >
            <Input placeholder="https://example.com" className={inputClassName} />
          </Form.Item>
          <Form.Item
            label="Logo"
            name="logo"
            rules={[{ required: true, message: "请填写 logo 链接" }]}
          >
            <Input
              placeholder="https://example.com/logo.png"
              className={inputClassName}
            />
          </Form.Item>
          <Form.Item
            label="描述"
            name="description"
            rules={[
              { required: true, message: "请填写描述" },
              { max: 80, message: "不超过 80 字" },
            ]}
          >
            <Input.TextArea
              placeholder="一句话介绍你的站点"
              rows={3}
              showCount
              maxLength={80}
              className={inputClassName}
            />
          </Form.Item>
          <Form.Item
            label="邮箱"
            name="email"
            rules={[
              { required: true, message: "请填写邮箱" },
              { type: "email", message: "请输入有效邮箱" },
            ]}
          >
            <Input
              placeholder="用于联系，审核结果会通知"
              className={inputClassName}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={submitting}>
              提交申请
            </Button>
          </Form.Item>
        </Form>
      </section>
    </div>
  );
};

export default FriendPageClient;
