<template>
  <div class="flex flex-col">
    <div class="flex"><a-input></a-input> <a-button>添加</a-button></div>
    <div
      class="flex p-3 border-b-1 border-b-solid border-b-gray-2"
      v-for="(item, index) in records"
      :key="index"
    >
      <div v-html="item"></div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'

const records = ref([
  '<a href="https://beian.miit.gov.cn/" target="_blank"\n' +
    '          style="height:20px;line-height:20px;margin: 0px 0px 0px 5px;">琼ICP备2023002354号</a>',
  '<a target="_blank" href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=44030602007624"\n' +
    '            style="display:inline-block;text-decoration:none;height:20px;line-height:20px;"><img src="/static/备案图标.png"\n' +
    '              style="float:left;" />\n' +
    '            <p style="float:left;height:20px;line-height:20px;margin: 0px 0px 0px 5px; color:#939393;">粤公网安备\n' +
    '              44030602007624号</p>\n' +
    '          </a>'
])
</script>
