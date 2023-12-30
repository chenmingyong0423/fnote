<template>
  <div ref="header"
       class="bg-#fff backdrop-blur-20 fixed top-0 dark:text-dtc w-full z-99 flex justify-between items-center p-1 dark_bg_gray duration-200 ease-linear max-h-[70px] select-none">
    <div>
      <NuxtLink to="/">
        <img :src="picture" alt=""
             class="w-15 h-15 border-rounded-50% mx5 cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0 select-none">
      </NuxtLink>
    </div>
    <div class="bg_transparent  lt-md:hidden">
      <Menu :items="homeStore.menuList"></Menu>
    </div>
    <div class="ml-auto pr-5 flex gap-x-4">
      <div
          class="i-ph-sun-dim-duotone dark:i-ph-moon-stars-fill cursor-pointer text-10 text-#86909c dark:text-dtc dark:hover:text-white"
          @click="homeStore.isBlackMode = !homeStore.isBlackMode"></div>
      <NuxtLink
          class="i-ph-list-magnifying-glass-duotone cursor-pointer dark:text-dtc text-10 text-#86909c dark:hover:text-white"
          to="/search?keyword=">
      </NuxtLink>
      <div
          class="i-grommet-icons:github text-12 dark:text-dtc cursor-pointer dark:hover:text-white active:bg-#e5e5e5  lt-md:hidden"/>
      <div class="i-ph:list text-10 text-#86909c dark:text-dtc cursor-pointer dark:hover:text-white active:bg-#e5e5e5  md:hidden" @click="homeStore.showSmallScreenMenu = true"></div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import {ref, onMounted} from 'vue';
import {useHomeStore} from '~/store/home';
import Menu from "~/components/Menu.vue";

const homeStore = useHomeStore()

const isBlackMode = computed(() => homeStore.isBlackMode)
const picture = ref<string>(homeStore.master_info.picture)


watch(isBlackMode, (newValue) => {
  localStorage.setItem("isBlackMode", newValue.toString())
  if (newValue) {

    document.querySelector('html')!.classList.add("dark");
  } else {
    document.querySelector('html')!.classList.remove("dark");
  }
})

const header = ref()
let scrollCount = ref(0)
const headerScroll = () => {
  if (document.documentElement.scrollTop > scrollCount.value) {
    header.value.setAttribute('style', `top:-${header.value.clientHeight}px`);
  } else {
    header.value.setAttribute('style', 'top:0px');
  }
  scrollCount.value = document.documentElement.scrollTop
}
onMounted(() => {
  window.addEventListener('scroll', headerScroll)
})
onBeforeUnmount(() => {
  window.removeEventListener('scroll', headerScroll)
})
</script>

<style scoped>
</style>
