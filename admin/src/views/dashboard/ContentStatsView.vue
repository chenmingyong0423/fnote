<template>
  <a-card title="总览">
    <a-flex gap="middle" horizontal>
      <a-statistic title="文章数量" :value="analysis.post_count" class="w-33%" />
      <a-statistic title="分类数量" :value="analysis.category_count" class="w-33%" />
      <a-statistic title="标签数量" :value="analysis.tag_count" class="w-33%" />
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
