<template>
  <a-tabs v-model:activeKey="activeKey">
    <a-tab-pane key="basic" tab="站点信息"><BasicView /></a-tab-pane>
    <a-tab-pane key="owner" tab="站长信息"><OwnerView /></a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置"><SeoView /></a-tab-pane>
    <a-tab-pane key="comment" tab="评论配置"><CommentSwitchView /></a-tab-pane>
    <a-tab-pane key="friend" tab="友链配置"><FriendSwitchView /></a-tab-pane>
    <a-tab-pane key="email" tab="邮件配置"><EmailView /></a-tab-pane>
    <a-tab-pane key="notice" tab="公告配置"><NoticeView /></a-tab-pane>
    <a-tab-pane key="front-post-count" tab="首页展示文章数量配置"
      ><FrontPostCountView
    /></a-tab-pane>
    <a-tab-pane key="records" tab="备案信息"><RecordView /></a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import BasicView from '@/views/setting/BasicView.vue'
import OwnerView from '@/views/setting/OwnerView.vue'
import SeoView from '@/views/setting/SeoView.vue'
import CommentSwitchView from '@/views/setting/CommentSwitchView.vue'
import FriendSwitchView from '@/views/setting/FriendSwitchView.vue'
import EmailView from '@/views/setting/EmailView.vue'
import NoticeView from '@/views/setting/NoticeView.vue'
import FrontPostCountView from '@/views/setting/FrontPostCountView.vue'
import RecordView from '@/views/setting/RecordView.vue'
const activeKey = ref('basic')
</script>
