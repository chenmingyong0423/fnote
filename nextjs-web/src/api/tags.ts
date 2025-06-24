import { request } from "../utils/http";
import type { Response } from "./types";

export interface TagVO {
  name: string;
  route: string;
  count: number;
}

export interface TagListResponse {
  list: TagVO[];
}

export async function getTags(): Promise<TagVO[]> {
  const res = await request<Response<TagListResponse>>("/api/tags");
  if (res.code !== 0 || !res.data) throw new Error(res.message);
  return res.data.list;
}

// 根据路由获取标签名称
export async function getTagNameByRoute(route: string): Promise<string> {
  const res = await request<{ code: number; message: string; data: { name: string } }>(`/api/tags/route/${route}`);
  if (res.code !== 0 || !res.data) throw new Error(res.message || "Failed to fetch tag name");
  return res.data.name;
}
