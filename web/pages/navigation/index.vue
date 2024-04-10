<template>
  <div class="text-black bg-white p-5 b-rounded-4 dark_bg_gray dark:text-dtc">
    <div class="mb-10">
      <div
        class="line-height-10 text-10 light_border_bottom p-b-5 dark_text_white"
      >
        分类
      </div>
      <div
        class="mt-5 flex items-center flex-wrap gap-x-10 gap-y-10 lt-md:gap-x-5 lt-md:gap-y-2"
      >
        <NuxtLink
          class="light_border p-5 w-90 h-40 b-rounded-4 dark_bg_gray dark:text-dtc cursor-pointer custom_cursor_flow custom_shadow lt-md:w-36%"
          v-for="(c, index) in categories"
          :key="index"
          :to="'/categories/' + c.route"
        >
          <div class="i-ph-columns-fill w-10 h-10 text-gray-5 p-y-1"></div>
          <div class="text-8 font-bold p-y-1 dark_text_white">{{ c.name }}</div>
          <div class="p-y-1 truncate">{{ c.description }}</div>
          <div class="p-y-1 flex items-center gap-x-1 h-8 text-5">
            <div class="i-ph-book dark:text-dtc"></div>
            <div>{{ c.count }}</div>
          </div>
        </NuxtLink>
      </div>
    </div>

    <div>
      <div
        class="line-height-10 text-10 light_border_bottom p-b-5 dark_text_white"
      >
        标签
      </div>
      <div
        class="mt-5 flex items-center flex-wrap gap-x-10 gap-y-10 lt-md:gap-3"
      >
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
import type { IResponse, IListData } from "~/api/http";
import type { ICategoryWithCount } from "~/api/category";
import type { ITagWithCount } from "~/api/tag";
import { getTagList } from "~/api/tag";
import { getCategoriesAndTags } from "~/api/category";
import { useConfigStore } from "~/store/config";

let categories = ref<ICategoryWithCount[]>([]);
let ct = ref<string>("");

const categoryAndTags = async () => {
  try {
    let httpRes: any = await getCategoriesAndTags();
    let res: IResponse<IListData<ICategoryWithCount>> = httpRes.data.value;
    categories.value = res.data?.list || [];
    categories.value.forEach((item: ICategoryWithCount) => {
      ct.value += item.name + "、";
    });
  } catch (error) {
    console.log(error);
  }
};

let tags = ref<ITagWithCount[]>([]);
const tagList = async () => {
  try {
    let httpRes: any = await getTagList();
    let res: IResponse<IListData<ITagWithCount>> = httpRes.data.value;
    tags.value = res.data?.list || [];
    tags.value.forEach((item: ITagWithCount) => {
      ct.value += item.name + "、";
    });
  } catch (error) {
    console.log(error);
  }
};

await categoryAndTags();
await tagList();
ct.value = ct.value.substring(0, ct.value.length - 1);

const configStore = useConfigStore();
useHead({
  title: `文章分类和标签 - ${configStore.seo_meta_config.title === '' ? configStore.website_info.website_name : configStore.seo_meta_config.title}`,
  meta: [
    {
      name: "description",
      content: `所有的文章分类和标签，包括${ct.value}等不同主题。`,
    },
    { name: "keywords", content: configStore.seo_meta_config.keywords },
    { name: "author", content: configStore.seo_meta_config.author },
    { name: "robots", content: configStore.seo_meta_config.robots },
  ],
  link: [
    { rel: "icon", type: "image/x-icon", href: configStore.website_info.website_icon },
  ],
});
useSeoMeta({
  ogTitle: `文章分类和标签 - ${configStore.seo_meta_config.og_title === '' ? configStore.website_info.website_name : configStore.seo_meta_config.og_title}`,
  ogDescription: `所有的文章分类和标签，包括${ct.value}等不同主题。`,
  ogImage: configStore.seo_meta_config.og_image,
  twitterCard: "summary",
});
</script>

<style scoped></style>
