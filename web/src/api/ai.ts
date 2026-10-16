import type { ArticleReference, ChatConversation } from "../components/ai/types";

export interface ChatEvent {
  type: "meta" | "delta" | "references" | "tool" | "done";
  conversation_id?: string;
  request_id?: string;
  message_id?: string;
  delta?: string;
  references?: ArticleReference[];
  tool?: string;
  status?: "complete" | "stopped" | "failed" | "interrupted";
  message?: string;
}
export interface ChatHistory {
  conversations: ChatConversation[];
  enabled: boolean;
  assistant_name: string;
}

async function checkedFetch(path: string, init: RequestInit = {}) {
  const response = await fetch(`/api/ai${path}`, {
    ...init, credentials: "same-origin", cache: "no-store",
    headers: { "Content-Type": "application/json", ...init.headers },
  });
  if (!response.ok) {
    const body = await response.json().catch(() => null);
    throw new Error(typeof body?.message === "string" ? body.message : "对话服务暂不可用，请稍后重试");
  }
  return response;
}
async function jsonRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const response = await checkedFetch(path, { ...init, signal: init.signal
    ? AbortSignal.any([init.signal, AbortSignal.timeout(15000)]) : AbortSignal.timeout(15000) });
  const body = await response.json();
  if (body.code !== 0 || !body.data) throw new Error(body.message || "对话请求失败");
  return body.data as T;
}
export const getChatHistory = (signal?: AbortSignal) => jsonRequest<ChatHistory>("/conversations", { signal });
export const getChatConversation = (id: string, signal?: AbortSignal) =>
  jsonRequest<ChatConversation>(`/conversations/${encodeURIComponent(id)}`, { signal });
export const createChatConversation = (title: string, signal?: AbortSignal) =>
  jsonRequest<ChatConversation>("/conversations", { method: "POST", body: JSON.stringify({ title }), signal });
export const stopChat = async (id: string, requestID: string) => {
  await checkedFetch(`/conversations/${encodeURIComponent(id)}/stop`, {
    method: "POST", body: JSON.stringify({ request_id: requestID }), signal: AbortSignal.timeout(15000),
  });
};
export const sendChat = (id: string, requestID: string, content: string, signal: AbortSignal) =>
  checkedFetch(`/conversations/${encodeURIComponent(id)}/messages`, {
    method: "POST", body: JSON.stringify({ request_id: requestID, content }), signal,
  });

// Fetch streams may split UTF-8 characters, CRLFs and SSE events at any byte.
// Keep only one bounded event, ignore heartbeat comments and require a terminal event.
export async function consumeChatStream(response: Response, receive: (event: ChatEvent) => void | Promise<void>) {
  if (!response.headers.get("content-type")?.includes("text/event-stream") || !response.body)
    throw new Error("对话服务未返回有效的数据流");
  const reader = response.body.getReader();
  const decoder = new TextDecoder();
  let buffer = "", done = false;
  const process = async () => {
    let boundary: number;
    while ((boundary = buffer.indexOf("\n\n")) >= 0) {
      if (boundary > 128 * 1024) throw new Error("对话数据过长");
      const block = buffer.slice(0, boundary); buffer = buffer.slice(boundary + 2);
      const data = block.split("\n").filter(line => line.startsWith("data:")).map(line => line.slice(5).trimStart()).join("\n");
      if (!data) continue;
      const event = JSON.parse(data) as ChatEvent;
      if (!event || !["meta", "delta", "references", "tool", "done"].includes(event.type))
        throw new Error("对话数据格式不正确");
      await receive(event);
      if (event.type === "done") done = true;
    }
    if (buffer.length > 128 * 1024) throw new Error("对话数据过长");
  };
  try {
    for (;;) {
      const chunk = await reader.read();
      buffer += decoder.decode(chunk.value, { stream: !chunk.done });
      // A trailing CR is retained until its LF arrives in the next frame.
      buffer = buffer.replace(/\r\n/g, "\n");
      await process();
      if (done) { await reader.cancel(); break; }
      if (chunk.done) break;
    }
    if (!done) throw new Error("回答连接中断，已生成内容可从历史中恢复");
  } finally { reader.releaseLock(); }
}
