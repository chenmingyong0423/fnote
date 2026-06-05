import { request } from "../utils/http";

export interface CarouselItemVO {
  id: string;
  title: string;
  summary: string;
  cover_img: string;
  show: boolean;
  color: string;
  created_at: number;
  updated_at: number;
}

export interface CarouselListResponse {
  list: CarouselItemVO[];
}

export interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
}

export async function getCarouselList(): Promise<CarouselItemVO[]> {
  const res = await request<ApiResponse<CarouselListResponse>>("/api/configs/index/carousel");
  if (res.code !== 0 || !res.data) throw new Error(res.message || "Failed to fetch carousel list");
  // 只返回 show 为 true 的项
  return res.data.list;
}
