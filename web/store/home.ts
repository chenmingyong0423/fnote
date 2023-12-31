import type {ISocialInfoItem, IWebMaster, INotice, IPayInfo, SeoMetaConfigVO} from "@/api/config"
import type {IMenu} from "@/api/category"
import {defineStore} from 'pinia'

export const useHomeStore = defineStore("home", {
    state: () => ({
        searchVisible: false,//搜索弹窗状态
        menuList: [] as IMenu[],//菜单列表
        classification: {} as IMenu | undefined,//当前分类信息
        master_info: {
            name: "",
            post_count: 0,
            category_count: 0,
            website_views: 0,
            website_live_time: 0,
            profile: "",
            picture: "",
            website_icon: "",
            domain: "",
            records: [],
        } as IWebMaster,
        notice_info: {
            content: ""
        } as INotice,
        pay_info: [] as IPayInfo[],
        social_info_list: [] as ISocialInfoItem[],
        seo_meta_config: {
            title: "fnote blog",
            description: "fnote blog",
            ogTitle: "fnote blog",
        } as SeoMetaConfigVO,
        isBlackMode: false,
        showSmallScreenMenu: false,
        apiDomain: '',
    }),
    // 持久化存储
    persist: process.client && {
        storage: localStorage,//存储模式：localStorage || sessionStorage
        paths: ['classification']//要存储的数据
    }

});
