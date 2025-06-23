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
