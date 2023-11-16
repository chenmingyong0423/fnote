<template>
    <div>
        <div ref="myDom">
            <MyHeader />
        </div>
        <!-- <el-scrollbar max-height="100vh" @scroll="headerScroll"> -->
        <div class="bg-#F0F2F5 dark_bg_black">
            <slot></slot>
        </div>
        <MyFooter />
        <!-- </el-scrollbar> -->
        <el-backtop :right="100" :bottom="100" class="lt-lg:important:right-20 lt-lg:important:bottom-50" />
        <!-- <SearchDialog /> -->
    </div>
</template>
<script lang="ts" setup>
const myDom = ref()
let scrollCount = ref(0)
const headerScroll = (res: any) => {
    // const scrollNow = ref(document.documentElement.scrollTop)
    // console.log(Math.abs(scrollNow.value - scrollCount.value));
    if (res.scrollTop > scrollCount.value) {
        myDom.value.firstChild.setAttribute('style', `top:-50px`);
    }
    else {
        myDom.value.firstChild.setAttribute('style', 'top:0px');
    }
    scrollCount.value = res.scrollTop
}

import { useHomeStore } from '../store/home';
import { getWebMaster, IWebmaster } from "~/api/config"
import { IResponse } from "../api/http";
const info = useHomeStore()
const webMaster = async () => {
    try {
        let postRes: any = await getWebMaster()
        let res: IResponse<IWebmaster> = postRes.data.value
        if (res.data) {
            info.masterInfo = res.data


            // const newLink = document.createElement('link');
            // newLink.rel = 'icon';
            // newLink.type = 'image/x-icon';
            // newLink.href = info.masterInfo.website_icon;
            // const oldLink = document.querySelector("link[rel*='icon']");

            // if (oldLink) {
            //     document.head.removeChild(oldLink);
            // }
            // document.head.appendChild(newLink);
        }
    } catch (error) {
        console.log(error);
    }
};
webMaster()


</script>

