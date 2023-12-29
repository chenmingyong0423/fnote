<template>
  <div class="bg-white">
    <div class="h-[50px] line-height-[50px] pl-3 custom_bottom_border_gray mb-2">
      2 评论
    </div>
    <div class="flex">
      <div class="w-8% min-h-[180px] ml-1% flex justify-center">
        <img :src="1" alt=""
             class="w-12 h-12 border-rounded-50%  cursor-pointer hover:rotate-360 ease-out duration-1000 lt-lg:mr0"
             v-if="false">
        <div class="i-ph-user-circle-duotone w-12 h-12 border-rounded-50%  lt-lg:mr0 text-gray-4" v-else></div>
      </div>
      <div class="w-91% flex flex-col ">
        <div class="text-gray-4 h-[55px] line-height-[35px]">
          <span>陈明勇</span> <span>发表于 2023-06-14 19:55:03</span>
        </div>
        <div>
          <v-md-preview :text="`# 123 \n123`"
                        class="lt-lg:important:p0" :class="{'dark': isBlackMode}"
          ></v-md-preview>
        </div>
        <div>
          <div class="w-15 h-8 custom_border_gray text-center line-height-8 b-rounded-2 cursor-pointer hover:bg-gray-1">回复</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">

import {useHomeStore} from "~/store/home";
const info = useHomeStore()

const isBlackMode = computed(() => info.isBlackMode)
</script>

<style scoped>
 /deep/ .github-markdown-body {
   padding: 0;
 }
</style>