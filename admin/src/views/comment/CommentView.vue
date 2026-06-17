<template>
  <a-card title="评论列表">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="refresh" />
      </a-tooltip>
    </template>

    <div class="comment-toolbar">
      <div class="comment-toolbar-actions">
        <a-button v-if="data.length > 0" @click="expandOrHideRows">
          {{ expandedRowKeys.length === 0 ? '全部展开' : '全部折叠' }}
        </a-button>
        <template v-if="selectedRowKeys.length > 0">
          <span class="selection-count">已选 {{ selectedRowKeys.length }} 条</span>
          <a-button type="primary" :loading="batchApproving" @click="batchApproveComment">
            通过所选
          </a-button>
          <a-button type="primary" danger :loading="batchDeleting" @click="batchDeleteComment">
            删除所选
          </a-button>
        </template>
      </div>

      <div class="comment-filter">
        <span>状态：</span>
        <a-select
          v-model:value="approveStatus"
          class="status-select"
          :options="statusList"
          @change="handleChange"
        />
      </div>
    </div>

    <a-spin :spinning="loading">
      <a-table
        :columns="columns"
        :data-source="data"
        :pagination="pagination"
        :row-key="getRowKey"
        :row-selection="selection"
        :scroll="{ x: 1280 }"
        children-column-name="replies"
        v-model:expandedRowKeys="expandedRowKeys"
        bordered
        @change="change"
        @expandedRowsChange="expandedRowsChange"
      >
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.key === 'user_info'">
            <div class="comment-user">
              <a-avatar :src="record.user_info.picture" />
              <div class="comment-user-meta">
                <a
                  v-if="record.user_info.website"
                  :href="record.user_info.website"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {{ record.user_info.name }}
                </a>
                <span v-else class="comment-user-name">{{ record.user_info.name }}</span>
                <span class="comment-user-email">{{ record.user_info.email }}</span>
                <span v-if="record.user_info.ip" class="comment-user-ip">{{
                  record.user_info.ip
                }}</span>
              </div>
            </div>
          </template>

          <template v-else-if="column.key === 'post'">
            <a :href="record.post_info.post_url" target="_blank" rel="noopener noreferrer">
              {{ record.post_info.post_title }}
            </a>
          </template>

          <template v-else-if="column.key === 'content'">
            <a-tooltip :title="record.content">
              <div class="comment-content">{{ record.content }}</div>
            </a-tooltip>
          </template>

          <template v-else-if="column.key === 'reply_count'">
            <span>{{ record.type === 'comment' ? record.reply_count : '-' }}</span>
          </template>

          <template v-else-if="column.key === 'approval_status'">
            <a-tag :color="record.approval_status ? 'success' : 'processing'">
              {{ record.approval_status ? '审核通过' : '待审核' }}
            </a-tag>
          </template>

          <template v-else-if="column.key === 'type'">
            <a-tag :color="record.type === 'comment' ? 'blue' : 'cyan'">
              {{ record.type === 'comment' ? '评论' : '回复' }}
            </a-tag>
          </template>

          <template v-else-if="column.key === 'created_at' || column.key === 'updated_at'">
            {{ formatTime(Number(text)) }}
          </template>

          <template v-else-if="column.key === 'operation'">
            <div class="comment-actions">
              <a-popconfirm
                v-if="!record.approval_status"
                title="确认通过？"
                @confirm="approveComment(record)"
              >
                <a-button type="link" size="small" :loading="isApproving(record)"> 通过 </a-button>
              </a-popconfirm>
              <a-popconfirm title="确认删除？" @confirm="deleteById(record)">
                <a-button type="link" size="small" danger :loading="isDeleting(record)">
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

<script setup lang="ts">
import dayjs from 'dayjs'
import { computed, h, ref } from 'vue'
import type { IPageData, IResponse, PageRequest } from '@/interfaces/Common'
import {
  ApproveCommentById,
  ApproveReplyById,
  type AdminCommentVO,
  DeleteCommentById,
  DeleteReplyById,
  GetComments,
  batchApproved,
  type BatchApprovedCommentRequest,
  batchDelete
} from '@/interfaces/Comment'
import { message, Table } from 'ant-design-vue'
import type { TableColumnType, TableProps } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { toErrorMessage } from '@/utils/error'

document.title = '评论列表 - 后台管理'

type Key = string | number

const showSorterTooltip = ref('点击升序排序')
const columns = computed<TableColumnType<AdminCommentVO>[]>(() => [
  {
    title: '用户',
    dataIndex: 'user_info',
    key: 'user_info',
    width: 240
  },
  {
    title: '文章',
    dataIndex: ['post_info', 'post_url'],
    key: 'post',
    width: 220,
    ellipsis: true
  },
  {
    title: '内容',
    dataIndex: 'content',
    key: 'content',
    width: 300
  },
  {
    title: '回复数',
    dataIndex: 'reply_count',
    key: 'reply_count',
    width: 90
  },
  {
    title: '状态',
    key: 'approval_status',
    dataIndex: 'approval_status',
    width: 110
  },
  {
    title: '类型',
    key: 'type',
    dataIndex: 'type',
    width: 90
  },
  {
    title: '提交时间',
    key: 'created_at',
    dataIndex: 'created_at',
    sorter: (c1: AdminCommentVO, c2: AdminCommentVO) => c1.created_at - c2.created_at,
    defaultSortOrder: 'descend',
    sortDirections: ['descend', 'ascend'],
    showSorterTooltip: { title: showSorterTooltip.value },
    width: 170
  },
  {
    title: '更新时间',
    key: 'updated_at',
    dataIndex: 'updated_at',
    width: 170
  },
  {
    title: '操作',
    key: 'operation',
    dataIndex: 'operation',
    fixed: 'right',
    width: 120
  }
])

const data = ref<AdminCommentVO[]>([])
const pageReq = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5
} as PageRequest)
const total = ref(0)
const loading = ref(false)
const batchApproving = ref(false)
const batchDeleting = ref(false)
const approvingKeys = ref<Set<string>>(new Set())
const deletingKeys = ref<Set<string>>(new Set())
const commentMap = new Map<string, AdminCommentVO>()
const selectedComments = ref<BatchApprovedCommentRequest>({
  comment_ids: [],
  replies: {}
})
const selectedRowKeys = ref<Key[]>([])
const expandedRowKeys = ref<Key[]>([])
const approveStatus = ref(-1)

const statusList = [
  { label: '全部', value: -1 },
  { label: '待审核', value: 0 },
  { label: '已审核', value: 1 }
]

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const generateDefaultSort = () => {
  if (!pageReq.value.sort || pageReq.value.sort.length === 0) {
    pageReq.value.sort = '-created_at'
  }
}

const getRowKey = (record: AdminCommentVO) => record.key || record.id

const formatTime = (timestamp?: number) => {
  if (!timestamp) {
    return '-'
  }
  return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

const toCommentKey = (record: AdminCommentVO) => record.key || record.id

const setPending = (source: typeof approvingKeys, key: string, pending: boolean) => {
  const next = new Set(source.value)
  if (pending) {
    next.add(key)
  } else {
    next.delete(key)
  }
  source.value = next
}

const isApproving = (record: AdminCommentVO) => approvingKeys.value.has(toCommentKey(record))

const isDeleting = (record: AdminCommentVO) => deletingKeys.value.has(toCommentKey(record))

const clearTableState = () => {
  selectedRowKeys.value = []
  selectedComments.value = {
    comment_ids: [],
    replies: {}
  }
  expandedRowKeys.value = []
}

const normalizeComments = (comments: AdminCommentVO[]) => {
  commentMap.clear()
  comments.forEach((commentVO) => {
    commentVO.key = commentVO.id
    commentVO.replies?.forEach((replyVO) => {
      replyVO.fid = commentVO.id
      replyVO.key = `${commentVO.id}~${replyVO.id}`
      commentMap.set(replyVO.key, replyVO)
    })
    commentMap.set(commentVO.key, commentVO)
  })
  return comments
}

const buildSelectedComments = (keys: Key[]) => {
  const next: BatchApprovedCommentRequest = {
    comment_ids: [],
    replies: {}
  }

  keys.forEach((key) => {
    const stringKey = String(key)
    const record = commentMap.get(stringKey)
    if (!record) {
      return
    }

    if (stringKey.includes('~')) {
      const [commentId, replyId] = stringKey.split('~')
      next.replies[commentId] = next.replies[commentId] || []
      next.replies[commentId].push(replyId)
    } else {
      next.comment_ids.push(stringKey)
    }
  })

  return next
}

const updateSelectedRows = (keys: Key[]) => {
  selectedRowKeys.value = keys
  selectedComments.value = buildSelectedComments(keys)
}

const get = async () => {
  try {
    generateDefaultSort()
    loading.value = true
    const response: any = await GetComments(pageReq.value)
    const result: IResponse<IPageData<AdminCommentVO>> = response.data
    if (result.code !== 0) {
      message.error(result.message || '评论列表加载失败')
      return
    }
    data.value = normalizeComments(result.data?.list || [])
    total.value = result.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '评论列表加载失败'))
  } finally {
    loading.value = false
  }
}

get()

const change: TableProps<AdminCommentVO>['onChange'] = (pagination, _filters, sorter: any) => {
  pageReq.value.pageNo = Number(pagination.current || 1)
  pageReq.value.pageSize = Number(pagination.pageSize || 5)
  clearTableState()
  switch (sorter.order) {
    case 'ascend':
      pageReq.value.sort = '+created_at'
      showSorterTooltip.value = '点击默认排序'
      break
    case 'descend':
      pageReq.value.sort = '-created_at'
      showSorterTooltip.value = '点击升序排序'
      break
    default:
      pageReq.value.sort = '-created_at'
      showSorterTooltip.value = '点击降序排序'
  }
  get()
}

const approveComment = async (record: AdminCommentVO) => {
  const key = toCommentKey(record)
  if (approvingKeys.value.has(key)) {
    return
  }

  try {
    setPending(approvingKeys, key, true)
    const response: any =
      record.type === 'comment'
        ? await ApproveCommentById(record.id)
        : await ApproveReplyById(record.fid || '', record.id)

    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('审核成功')
    clearTableState()
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '审核失败'))
  } finally {
    setPending(approvingKeys, key, false)
  }
}

const deleteById = async (record: AdminCommentVO) => {
  const key = toCommentKey(record)
  if (deletingKeys.value.has(key)) {
    return
  }

  try {
    setPending(deletingKeys, key, true)
    const response: any =
      record.type === 'comment'
        ? await DeleteCommentById(record.id)
        : await DeleteReplyById(record.fid || '', record.id)

    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    clearTableState()
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    setPending(deletingKeys, key, false)
  }
}

const onChange = (keys: Key[]) => {
  updateSelectedRows(keys)
}

const buildDeleteRequest = () => {
  const deleteRequest: BatchApprovedCommentRequest = JSON.parse(
    JSON.stringify(selectedComments.value)
  )
  deleteRequest.comment_ids.forEach((commentId) => {
    if (deleteRequest.replies[commentId]) {
      delete deleteRequest.replies[commentId]
    }
  })
  return deleteRequest
}

const batchApproveComment = async () => {
  if (selectedRowKeys.value.length === 0 || batchApproving.value) {
    return
  }

  try {
    batchApproving.value = true
    const apiResponse = await batchApproved(selectedComments.value)
    if (apiResponse.data?.code !== 0) {
      message.error(apiResponse.data?.message || '批量审核失败')
      return
    }
    message.success('批量审核成功')
    clearTableState()
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '批量审核失败'))
  } finally {
    batchApproving.value = false
  }
}

const batchDeleteComment = async () => {
  if (selectedRowKeys.value.length === 0 || batchDeleting.value) {
    return
  }

  try {
    batchDeleting.value = true
    const apiResponse = await batchDelete(buildDeleteRequest())
    if (apiResponse.data?.code !== 0) {
      message.error(apiResponse.data?.message || '批量删除失败')
      return
    }
    message.success('批量删除成功')
    clearTableState()
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '批量删除失败'))
  } finally {
    batchDeleting.value = false
  }
}

const handleChange = (value: number) => {
  pageReq.value.pageNo = 1
  pageReq.value.approvalStatus = value === -1 ? undefined : value !== 0
  clearTableState()
  get()
}

const selection = computed<TableProps<AdminCommentVO>['rowSelection']>(() => {
  return {
    selectedRowKeys: selectedRowKeys.value,
    onChange,
    hideDefaultSelections: true,
    selections: [
      Table.SELECTION_ALL,
      Table.SELECTION_NONE,
      {
        key: 'approved',
        text: '选中已审核',
        onSelect: (changeableRowKeys: Key[]) => {
          updateSelectedRows(
            changeableRowKeys.filter((key) => commentMap.get(String(key))?.approval_status)
          )
        }
      },
      {
        key: 'unapproved',
        text: '选中未审核',
        onSelect: (changeableRowKeys: Key[]) => {
          updateSelectedRows(
            changeableRowKeys.filter((key) => !commentMap.get(String(key))?.approval_status)
          )
        }
      }
    ]
  }
})

const expandedRowsChange = (rowKeys: Key[]) => {
  expandedRowKeys.value = rowKeys
}

const expandOrHideRows = () => {
  if (expandedRowKeys.value.length === 0) {
    expandedRowKeys.value = data.value.map((item) => item.key || item.id)
  } else {
    expandedRowKeys.value = []
  }
}

const refresh = async () => {
  clearTableState()
  await get()
}
</script>

<style scoped>
.comment-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.comment-toolbar-actions,
.comment-filter,
.comment-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-filter {
  margin-left: auto;
}

.status-select {
  width: 120px;
}

.selection-count,
.comment-user-email,
.comment-user-ip {
  color: rgba(0, 0, 0, 0.45);
}

.comment-user {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  min-width: 0;
}

.comment-user-meta {
  display: flex;
  flex-direction: column;
  min-width: 0;
  line-height: 1.5;
}

.comment-user-name {
  font-weight: 600;
}

.comment-user-email,
.comment-user-ip {
  overflow: hidden;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-content {
  display: -webkit-box;
  overflow: hidden;
  line-height: 1.5;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.comment-actions :deep(.ant-btn) {
  padding: 0;
}
</style>
