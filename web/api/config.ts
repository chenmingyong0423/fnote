import httpRequest from "./http";

export interface IWebsiteInfo {
  website_config: IWebsite;
  notice_config: INotice;
  social_info_config: ISocialInfo;
  pay_info_config: IPayInfo[];
  seo_meta_config: SeoMetaConfigVO;
}

export interface ISocialInfo {
  social_info_list: ISocialInfoItem[];
}

export interface ISocialInfoItem {
  social_name: string;
  social_value: string;
  css_class: string;
  is_link: boolean;
}

export interface IWebsite {
  website_name: string;
  icon: string;
  live_time: number;
  records: string[];
  owner_name: string;
  owner_profile: string;
  owner_picture: string;
}

export interface INotice {
  title: string;
  content: string;
  publish_time: number;
}

export interface IPayInfo {
  name: string;
  image: string;
}

export interface SeoMetaConfigVO {
  title: string;
  description: string;
  og_title: string;
  og_image: string;
  twitter_card: string;
  baiduSite_verification: string;
  keywords: string;
  author: string;
  robots: string;
}

const prefix = "/configs";

export const getWebsiteInfo = () => httpRequest.get(`${prefix}/index`);

export interface InitializationStatusVO {
  initStatus: boolean;
}
export const getInitializationStatus = () => httpRequest.get(`${prefix}/check-initialization`);