<template>
  <div class="dashboard-stat-grid" :style="{ '--stat-min-width': minWidth }">
    <div v-for="item in items" :key="item.key" class="dashboard-stat-item">
      <span class="dashboard-stat-label">{{ item.label }}</span>
      <a-statistic
        :value="item.value"
        :precision="item.precision"
        :suffix="item.suffix"
        class="dashboard-stat-value"
      />
      <span v-if="item.description" class="dashboard-stat-description">
        {{ item.description }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
export interface DashboardStatItem {
  key: string
  label: string
  value: number
  suffix?: string
  precision?: number
  description?: string
}

withDefaults(
  defineProps<{
    items: DashboardStatItem[]
    minWidth?: string
  }>(),
  {
    minWidth: '160px'
  }
)
</script>

<style scoped>
.dashboard-stat-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(var(--stat-min-width), 1fr));
  gap: 12px;
}

.dashboard-stat-item {
  min-width: 0;
  padding: 14px 16px;
  background: #fafafa;
  border: 1px solid #f0f0f0;
  border-radius: 6px;
}

.dashboard-stat-label,
.dashboard-stat-description {
  display: block;
  color: rgba(0, 0, 0, 0.45);
  line-height: 1.5;
}

.dashboard-stat-label {
  margin-bottom: 6px;
  font-size: 13px;
}

.dashboard-stat-description {
  margin-top: 4px;
  font-size: 12px;
}

.dashboard-stat-value :deep(.ant-statistic-content) {
  color: rgba(0, 0, 0, 0.88);
  font-size: 24px;
  line-height: 1.2;
}
</style>
