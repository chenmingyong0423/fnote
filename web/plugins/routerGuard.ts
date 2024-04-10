import {defineNuxtPlugin} from '#app';
import {useRouter} from 'vue-router';
import {getInitializationStatus, type InitializationStatusVO} from "~/api/config";
import {useConfigStore} from "~/store/config";
import type {IPageData, IResponse} from "~/api/http";
import type {IComment} from "~/api/comment";

export default defineNuxtPlugin(nuxtApp => {
    const router = useRouter();
    let flag = true
    router.beforeEach(async (to, from, next) => {
        const cfg = useConfigStore()
        if (flag && process.client) {
            try {
                const res: any = await getInitializationStatus()
                const apiRes: IResponse<InitializationStatusVO> = res.data.value;
                if (apiRes) {
                    if (apiRes.code === 0) {
                        cfg.initialization = apiRes.data?.initStatus || false
                    }
                }
            } catch (e) {
                console.log(e)
            } finally {
                flag = false
            }
        }
        if (!cfg.initialization) {
            const host = process.env.WEBSITE_DOMAIN || "http://localhost:5173";
            if (process.client) {
                // 客户端重定向
                window.location.href = host + "/init";
                return;
            }
        }
        next()
    });
});
