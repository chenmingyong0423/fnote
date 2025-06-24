import { request } from "../utils/http";
import type { Response } from "./types";

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

export async function getLatestPosts(): Promise<LatestPostVO[]> {
  const res = await request<Response<LatestPostsResponse>>("/api/posts/latest");
  if (res.code !== 0 || !res.data) throw new Error(res.message || "Failed to fetch latest posts");
  return res.data.list;
}

// 分页查询文章列表
export interface PostListParams {
  pageNo: number;
  pageSize: number;
  field?: string;
  order?: string;
  keyword?: string;
}

export interface PostListResponse {
  PageNo: number;
  PageSize: number;
  totalPages: number;
  totalCount: number;
  list: LatestPostVO[];
}

export async function getPostList(params: PostListParams): Promise<PostListResponse> {
  const searchParams = new URLSearchParams();
  searchParams.set("pageNo", String(params.pageNo));
  searchParams.set("pageSize", String(params.pageSize));
  if (params.field) searchParams.set("field", params.field);
  if (params.order) searchParams.set("order", params.order);
  if (params.keyword) searchParams.set("keyword", params.keyword);

  const res = await request<Response<PostListResponse>>(`/api/posts?${searchParams.toString()}`);
  if (res.code !== 0 || !res.data) throw new Error(res.message || "Failed to fetch post list");
  return res.data;
}
