<template>
  <div :id="props.id"></div>
</template>

<script lang="ts" setup>
import type { EChartsType } from 'echarts/core'
import { echarts } from '@/utils/echarts-setup.js'
import { onBeforeUnmount, onMounted, type PropType, watch } from 'vue'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const chartTextColor = () => (themeStore.mode === 'dark' ? 'rgba(255, 255, 255, 0.65)' : '#666')
const chartLineColor = () => (themeStore.mode === 'dark' ? '#303030' : '#e8e8e8')

const props = defineProps({
  id: {
    type: String,
    required: true
  },
  data: {
    type: Object as PropType<{
      title: string
      series: {
        name: string
        type: string
        data: number[][]
      }[]
    }>,
    required: true
  }
})

watch(
  () => props.data,
  () => {
    setOption()
  },
  { deep: true }
)

const legend = () => {
  return props.data.series.map((item) => item.name)
}

let chart: EChartsType | null = null
const resizeChart = () => {
  chart?.resize()
}

const setOption = () => {
  const textColor = chartTextColor()
  const lineColor = chartLineColor()
  chart?.setOption({
    textStyle: { color: textColor },
    tooltip: {
      trigger: 'axis',
      position: function (pt: any) {
        return [pt[0], '10%']
      }
    },
    legend: {
      data: legend(),
      textStyle: { color: textColor }
    },
    title: {
      text: props.data.title,
      left: '5%',
      textStyle: { color: textColor }
    },
    xAxis: {
      type: 'time',
      boundaryGap: false,
      axisLabel: {
        color: textColor,
        formatter: function (value: number) {
          const date = new Date(value)
          const year = date.getFullYear()
          const month = ('0' + (date.getMonth() + 1)).slice(-2)
          const day = ('0' + date.getDate()).slice(-2)
          return year + '-' + month + '-' + day
        }
      },
      axisLine: { lineStyle: { color: lineColor } }
    },
    yAxis: {
      type: 'value',
      boundaryGap: [0, '100%'],
      axisLabel: { color: textColor },
      splitLine: { lineStyle: { color: lineColor } }
    },
    series: props.data.series
  })
}

watch(() => themeStore.mode, setOption)

onMounted(() => {
  const chartElement = document.getElementById(props.id)
  if (!chartElement) {
    return
  }
  chart = echarts.init(chartElement)
  setOption()
  window.addEventListener('resize', resizeChart)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', resizeChart)
  chart?.dispose()
  chart = null
})
</script>
