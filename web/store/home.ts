import type {IWebmasterInfo, IWebMaster, INotice} from "@/server/api/config"
import type {IMenu} from "@/server/api/category"
import {defineStore} from 'pinia'

export const useHomeStore = defineStore("home", {
    state: () => ({
        searchVisible: false,//搜索弹窗状态
        is_black_mode: false,//暗黑模式状态
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
        notice_config: {
            content: ""
        } as INotice
    }),
    // 持久化存储
    persist: process.client && {
        storage: localStorage,//存储模式：localStorage || sessionStorage
        paths: ['classification']//要存储的数据
    }

});
