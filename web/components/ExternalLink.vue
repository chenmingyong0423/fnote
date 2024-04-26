<template>
  <div
      class="flex items-center justify-center gap-x-3 text-gray-5 text-6 p-t-5 w-full"
  >
    <div
        class="group w-12 h-12 bg-white rounded-full flex items-center overflow-hidden md:hover:w-50 cursor-pointer duration-500 dark:text-dtc dark_bg_gray"
        v-for="(icon, index) in configStore.social_info_list"
        :key="index"
    >
      <div
          class="bg-white rounded-full dark:text-dtc dark_bg_gray p-3 w-10 flex items-center justify-center shadow-2xl shadow-black group-hover:bg-#1E80FF group-hover:text-#fff duration-300 dark:group-hover:bg-#1E80FF/50"
      >
        <a
            v-if="icon.is_link"
            :class="icon.css_class"
            :href="icon.social_value"
            target="_blank"
        ></a>
        <span
            v-else
            :class="icon.css_class"
            @click="copyExternalLink(icon.social_name + ': ' + icon.social_value)"
        ></span>
      </div>
      <div class="text-4 w-full text-center truncate select-none dark:text-dtc">
        <a v-if="icon.is_link" :href="icon.social_value" target="_blank">{{
            icon.social_name
          }}</a>
        <span
            v-else
            @click="copyExternalLink(icon.social_name + ': ' + icon.social_value)"
        >
          {{ icon.social_name }}
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {useAlertStore} from "~/store/toast";
import {useConfigStore} from "~/store/config";

const toast = useAlertStore();

const copyExternalLink = async (content: string) => {
  await navigator.clipboard.writeText(content);
  toast.showToast("复制成功！", 2000);
};

const configStore = useConfigStore();
</script>

<style scoped></style>
