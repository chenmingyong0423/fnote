<template>
  <div>
    <a
        class="slide-up item group flex p-5 bg-#fff b-rounded-4 h-50 cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--2 dark:text-dtc dark_bg_gray mb-5"
        v-for="(item, index) in props.posts"
        :key="index"
        :href="baseUrl + '/posts/' +item.sug"
        target="_blank"
        :title="item.title"
    >
      <div class="w-1/3 overflow-hidden relative">
        <img
            class="object-contain max-w-full h-full"
            :src="serverHost + item.cover_img"
            :alt="item.title"
        />
        <div
            class="flex flex-wrap gap-x-3 w-full z-99 absolute top-3 left--100% group-hover:left-1% ease-linear duration-200"
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
      <div class="w-2/3 flex flex-col">
        <div class="mb-2 text-6 h-50">
          {{ item.title }}
        </div>
        <div class="flex-grow h-100 line-clamp-4">
          <p class="line-height-loose text-gray-5">{{ item.summary }}</p>
        </div>
        <div class="flex gap-x-3 mt-2 h-20">
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
    </a>
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
});
const runtimeConfig = useRuntimeConfig()
const baseUrl = runtimeConfig.public.domain;
const serverHost = runtimeConfig.public.serverHost;
</script>

<style scoped>
.line-clamp-4 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3; /* 可以调整这个数值来设置行数 */
  overflow: hidden;
}

.item::before {
  content: "";
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

@keyframes slideUp {
  0% {
    transform: translateY(+100%);
  }
  100% {
    transform: translateY(0);
  }
}

.slide-up {
  animation: slideUp 0.5s ease;
  animation-iteration-count: 1;
}
</style>
