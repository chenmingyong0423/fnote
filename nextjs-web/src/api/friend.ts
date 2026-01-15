import { request } from "../utils/http";
import type { Response } from "./types";

export interface FriendItem {
  name: string;
  url: string;
  logo: string;
  description: string;
}

export interface FriendListResponse {
  list: FriendItem[];
}

export interface FriendSummaryResponse {
  introduction: string; // markdown 文本
}

export interface CreateFriendBody {
  name: string;
  url: string;
  logo: string;
  description: string;
  email: string;
}

export async function getFriends(): Promise<FriendItem[]> {
  const res = await request<Response<FriendListResponse>>("/api/friends");
  if (res.code !== 0 || !res.data) throw new Error(res.message || "获取友链失败");
  return res.data.list;
}

export async function getFriendSummary(): Promise<string> {
  const res = await request<Response<FriendSummaryResponse>>("/api/friends/summary");
  if (res.code !== 0 || !res.data) throw new Error(res.message || "获取友链说明失败");
  return res.data.introduction;
}

export interface CreateFriendResponse {
  code: number;
  message: string;
}

export async function createFriend(body: CreateFriendBody): Promise<CreateFriendResponse> {
  const res = await request<Response<null>>("/api/friends", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  return { code: res.code, message: res.message };
}

