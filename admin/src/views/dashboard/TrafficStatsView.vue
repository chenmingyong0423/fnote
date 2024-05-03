<template>
  <a-card title="今日数据">
    <a-flex gap="middle" horizontal>
      <a-statistic title="今日访问量（PV）" :value="todayTrafficStatsVO.view_count" class="w-25%" />
      <a-statistic
        title="今日访问用户"
        :value="todayTrafficStatsVO.user_view_count"
        class="w-25%"
      />
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
  <a-card title="趋势图" class="mt-5">
    <div>
      <div id="user-analysis-4-week" class="w-full h-120" />
      <div id="user-analysis-4-month" class="w-full h-120" />
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { message } from 'ant-design-vue'
import {
  GetTendencyStats,
  GetTodayTrafficStats,
  GetTrafficStats,
  type TendencyData,
  type TendencyDataVO,
  type TodayTrafficStatsVO,
  type TrafficStatsVO
} from '@/interfaces/DataAnalysis'
import { onMounted, reactive, ref, watch } from 'vue'

document.title = '流量统计 - 后台管理'

const todayTrafficStatsVO = ref<TodayTrafficStatsVO>({
  view_count: 0,
  user_view_count: 0,
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
  } catch (error) {
    console.log(error)
  }
}
getTodayTrafficStats()

const trafficStatsVO = ref<TrafficStatsVO>({
  view_count: 0,
  comment_count: 0,
  like_count: 0
})

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

import { echarts } from '@/utils/echarts-setup'
import type { IResponse } from '@/interfaces/Common' // 确保路径正确

onMounted(() => {
  initUserAnalysis4Week()
  initUserAnalysis4Month()
})

const tendencyData4Week = reactive<{ pv: number[][]; uv: number[][] }>({
  pv: [],
  uv: []
})
// 计算属性
watch(
  () => tendencyData4Week,
  () => {
    initUserAnalysis4Week()
  },
  { deep: true }
)

const getTendencyStats4Week = async () => {
  try {
    const response: any = await GetTendencyStats('week')
    const apiResponse: IResponse<TendencyDataVO> = response.data
    if (apiResponse.code !== 0) {
      message.error(apiResponse.message)
      return
    }
    tendencyData4Week.pv = []
    tendencyData4Week.uv = []
    apiResponse.data?.pv.forEach((item: TendencyData) => {
      tendencyData4Week.pv.push([item.timestamp * 1000, item.view_count])
    })
    apiResponse.data?.uv.forEach((item: TendencyData) => {
      tendencyData4Week.uv.push([item.timestamp * 1000, item.view_count])
    })
  } catch (error) {
    console.log(error)
  }
}

getTendencyStats4Week()

const tendencyData4Month = reactive<{ pv: number[][]; uv: number[][] }>({
  pv: [],
  uv: []
})

// 计算属性
watch(
  () => tendencyData4Month,
  () => {
    initUserAnalysis4Month()
  },
  { deep: true }
)

const getTendencyStats4Month = async () => {
  try {
    const response: any = await GetTendencyStats('month')
    const apiResponse: IResponse<TendencyDataVO> = response.data
    if (apiResponse.code !== 0) {
      message.error(apiResponse.message)
      return
    }
    tendencyData4Month.pv = []
    tendencyData4Month.uv = []
    apiResponse.data?.pv.forEach((item: TendencyData) => {
      tendencyData4Month.pv.push([item.timestamp * 1000, item.view_count])
    })
    apiResponse.data?.uv.forEach((item: TendencyData) => {
      tendencyData4Month.uv.push([item.timestamp * 1000, item.view_count])
    })
  } catch (error) {
    console.log(error)
  }
}

getTendencyStats4Month()

const initUserAnalysis4Week = () => {
  const myChart = echarts.init(document.getElementById('user-analysis-4-week'))
  myChart.setOption({
    tooltip: {
      trigger: 'axis',
      position: function (pt: any) {
        return [pt[0], '10%']
      }
    },
    legend: {
      data: ['浏览量', '用户访问量']
    },
    title: {
      text: '最近 7 天',
      left: '5%'
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        formatter: function (value: number) {
          const date = new Date(value)
          const year = date.getFullYear()
          const month = ('0' + (date.getMonth() + 1)).slice(-2)
          const day = ('0' + date.getDate()).slice(-2)
          return year + '-' + month + '-' + day
        }
      }
    },
    yAxis: {
      type: 'value',
      boundaryGap: [0, '100%']
    },
    series: [
      {
        name: '浏览量',
        type: 'line',
        data: tendencyData4Week.pv
      },
      {
        name: '用户访问量',
        type: 'line',
        data: tendencyData4Week.uv
      }
    ]
  })
}

const initUserAnalysis4Month = () => {
  const myChart = echarts.init(document.getElementById('user-analysis-4-month'))
  myChart.setOption({
    tooltip: {
      trigger: 'axis',
      position: function (pt: any) {
        return [pt[0], '10%']
      }
    },
    legend: {
      data: ['浏览量', '用户访问量']
    },
    title: {
      text: '最近 30 天',
      left: '5%'
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        formatter: function (value: number) {
          const date = new Date(value)
          const year = date.getFullYear()
          const month = ('0' + (date.getMonth() + 1)).slice(-2)
          const day = ('0' + date.getDate()).slice(-2)
          return year + '-' + month + '-' + day
        }
      }
    },
    yAxis: {
      type: 'value',
      boundaryGap: [0, '100%']
    },
    series: [
      {
        name: '浏览量',
        type: 'line',
        data: tendencyData4Month.pv
      },
      {
        name: '用户访问量',
        type: 'line',
        data: tendencyData4Month.uv
      }
    ]
  })
}
</script>
