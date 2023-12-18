<template>
  <div>
    <div ref="myDom">
      <Header/>
    </div>
    <!-- <el-scrollbar max-height="100vh" @scroll="headerScroll" /> -->
    <div class="bg-#F0F2F5 dark:text dark:bg-black pt220">
      <slot></slot>
    </div>
    <div>
      <Footer/>
    </div>
  </div>
</template>
<script lang="ts" setup>
import {useHomeStore} from '~/store/home';
import {getWebMaster} from "~/server/api/config"
import type {IWebmasterInfo} from "~/server/api/config"
import type {IResponse} from "~/server/api/http";
import {onMounted} from "vue";

const info = useHomeStore()
onMounted(() => {
  let isBlackMode = localStorage.getItem("isBlackMode")
  if (isBlackMode === '') {
    localStorage.setItem("isBlackMode", info.is_black_mode.toString())
  } else {
    info.is_black_mode = isBlackMode === 'true'
  }
})

const webMaster = async () => {
  try {
    let postRes: any = await getWebMaster()
    let res: IResponse<IWebmasterInfo> = postRes.data.value
    if (res && res.data) {
      info.master_info = res.data.web_master_config
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