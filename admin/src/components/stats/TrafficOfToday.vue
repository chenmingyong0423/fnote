<template>
  <a-card title="今日数据">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="todayTrafficLoading"
          @click="getTodayTrafficStats"
        />
      </a-tooltip>
    </template>
    <a-spin :spinning="todayTrafficLoading">
      <a-flex gap="middle" horizontal>
        <a-statistic
          title="今日访问量（PV）"
          :value="todayTrafficStatsVO.view_count"
          class="w-25%"
        />
        <a-statistic
          title="今日访问用户"
          :value="todayTrafficStatsVO.user_view_count"
          class="w-25%"
        />
        <a-statistic title="今日评论数" :value="todayTrafficStatsVO.comment_count" class="w-25%" />
        <a-statistic title="今日点赞数" :value="todayTrafficStatsVO.like_count" class="w-25%" />
      </a-flex>
    </a-spin>
  </a-card>
</template>
<script setup lang="ts">
import { h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { GetTodayTrafficStats, type TodayTrafficStatsVO } from '@/interfaces/DataAnalysis'
import { message } from 'ant-design-vue'

const todayTrafficLoading = ref<boolean>(false)
const todayTrafficStatsVO = ref<TodayTrafficStatsVO>({
  view_count: 0,
  user_view_count: 0,
  comment_count: 0,
  like_count: 0
})

const getTodayTrafficStats = async () => {
  try {
    todayTrafficLoading.value = true
    const response: any = await GetTodayTrafficStats()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    todayTrafficStatsVO.value = response.data.data || todayTrafficStatsVO.value
  } catch (error) {
    console.log(error)
  } finally {
    todayTrafficLoading.value = false
  }
}
getTodayTrafficStats()
</script>
