import { request } from "../utils/http";

export interface WebsiteCountStatsVO {
  post_count: number;
  category_count: number;
  tag_count: number;
  comment_count: number;
  like_count: number;
  website_view_count: number;
}

export interface ResponseBody<T> {
  code: number;
  message: string;
  data?: T;
}

export async function getWebsiteStats(): Promise<WebsiteCountStatsVO> {
  const res = await request<ResponseBody<WebsiteCountStatsVO>>("/api/stats");
  if (res.code !== 0 || !res.data) {
    throw new Error(res.message || "Failed to fetch website stats");
  }
  return res.data;
}
