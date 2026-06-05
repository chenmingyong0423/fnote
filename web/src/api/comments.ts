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

export interface CommentReply {
  id: string;
  comment_id: string;
  content: string;
  name: string;
  picture: string;
  website: string;
  reply_to_id: string;
  reply_to: string;
  reply_time: number;
}

export interface CommentItem {
  id: string;
  content: string;
  username: string;
  picture: string;
  website: string;
  comment_time: number;
  replies: CommentReply[];
}

export interface CommentListResponse {
  list: CommentItem[];
}

// 获取指定文章的评论列表
export async function getCommentsByPostId(postId: string): Promise<CommentItem[]> {
  const res = await request<Response<CommentListResponse>>(`/api/comments/id/${postId}`);
  if (res.code !== 0 || !res.data) throw new Error(res.message || 'Failed to fetch comments');
  return res.data.list;
}

export interface AddCommentBody {
  postId: string;
  username: string;
  email: string;
  website?: string;
  content: string;
}

export interface AddCommentResponse {
  id: string;
}

// 发表评论
export async function addComment(body: AddCommentBody): Promise<AddCommentResponse> {
  const res = await request<Response<{ id: string }>>('/api/comments', {
    method: 'POST',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (res.code !== 0 || !res.data) throw new Error(res.message || '评论失败');
  return res.data;
}

export interface AddReplyBody {
  postId: string;
  username: string;
  email: string;
  website?: string;
  content: string;
  replyToId?: string;
}

export interface AddReplyResponse {
  id: string;
}

// 评论回复
export async function addReply(commentId: string, body: AddReplyBody): Promise<AddReplyResponse> {
  const res = await request<Response<{ id: string }>>(`/api/comments/${commentId}/replies`, {
    method: 'POST',
    body: JSON.stringify(body),
    headers: { 'Content-Type': 'application/json' },
  });
  if (res.code !== 0 || !res.data) throw new Error(res.message || '回复失败');
  return res.data;
}
