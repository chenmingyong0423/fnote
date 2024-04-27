<template>
  <a-tabs v-model:activeKey="activeKey" :destroyInactiveTabPane="true">
    <a-tab-pane key="basic" tab="站点信息"><BasicView /></a-tab-pane>
    <a-tab-pane key="carousel" tab="轮播图配置"><CarouselView /></a-tab-pane>
    <a-tab-pane key="seo" tab="seo 配置"><SeoView /></a-tab-pane>
    <a-tab-pane key="verification" tab="站点验证"><VerificationView /></a-tab-pane>
    <a-tab-pane key="push" tab="文章推送配置"><PushView /></a-tab-pane>
    <a-tab-pane key="comment" tab="评论配置"><CommentSwitchView /></a-tab-pane>
    <a-tab-pane key="friend" tab="友链配置"><FriendSwitchView /></a-tab-pane>
    <a-tab-pane key="email" tab="邮件配置"><EmailView /></a-tab-pane>
    <a-tab-pane key="notice" tab="公告配置"><NoticeView /></a-tab-pane>
    <a-tab-pane key="front-post-count" tab="首页展示文章数量配置"
      ><FrontPostCountView
    /></a-tab-pane>
    <a-tab-pane key="pay" tab="支付二维码配置"><RecordView /></a-tab-pane>
    <a-tab-pane key="social" tab="社交信息配置"><SocialView /></a-tab-pane>
  </a-tabs>
</template>
<script lang="ts" setup>
import { ref } from 'vue'
import BasicView from '@/views/setting/BasicView.vue'
import SeoView from '@/views/setting/SeoView.vue'
import CommentSwitchView from '@/views/setting/CommentSwitchView.vue'
import FriendSwitchView from '@/views/setting/FriendSwitchView.vue'
import EmailView from '@/views/setting/EmailView.vue'
import NoticeView from '@/views/setting/NoticeView.vue'
import FrontPostCountView from '@/views/setting/FrontPostCountView.vue'
import RecordView from '@/views/setting/PayView.vue'
import SocialView from '@/views/setting/SocialView.vue'
import VerificationView from '@/views/setting/VerificationView.vue'
import PushView from '@/views/setting/PushView.vue'
import CarouselView from '@/views/setting/CarouselView.vue'
const activeKey = ref('basic')
</script>
