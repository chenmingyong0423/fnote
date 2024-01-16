<template>
  <div class="flex flex-wrap justify-between">
    <NuxtLink
      class="item group flex flex-col items-center box-border p-5 bg-#fff b-rounded-4 w-49% lt-md:w-100% h-100 cursor-pointer duration-100 custom_shadow dark:text-dtc dark_bg_gray mb-5"
      v-for="(item, index) in props.posts"
      :key="index"
      :to="'/posts/' + item.sug"
      :title="item.title"
    >
      <div class="h-2/3 overflow-hidden relative w-full">
        <img
          class="object-contain w-full h-full group-hover:scale-110 duration-500"
          :src="item.cover_img"
          :alt="item.title"
        />
        <div
          class="flex flex-wrap gap-y-2 gap-x-3 z-99 w-auto absolute top-3 ease-linear duration-200"
        >
          <span
            class="bg-#22c55e rounded-3 text-white py-0.2em px-0.8em"
            v-if="item.sticky_weight == 1"
            >↑置顶</span
          >
          <span
            class="bg-#2db7f5 rounded-3 text-white py-0.2em px-0.8em"
            v-for="(category, idx) in item.categories"
            :key="idx"
            >{{ category }}</span
          >
          <span
            class="bg-orange rounded-3 text-white py-0.2em px-0.8em"
            v-for="(tag, idx) in item.tags"
            :key="idx"
            >{{ tag }}</span
          >
        </div>
      </div>
      <div class="h-1/3 overflow-hidden relative w-full flex flex-col">
        <div class="mb-2 text-6 h-15">
          {{ item.title }}
        </div>
        <div class="flex-grow">
          <p class="line-height-loose text-gray-5 truncate">
            {{ item.summary }}
          </p>
        </div>
        <div class="flex gap-x-3 h-10 mt-auto">
          <div class="flex gap-x-1 items-center">
            <span class="i-ph-eye"></span><span>{{ item.visit_count }}</span>
          </div>
          <div class="flex gap-x-1 items-center">
            <span class="i-ph-thumbs-up-duotone"></span
            ><span>{{ item.like_count }}</span>
          </div>
          <div class="flex gap-x-1 items-center">
            <span class="i-ph-chats-duotone"></span
            ><span>{{ item.comment_count }}</span>
          </div>
          <div class="ml-auto flex gap-x-1 items-center">
            <span>{{
              $dayjs(item.create_time * 1000).format("YYYY-MM-DD")
            }}</span>
          </div>
        </div>
      </div>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import type { PropType } from "vue";
import type { IPost } from "~/api/post";

const props = defineProps({
  posts: {
    type: Array as PropType<IPost[]>,
    default: () => [],
    required: true,
  },
});
</script>

<style scoped></style>
