<template>
  <div>
    <ul class="list-none flex gap-x-10">
      <li class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
          :class=" route.path === '/' ? 'current' : '' "
          @click="router.push('/')">
        <span class="i-ph-house"/> <span>主页</span></li>
      <li class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue" :class=" route.path === '/category' ? 'current' : '' "
          @click="router.push('/category')"><span class="i-ph-list"/> <span>文章列表</span></li>
      <li class="flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue" v-for="(item, index) in props.items"
          :key="index" :class=" route.path === item.route ? 'current' : '' " @click="router.push(item.route)"><span
          class="i-ph-article-light"/> <span>{{ item.name }}</span></li>
      <li class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue" :class=" route.path === '/friend' ? 'current' : '' "
          @click="router.push('/friend')"><span class="i-ph-users"/> <span>友链</span></li>
      <li class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue" :class=" route.path === '/about-me' ? 'current' : '' "
          @click="router.push('/about-me')"><span class="i-ph-user"/> <span>关于我</span></li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import type {PropType} from "vue";
import type {IMenu} from "~/server/api/category"

const route = useRoute()
const router = useRouter()
const props = defineProps({
  items: {
    type: Array as PropType<IMenu[]>,
    default: () => [],
    required: true,
  },
})
</script>

<style scoped>
.current {
  color: #1e80ff;
}
</style>
