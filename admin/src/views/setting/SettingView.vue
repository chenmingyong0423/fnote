<template>
  <div class="setting-page">
    <div class="setting-heading">
      <div>
        <h1>博客设置</h1>
        <span>{{ activeSetting.label }}</span>
      </div>
      <a-tooltip title="重新检查配置">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="healthLoading"
          @click="loadConfigHealth"
        />
      </a-tooltip>
    </div>

    <a-alert
      v-if="health"
      class="health-alert"
      :type="health.required.completed ? (activeItems.length ? 'warning' : 'success') : 'error'"
      show-icon
    >
      <template #message>
        <div class="health-summary">
          <span class="health-score">配置健康分 {{ health.score }}</span>
          <a-space wrap size="small">
            <a-tag :color="health.required.completed ? 'success' : 'error'">
              必须 {{ health.required.done }}/{{ health.required.total }}
            </a-tag>
            <a-tag :color="health.recommended.completed ? 'success' : 'warning'">
              推荐 {{ health.recommended.done }}/{{ health.recommended.total }}
            </a-tag>
            <a-tag v-if="activeItems.length" color="warning">待完善 {{ activeItems.length }}</a-tag>
          </a-space>
          <a-button
            v-if="activeItems.length || quietItems.length"
            type="link"
            size="small"
            @click="healthExpanded = !healthExpanded"
          >
            {{ healthExpanded ? '收起' : '查看详情' }}
            <template #icon>
              <DownOutlined :class="{ 'is-expanded': healthExpanded }" />
            </template>
          </a-button>
        </div>
      </template>
      <template v-if="healthExpanded" #description>
        <div class="health-details">
          <div v-if="activeItems.length" class="health-list">
            <div v-for="item in activeItems" :key="item.key" class="health-item">
              <div class="health-item-content">
                <a-tag :color="levelColor(item.level)">{{ levelText(item.level) }}</a-tag>
                <div>
                  <div class="font-medium">{{ item.label }}</div>
                  <div v-if="item.missing_fields?.length" class="health-item-note">
                    缺少：{{ item.missing_fields.join('、') }}
                  </div>
                </div>
              </div>
              <a-space wrap>
                <a-button size="small" type="link" @click="openItem(item)">去配置</a-button>
                <a-button
                  size="small"
                  type="link"
                  :loading="updatingKey === item.key"
                  @click="snoozeItem(item)"
                >
                  7 天后提醒
                </a-button>
                <a-button
                  size="small"
                  type="link"
                  danger
                  :loading="updatingKey === item.key"
                  @click="ignoreItem(item)"
                >
                  忽略
                </a-button>
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

    <a-select
      class="setting-select"
      :value="activeKey"
      :options="settingOptions"
      @change="switchTab"
    />

    <div class="setting-layout">
      <nav class="setting-nav" aria-label="博客设置导航">
        <a-menu
          mode="inline"
          :selected-keys="[activeKey]"
          :inline-indent="16"
          @click="handleMenuClick"
        >
          <a-menu-item-group v-for="group in settingGroups" :key="group.key" :title="group.label">
            <a-menu-item v-for="item in group.items" :key="item.key">
              <template #icon><component :is="item.icon" /></template>
              {{ item.label }}
            </a-menu-item>
          </a-menu-item-group>
        </a-menu>
      </nav>

      <main class="setting-content">
        <component :is="activeSetting.component" :key="activeSetting.key" />
      </main>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, h, onMounted, ref, watch, type Component } from 'vue'
import { message, type MenuProps } from 'ant-design-vue'
import {
  CommentOutlined,
  DownOutlined,
  FileSearchOutlined,
  FileTextOutlined,
  GlobalOutlined,
  LinkOutlined,
  MailOutlined,
  NotificationOutlined,
  OrderedListOutlined,
  PictureOutlined,
  QrcodeOutlined,
  ReloadOutlined,
  SafetyCertificateOutlined,
  SearchOutlined,
  SendOutlined,
  ShareAltOutlined
} from '@ant-design/icons-vue'
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
import RobotsView from '@/views/setting/RobotsView.vue'
import {
  GetConfigHealth,
  UpdateConfigCheckState,
  type ConfigHealth,
  type ConfigHealthItem,
  type ConfigHealthLevel,
  type ConfigHealthStatus
} from '@/interfaces/Config'

interface SettingItem {
  key: string
  label: string
  icon: Component
  component: Component
}

interface SettingGroup {
  key: string
  label: string
  items: SettingItem[]
}

const settingGroups: SettingGroup[] = [
  {
    key: 'site',
    label: '站点与展示',
    items: [
      { key: 'basic', label: '站点信息', icon: GlobalOutlined, component: BasicView },
      { key: 'carousel', label: '首页轮播', icon: PictureOutlined, component: CarouselView },
      { key: 'notice', label: '站点公告', icon: NotificationOutlined, component: NoticeView },
      {
        key: 'front-post-count',
        label: '首页文章数量',
        icon: OrderedListOutlined,
        component: FrontPostCountView
      }
    ]
  },
  {
    key: 'discovery',
    label: '搜索与分发',
    items: [
      { key: 'seo', label: 'SEO 元信息', icon: SearchOutlined, component: SeoView },
      { key: 'sitemap', label: '站点地图', icon: FileSearchOutlined, component: SiteMapVIew },
      { key: 'robots', label: 'Robots 文件', icon: FileTextOutlined, component: RobotsView },
      {
        key: 'verification',
        label: '站点验证',
        icon: SafetyCertificateOutlined,
        component: VerificationView
      },
      { key: 'push', label: '文章推送', icon: SendOutlined, component: PushView }
    ]
  },
  {
    key: 'interaction',
    label: '互动与联系',
    items: [
      { key: 'comment', label: '评论设置', icon: CommentOutlined, component: CommentSwitchView },
      { key: 'friend', label: '友链申请', icon: LinkOutlined, component: FriendSwitchView },
      { key: 'email', label: '邮件通知', icon: MailOutlined, component: EmailView },
      { key: 'social', label: '社交信息', icon: ShareAltOutlined, component: SocialView },
      { key: 'pay', label: '支付二维码', icon: QrcodeOutlined, component: RecordView }
    ]
  }
]

const settingItems = settingGroups.flatMap((group) => group.items)
const settingOptions = settingGroups.map((group) => ({
  label: group.label,
  options: group.items.map((item) => ({ label: item.label, value: item.key }))
}))
const defaultKey = 'basic'
const route = useRoute()
const router = useRouter()

const resolveKey = (tab: unknown) =>
  typeof tab === 'string' && settingItems.some((item) => item.key === tab) ? tab : defaultKey

const activeKey = ref(resolveKey(route.query.tab))
const activeSetting = computed(
  () => settingItems.find((item) => item.key === activeKey.value) || settingItems[0]
)
const health = ref<ConfigHealth>()
const healthLoading = ref(false)
const healthExpanded = ref(false)
const updatingKey = ref('')

const activeItems = computed(
  () => health.value?.items.filter((item) => item.status === 'missing') || []
)
const quietItems = computed(
  () =>
    health.value?.items.filter((item) => item.status === 'ignored' || item.status === 'snoozed') ||
    []
)

const switchTab = (key: string | number) => {
  const nextKey = resolveKey(String(key))
  activeKey.value = nextKey
  if (route.query.tab !== nextKey) {
    router.replace({ path: route.path, query: { ...route.query, tab: nextKey } })
  }
}

const handleMenuClick: MenuProps['onClick'] = ({ key }) => switchTab(String(key))

const openItem = (item: ConfigHealthItem) => {
  const tab = item.href ? new URL(item.href, window.location.origin).searchParams.get('tab') : null
  switchTab(tab || activeKey.value)
}

const loadConfigHealth = async () => {
  try {
    healthLoading.value = true
    const response = await GetConfigHealth()
    health.value = response.data.data
  } catch (error) {
    console.error(error)
    message.error('配置检查加载失败')
  } finally {
    healthLoading.value = false
  }
}

const updateState = async (
  item: ConfigHealthItem,
  status: 'active' | 'ignored' | 'snoozed',
  snoozeDays?: number
) => {
  try {
    updatingKey.value = item.key
    await UpdateConfigCheckState(item.key, { status, snooze_days: snoozeDays })
    await loadConfigHealth()
  } finally {
    updatingKey.value = ''
  }
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

const levelText = (level: ConfigHealthLevel) =>
  ({ required: '必须', recommended: '推荐', optional: '可选' })[level]

const levelColor = (level: ConfigHealthLevel) =>
  ({ required: 'error', recommended: 'warning', optional: 'processing' })[level]

const statusColor = (status: ConfigHealthStatus) =>
  ({ ok: 'success', missing: 'warning', ignored: 'default', snoozed: 'processing' })[status]

const statusText = (item: ConfigHealthItem) => {
  if (item.status === 'ignored') return '已忽略'
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
    const nextKey = resolveKey(tab)
    activeKey.value = nextKey
    if (tab !== nextKey) {
      router.replace({ path: route.path, query: { ...route.query, tab: nextKey } })
    }
  }
)

onMounted(() => {
  document.title = '博客设置 - 后台管理'
  if (route.query.tab !== activeKey.value) {
    router.replace({ path: route.path, query: { ...route.query, tab: activeKey.value } })
  }
  loadConfigHealth()
})
</script>

<style scoped>
.setting-page {
  min-width: 0;
}

.setting-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.setting-heading h1 {
  margin: 0;
  font-size: 22px;
  line-height: 1.4;
  letter-spacing: 0;
}

.setting-heading span {
  color: var(--app-text-secondary);
  font-size: 13px;
}

.health-alert {
  margin-bottom: 16px;
}

.health-summary,
.health-item,
.health-item-content {
  display: flex;
  align-items: center;
}

.health-summary {
  flex-wrap: wrap;
  gap: 10px;
}

.health-summary > :last-child {
  margin-left: auto;
}

.health-score {
  font-weight: 600;
}

.health-details,
.health-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.health-details {
  padding-top: 8px;
}

.health-item {
  justify-content: space-between;
  gap: 12px;
  padding: 10px 12px;
  background: var(--app-surface-translucent);
  border-radius: 6px;
}

.health-item-content {
  min-width: 0;
  gap: 6px;
}

.health-item-note {
  color: var(--app-text-secondary);
  font-size: 12px;
}

.health-summary :deep(.anticon-down) {
  transition: transform 0.2s;
}

.health-summary :deep(.anticon-down.is-expanded) {
  transform: rotate(180deg);
}

.setting-layout {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  align-items: start;
  gap: 20px;
}

.setting-nav {
  position: sticky;
  top: 16px;
  overflow: hidden;
  border: 1px solid var(--app-border);
  border-radius: 6px;
}

.setting-nav :deep(.ant-menu) {
  border-inline-end: 0;
}

.setting-content {
  min-width: 0;
}

.setting-select {
  display: none;
  width: 100%;
  margin-bottom: 16px;
}

@media (max-width: 768px) {
  .setting-layout {
    display: block;
  }

  .setting-nav {
    display: none;
  }

  .setting-select {
    display: block;
  }

  .health-item {
    align-items: flex-start;
    flex-direction: column;
  }

  .health-summary > :last-child {
    margin-left: 0;
  }
}
</style>
