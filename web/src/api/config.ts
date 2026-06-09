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

export const DEFAULT_COMMON_CONFIG: CommonConfigVO = {
  website_meta: {
    website_name: "Fnote",
    website_icon: "",
    website_owner: "Fnote",
    website_owner_profile: "",
    website_owner_avatar: "",
    website_runtime: 0,
  },
  seo_meta: {
    title: "Fnote",
    description: "站点数据暂时不可用",
    og_title: "Fnote",
    og_image: "",
    baidu_site_verification: "",
    keywords: "",
    author: "Fnote",
    robots: "noindex, nofollow",
  },
  third_party_site_verification: [],
  records: [],
};

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

export const DEFAULT_WEBSITE_OWNER_CONFIG: WebsiteOwnerConfigVO = {
  website_owner: "Fnote",
  website_owner_profile: "站点信息暂时无法加载",
  website_owner_avatar: "",
};

// 获取网站主信息
export async function getWebsiteOwnerConfig(): Promise<WebsiteOwnerConfigVO> {
  const res = await request<Response<WebsiteOwnerConfigVO>>("/api/configs/owner");
  if (res.code !== 0) throw new Error(res.message);
  return res.data;
}
