import httpRequest from "./http";

export interface IWebsiteInfo {
    website_config: IWebsite,
    owner_config: IOwner,
    notice_config: INotice,
    social_info_config: ISocialInfo,
    pay_info_config: IPayInfo[],
    seo_meta_config: SeoMetaConfigVO,
}

export interface IOwner {
    name: string;
    profile: string;
    picture: string;
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
    name: string;
    icon: string;
    post_count: number;
    category_count: number;
    view_count: number;
    live_time: number;
    domain: string;
    records: string[];
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
    ogTitle: string;
    ogImage: string;
    twitterCard: string;
    baiduSiteVerification: string;
    keywords: string;
    author: string;
    robots: string;
}

const prefix = "/configs"

export const getWebsiteInfo = () => httpRequest.get(`${prefix}/index`)


