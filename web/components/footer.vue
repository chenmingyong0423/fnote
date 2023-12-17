<template>
  <div
      class="footer_shadow flex flex-col justify-evenly items-center text-16 h120 dark_bg_gray dark:text-dtc shadow-xl">
    <div>
      Copyright © {{ year }} - Designed by
      <a class="decoration-none text-#000 font-bold hover:underline dark_text_white" href="github.com/chenmingyong0423/fnote">Fnote</a>
    </div>
    <div class="flex items-center lt-lg:flex-col">
            <span v-for="(item, index) in records" :key="index" v-html="item">
            </span>
    </div>
    <div>
    </div>
  </div>
</template>


<script lang="ts" setup>
import {useHomeStore} from '~/store/home';
const homeStore = useHomeStore()
const date = new Date(homeStore.master_info.website_live_time)
const year = date.getFullYear()

const records = homeStore.master_info.records
</script>

<style scoped></style>