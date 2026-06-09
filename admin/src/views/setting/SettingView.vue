<template>
  <div>
    <a-alert
      v-if="health"
      class="mb-4"
      :type="health.required.completed ? (activeItems.length ? 'warning' : 'success') : 'error'"
      show-icon
    >
      <template #message>
        配置健康分：{{ health.score }}，必须项 {{ health.required.done }}/{{
          health.required.total
        }}，推荐项 {{ health.recommended.done }}/{{ health.recommended.total }}
      </template>
      <template #description>
        <div class="flex flex-col gap-3">
          <a-space wrap>
            <a-tag :color="health.required.completed ? 'success' : 'error'">
              必须项：{{ health.required.done }}/{{ health.required.total }}
            </a-tag>
            <a-tag :color="health.recommended.completed ? 'success' : 'warning'">
              推荐项：{{ health.recommended.done }}/{{ health.recommended.total }}
            </a-tag>
            <a-tag color="processing">
              可选增强：{{ health.optional.done }}/{{ health.optional.total }}
            </a-tag>
          </a-space>

          <div v-if="activeItems.length" class="flex flex-col gap-2">
            <div
              v-for="item in activeItems"
              :key="item.key"
              class="flex flex-wrap items-center justify-between gap-2 rounded bg-white/70 px-3 py-2"
            >
              <div>
                <a-tag :color="levelColor(item.level)">{{ levelText(item.level) }}</a-tag>
                <span class="font-medium">{{ item.label }}</span>
                <span class="text-gray-500">
                  {{ item.missing_fields?.length ? `缺少：${item.missing_fields.join('、')}` : '' }}
                </span>
              </div>
              <a-space>
                <a-button size="small" type="link" @click="openItem(item)">
                  <template #icon><ArrowRightOutlined /></template>
                  去配置
                </a-button>
                <a-button size="small" type="link" @click="snoozeItem(item)">7 天后提醒</a-button>
                <a-button size="small" type="link" danger @click="ignoreItem(item)">忽略</a-button>
              </a-space>
            </div>
          </div>

          <a-space v-if="quietItems.length" wrap>
            <a-tag v-for="item in quietItems" :key="item.key" :color="statusColor(item.status)">
              {{ item.label }}：{{ statusText(item) }}
              <a class="ml-2" @click="activateItem(item)">重新提醒</a>
            </a-tag>
          </a-space>
        </div>
      </template>
    </a-alert>

    <a-tabs v-model:activeKey="activeKey" :destroyInactiveTabPane="true" @change="switchTab">
      <a-tab-pane key="basic" tab="站点信息"><BasicView /></a-tab-pane>
      <a-tab-pane key="carousel" tab="轮播图配置"><CarouselView /></a-tab-pane>
      <a-tab-pane key="seo" tab="seo 配置"><SeoView /></a-tab-pane>
      <a-tab-pane key="sitemap" tab="sitemap 配置"><SiteMapVIew /></a-tab-pane>
      <a-tab-pane key="verification" tab="站点验证"><VerificationView /></a-tab-pane>
      <a-tab-pane key="push" tab="文章推送配置"><PushView /></a-tab-pane>
      <a-tab-pane key="comment" tab="评论配置"><CommentSwitchView /></a-tab-pane>
      <a-tab-pane key="friend" tab="友链配置"><FriendSwitchView /></a-tab-pane>
      <a-tab-pane key="email" tab="邮件配置"><EmailView /></a-tab-pane>
      <a-tab-pane key="notice" tab="公告配置"><NoticeView /></a-tab-pane>
      <a-tab-pane key="front-post-count" tab="首页展示文章数量配置"><FrontPostCountView /></a-tab-pane>
      <a-tab-pane key="pay" tab="支付二维码配置"><RecordView /></a-tab-pane>
      <a-tab-pane key="social" tab="社交信息配置"><SocialView /></a-tab-pane>
    </a-tabs>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, watch } from 'vue'
import { message } from 'ant-design-vue'
import { ArrowRightOutlined } from '@ant-design/icons-vue'
import { useRoute, useRouter } from 'vue-router'
import BasicView from '@/views/setting/BasicView.vue'
import SeoView from '@/views/setting/SeoView.vue'
import CommentSwitchView from '@/views/setting/CommentSwitchView.vue'
import FriendSwitchView from '@/views/setting/FriendView.vue'
import EmailView from '@/views/setting/EmailView.vue'
import NoticeView from '@/views/setting/NoticeView.vue'
import FrontPostCountView from '@/views/setting/FrontPostCountView.vue'
import RecordView from '@/views/setting/PayView.vue'
import SocialView from '@/views/setting/SocialView.vue'
import VerificationView from '@/views/setting/VerificationView.vue'
import PushView from '@/views/setting/PushView.vue'
import CarouselView from '@/views/setting/CarouselView.vue'
import SiteMapVIew from '@/views/setting/SiteMapVIew.vue'
import {
  GetConfigHealth,
  UpdateConfigCheckState,
  type ConfigHealth,
  type ConfigHealthItem,
  type ConfigHealthLevel,
  type ConfigHealthStatus
} from '@/interfaces/Config'

const route = useRoute()
const router = useRouter()
const activeKey = ref(typeof route.query.tab === 'string' ? route.query.tab : 'basic')
const health = ref<ConfigHealth>()

const activeItems = computed(
  () => health.value?.items.filter((item) => item.status === 'missing') || []
)
const quietItems = computed(
  () =>
    health.value?.items.filter((item) => item.status === 'ignored' || item.status === 'snoozed') ||
    []
)

const switchTab = (key: string | number) => {
  activeKey.value = String(key)
  router.replace({ path: route.path, query: { ...route.query, tab: activeKey.value } })
}

const openItem = (item: ConfigHealthItem) => {
  const tab = item.href?.split('tab=')[1]
  switchTab(tab || activeKey.value)
}

const loadConfigHealth = async () => {
  const response = await GetConfigHealth()
  health.value = response.data.data
}

const updateState = async (
  item: ConfigHealthItem,
  status: 'active' | 'ignored' | 'snoozed',
  snoozeDays?: number
) => {
  await UpdateConfigCheckState(item.key, { status, snooze_days: snoozeDays })
  await loadConfigHealth()
}

const snoozeItem = async (item: ConfigHealthItem) => {
  await updateState(item, 'snoozed', 7)
  message.success('已设置 7 天后提醒')
}

const ignoreItem = async (item: ConfigHealthItem) => {
  await updateState(item, 'ignored')
  message.success('已忽略该提示')
}

const activateItem = async (item: ConfigHealthItem) => {
  await updateState(item, 'active')
  message.success('已恢复提醒')
}

const levelText = (level: ConfigHealthLevel) => {
  const mapping = {
    required: '必须',
    recommended: '推荐',
    optional: '可选'
  }
  return mapping[level]
}

const levelColor = (level: ConfigHealthLevel) => {
  const mapping = {
    required: 'error',
    recommended: 'warning',
    optional: 'processing'
  }
  return mapping[level]
}

const statusColor = (status: ConfigHealthStatus) => {
  const mapping = {
    ok: 'success',
    missing: 'warning',
    ignored: 'default',
    snoozed: 'processing'
  }
  return mapping[status]
}

const statusText = (item: ConfigHealthItem) => {
  if (item.status === 'ignored') {
    return '已忽略'
  }
  if (item.status === 'snoozed') {
    return item.snoozed_until
      ? `${new Date(item.snoozed_until * 1000).toLocaleDateString()} 后提醒`
      : '稍后提醒'
  }
  return item.configured ? '已配置' : '未配置'
}

watch(
  () => route.query.tab,
  (tab) => {
    if (typeof tab === 'string') {
      activeKey.value = tab
    }
  }
)

onMounted(() => {
  loadConfigHealth()
})
</script>
