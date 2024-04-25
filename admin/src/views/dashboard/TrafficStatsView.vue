<template>
  <a-card title="今日数据">
    <a-flex gap="middle" horizontal>
      <a-statistic title="今日访问量（PV）" :value="0" class="w-25%" />
      <a-statistic title="今日访问用户" :value="0" class="w-25%" />
      <a-statistic title="今日评论数" :value="0" class="w-25%" />
      <a-statistic title="今日点赞数" :value="0" class="w-25%" />
    </a-flex>
  </a-card>
  <a-card title="总数据" class="mt-5">
    <a-flex gap="middle" horizontal>
      <a-statistic title="总访问量" :value="analysis.total_view_count" class="w-33%" />
      <a-statistic title="总评论数" :value="analysis.comment_count" class="w-33%" />
      <a-statistic title="总点赞数" :value="analysis.like_count" class="w-33%" />
    </a-flex>
  </a-card>
</template>

<script lang="ts" setup>
import { message } from 'ant-design-vue'
import { type DataAnalysisVO, GetDataAnalysis } from '@/interfaces/DataAnalysis'
import { ref } from 'vue'

const analysis = ref<DataAnalysisVO>({
  post_count: 0,
  category_count: 0,
  tag_count: 0,
  comment_count: 0,
  like_count: 0,
  today_view_count: 0,
  total_view_count: 0,
  today_user_visit_count: 0
})
const getAnalysis = async () => {
  try {
    const response: any = await GetDataAnalysis()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    analysis.value = response.data.data || analysis.value
  } catch (error) {
    console.log(error)
  }
}
getAnalysis()
</script>
