<template>
  <div>
    <div
      class="duration-400 fixed z-999 backdrop-blur-md bg-white bg-opacity-40 md:hidden flex items-center justify-end custom_shadow"
      :class="info.showSmallScreenMenu ? 'inset-0' : 'inset-full'"
      @click="info.showSmallScreenMenu = false"
    ></div>

    <div
      class="fixed z-999 right-0 bg-white w-70% h-full bg-#f7f9fe custom_border_gray duration-300 md:hidden dark:text-dtc dark_bg_full_black"
      :class="info.showSmallScreenMenu ? ' translate-x-0' : ' translate-x-100%'"
    >
      <Profile></Profile>
      <Menu :items="homeStore.menuList" @click="info.showSmallScreenMenu = false"></Menu>
    </div>
  </div>
</template>
<script setup lang="ts">
import { useHomeStore } from "~/store/home";

const info = useHomeStore();
const homeStore = useHomeStore();
</script>
