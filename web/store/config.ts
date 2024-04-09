import type {
  ISocialInfoItem,
  INotice,
  IPayInfo,
  SeoMetaConfigVO,
  IWebsite,
} from "@/api/config";
import type { IMenu } from "@/api/category";
import { defineStore } from "pinia";
import type { WebsiteCountStats } from "~/api/statiscs";

export const useConfigStore = defineStore("config", {
  state: () => ({
    website_info: {
      website_name: "",
      icon: "",
      live_time: 0,
      records: [],
      owner_name: "",
      owner_profile: "",
      owner_picture: "",
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
    initialization: false,
  }),
});
