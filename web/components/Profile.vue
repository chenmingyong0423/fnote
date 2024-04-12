<template>
  <div
      class="flex flex-col items-center justify-center bg-#fff p-10 b-rounded-4 dark:text-dtc md:dark_bg_gray ease-linear duration-100 md:custom_shadow md:hover:translate-y--2 lt-md:p-5 lt-md:bg_transparent"
  >
    <div class="avatar">
      <img
          :src="apiHost + configStore.website_info.website_owner_avatar"
          alt=""
          class="w-25 h-25 border-rounded-50% mx5 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
      />
    </div>
    <div class="introduction flex flex-col items-center justify-center p-5">
      <span class="text-1.5em mb-2">{{
          configStore.website_info.website_owner
        }}</span>
      <span class="text-gray-5 mb-2">{{
          configStore.website_info.website_owner_profile
        }}</span>
    </div>
    <div
        class="flex items-center justify-between border-t w-full text-gray-5 border-t-1 border-t-gray-2 border-t-solid p-t-5 mb-5"
    >
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{
            configStore.website_count_stats.post_count
          }}</span>
        <span class="">文章</span>
      </div>
      <div
          class="flex flex-col items-center justify-center w-33% border-x-1 border-x-gray-2 border-x-solid"
      >
        <span class="mb-1">{{
            configStore.website_count_stats.category_count
          }}</span>
        <span class="">分类</span>
      </div>
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{
            configStore.website_count_stats.tag_count
          }}</span>
        <span class="">标签</span>
      </div>
    </div>
    <div
        class="flex items-center justify-between border-t w-full text-gray-5 p-b-5 border-b-1 border-b-gray-2 border-b-solid"
    >
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{
            configStore.website_count_stats.comment_count
          }}</span>
        <span class="">评论</span>
      </div>
      <div
          class="flex flex-col items-center justify-center w-33% border-x-1 border-x-gray-2 border-x-solid"
      >
        <span class="mb-1">{{
            configStore.website_count_stats.like_count
          }}</span>
        <span class="">点赞</span>
      </div>
      <div class="flex flex-col items-center justify-center w-33%">
        <span class="mb-1">{{
            configStore.website_count_stats.website_view_count
          }}</span>
        <span class="">浏览量</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {useConfigStore} from "~/store/config";

const configStore = useConfigStore();
const runtimeConfig = useRuntimeConfig()
const apiHost = runtimeConfig.public.apiHost;
const get = (icon: string): string => {
  switch (icon) {
    case "i-fa6-brands:x-twitter":
      return "i-fa6-brands:x-twitter";
    case "i-fa6-brands:facebook":
      return "i-fa6-brands:facebook";
    case "i-fa6-brands:instagram":
      return "i-fa6-brands:instagram";
    case "i-fa6-brands:youtube":
      return "i-fa6-brands:youtube";
    case "i-fa6-brands:bilibili":
      return "i-fa6-brands:bilibili";
    case "i-fa6-brands:qq":
      return "i-fa6-brands:qq";
    case "i-fa6-brands:github":
      return "i-fa6-brands:github";
    case "i-fa6-brands:square-git":
      return "i-fa6-brands:square-git";
    case "i-fa6-brands:weixin":
      return "i-fa6-brands:weixin";
    case "i-fa6-brands:zhihu":
      return "i-fa6-brands:zhihu";
    case "i-bi:link-45deg":
      return "i-bi:link-45deg";
  }
  return "";
};
</script>

<style scoped></style>
