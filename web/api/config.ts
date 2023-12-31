import httpRequest from "./http";

export interface IWebmasterInfo {
    web_master_config: IWebMaster,
    notice_config: INotice,
    social_info_config: ISocialInfo,
    pay_info_config: IPayInfo[],
    seo_meta_config: SeoMetaConfigVO,
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

export interface IWebMaster {
    name: string;
    post_count: number;
    category_count: number;
    website_views: number;
    website_live_time: number;
    profile: string;
    picture: string;
    website_icon: string;
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

const prefix = "configs"

export const getWebMaster = (apiDomain: string) => httpRequest.get(`${apiDomain}/${prefix}/index`)


