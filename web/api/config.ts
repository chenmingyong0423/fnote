import httpRequest from "./http";

export interface IWebsiteInfo {
    website_config: IWebsite;
    notice_config: INotice;
    social_info_config: ISocialInfo;
    pay_info_config: IPayInfo[];
    seo_meta_config: SeoMetaConfigVO;
    third_party_site_verification: TPSVVO[];
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
    website_icon: string;
    website_owner: string;
    website_owner_profile: string;
    website_owner_avatar: string;
    website_owner_email: string;
    website_runtime: number;
    website_records: string[];
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

export interface TPSVVO {
    key: string;
    value: string;
    description: string;
}

const prefix = "/configs";

export const getWebsiteInfo = () => httpRequest.get(`${prefix}/index`);

export interface InitializationStatusVO {
    initStatus: boolean;
}

export const getInitializationStatus = () => httpRequest.get(`${prefix}/check-initialization`);

export interface CarouselVO {
    id: string;
    title: string;
    summary: string;
    cover_img: string;
    created_at: number;
}

export const GetCarousel = () => httpRequest.get(`${prefix}/index/carousel`);
