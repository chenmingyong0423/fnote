import { request } from "../utils/http";
import type { Response } from "./types";

export interface WebsiteCountStatsVO {
  post_count: number;
  category_count: number;
  tag_count: number;
  comment_count: number;
  like_count: number;
  website_view_count: number;
}

export async function getWebsiteStats(): Promise<WebsiteCountStatsVO> {
  const res = await request<Response<WebsiteCountStatsVO>>("/api/stats");
  if (res.code !== 0 || !res.data) {
    throw new Error(res.message || "Failed to fetch website stats");
  }
  return res.data;
}
