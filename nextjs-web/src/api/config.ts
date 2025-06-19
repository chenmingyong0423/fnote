import { request } from "../utils/http";

export interface WebsiteConfigVO {
  website_name: string;
  website_icon: string;
  website_owner: string;
  website_owner_profile: string;
  website_owner_avatar: string;
  website_runtime?: number;
  website_records?: string[];
}

export interface NoticeConfigVO {
  title: string;
  content: string;
  publish_time: number;
  enabled: boolean;
}

export interface SocialInfoVO {
  name: string;
  url: string;
  icon?: string;
}

export interface SocialInfoConfigVO {
  social_info_list: SocialInfoVO[];
}

export interface PayInfoConfigVO {
  name: string;
  image: string;
}

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

export interface IndexConfigVO {
  website_config: WebsiteConfigVO;
  notice_config: NoticeConfigVO;
  social_info_config: SocialInfoConfigVO;
  pay_info_config: PayInfoConfigVO[];
  seo_meta_config: SeoMetaConfigVO;
  third_party_site_verification: TPSVVO[];
}

export interface Response<T> {
  code: number;
  message: string;
  data: T;
}

// 获取网站相关配置信息
export async function getIndexConfig(): Promise<IndexConfigVO> {
  const res = await request<Response<IndexConfigVO>>("/api/configs/index");
  if (res.code !== 0) throw new Error(res.message);
  return res.data;
}
