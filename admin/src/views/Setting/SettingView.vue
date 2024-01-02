<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="basic" tab="站点信息"><BasicView /></a-tab-pane>
    <a-tab-pane key="email" tab="站长信息"><OnwerView /></a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置">Content of Tab Pane 3</a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置">Content of Tab Pane 3</a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import BasicView from '@/views/Setting/BasicView.vue'
import OnwerView from '@/views/Setting/OnwerView.vue'
const activeKey = ref('basic')
</script>
