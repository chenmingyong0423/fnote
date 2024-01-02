<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="basic" tab="站点设置"><BasicView /></a-tab-pane>
    <a-tab-pane key="email" tab="邮件配置" force-render>Content of Tab Pane 2</a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置">Content of Tab Pane 3</a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置">Content of Tab Pane 3</a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import { ref } from 'vue';
import BasicView from '@/views/Setting/BasicView.vue'
const activeKey = ref('basic');
</script>