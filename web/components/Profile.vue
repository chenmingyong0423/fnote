<template>
  <div
      class="flex flex-col items-center justify-center bg-#fff p-10 b-rounded-4 cursor-pointer dark:text-dtc md:dark_bg_gray ease-linear duration-100 custom_shadow hover:translate-y--2 lt-md:p-5 lt-md:bg_transparent"
      @mouseleave="isTooltipVisible=false"
  >
    <div class="avatar">
      <img :src="homeStore.master_info.picture"
           alt=""
           class="w-25 h-25 border-rounded-50% mx5 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0">
    </div>
    <div class="introduction flex flex-col items-center justify-center p-5">
      <span class="text-1.5em mb-2">{{ homeStore.master_info.name }}</span>
      <span class="text-gray-5 mb-2">{{ homeStore.master_info.profile }}</span>
    </div>
    <div
        class="flex items-center justify-between border-t w-full text-gray-5 border-t-1 border-t-gray-2 border-t-solid p-t-5 mb-5">
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{ homeStore.master_info.post_count }}</span>
        <span class="">文章</span>
      </div>
      <div class="flex flex-col items-center justify-center w-33% border-x-1 border-x-gray-2 border-x-solid">
        <span class="mb-1">{{ homeStore.master_info.category_count }}</span>
        <span class="">分类</span>
      </div>
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{ homeStore.master_info.website_views }}</span>
        <span class="">浏览量</span>
      </div>
    </div>
    <div
        class="flex items-center justify-center gap-x-3 text-gray-5 text-6 border-t-1 border-t-gray-2 border-t-solid p-t-5 w-full lt-md:hidden">
      <template v-for="(icon, index) in homeStore.social_info_list" :key="index">
        <a v-if="icon.is_link" :class="get(icon.css_class)" class="custom_icon text-6" :title="icon.social_name"
           :href="icon.social_value" target="_blank"></a>
        <span v-else :class="get(icon.css_class)" class="custom_icon text-6" :title="icon.social_name"
              @click="copyExternalLink(icon.social_name + ': ' +icon.social_value)"></span>
      </template>
    </div>
  </div>
</template>

<script lang="ts" setup>

import {useHomeStore} from "~/store/home";

const isTooltipVisible = ref(false);

import {useAlertStore} from '~/store/toast';

const toast = useAlertStore();

const copyExternalLink = async (content: string) => {
  await navigator.clipboard.writeText(content);
  toast.showToast('复制成功！', 2000);
}

const tooltipVisibleChanged = (visible: boolean) => {
  isTooltipVisible.value = visible;
}

const homeStore = useHomeStore()
const get = (icon: string): string => {
  switch (icon) {
    case "i-fa6-brands:x-twitter":
      return "i-fa6-brands:x-twitter"
    case "i-fa6-brands:facebook":
      return "i-fa6-brands:facebook"
    case "i-fa6-brands:instagram":
      return "i-fa6-brands:instagram"
    case "i-fa6-brands:youtube":
      return "i-fa6-brands:youtube"
    case "i-fa6-brands:bilibili":
      return "i-fa6-brands:bilibili"
    case "i-fa6-brands:qq":
      return "i-fa6-brands:qq"
    case "i-fa6-brands:github":
      return "i-fa6-brands:github"
    case "i-fa6-brands:square-git":
      return "i-fa6-brands:square-git"
    case "i-fa6-brands:weixin":
      return "i-fa6-brands:weixin"
    case "i-fa6-brands:zhihu":
      return "i-fa6-brands:zhihu"
    case "fa6-brands:internet-explorer":
      return "fa6-brands:internet-explorer"
  }
  return ""
}
</script>

<style scoped>
</style>