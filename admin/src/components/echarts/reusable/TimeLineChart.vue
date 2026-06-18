<template>
  <div :id="props.id"></div>
</template>

<script lang="ts" setup>
import type { EChartsType } from 'echarts/core'
import { echarts } from '@/utils/echarts-setup.js'
import { onBeforeUnmount, onMounted, type PropType, watch } from 'vue'

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
  chart?.setOption({
    tooltip: {
      trigger: 'axis',
      position: function (pt: any) {
        return [pt[0], '10%']
      }
    },
    legend: {
      data: legend()
    },
    title: {
      text: props.data.title,
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
    series: props.data.series
  })
}

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
