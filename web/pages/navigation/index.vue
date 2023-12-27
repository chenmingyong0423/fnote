<template>
  <div class="text-black bg-white p-5 b-rounded-4 dark_bg_gray dark:text-dtc">
    <div class="mb-10">
      <div class="line-height-10 text-10 light_border_bottom p-b-5 dark_text_white">分类</div>
      <div class="mt-5 flex items-center flex-wrap gap-x-10 gap-y-10">
        <NuxtLink
            class=" light_border p-5 w-90 h-40 b-rounded-4 dark_bg_gray dark:text-dtc cursor-pointer custom_cursor_flow"
            v-for="(c, index) in categories"
            :key="index" :to="'/categories/' + c.route">
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
        <NuxtLink
            class="flex text-5 gap-x-2 b-rounded-4 border-solid border-gray-2 border-1 p-2 cursor-pointer custom_cursor_flow"
            v-for="(t, index) in tags"
            :key="index"
            :to="'/tags/' + t.route"
        >
          <span># {{ t.name }}</span>
          <span>{{ t.count }}</span>
        </NuxtLink>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import type {IResponse, IListData} from "~/api/http";
import type {ICategoryWithCount} from '~/api/category'
import type {ITagWithCount} from '~/api/tag'
import {getTagList} from '~/api/tag'
import {getCategoriesAndTags} from '~/api/category'

let categories = ref<ICategoryWithCount[]>([]);

const categoryAndTags = async () => {
  try {
    let httpRes: any = await getCategoriesAndTags()
    let res: IResponse<IListData<ICategoryWithCount>> = httpRes.data.value
    categories.value = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
categoryAndTags()

let tags = ref<ITagWithCount[]>([]);
const tagList = async () => {
  try {
    let httpRes: any = await getTagList()
    let res: IResponse<IListData<ITagWithCount>> = httpRes.data.value
    tags.value = res.data?.list || []
  } catch (error) {
    console.log(error);
  }
};
tagList()
</script>

<style scoped>
</style>