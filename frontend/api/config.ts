import httpRequest from "./http";

export interface IWebmaster {
    name: string;
    post_count: number;
    category_count: number;
    website_views: number;
    website_live_time: number;
    profile: string;
    picture: string;
    website_icon: string;
    domain: string;
}
const prefix = "/configs"

const getWebMaster = () => {
    return httpRequest.get(prefix + "/webmaster")
};

export {
    getWebMaster
}