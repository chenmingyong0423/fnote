<template>
  <a-card title="今日数据">
    <a-flex gap="middle" horizontal>
      <a-statistic title="今日访问量（PV）" :value="todayTrafficStatsVO.view_count" class="w-25%" />
      <a-statistic title="今日访问用户" :value="todayTrafficStatsVO.user_view_count" class="w-25%" />
      <a-statistic title="今日评论数" :value="todayTrafficStatsVO.comment_count" class="w-25%" />
      <a-statistic title="今日点赞数" :value="todayTrafficStatsVO.like_count" class="w-25%" />
    </a-flex>
  </a-card>
  <a-card title="总数据" class="mt-5">
    <a-flex gap="middle" horizontal>
      <a-statistic title="总访问量" :value="trafficStatsVO.view_count" class="w-33%" />
      <a-statistic title="总评论数" :value="trafficStatsVO.comment_count" class="w-33%" />
      <a-statistic title="总点赞数" :value="trafficStatsVO.like_count" class="w-33%" />
    </a-flex>
  </a-card>
</template>

<script lang="ts" setup>
import { message } from 'ant-design-vue'
import {
  GetTodayTrafficStats, GetTrafficStats,
  type TodayTrafficStatsVO,
  type TrafficStatsVO
} from '@/interfaces/DataAnalysis'
import { ref } from 'vue'


const todayTrafficStatsVO = ref<TodayTrafficStatsVO>({
  view_count: 0,
  user_view_count: 0,
  comment_count: 0,
  like_count: 0
})

const trafficStatsVO = ref<TrafficStatsVO>({
  view_count: 0,
  comment_count: 0,
  like_count: 0
})
const getTodayTrafficStats = async () => {
  try {
    const response: any = await GetTodayTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    todayTrafficStatsVO.value = response.data.data || todayTrafficStatsVO.value
    console.log(trafficStatsVO)
  } catch (error) {
    console.log(error)
  }
}
getTodayTrafficStats()

const getTrafficStats = async () => {
  try {
    const response: any = await GetTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    trafficStatsVO.value = response.data.data || trafficStatsVO.value
  } catch (error) {
    console.log(error)
  }
}
getTrafficStats()
</script>
