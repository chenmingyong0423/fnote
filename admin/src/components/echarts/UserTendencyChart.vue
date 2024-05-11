<template>
  <a-card title="趋势图" class="mt-5">
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
      <TimeLineChart id="tendencySeries4Week" :data="tendencySeries4Week" class="w-full h-120" />
      <TimeLineChart id="tendencySeries4Month" :data="tendencySeries4Month" class="w-full h-120" />
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

const tendencySeries4Week = reactive<{
  title: string
  series: {
    name: String
    type: String
    data: number[][]
  }[]
}>({
  title: '最近 7 天',
  series: [
    {
      name: '浏览量',
      type: 'line',
      data: []
    },
    {
      name: '用户访问量',
      type: 'line',
      data: []
    }
  ]
})

const tendencyLoading = ref<boolean>(false)

const getTendencyStats4Week = async () => {
  try {
    const response: any = await GetTendencyStats('week')
    const apiResponse: IResponse<TendencyDataVO> = response.data
    if (apiResponse.code !== 0) {
      message.error(apiResponse.message)
      return
    }
    tendencySeries4Week.series[0].data = []
    tendencySeries4Week.series[1].data = []
    apiResponse.data?.pv.forEach((item: TendencyData) => {
      tendencySeries4Week.series[0].data.push([item.timestamp * 1000, item.view_count])
    })
    apiResponse.data?.uv.forEach((item: TendencyData) => {
      tendencySeries4Week.series[1].data.push([item.timestamp * 1000, item.view_count])
    })
  } catch (error) {
    console.log(error)
  }
}

getTendencyStats4Week()

const tendencySeries4Month = reactive<{
  title: string
  series: {
    name: String
    type: String
    data: number[][]
  }[]
}>({
  title: '最近 30 天',
  series: [
    {
      name: '浏览量',
      type: 'line',
      data: []
    },
    {
      name: '用户访问量',
      type: 'line',
      data: []
    }
  ]
})

const getTendencyStats4Month = async () => {
  try {
    const response: any = await GetTendencyStats('month')
    const apiResponse: IResponse<TendencyDataVO> = response.data
    if (apiResponse.code !== 0) {
      message.error(apiResponse.message)
      return
    }
    tendencySeries4Month.series[0].data = []
    tendencySeries4Month.series[1].data = []
    apiResponse.data?.pv.forEach((item: TendencyData) => {
      tendencySeries4Month.series[0].data.push([item.timestamp * 1000, item.view_count])
    })
    apiResponse.data?.uv.forEach((item: TendencyData) => {
      tendencySeries4Month.series[1].data.push([item.timestamp * 1000, item.view_count])
    })
  } catch (error) {
    console.log(error)
  }
}

getTendencyStats4Month()

const refresh = async () => {
  tendencyLoading.value = true
  await getTendencyStats4Week()
  await getTendencyStats4Month()
  tendencyLoading.value = false
}
</script>
