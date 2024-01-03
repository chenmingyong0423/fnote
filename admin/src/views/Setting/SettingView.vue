<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="basic" tab="站点信息"><BasicView /></a-tab-pane>
    <a-tab-pane key="email" tab="站长信息"><OwnerView /></a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置"><SeoView /></a-tab-pane>
    <a-tab-pane key="comment switch" tab="评论配置"><CommentSwitchView /></a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import BasicView from '@/views/Setting/BasicView.vue'
import OwnerView from '@/views/Setting/OwnerView.vue'
import SeoView from '@/views/Setting/SeoView.vue'
import CommentSwitchView from '@/views/Setting/CommentSwitchView.vue'
const activeKey = ref('basic')
</script>
