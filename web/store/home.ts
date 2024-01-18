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

export const useHomeStore = defineStore("home", {
  state: () => ({
    searchVisible: false, //搜索弹窗状态
    menuList: [] as IMenu[], //菜单列表
    classification: {} as IMenu | undefined, //当前分类信息
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
    isBlackMode: false,
    showSmallScreenMenu: false,
  }),
  // 持久化存储
  persist: process.client && {
    storage: localStorage, //存储模式：localStorage || sessionStorage
    paths: ["classification"], //要存储的数据
  },
});
