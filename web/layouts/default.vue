<template>
  <div class="dark:bg-#03080c">
    <SmallMenu></SmallMenu>
    <MyToast></MyToast>
    <div ref="myDom">
      <Header class="slide-down"/>
    </div>
    <div class="bg-#F0F2F5 dark:bg-#03080c pt-25 p-5">
      <div class="slide-up w-90% m-auto lt-md:w-99%">
        <slot></slot>
      </div>
    </div>
    <div>
      <Footer/>
    </div>
    <div>
      <div ref="scrollToTop"
           class="i-ph:caret-circle-up-fill text-#1e80ff dark:text-dtc w-15 h-15 fixed bottom-[100px] right-[20px] cursor-pointer"
           style="display: none;" @click="handleScrollToTop"></div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import {useHomeStore} from '~/store/home';
import {type IWebsiteInfo, getWebsiteInfo} from "~/api/config"
import type {IListData, IResponse} from "~/api/http";
import {getMenus, type IMenu} from "~/api/category";
import SmallMenu from "~/components/SmallMenu.vue";
import {collectVisitLog, type VisitLogRequest} from "~/api/statiscs";

const myDom = ref()
const homeStore = useHomeStore()
onMounted(() => {
  let isBlackMode = localStorage.getItem("isBlackMode")
  if (isBlackMode === '') {
    localStorage.setItem("isBlackMode", homeStore.isBlackMode.toString())
  } else {
    homeStore.isBlackMode = isBlackMode === 'true'
  }
})

const webMaster = async () => {
  try {
    let postRes: any = await getWebsiteInfo()
    let res: IResponse<IWebsiteInfo> = postRes.data.value
    if (res && res.data) {
      homeStore.website_info = res.data.website_config
      homeStore.owner_info = res.data.owner_config
      homeStore.notice_info = res.data.notice_config
      homeStore.social_info_list = res.data.social_info_config.social_info_list
      homeStore.pay_info = res.data.pay_info_config
      homeStore.seo_meta_config = res.data.seo_meta_config
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

const scrollToTop = ref()

const scrollEvent = () => {
  if (document.body.scrollTop > 50 || document.documentElement.scrollTop > 20) {
    scrollToTop.value.style.display = "block";
  } else {
    scrollToTop.value.style.display = "none";
  }
}
const handleScrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth',
  });
}

onMounted(() => {
  window.addEventListener('scroll', scrollEvent)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', scrollEvent)
})

const collect = async () => {
  try {
    const req = {
      url: window.location.href
    } as VisitLogRequest
    await collectVisitLog(req)
  } catch (error) {
    console.log(error);
  }
}
collect()
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