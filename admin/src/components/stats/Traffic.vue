<template>
  <a-card title="总数据" class="mt-5">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="trafficStatsLoading"
          @click="getTrafficStats"
        />
      </a-tooltip>
    </template>
    <a-spin :spinning="trafficStatsLoading">
      <a-flex gap="middle" horizontal>
        <a-statistic title="总访问量" :value="trafficStatsVO.view_count" class="w-33%" />
        <a-statistic title="总评论数" :value="trafficStatsVO.comment_count" class="w-33%" />
        <a-statistic title="总点赞数" :value="trafficStatsVO.like_count" class="w-33%" />
      </a-flex>
    </a-spin>
  </a-card>
</template>
<script setup lang="ts">
import { h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { GetTrafficStats, type TrafficStatsVO } from '@/interfaces/DataAnalysis'
import { message } from 'ant-design-vue'

const trafficStatsVO = ref<TrafficStatsVO>({
  view_count: 0,
  comment_count: 0,
  like_count: 0
})

const trafficStatsLoading = ref<boolean>(false)

const getTrafficStats = async () => {
  try {
    trafficStatsLoading.value = true
    const response: any = await GetTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    trafficStatsVO.value = response.data.data || trafficStatsVO.value
  } catch (error) {
    console.log(error)
  } finally {
    trafficStatsLoading.value = false
  }
}
getTrafficStats()
</script>
