<template>
  <div class="flex flex-col items-center justify-center bg-#fff p-10 b-rounded-4 cursor-pointer dark:text-dtc dark_bg_gray ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--2">
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
        class="flex items-center justify-between border-t w-full text-gray-5 border-t-1 border-t-gray-2 border-t-solid p-t-5">
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
  </div>
</template>

<script lang="ts" setup>

import {useHomeStore} from "~/store/home";

const homeStore = useHomeStore()

</script>

<style scoped>
</style>