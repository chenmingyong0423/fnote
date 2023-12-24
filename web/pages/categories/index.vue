<template>
  <div class="text-black bg-white p-5 b-rounded-4 dark_bg_gray dark:text-dtc">
    <div class="mb-10">
      <div class="line-height-10 text-10 light_border_bottom p-b-5 dark_text_white">分类</div>
      <div class="mt-5 flex items-center flex-wrap gap-x-10 gap-y-10">
        <NuxtLink class=" light_border p-5 w-90 h-40 b-rounded-4 dark_bg_gray dark:text-dtc cursor-pointer custom_cursor_flow" v-for="(c, index) in categories"
             :key="index" :to="c.route">
          <div class="i-ph-columns-fill w-10 h-10 text-gray-5 p-y-1 "></div>
          <div class="text-8 font-bold p-y-1 dark_text_white">{{ c.name }}</div>
          <div class="p-y-1">{{ c.description }}</div>
          <div class="p-y-1 flex items-center gap-x-1 h-8 text-5">
            <div class="i-ph-book dark:text-dtc"></div>
            <div>{{ c.count }}</div>
          </div>
        </NuxtLink>
      </div>
    </div>

    <div>
      <div class="line-height-10 text-10 light_border_bottom p-b-5 dark_text_white">标签</div>
      <div class="mt-5 flex items-center flex-wrap gap-x-10 gap-y-10">
        <NuxtLink class="flex text-5 gap-x-2 b-rounded-4 border-solid border-gray-2 border-1 p-2 cursor-pointer custom_cursor_flow" v-for="(t, index) in tags"
             :key="index">
          <span># {{ t.name }}</span>
          <span>{{ t.count }}</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import type {IResponse} from "~/server/api/http";
import type {ICategoryAndTags, ICategory, ITag} from '~/server/api/category'
import {getCategoriesAndTags} from '~/server/api/category'

let categories = ref<ICategory[]>([]);
let tags = ref<ITag[]>([]);

const categoryAndTags = async () => {
  try {
    let httpRes: any = await getCategoriesAndTags()
    let res: IResponse<ICategoryAndTags> = httpRes.data.value
    categories.value = res.data?.categories || []
    tags.value = res.data?.tags || []
  } catch (error) {
    console.log(error);
  }
};
categoryAndTags()

</script>

<style scoped>
</style>