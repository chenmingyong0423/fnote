<template>
    <div
        class="footer_shadow flex flex-col justify-evenly items-center text-16 h120 dark_bg_gray dark_text_white shadow-xl">
        <div>
            Copyright © {{ year }} - Designed by
            <a class="decoration-none text-#000 font-bold hover:underline" href="github.com/chenmingyong0423/fnote">Fnote</a>
        </div>
        <div class="flex items-center lt-lg:flex-col">
            <span v-for="(item, index) in recoreds" :key="index"  v-html="item" >
            </span>
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

const recoreds = homeStore.masterInfo.records

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