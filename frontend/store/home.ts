import { IWebmaster } from "~/api/config"
import { IMenu } from "~/api/category"

export const useHomeStore = defineStore("home", {
    state: () => ({
        isBlackMode: false,
        // myHeaderBg: 'bg-#000/20 backdrop-blur-20',
        headerTextColor: '#fff',
        menuList: [] as IMenu[],
        masterInfo: {
            name: '',
            post_count: 0,
            category_count: 0,
            website_views: 0,
            website_live_time: 0,
            profile: '',
            picture: '',
            website_icon: '',
            domain: 'localhost:8080'
        } as IWebmaster
    })
});