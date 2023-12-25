<template>
  <div class="dark:bg-#03080c">
    <div ref="myDom">
      <Header class="slide-down"/>
    </div>
    <!-- <el-scrollbar max-height="100vh" @scroll="headerScroll" /> -->
    <div class="bg-#F0F2F5 dark:bg-#03080c pt-25 p-5">
      <div class="slide-up w-90% m-auto">
        <slot></slot>
      </div>
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
import type {IListData, IResponse} from "~/server/api/http";
import {onMounted, ref} from "vue";
import {getMenus, type IMenu} from "~/server/api/category";

const info = useHomeStore()
const myDom = ref()
const homeStore = useHomeStore()

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
      info.notice_info = res.data.notice_config
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

const menus = async () => {
  try {
    let postRes: any = await getMenus()
    let res: IResponse<IListData<IMenu>> = postRes.data.value
    homeStore.menuList = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
menus()
</script>

<style scoped>
@keyframes slideDown {
  0% {
    transform: translateY(-100%);
  }
  100% {
    transform: translateY(0);
  }
}

.slide-down {
  animation: slideDown 1s ease;
  animation-iteration-count: 1;
}

@keyframes slideUp {
  0% {
    transform: translateY(+100%);
  }
  100% {
    transform: translateY(0);
  }
}

.slide-up {
  animation: slideUp 0.5s ease;
  animation-iteration-count: 1;
}
</style>