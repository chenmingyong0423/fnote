<template>
  <div class="dashboard-page">
    <a-card title="内容发布总览">
      <template #extra>
        <div class="dashboard-card-extra">
          <span v-if="lastUpdatedText" class="dashboard-updated-at">{{ lastUpdatedText }}</span>
          <a-tooltip title="刷新数据">
            <a-button
              shape="circle"
              :icon="h(ReloadOutlined)"
              :loading="loading"
              @click="getContentStatsVO"
            />
          </a-tooltip>
        </div>
      </template>
      <a-spin :spinning="loading">
        <DashboardStatGrid :items="contentStatItems" />
      </a-spin>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { type ContentStatsVO, GetContentStatsVO } from '@/interfaces/DataAnalysis'
import { GetPostDraft } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import { computed, h, ref } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import DashboardStatGrid from '@/components/stats/DashboardStatGrid.vue'
import dayjs from 'dayjs'
import { toErrorMessage } from '@/utils/error'

document.title = '内容发布统计 - 后台管理'

const contentStatsVO = ref<ContentStatsVO>({
  post_count: 0,
  category_count: 0,
  tag_count: 0
})
const draftCount = ref(0)

const loading = ref(false)
const lastUpdatedAt = ref<number>()

const contentStatItems = computed(() => [
  {
    key: 'post_count',
    label: '已发布文章',
    value: contentStatsVO.value.post_count
  },
  {
    key: 'draft_count',
    label: '草稿数量',
    value: draftCount.value
  },
  {
    key: 'category_count',
    label: '分类数量',
    value: contentStatsVO.value.category_count
  },
  {
    key: 'tag_count',
    label: '标签数量',
    value: contentStatsVO.value.tag_count
  }
])

const lastUpdatedText = computed(() => {
  if (!lastUpdatedAt.value) {
    return ''
  }
  return `更新于 ${dayjs(lastUpdatedAt.value).format('HH:mm:ss')}`
})

const getContentStatsVO = async () => {
  try {
    loading.value = true
    const [contentStatsResponse, draftResponse]: any[] = await Promise.all([
      GetContentStatsVO(),
      GetPostDraft({
        pageNo: 1,
        pageSize: 1,
        sortField: 'updated_at',
        sortOrder: 'DESC'
      })
    ])

    if (contentStatsResponse.data.code !== 0) {
      message.error(contentStatsResponse.data.message || '内容统计加载失败')
      return
    }
    if (draftResponse.data.code !== 0) {
      message.error(draftResponse.data.message || '草稿统计加载失败')
      return
    }

    contentStatsVO.value = contentStatsResponse.data.data || contentStatsVO.value
    draftCount.value = draftResponse.data.data?.totalCount || 0
    lastUpdatedAt.value = Date.now()
  } catch (error) {
    message.error(toErrorMessage(error, '内容发布统计加载失败'))
  } finally {
    loading.value = false
  }
}

getContentStatsVO()
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dashboard-card-extra {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.dashboard-updated-at {
  color: var(--app-text-secondary);
  font-size: 12px;
  white-space: nowrap;
}
</style>
