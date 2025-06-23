import { request } from "../utils/http";
import type { Response } from "./types";

export interface CategoryResponse {
  list: CategoryWithCountVO[];
}

// 分类及文章数
export interface CategoryWithCountVO {
  name: string;
  route: string;
  description: string;
  count: number;
}

// 分类名称
export interface CategoryNameVO {
  name: string;
}

// 菜单
export interface MenuVO {
  name: string;
  route: string;
}

const API_PREFIX = "/api/categories";

// 获取所有已启用分类及文章数
export async function getCategories(): Promise<CategoryWithCountVO[]> {
  const res = await request<Response<CategoryResponse>>(`${API_PREFIX}`);
  if (res.code !== 0 || !res.data) throw new Error(res.message);
  return res.data.list;
}

// 根据路由获取分类名称
export async function getCategoryNameByRoute(route: string): Promise<CategoryNameVO> {
  const res = await request<Response<CategoryNameVO>>(`${API_PREFIX}/route/${route}`);
  if (res.code !== 0) throw new Error(res.message);
  return res.data;
}

// 获取导航菜单（分类）
export async function getMenus(): Promise<MenuVO[]> {
  const res = await request<Response<{ list: MenuVO[] }>>(`/api/menus`);
  if (res.code !== 0) throw new Error(res.message);
  return res.data.list;
}
