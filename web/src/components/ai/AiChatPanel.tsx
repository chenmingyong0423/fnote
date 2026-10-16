"use client";

import { useMemo, useRef, useState } from "react";
import { Alert, App, Avatar, Button, Dropdown, Modal, Popconfirm, Tooltip } from "antd";
import { Bubble, Prompts, Sender, Welcome } from "@ant-design/x";
import type { BubbleProps } from "@ant-design/x";
import type { SenderRef } from "@ant-design/x/es/sender";
import {
  ArrowUpOutlined, BookOutlined, CheckOutlined, CloseOutlined,
  CompassOutlined, CopyOutlined, HistoryOutlined, PlusOutlined, ReadOutlined,
} from "@ant-design/icons";
import ReactMarkdown from "react-markdown";
import type { ChatMessage } from "./types";
import { useAiChat } from "./useAiChat";
import styles from "./AiAssistant.module.css";

const suggestions = [
  { key: "overview", icon: <ReadOutlined />, label: "发现博客", description: "博客里有哪些值得探索的主题？" },
  { key: "reading", icon: <CompassOutlined />, label: "找点灵感", description: "推荐一个适合入门的阅读路线" },
  { key: "summary", icon: <BookOutlined />, label: "提炼重点", description: "可以帮我总结文章的核心内容吗？" },
];

export default function AiChatPanel({ open, onClose, afterClose, logoSrc, assistantName }: {
  open: boolean;
  onClose: () => void;
  afterClose: () => void;
  logoSrc: string;
  assistantName: string;
}) {
  const { message: toast } = App.useApp();
  const { messages, conversations, busy, loading, ready, enabled, assistantName: configuredName, activity, error, setError, canRetry, retry, submit: sendMessage, stop, reset: resetConversation, selectConversation } = useAiChat();
  const displayName = configuredName || assistantName;
  const inputDisabled = loading || !ready || !enabled;
  const [input, setInput] = useState("");
  const [copiedId, setCopiedId] = useState<string>();
  const sender = useRef<SenderRef>(null);
  const submit = (value: string) => {
    if (sendMessage(value)) setInput("");
  };
  const reset = () => {
    if (!resetConversation()) return;
    setInput(""); setCopiedId(undefined); sender.current?.focus();
  };
  const roles = useMemo(() => ({
    assistant: { placement: "start" as const, avatar: <Avatar src={logoSrc} alt={displayName} size={28} className={styles.avatar}>勇</Avatar>, styles: { content: { background: "var(--chat-answer)", borderRadius: 16 } } },
    user: { placement: "end" as const, styles: { content: { background: "var(--chat-question)", borderRadius: 16 } } },
  }), [logoSrc, displayName]);

  const copy = async (item: ChatMessage) => {
    try {
      await navigator.clipboard.writeText(item.content);
      setCopiedId(item.id);
      toast.success("回复已复制");
    } catch {
      toast.info("暂时无法复制，请选中文本手动复制。");
    }
  };

  return (
    <Modal
      open={open}
      onCancel={onClose}
      afterClose={afterClose}
      afterOpenChange={(visible) => { if (visible && !window.matchMedia("(max-width: 600px)").matches) sender.current?.focus(); }}
      rootClassName={styles.modal}
      width={460}
      footer={null}
      closable={false}
      styles={{ content: { padding: 0, overflow: "hidden" }, header: { margin: 0, background: "transparent" }, mask: { backgroundColor: "rgba(15, 23, 42, 0.16)" } }}
      title={
        <div className={styles.header}>
          <Avatar src={logoSrc} alt={displayName} size={40} className={styles.brandIcon}>勇</Avatar>
          <div className={styles.heading}>{displayName}</div>
          <Dropdown trigger={["click"]} disabled={busy || loading} menu={{ items: conversations.length ? conversations.map(item => ({ key: item.id, label: item.title })) : [{ key: "empty", label: "暂无历史对话", disabled: true }], onClick: ({ key }) => void selectConversation(key) }}>
            <Button type="text" shape="circle" icon={<HistoryOutlined />} aria-label="历史对话" title="历史对话" disabled={busy || loading} />
          </Dropdown>
          <Tooltip title="新建对话">
            <Popconfirm title="开始新的对话？" description="之前的对话会保留在历史记录中。" onConfirm={reset} okText="新建" cancelText="保留" disabled={busy || loading || (!messages.length && !input)}>
              <Button type="text" shape="circle" icon={<PlusOutlined />} aria-label="新建对话" disabled={busy || loading || (!messages.length && !input)} />
            </Popconfirm>
          </Tooltip>
          <Button type="text" shape="circle" icon={<CloseOutlined />} aria-label="关闭 AI 助手" onClick={onClose} />
        </div>
      }
    >
      <div className={styles.body}>
        <div className={styles.scope}><span><BookOutlined /> 博客问答</span><span role="status">{loading ? "正在加载对话" : busy ? activity || "正在生成回答" : ""}</span></div>
        {error && <Alert type="warning" message={error} showIcon closable onClose={() => setError("")} action={<Button size="small" disabled={busy || loading} onClick={retry}>{canRetry ? "重试发送" : "重新加载"}</Button>} />}
        {messages.length === 0 ? (
          <div className={styles.welcomeArea}>
            <div className={styles.eyebrow}>探索博客里的知识</div>
            <Welcome
              variant="borderless"
              className={styles.welcome}
              title={<>你好，有什么想了解的？<span className={styles.sparkle}>✦</span></>}
              description="从一个问题开始，探索文章、梳理观点，或寻找下一篇值得读的内容。"
            />
            <Prompts
              title="试着这样问"
              vertical
              items={suggestions.map((item) => ({
                ...item,
                description: <button type="button" className={styles.promptButton} disabled={busy || inputDisabled} onClick={() => submit(item.description)}>{item.description}</button>,
              }))}
              className={styles.prompts}
            />
            <p className={styles.welcomeNote}>{ready && !enabled ? "AI 助手暂未开启。" : "围绕博客公开文章提问，回答会附上可查阅的来源。"}</p>
          </div>
        ) : (
          <Bubble.List
            className={styles.messages}
            roles={roles}
            autoScroll
            items={messages.map((item): BubbleProps<string> & { key: string; role: string } => ({
              key: item.id,
              role: item.role,
              content: item.content,
              loading: item.status === "typing" && !item.content,
              messageRender: (content: string) => <div className={styles.markdown}><ReactMarkdown skipHtml components={{ img: ({ alt }) => <span>{alt || "图片"}</span>, a: ({ href, children }) => <a href={href} target="_blank" rel="noopener noreferrer">{children}</a> }}>{content || (item.status !== "typing" ? item.error || "本次未生成回答" : "")}</ReactMarkdown></div>,
              footer: item.role === "assistant" ? (
                <div>
                {!!item.references?.length && <div className={styles.sources}>{item.references.map(reference => <a key={reference.href} href={reference.href} target="_blank" rel="noopener noreferrer">{reference.title}</a>)}</div>}
                <div className={styles.messageFooter}>
                  <span>{item.status === "stopped" ? "已停止" : item.status === "typing" ? "正在生成" : item.status === "failed" ? "生成失败" : item.status === "interrupted" ? "回答已中断" : "AI 回答"}</span>
                  <Tooltip title={copiedId === item.id ? "已复制" : "复制回复"}>
                    <Button size="small" type="text" aria-label="复制回复" disabled={!item.content} icon={copiedId === item.id ? <CheckOutlined /> : <CopyOutlined />} onClick={() => void copy(item)} />
                  </Tooltip>
                </div>
                </div>
              ) : undefined,
            }))}
          />
        )}
        <div className={styles.composer}>
          <p className={styles.previewNote}>AI 回答可能存在偏差，请结合引用文章核实。</p>
          <Sender
            ref={sender}
            value={input}
            onChange={setInput}
            onSubmit={(value) => { void submit(value); }}
            onCancel={() => void stop()}
            disabled={inputDisabled}
            loading={busy}
            placeholder="问问博客里的知识…"
            autoSize={{ minRows: 2, maxRows: 5 }}
            actions={(_, { components: { SendButton, LoadingButton } }) => busy
              ? <LoadingButton aria-label="停止生成" title="停止生成" />
              : <SendButton icon={<ArrowUpOutlined />} aria-label="发送问题" title="发送问题" disabled={inputDisabled || !input.trim()} />}
          />
          <div className={styles.inputHint}>
            <span>Enter 发送 · Shift + Enter 换行</span>
          </div>
          <span className={styles.srOnly} role="status" aria-live="polite">{busy ? "正在生成回答，可点击停止。" : "可以输入问题。"}</span>
        </div>
      </div>
    </Modal>
  );
}
