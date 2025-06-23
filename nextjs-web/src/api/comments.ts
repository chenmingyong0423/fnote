import { request } from "../utils/http";
import type { LatestComment } from "../components/LatestComments";
import type { Response } from "./types";

export interface LatestCommentApiItem {
  post_id: string;
  post_title: string;
  post_url: string;
  picture: string;
  name: string;
  content: string;
  created_at: number;
}

export interface LatestCommentsResponse {
  list: LatestCommentApiItem[];
}

// 获取最新评论
export async function getLatestComments(): Promise<LatestComment[]> {
  const res = await request<Response<LatestCommentsResponse>>("/api/comments/latest");
  if (res.code !== 0) throw new Error(res.message);
  return res.data.list.map((item, idx) => ({
    id: idx + 1, // 或 item.post_id
    user: item.name,
    avatar: item.picture,
    content: item.content,
    article: { title: item.post_title, link: item.post_url },
    created_at: item.created_at,
  }));
}
