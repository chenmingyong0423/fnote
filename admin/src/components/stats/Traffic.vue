<template>
  <a-card title="总数据">
    <template #extra>
      <div class="dashboard-card-extra">
        <span v-if="lastUpdatedText" class="dashboard-updated-at">{{ lastUpdatedText }}</span>
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="trafficStatsLoading"
            @click="getTrafficStats"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="trafficStatsLoading">
      <DashboardStatGrid :items="trafficItems" />
    </a-spin>
  </a-card>
</template>

<script setup lang="ts">
import { computed, h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { GetTrafficStats, type TrafficStatsVO } from '@/interfaces/DataAnalysis'
import { message } from 'ant-design-vue'
import DashboardStatGrid from '@/components/stats/DashboardStatGrid.vue'
import dayjs from 'dayjs'
import { toErrorMessage } from '@/utils/error'

const trafficStatsVO = ref<TrafficStatsVO>({
  view_count: 0,
  comment_count: 0,
  like_count: 0
})

const trafficStatsLoading = ref(false)
const lastUpdatedAt = ref<number>()

const trafficItems = computed(() => [
  {
    key: 'view_count',
    label: '总访问量',
    value: trafficStatsVO.value.view_count
  },
  {
    key: 'comment_count',
    label: '总评论数',
    value: trafficStatsVO.value.comment_count
  },
  {
    key: 'like_count',
    label: '总点赞数',
    value: trafficStatsVO.value.like_count
  }
])

const lastUpdatedText = computed(() => {
  if (!lastUpdatedAt.value) {
    return ''
  }
  return `更新于 ${dayjs(lastUpdatedAt.value).format('HH:mm:ss')}`
})

const getTrafficStats = async () => {
  try {
    trafficStatsLoading.value = true
    const response: any = await GetTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.message || '总数据加载失败')
      return
    }
    trafficStatsVO.value = response.data.data || trafficStatsVO.value
    lastUpdatedAt.value = Date.now()
  } catch (error) {
    message.error(toErrorMessage(error, '总数据加载失败'))
  } finally {
    trafficStatsLoading.value = false
  }
}

getTrafficStats()
</script>

<style scoped>
.dashboard-card-extra {
  display: flex;
  align-items: center;
  gap: 10px;
}

.dashboard-updated-at {
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  white-space: nowrap;
}
</style>
