<template>
  <div>
    <div
        class="item group flex p-5 bg-#fff b-rounded-4 h-50 cursor-pointer ease-linear duration-100 hover:drop-shadow-xl hover:translate-y--2 dark:text-dtc dark_bg_gray relative mb-5 hover:bg-gray-1 dark:hover:bg-#fff/20" v-for="(item, index) in props.posts" :key="index" @click="router.push('/posts/' + item.sug)">
      <div
          class="flex justify-between gap-x-3 z-99 w-auto absolute top-3 left--100% group-hover:left-1% ease-linear duration-200">
        <span class="bg-#2db7f5 rounded-3 text-white py-0.2em px-0.8em">{{ item.category }}</span>
        <span class="bg-orange rounded-3 text-white py-0.2em px-0.8em" v-for="(tag, idx) in item.tags" :key="idx">{{ tag }}</span>
      </div>
      <div class="w-1/3">
        <img class="object-contain max-w-full h-full"
             :src="item.cover_img"
             :alt="item.title" />
      </div>
      <div class="w-2/3 flex flex-col">
        <div class="mb-2 text-10 h-50 truncate">
          {{ item.title }}
        </div>
        <div class="flex-grow h-100 line-clamp-4">
          <p class="line-height-loose"> {{ item.summary}} </p>
        </div>
        <div class="flex gap-x-3 mt-2 h-20">
          <div class="flex gap-x-1 items-center"><span class="i-ph-eye"></span><span>{{ item.visit_count }}</span></div>
          <div class="flex gap-x-1 items-center"><span class="i-ph-thumbs-up-duotone"></span><span>{{ item.like_count }}</span></div>
          <div class="flex gap-x-1 items-center"><span class="i-ph-chats-duotone"></span><span>{{ item.comment_count }}</span></div>
          <div class="ml-auto flex gap-x-1 items-center"><span>{{ $dayjs(item.create_time).format('YYYY-MM-DD HH:mm:ss') }}</span></div>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import type {PropType} from "vue";
import type {IPost} from "~/server/api/post";
const router = useRouter()

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