<template>
  <a-card title="今日数据">
    <template #extra>
      <div class="dashboard-card-extra">
        <span v-if="lastUpdatedText" class="dashboard-updated-at">{{ lastUpdatedText }}</span>
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="todayTrafficLoading"
            @click="getTodayTrafficStats"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="todayTrafficLoading">
      <DashboardStatGrid :items="todayTrafficItems" min-width="150px" />
    </a-spin>
  </a-card>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { GetTodayTrafficStats, type TodayTrafficStatsVO } from '@/interfaces/DataAnalysis'
import { message } from 'ant-design-vue'
import DashboardStatGrid from '@/components/stats/DashboardStatGrid.vue'
import dayjs from 'dayjs'
import { toErrorMessage } from '@/utils/error'

const todayTrafficLoading = ref(false)
const lastUpdatedAt = ref<number>()
const todayTrafficStatsVO = ref<TodayTrafficStatsVO>({
  view_count: 0,
  user_view_count: 0,
  comment_count: 0,
  like_count: 0
})

const todayTrafficItems = computed(() => [
  {
    key: 'view_count',
    label: '访问量（PV）',
    value: todayTrafficStatsVO.value.view_count
  },
  {
    key: 'user_view_count',
    label: '访问用户',
    value: todayTrafficStatsVO.value.user_view_count
  },
  {
    key: 'comment_count',
    label: '评论数',
    value: todayTrafficStatsVO.value.comment_count
  },
  {
    key: 'like_count',
    label: '点赞数',
    value: todayTrafficStatsVO.value.like_count
  }
])

const lastUpdatedText = computed(() => {
  if (!lastUpdatedAt.value) {
    return ''
  }
  return `更新于 ${dayjs(lastUpdatedAt.value).format('HH:mm:ss')}`
})

const getTodayTrafficStats = async () => {
  try {
    todayTrafficLoading.value = true
    const response: any = await GetTodayTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.message || '今日数据加载失败')
      return
    }
    todayTrafficStatsVO.value = response.data.data || todayTrafficStatsVO.value
    lastUpdatedAt.value = Date.now()
  } catch (error) {
    message.error(toErrorMessage(error, '今日数据加载失败'))
  } finally {
    todayTrafficLoading.value = false
  }
}

getTodayTrafficStats()
</script>

<style scoped>
.dashboard-card-extra {
  display: flex;
  align-items: center;
  gap: 10px;
}

.dashboard-updated-at {
  color: var(--app-text-secondary);
  font-size: 12px;
  white-space: nowrap;
}
</style>
