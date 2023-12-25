<template>
  <div>
    <div class="list-none flex gap-x-10">
      <NuxtLink class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
                :class=" route.path === '/' ? 'current' : '' "
                to="/">
        <span class="i-ph-house"/> <span>主页</span>
      </NuxtLink>
      <NuxtLink class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
                :class=" route.path === '/category' ? 'current' : '' "
                to="/categories"><span class="i-ph-list"/> <span>文章导航</span></NuxtLink>
      <NuxtLink class="flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
                v-for="(item, index) in props.items"
                :key="index" :class=" route.path === item.route || route.path.startsWith(item.route) ? 'current' : '' " :to="item.route"><span
          class="i-ph-article-light"/> <span>{{ item.name }}</span></NuxtLink>
      <NuxtLink class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
                :class=" route.path === '/friend' ? 'current' : '' "
                to="/friend"><span class="i-ph-users"/> <span>友链</span></NuxtLink>
      <NuxtLink class=" flex items-center gap-x-4 cursor-pointer p-2 hover:border_bottom_blue"
                :class=" route.path === '/about-me' ? 'current' : '' "
                to="/about-me"><span class="i-ph-user"/> <span>关于我</span></NuxtLink>
    </div>
  </div>
</template>

<script setup lang="ts">
import type {PropType} from "vue";
import type {IMenu} from "~/api/category"

const route = useRoute()
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
