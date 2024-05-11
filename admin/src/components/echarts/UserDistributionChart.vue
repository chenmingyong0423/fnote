<template>
  <a-card title="用户分布图" class="mt-5">
    <template #extra>
      <div class="flex gap-x-3">
        <a-range-picker v-model:value="datetime" show-time @change="datetimeChanged" />
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="userDistributionLoading"
            @click="getUserDistribution"
          />
        </a-tooltip>
      </div>
    </template>
    <div>
      <a-spin :spinning="userDistributionLoading">
        <div id="user-distribution" class="w-full h-120" />
      </a-spin>
    </div>
  </a-card>
</template>
<script setup lang="ts">
import { h, onMounted, reactive, ref, watch } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import type { EChartsType } from 'echarts/core'
import dayjs, { type Dayjs } from 'dayjs'
import { GetUserDistributionStats, type UserDistributionVO } from '@/interfaces/DataAnalysis'
import type { IListData, IResponse } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import { echarts } from '@/utils/echarts-setup'

const userDistributionData = reactive<{
  seriesData: { name: string; value: number }[]
  legendData: string[]
  totalUsers: number
}>({ seriesData: [], legendData: [], totalUsers: 0 })

let userDistributionChart: EChartsType | null = null
const setUserDistributionChart = () => {
  userDistributionChart?.setOption({
    title: {
      text: '用户分布图',
      subtext: `总用户：${userDistributionData.totalUsers}, 地区个数：${userDistributionData.legendData.length}`,
      left: 'center'
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b} : {c} ({d}%)'
    },
    legend: {
      type: 'scroll',
      orient: 'vertical',
      right: 10,
      top: 20,
      bottom: 20,
      data: userDistributionData.legendData
    },
    series: [
      {
        name: '地区',
        type: 'pie',
        radius: '50%',
        data: userDistributionData.seriesData,
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  })
}

type RangeValue = [Dayjs, Dayjs]

const datetimeFormat = 'YYYY-MM-DD HH:mm:ss'

const datetime = ref<RangeValue>([
  dayjs().startOf('day'), // 当天的00:00:00
  dayjs().endOf('day') // 当天的23:59:59
])

const datetimeChanged = () => {
  getUserDistribution()
}

const userDistributionLoading = ref(false)

const getUserDistribution = async () => {
  try {
    userDistributionLoading.value = true
    const response: any = await GetUserDistributionStats(
      datetime.value[0].format(datetimeFormat),
      datetime.value[1].format(datetimeFormat)
    )
    const apiResponse: IResponse<IListData<UserDistributionVO>> = response.data
    if (apiResponse.code !== 0) {
      message.error(apiResponse.message)
      return
    }
    userDistributionData.seriesData = []
    userDistributionData.legendData = []
    userDistributionData.totalUsers = 0
    apiResponse.data?.list.forEach((item: UserDistributionVO) => {
      userDistributionData.seriesData.push({ name: item.location, value: item.user_count })
      userDistributionData.legendData.push(item.location)
      userDistributionData.totalUsers += item.user_count
    })
  } catch (error) {
    console.log(error)
  } finally {
    userDistributionLoading.value = false
  }
}

getUserDistribution()

// 计算属性
watch(
  () => userDistributionData,
  () => {
    setUserDistributionChart()
  },
  { deep: true }
)

onMounted(() => {
  userDistributionChart = echarts.init(document.getElementById('user-distribution'))
  setUserDistributionChart()
})
</script>
