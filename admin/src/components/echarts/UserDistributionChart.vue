<template>
  <a-card title="用户分布">
    <template #extra>
      <div class="dashboard-chart-extra">
        <a-range-picker
          v-model:value="datetime"
          show-time
          :allow-clear="false"
          @change="datetimeChanged"
        />
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
    <a-spin :spinning="userDistributionLoading">
      <div id="user-distribution" class="user-distribution-chart" />
    </a-spin>
  </a-card>
</template>

<script setup lang="ts">
import { h, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import type { EChartsType } from 'echarts/core'
import dayjs, { type Dayjs } from 'dayjs'
import { GetUserDistributionStats, type UserDistributionVO } from '@/interfaces/DataAnalysis'
import type { IListData, IResponse } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import { echarts } from '@/utils/echarts-setup'
import { toErrorMessage } from '@/utils/error'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()

const userDistributionData = reactive<{
  seriesData: { name: string; value: number }[]
  legendData: string[]
  totalUsers: number
}>({ seriesData: [], legendData: [], totalUsers: 0 })

let userDistributionChart: EChartsType | null = null

const setUserDistributionChart = () => {
  const textColor = themeStore.mode === 'dark' ? 'rgba(255, 255, 255, 0.65)' : '#666'
  userDistributionChart?.setOption({
    textStyle: { color: textColor },
    title: {
      text: '用户分布',
      subtext: `总用户：${userDistributionData.totalUsers}，地区数：${userDistributionData.legendData.length}`,
      left: 'center',
      textStyle: { color: textColor },
      subtextStyle: { color: textColor }
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
      data: userDistributionData.legendData,
      textStyle: { color: textColor }
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

const datetime = ref<RangeValue>([dayjs().startOf('day'), dayjs().endOf('day')])

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
      message.error(apiResponse.message || '用户分布加载失败')
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
    message.error(toErrorMessage(error, '用户分布加载失败'))
  } finally {
    userDistributionLoading.value = false
  }
}

const resizeChart = () => {
  userDistributionChart?.resize()
}

watch(
  () => userDistributionData,
  () => {
    setUserDistributionChart()
  },
  { deep: true }
)

watch(() => themeStore.mode, setUserDistributionChart)

onMounted(() => {
  const chartElement = document.getElementById('user-distribution')
  if (!chartElement) {
    return
  }
  userDistributionChart = echarts.init(chartElement)
  setUserDistributionChart()
  getUserDistribution()
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)
  userDistributionChart?.dispose()
  userDistributionChart = null
})
</script>

<style scoped>
.dashboard-chart-extra {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.user-distribution-chart {
  width: 100%;
  height: 420px;
}
</style>
