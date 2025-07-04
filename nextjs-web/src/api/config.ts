import { request } from "../utils/http";
import type { Response } from "./types";

export interface SeoMetaConfigVO {
  title: string;
  description: string;
  og_title: string;
  og_image: string;
  baidu_site_verification: string;
  keywords: string;
  author: string;
  robots: string;
}

export interface TPSVVO {
  key: string;
  value: string;
  description: string;
}

export interface WebsiteMetaVO {
  website_name: string;
  website_icon: string;
  website_owner: string;
  website_owner_profile: string;
  website_owner_avatar: string;
  website_runtime: number;
}

export interface CommonConfigVO {
  website_meta: WebsiteMetaVO;
  seo_meta: SeoMetaConfigVO;
  third_party_site_verification: TPSVVO[];
  records: string[];
}

// 获取通用配置信息
export async function getCommonConfig(): Promise<CommonConfigVO> {
  const res = await request<Response<CommonConfigVO>>("/api/configs/common");
  if (res.code !== 0) throw new Error(res.message);
  return res.data;
}

export interface WebsiteOwnerConfigVO {
  website_owner: string;
  website_owner_profile: string;
  website_owner_avatar: string;
}

// 获取网站主信息
export async function getWebsiteOwnerConfig(): Promise<WebsiteOwnerConfigVO> {
  const res = await request<Response<WebsiteOwnerConfigVO>>("/api/configs/owner");
  if (res.code !== 0) throw new Error(res.message);
  return res.data;
}
