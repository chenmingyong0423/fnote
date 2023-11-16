<template>
    <div
        class="footer_shadow flex flex-col justify-evenly items-center text-16 h120 dark_bg_gray dark_text_white shadow-xl">
        <div>
            Copyright © {{ year }} - Designed by
            <a class="decoration-none text-#939393" href="github.com/chenmingyong0423/fnote">Fnote</a>
        </div>
        <div class="flex items-center lt-lg:flex-col">
            <a href="https://beian.miit.gov.cn/" target="_blank" class="decoration-none text-#939393">
                琼ICP备2023002354号
            </a>


            <span class="mx5 lt-lg:display-none">-</span>

            <a target="_blank" href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=44030602007624"
                class="decoration-none lt-lg:mt5">
                <img src="https://chenmingyong.cn/assets/%E5%A4%87%E6%A1%88%E5%9B%BE%E6%A0%87-a20583c8.png"
                    style="float: left;">
                <span class="text-#939393">
                    粤公安网备44030602007624号
                </span>
            </a>

        </div>
        <div>
        </div>
        <div>{{ time }}</div>
    </div>
</template>


<script lang="ts" setup>
import { useHomeStore } from '~/store/home';
const homeStore = useHomeStore()
const date = new Date(homeStore.masterInfo.website_live_time)
const year = date.getFullYear()
const time = ref("本站已运行：")

// 更新网站运营时间
setInterval(() => {
    // 计算当前时间与上线时间之间的毫秒数差
    // 计算当前时间与上线时间之间的毫秒数差
    const millisecondsDiff = Date.now() - date.getTime();

    // 将毫秒数转换成天、小时、分钟和秒
    let seconds = Math.floor(millisecondsDiff / 1000);
    let minutes = Math.floor(seconds / 60);
    let hours = Math.floor(minutes / 60);
    let days = Math.floor(hours / 24);

    // 将时间数值归零以计算余数
    seconds %= 60;
    minutes %= 60;
    hours %= 24;
    // 更新网站已经运营的时间
    time.value = `本站已运行：${days} 天 ${hours} 时 ${minutes} 分 ${seconds} 秒`
}, 1000);
</script>

<style scoped></style>    