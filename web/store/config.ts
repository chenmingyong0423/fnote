import type {
    ISocialInfoItem,
    INotice,
    IPayInfo,
    SeoMetaConfigVO,
    IWebsite, TPSVVO,
} from "@/api/config";
import {defineStore} from "pinia";
import type {WebsiteCountStats} from "~/api/statiscs";

export const useConfigStore = defineStore("config", {
    state: () => ({
        website_info: {
            website_name: '',
            website_icon: '',
            website_owner: '',
            website_owner_profile: '',
            website_owner_avatar: '',
            website_owner_email: '',
            website_runtime: 0,
            website_records: []
        } as IWebsite,
        website_count_stats: {
            post_count: 0,
            category_count: 0,
            tag_count: 0,
            comment_count: 0,
            like_count: 0,
            website_view_count: 0,
        } as WebsiteCountStats,
        notice_info: {
            content: "",
        } as INotice,
        pay_info: [] as IPayInfo[],
        social_info_list: [] as ISocialInfoItem[],
        seo_meta_config: {
            title: "fnote blog",
            description: "fnote blog",
            og_title: "fnote blog",
        } as SeoMetaConfigVO,
        tpsv_list: [] as TPSVVO[],
        initialization: false,
    }),
});
