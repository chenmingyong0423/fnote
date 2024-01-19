<template>
  <div
      class="flex items-center justify-center gap-x-3 text-gray-5 text-6 p-t-5 w-full"
  >
    <div class="group w-12 h-12 bg-white rounded-full flex items-center overflow-hidden md:hover:w-50 cursor-pointer duration-500" v-for="(icon, index) in homeStore.social_info_list" :key="index">
      <div
          class="bg-white rounded-full p-3 w-10 flex items-center justify-center shadow-2xl shadow-black group-hover:bg-#1E80FF group-hover:text-#fff duration-300"
      >
        <a
            v-if="icon.is_link"
            :class="get(icon.css_class)"
            :title="icon.social_name"
            :href="icon.social_value"
            target="_blank"
        ></a>
        <span
            v-else
            :class="get(icon.css_class)"
            @click="copyExternalLink(icon.social_name + ': ' + icon.social_value)"
        ></span>
      </div>
      <div class="text-4 w-full text-center truncate select-none">
        <a
            v-if="icon.is_link"
            :href="icon.social_value"
            target="_blank"
        >{{icon.social_name}}</a>
        <span v-else @click="copyExternalLink(icon.social_name + ': ' + icon.social_value)">
          {{icon.social_name}}
        </span>
      </div>
    </div>

  </div>
</template>

<script lang="ts" setup>
import { useHomeStore } from "~/store/home";

import { useAlertStore } from "~/store/toast";

const toast = useAlertStore();

const copyExternalLink = async (content: string) => {
  await navigator.clipboard.writeText(content);
  toast.showToast("复制成功！", 2000);
};

const homeStore = useHomeStore();
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
