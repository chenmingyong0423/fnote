"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { consumeChatStream, createChatConversation, getChatConversation, getChatHistory, sendChat, stopChat } from "@/src/api/ai";
import type { ChatConversation, ChatMessage } from "./types";

type Pending = { conversationId: string; requestId: string; content: string; localId: string; stopRequested: boolean };
const selectedKey = "fnote-ai-conversation";
const toolLabels: Record<string, string> = { search_posts: "正在搜索文章", read_post: "正在读取文章", get_blog_overview: "正在了解博客", get_author_profile: "正在读取博主公开信息" };
const errorText = (error: unknown) => {
  if (error instanceof DOMException && ["AbortError", "TimeoutError"].includes(error.name)) return "连接超时或中断，可从历史中恢复回答";
  if (error instanceof TypeError) return "无法连接对话服务，请检查网络后重试";
  return error instanceof Error ? error.message : "对话请求失败，请稍后重试";
};
const remember = (id: string) => { try { if (id) localStorage.setItem(selectedKey, id); else localStorage.removeItem(selectedKey); } catch {} };
const hydrate = (c: ChatConversation): ChatConversation => ({ ...c, messages: (c.messages || []).map((m, i, all) => ({ ...m,
  request_id: m.request_id || (m.role === "assistant" ? all[i - 1]?.request_id : undefined),
})) });

export function useAiChat() {
  const [conversations, setConversations] = useState<ChatConversation[]>([]);
  const [selectedId, setSelectedId] = useState("");
  const [sending, setSending] = useState(false);
  const [loading, setLoading] = useState(true);
  const [ready, setReady] = useState(false);
  const [enabled, setEnabled] = useState(false);
  const [assistantName, setAssistantName] = useState("");
  const [error, setError] = useState("");
  const [activity, setActivity] = useState("");
  const [canRetry, setCanRetry] = useState(false);
  const selected = useRef("");
  const pending = useRef<Pending | null>(null);
  const running = useRef<AbortController | null>(null);
  const selection = useRef<AbortController | null>(null);
  const mounted = useRef(false);

  const storeConversation = useCallback((conversation: ChatConversation) => {
    setConversations(list => [hydrate(conversation), ...list.filter(c => c.id !== conversation.id)]);
  }, []);
  const selectConversation = useCallback(async (id: string) => {
    if (running.current) return;
    selection.current?.abort();
    const controller = new AbortController(); selection.current = controller;
    selected.current = id; setSelectedId(id); remember(id); setError(""); setCanRetry(false); pending.current = null; setLoading(true);
    try { const c = await getChatConversation(id, controller.signal); if (!controller.signal.aborted) storeConversation(c); }
    catch (e) { if (!controller.signal.aborted) setError(errorText(e)); }
    finally { if (!controller.signal.aborted) setLoading(false); }
  }, [storeConversation]);

  const reloadHistory = useCallback(async () => {
    selection.current?.abort();
    const controller = new AbortController(); selection.current = controller; setLoading(true); setError("");
    try {
      const history = await getChatHistory(controller.signal);
      if (controller.signal.aborted) return;
      setEnabled(history.enabled); setAssistantName(history.assistant_name); setReady(true);
      setConversations(history.conversations.map(hydrate));
      let saved = ""; try { saved = localStorage.getItem(selectedKey) || ""; } catch {}
      if (history.conversations.some(c => c.id === saved)) await selectConversation(saved);
      else { selected.current = ""; setSelectedId(""); remember(""); }
    } catch (e) { if (!controller.signal.aborted) setError(errorText(e)); }
    finally { if (!controller.signal.aborted) setLoading(false); }
  }, [selectConversation]);

  useEffect(() => {
    mounted.current = true;
    void reloadHistory();
    return () => { mounted.current = false; selection.current?.abort(); running.current?.abort(); };
  }, [reloadHistory]);

  const currentConversation = conversations.find(c => c.id === selectedId);
  const busy = sending || currentConversation?.generation?.status === "running";
  // A refreshed page or another browser tab can observe a running answer. Use
  // ordinary detail reads to recover it; no separate SSE subscription is needed.
  useEffect(() => {
    if (!selectedId || sending || currentConversation?.generation?.status !== "running") return;
    const controller = new AbortController();
    const timer = setInterval(() => { void getChatConversation(selectedId, controller.signal).then(c => {
      if (!controller.signal.aborted) storeConversation(c);
    }).catch(() => {}); }, 1500);
    return () => { clearInterval(timer); controller.abort(); };
  }, [selectedId, sending, currentConversation?.generation?.status, storeConversation]);

  const run = async (item: Pending) => {
    const controller = new AbortController(); running.current = controller;
    const timer = setTimeout(() => controller.abort(), 105000);
    setSending(true); setCanRetry(false); setError(""); setActivity("正在准备回答");
    let answerID = `pending:${item.requestId}`;
    const updateAnswer = (update: (answer: ChatMessage) => ChatMessage) => setConversations(list => list.map(c => c.id === item.localId
      ? { ...c, messages: c.messages.map(m => m.role === "assistant" && m.request_id === item.requestId ? update(m) : m) } : c));
    try {
      if (!item.conversationId) {
        const conversation = await createChatConversation(item.content, controller.signal);
        const oldID = item.localId; item.conversationId = conversation.id; item.localId = conversation.id;
        setConversations(list => list.map(c => c.id === oldID ? { ...c, id: conversation.id } : c));
        selected.current = conversation.id; setSelectedId(conversation.id); remember(conversation.id);
      }
      if (item.stopRequested) { updateAnswer(m => ({ ...m, status: "stopped" })); pending.current = null; return; }
      const response = await sendChat(item.conversationId, item.requestId, item.content, controller.signal);
      await consumeChatStream(response, async event => {
        if (!mounted.current) return;
        if (event.type === "meta") {
          answerID = event.message_id || answerID;
          updateAnswer(m => ({ ...m, id: answerID }));
          if (item.stopRequested) await stopChat(item.conversationId, item.requestId);
        } else if (event.type === "delta") {
          setActivity("正在回答"); updateAnswer(m => ({ ...m, content: m.content + (event.delta || "") }));
        } else if (event.type === "references") updateAnswer(m => ({ ...m, references: event.references || [] }));
        else if (event.type === "tool") setActivity(toolLabels[event.tool || ""] || "正在查阅博客资料");
        else if (event.type === "done") {
          updateAnswer(m => ({ ...m, status: event.status || "failed", error: event.message }));
          if (event.message) setError(event.message);
          pending.current = null; setCanRetry(false);
        }
      });
    } catch (e) {
      if (!mounted.current) return;
      setError(errorText(e)); setCanRetry(true);
      updateAnswer(m => ({ ...m, status: "interrupted" }));
    } finally {
      clearTimeout(timer);
      if (mounted.current && item.conversationId) {
        try { storeConversation(await getChatConversation(item.conversationId)); } catch {}
      }
      if (running.current === controller) { running.current = null; if (mounted.current) { setSending(false); setActivity(""); } }
    }
  };

  const begin = (item: Pending) => {
    pending.current = item;
    const additions: ChatMessage[] = [
      { id: item.requestId, request_id: item.requestId, role: "user", content: item.content, status: "complete" },
      { id: `pending:${item.requestId}`, request_id: item.requestId, role: "assistant", content: "", status: "typing" },
    ];
    setConversations(list => {
      const current = list.find(c => c.id === item.localId);
      const conversation: ChatConversation = current
        ? { ...current, messages: [...current.messages.filter(m => m.request_id !== item.requestId), ...additions] }
        : { id: item.localId, title: Array.from(item.content).slice(0, 40).join(""), messages: additions };
      return [conversation, ...list.filter(c => c.id !== conversation.id)];
    });
    selected.current = item.localId; setSelectedId(item.localId);
    void run(item);
  };
  const submit = (value: string) => {
    if (running.current || busy || loading) return false;
    if (!ready || !enabled) { setError("AI 助手暂未开启或服务不可用"); return false; }
    const content = value.trim(); if (!content) return false;
    if (Array.from(content).length > 8000) { setError("问题过长，请缩短到 8000 字以内"); return false; }
    const requestId = crypto.randomUUID();
    const conversationId = selected.current.startsWith("draft:") ? "" : selected.current;
    if (!conversationId) setConversations(list => list.filter(c => !c.id.startsWith("draft:")));
    begin({ conversationId, requestId, content, localId: conversationId || `draft:${requestId}`, stopRequested: false });
    return true;
  };
  const retry = () => {
    if (running.current || busy || loading) return;
    if (pending.current) { pending.current.stopRequested = false; begin(pending.current); }
    else void reloadHistory();
  };
  const stop = async () => {
    const item = pending.current;
    if (running.current && item) {
      item.stopRequested = true; setActivity("正在停止并保存回答");
      if (!item.conversationId) return;
      try { await stopChat(item.conversationId, item.requestId); } catch (e) { setError(errorText(e)); }
    } else if (currentConversation?.generation?.status === "running") {
      try { await stopChat(currentConversation.id, currentConversation.generation.request_id); }
      catch (e) { setError(errorText(e)); }
    }
  };
  const reset = () => {
    if (running.current || busy || loading) return false;
    selected.current = ""; setSelectedId(""); remember(""); setError(""); setCanRetry(false); pending.current = null; return true;
  };
  return { messages: currentConversation?.messages || [], conversations, busy, loading, ready, enabled, assistantName, activity,
    error, setError, canRetry, retry, submit, stop, reset, selectConversation };
}
