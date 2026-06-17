<template>
  <a-card title="草稿箱">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getPostDrafts"
          />
        </a-tooltip>
      </div>
    </template>
    <div class="draft-toolbar">
      <a-input-search
        v-model:value="req.keyword"
        placeholder="搜索草稿标题"
        allow-clear
        class="draft-search"
        @search="searchDrafts"
        @pressEnter="searchDrafts"
      />
    </div>
    <a-spin :spinning="loading">
      <a-table
        :columns="columns"
        :data-source="listData"
        :pagination="pagination"
        row-key="id"
        bordered
        @change="change"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'title'">
            <span class="draft-title">{{ record.title || '未命名草稿' }}</span>
          </template>
          <template v-else-if="column.key === 'created_at'">
            <span>{{ formatTime(record.created_at) }}</span>
          </template>
          <template v-else-if="column.key === 'operation'">
            <div class="draft-actions">
              <a @click="router.push(`/home/post/draft/${record.id}`)">编辑</a>
              <a-popconfirm title="确认删除？" @confirm="deletePostDraft(record.id)">
                <a-button type="link" size="small" danger :loading="isDeleting(record.id)">
                  删除
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-spin>
  </a-card>
</template>
<script lang="ts" setup>
import { computed, h, ref } from 'vue'
import { DeletePostDraftById, GetPostDraft, type PageRequest } from '@/interfaces/Post'
import router from '@/router'
import dayjs from 'dayjs'
import { message } from 'ant-design-vue'
import type { TableColumnType, TableProps } from 'ant-design-vue'
import type { PostDraftBrief } from '@/interfaces/PostDraft'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { toErrorMessage } from '@/utils/error'

document.title = '草稿箱 - 后台管理'

const req = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'created_at',
  sortOrder: 'DESC',
  keyword: ''
} as PageRequest)

const listData = ref<PostDraftBrief[]>([])

const columns = computed<TableColumnType<PostDraftBrief>[]>(() => [
  {
    title: '标题',
    dataIndex: 'title',
    key: 'title',
    ellipsis: true
  },
  {
    title: '草稿 ID',
    dataIndex: 'id',
    key: 'id',
    width: 260,
    ellipsis: true
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    key: 'created_at',
    sorter: true,
    defaultSortOrder: 'descend',
    width: 180
  },
  {
    title: '操作',
    dataIndex: 'operation',
    key: 'operation',
    width: 120,
    fixed: 'right'
  }
])

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: req.value.pageNo,
  pageSize: req.value.pageSize
}))

const loading = ref(false)
const deletingIds = ref<Set<string>>(new Set())

const formatTime = (timestamp: number) => dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')

const setDeleting = (id: string, deleting: boolean) => {
  const next = new Set(deletingIds.value)
  if (deleting) {
    next.add(id)
  } else {
    next.delete(id)
  }
  deletingIds.value = next
}

const isDeleting = (id: string) => deletingIds.value.has(id)

const getPostDrafts = async () => {
  try {
    loading.value = true
    const response = await GetPostDraft(req.value)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    listData.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '草稿列表加载失败'))
  } finally {
    loading.value = false
  }
}

const deletePostDraft = async (id: string) => {
  if (isDeleting(id)) {
    return
  }
  try {
    setDeleting(id, true)
    const response: any = await DeletePostDraftById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getPostDrafts()
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    setDeleting(id, false)
  }
}

const searchDrafts = () => {
  req.value.pageNo = 1
  getPostDrafts()
}

const change: TableProps<PostDraftBrief>['onChange'] = (pagination, _filters, sorter: any) => {
  req.value.pageNo = Number(pagination.current || 1)
  req.value.pageSize = Number(pagination.pageSize || req.value.pageSize)
  req.value.sortField = sorter.field || 'created_at'
  req.value.sortOrder = sorter.order === 'ascend' ? 'ASC' : 'DESC'
  getPostDrafts()
}

getPostDrafts()
</script>

<style scoped>
.draft-toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 12px;
}

.draft-search {
  width: 240px;
}

.draft-title {
  font-weight: 500;
}

.draft-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  white-space: nowrap;
}

.draft-actions :deep(.ant-btn) {
  padding: 0;
}
</style>
