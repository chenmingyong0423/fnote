import { request } from "../utils/http";
import type { Response } from "./types";

export interface BaiduPushVO {
  remain?: number;
  success?: number;
  not_same_site?: string[];
  not_valid?: string[];
  error?: number;
  message?: string;
}

export interface PostIndexRequest {
  urls: string;
}

export async function baiduPushIndex(data: PostIndexRequest): Promise<Response<BaiduPushVO>> {
  return request<Response<BaiduPushVO>>("/api/post-index/baidu/push", {
    method: "POST",
    body: JSON.stringify(data),
    headers: { "Content-Type": "application/json" },
  });
}
