<template>
  <TrafficOfToday />
  <Traffic />
  <UserTendencyChart />
  <UserDistributionChart />
</template>

<script lang="ts" setup>
import UserDistributionChart from '@/components/echarts/UserDistributionChart.vue'
import UserTendencyChart from '@/components/echarts/UserTendencyChart.vue'
import TrafficOfToday from '@/components/stats/TrafficOfToday.vue'
import Traffic from '@/components/stats/Traffic.vue'
document.title = '流量统计 - 后台管理'
</script>
