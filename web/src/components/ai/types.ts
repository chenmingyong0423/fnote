export interface ArticleReference {
  title: string;
  href: string;
}

export interface ChatMessage {
  id: string;
  role: "user" | "assistant";
  content: string;
  status: "typing" | "complete" | "stopped" | "failed" | "interrupted";
  request_id?: string;
  error?: string;
  references?: ArticleReference[];
}

export interface ChatConversation {
  id: string;
  title: string;
  messages: ChatMessage[];
  generation?: { request_id: string; message_id: string; status: string };
  updated_at?: string;
}
