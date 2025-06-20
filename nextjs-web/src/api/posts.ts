import { request } from "../utils/http";

export interface LatestPostVO {
  sug: string;
  author: string;
  title: string;
  summary: string;
  cover_img: string;
  categories: string[];
  tags: string[];
  like_count: number;
  comment_count: number;
  visit_count: number;
  sticky_weight: number;
  created_at: number;
}

export interface LatestPostsResponse {
  list: LatestPostVO[];
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export async function getLatestPosts(): Promise<LatestPostVO[]> {
  const res = await request<ApiResponse<LatestPostsResponse>>("/api/posts/latest");
  if (res.code !== 0 || !res.data) throw new Error(res.message || "Failed to fetch latest posts");
  return res.data.list;
}
