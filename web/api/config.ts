import httpRequest from "./http";

export interface IWebmasterInfo {
    web_master_config: IWebMaster,
    notice_config: INotice,
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
const prefix = "/configs"

export const getWebMaster = () => httpRequest.get(prefix + "/index")


