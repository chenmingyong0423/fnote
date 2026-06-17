<template>
  <a-card title="访问趋势">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="tendencyLoading"
          @click="refresh"
        />
      </a-tooltip>
    </template>
    <a-spin :spinning="tendencyLoading">
      <div class="tendency-chart-grid">
        <TimeLineChart id="tendencySeries4Week" :data="tendencySeries4Week" class="chart-panel" />
        <TimeLineChart id="tendencySeries4Month" :data="tendencySeries4Month" class="chart-panel" />
      </div>
    </a-spin>
  </a-card>
</template>

<script setup lang="ts">
import { h, reactive, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { GetTendencyStats, type TendencyData, type TendencyDataVO } from '@/interfaces/DataAnalysis'
import type { IResponse } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import TimeLineChart from '@/components/echarts/reusable/TimeLineChart.vue'
import { toErrorMessage } from '@/utils/error'

type TendencySeries = {
  title: string
  series: {
    name: string
    type: string
    data: number[][]
  }[]
}

const createTendencySeries = (title: string): TendencySeries => ({
  title,
  series: [
    {
      name: '浏览量',
      type: 'line',
      data: []
    },
    {
      name: '访问用户',
      type: 'line',
      data: []
    }
  ]
})

const tendencySeries4Week = reactive<TendencySeries>(createTendencySeries('最近 7 天'))
const tendencySeries4Month = reactive<TendencySeries>(createTendencySeries('最近 30 天'))
const tendencyLoading = ref(false)

const toSeriesData = (items?: TendencyData[]) => {
  return items?.map((item) => [item.timestamp * 1000, item.view_count]) || []
}

const loadTendencyStats = async (
  period: 'week' | 'month',
  target: TendencySeries,
  fallback: string
) => {
  const response: any = await GetTendencyStats(period)
  const apiResponse: IResponse<TendencyDataVO> = response.data
  if (apiResponse.code !== 0) {
    throw new Error(apiResponse.message || fallback)
  }

  target.series[0].data = toSeriesData(apiResponse.data?.pv)
  target.series[1].data = toSeriesData(apiResponse.data?.uv)
}

const refresh = async () => {
  try {
    tendencyLoading.value = true
    await Promise.all([
      loadTendencyStats('week', tendencySeries4Week, '最近 7 天趋势加载失败'),
      loadTendencyStats('month', tendencySeries4Month, '最近 30 天趋势加载失败')
    ])
  } catch (error) {
    message.error(toErrorMessage(error, '趋势数据加载失败'))
  } finally {
    tendencyLoading.value = false
  }
}

refresh()
</script>

<style scoped>
.tendency-chart-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
}

.chart-panel {
  width: 100%;
  height: 420px;
}
</style>
