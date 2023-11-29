import { IWebmaster } from "~/api/config"
import { IMenu } from "~/api/category"

export const useHomeStore = defineStore("home", {
    state: () => ({
        searchVisible: false,//搜索弹窗状态
        isBlackMode: false,//暗黑模式状态
        menuList: [] as IMenu[],//菜单列表
        classification: {} as IMenu | undefined,//当前分类信息
        masterInfo: {
            name: '',
            post_count: 0,
            category_count: 0,
            website_views: 0,
            website_live_time: 0,
            profile: '',
            picture: '',
            website_icon: '',
            domain: 'localhost:8080',
            records:[],
        } as IWebmaster
    }),
    // 持久化存储
    persist: process.client && {
        storage: localStorage,//存储模式：localStorage || sessionStorage
        paths: ['classification']//要存储的数据
    }

});
