<template>
  <div class="flex justify-between mb-5">
    <a-card class="w-24%">
      <a-statistic title="文章数量" :value="analysis.post_count" style="margin-right: 50px" />
    </a-card>
    <a-card class="w-24%">
      <a-statistic title="分类数量" :value="analysis.category_count" style="margin-right: 50px" />
    </a-card>
    <a-card class="w-24%">
      <a-statistic title="标签总数" :value="analysis.tag_count" style="margin-right: 50px" />
    </a-card >
    <a-card class="w-24%">
      <a-statistic title="点赞数" :value="analysis.like_count" style="margin-right: 50px" />
    </a-card>
  </div>
  <div class="flex justify-between">
    <a-card class="">
      <a-statistic title="今日访问量" :value="analysis.today_view_count" style="margin-right: 50px" />
    </a-card>
    <a-card class="%">
      <a-statistic title="实际访问用户" :value="analysis.today_user_visit_count" style="margin-right: 50px" />
    </a-card>
    <a-card class="w-24%">
      <a-statistic title="总访问量" :value="analysis.total_view_count" style="margin-right: 50px" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>

import axios from '@/http/axios'
import { message } from 'ant-design-vue'
import type { DataAnalysisVO } from '@/interfaces/DataAnalysis'
import type { IResponse } from '@/interfaces/Common'
import { ref } from 'vue'

const analysis = ref<DataAnalysisVO>({
  totalVisits: 0,
  todayVisits: 0,
  totalUsers: 0,
  todayUsers: 0,
  totalPosts: 0,
  todayPosts: 0,
  totalComments: 0,
  todayComments: 0
})
const getAnalysis = async () => {
  try {
    const response = await axios.get<IResponse<DataAnalysisVO>>(`/admin/data-analysis`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    analysis.value = response.data.data
  } catch (error) {
    console.log(error)
    message.error('获取数据失败')
  }
};
getAnalysis()
</script>