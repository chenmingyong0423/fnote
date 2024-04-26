<template>
  <a-card title="总览">
    <a-flex gap="middle" horizontal>
      <a-statistic title="文章数量" :value="contentStatsVO.post_count" class="w-33%" />
      <a-statistic title="分类数量" :value="contentStatsVO.category_count" class="w-33%" />
      <a-statistic title="标签数量" :value="contentStatsVO.tag_count" class="w-33%" />
    </a-flex>
  </a-card>
</template>

<script lang="ts" setup>
import { type ContentStatsVO, GetContentStatsVO } from '@/interfaces/DataAnalysis'
import { message } from 'ant-design-vue'
import { ref } from 'vue'

document.title = '内容发布统计 - 后台管理'

const contentStatsVO = ref<ContentStatsVO>({
  post_count: 0,
  category_count: 0,
  tag_count: 0
})
const getContentStatsVO = async () => {
  try {
    const response: any = await GetContentStatsVO()
    if (response.data.code !== 0) {
      message.error(response.data.data.message)
      return
    }
    contentStatsVO.value = response.data.data || contentStatsVO.value
  } catch (error) {
    console.log(error)
  }
}
getContentStatsVO()

</script>
