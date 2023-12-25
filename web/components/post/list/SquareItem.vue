<template>
  <div class="flex flex-wrap justify-between">
    <NuxtLink
        class="item group flex flex-col items-center p-5 bg-#fff b-rounded-4 w-150 h-110 cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--2 dark:text-dtc dark_bg_gray mb-5"
        v-for="(item, index) in props.posts" :key="index" :to="'/posts/' + item.sug">
      <div class="h-2/3 overflow-hidden relative w-full">
        <img class="object-contain w-full h-full"
             :src="item.cover_img"
             :alt="item.title"/>
      </div>
      <div class="h-1/3 overflow-hidden relative w-full flex flex-col">
        <div class="mb-2 text-10 h-15 truncate">
          {{ item.title }}
        </div>
        <div class="flex-grow">
          <p class="line-height-loose text-gray-5"> {{ item.summary }} </p>
        </div>
        <div class="flex gap-x-3 h-10 mt-auto">
          <div class="flex gap-x-1 items-center"><span class="i-ph-eye"></span><span>{{ item.visit_count }}</span></div>
          <div class="flex gap-x-1 items-center"><span class="i-ph-thumbs-up-duotone"></span><span>{{
              item.like_count
            }}</span></div>
          <div class="flex gap-x-1 items-center"><span class="i-ph-chats-duotone"></span><span>{{
              item.comment_count
            }}</span></div>
          <div class="ml-auto flex gap-x-1 items-center"><span>{{
              $dayjs(item.create_time).format('YYYY-MM-DD HH:mm:ss')
            }}</span></div>
        </div>
      </div>
    </NuxtLink>
  </div>
</template>


<script setup lang="ts">
import type {PropType} from "vue";
import type {IPost} from "~/api/post";

const props = defineProps({
  posts: {
    type: Array as PropType<IPost[]>,
    default: () => [],
    required: true,
  },
})
</script>

<style scoped>
.line-clamp-4 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3; /* 可以调整这个数值来设置行数 */
  overflow: hidden;
}

.item::before {
  content: '';
  position: absolute;
  width: 99%;
  transform: scaleX(0);
  height: 2px;
  bottom: 0;
  left: 0.5%;
  background-color: #0087ca;
  transform-origin: bottom;
  transition: transform 0.3s ease-out;
}

.item:hover::before {
  transform: scaleX(1);
  transform-origin: right left;
}
</style>