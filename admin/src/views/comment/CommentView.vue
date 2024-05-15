<template>
  <a-card title="评论列表">
    <div class="flex mb-3 gap-x-2">
      <div class="flex gap-x-2" v-if="data.length > 0">
        <a-button @click="expandOrHideRows">{{
          expandedRowKeys.length === 0 ? '全部展开' : '全部折叠'
        }}</a-button>
      </div>
      <div class="flex gap-x-2" v-if="selectedRowKeys.length > 0">
        <a-button type="primary" @click="batchApproveComment">通过所选</a-button>
        <a-button type="primary" danger @click="batchDeleteComment">删除所选</a-button>
      </div>
      <div class="ml-auto">
        状态：
        <a-select
          ref="select"
          class="w-120px"
          v-model:value="approveStatus"
          :options="statusList"
          @change="handleChange"
        ></a-select>
      </div>
      <div>
        <a-tooltip title="刷新数据">
          <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="refresh" />
        </a-tooltip>
      </div>
    </div>
    <a-table
      :columns="columns"
      :data-source="data"
      :pagination="pagination"
      @change="change"
      :row-selection="selection"
      childrenColumnName="replies"
      v-model:expandedRowKeys="expandedRowKeys"
      @expandedRowsChange="expandedRowsChange"
    >
      <template #bodyCell="{ column, text, record }">
        <template v-if="column.dataIndex === 'user_info'">
          <div class="flex gap-x-3">
            <div>
              <a-avatar :src="record.user_info.picture" />
            </div>
            <div class="flex-col">
              <div>
                <a
                  :href="record.user_info.website"
                  target="_blank"
                  v-if="record.user_info.website"
                  >{{ record.user_info.name }}</a
                >
                <span class="font-bold" v-else>{{ record.user_info.name }}</span>
              </div>
              <div class="text-gray-5">{{ record.user_info.email }}</div>
            </div>
          </div>
        </template>
        <template v-if="column.dataIndex === 'post.post_url'">
          <a :href="record.post_info.post_url" target="_blank">{{ record.post_info.post_title }}</a>
        </template>
        <template v-if="column.dataIndex === 'content'">
          {{ text }}
        </template>
        <template v-if="column.dataIndex === 'approval_status'">
          <a-tag :color="record.approval_status ? 'success' : 'processing'"
            >{{ text ? '审核通过' : '未审核' }}
          </a-tag>
        </template>
        <template v-if="column.dataIndex === 'type'">
          <a-tag color="success">{{ record.type === 'comment' ? '评论' : '回复' }}</a-tag>
        </template>
        <template v-if="['created_at', 'updated_at'].includes(column.dataIndex)">
          {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
        <template v-else-if="column.dataIndex === 'operation'">
          <div class="editable-row-operations">
            <a-popconfirm
              v-if="data.length && !record.approval_status"
              title="确认通过？"
              @confirm="approveComment(record)"
            >
              <a>通过</a>
            </a-popconfirm>
            <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deleteById(record)">
              <a>删除</a>
            </a-popconfirm>
          </div>
        </template>
      </template>
    </a-table>
  </a-card>
</template>
<script setup lang="ts">
import dayjs from 'dayjs'
import { computed, ref } from 'vue'
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
import { message } from 'ant-design-vue'
import { Table } from 'ant-design-vue'
import { h } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import originalAxios from 'axios'

document.title = '评论列表 - 后台管理'
const showSorterTooltip = ref('点击升序排序')

const columns = [
  {
    title: '用户',
    dataIndex: 'user_info',
    key: 'user_info'
  },
  {
    title: '文章',
    dataIndex: 'post.post_url',
    key: 'post.post_url'
  },
  {
    title: '内容',
    dataIndex: 'content',
    key: 'content'
  },
  {
    title: '回复数',
    dataIndex: 'reply_count',
    key: 'reply_count'
  },
  {
    title: '状态',
    key: 'approval_status',
    dataIndex: 'approval_status'
  },
  {
    title: '类型',
    key: 'type',
    dataIndex: 'type'
  },
  {
    title: '提交时间',
    key: 'created_at',
    dataIndex: 'created_at',
    sorter: (c1: AdminCommentVO, c2: AdminCommentVO) => c1.created_at - c2.created_at,
    defaultSortOrder: 'descend',
    sortDirections: ['descend', 'ascend'],
    showSorterTooltip: { title: showSorterTooltip.value }
  },
  {
    title: '更新时间',
    key: 'updated_at',
    dataIndex: 'updated_at'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const data = ref<AdminCommentVO[]>([])
const pageReq = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5
} as PageRequest)

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const loading = ref(false)

const generateDefaultSort = () => {
  if (!pageReq.value.sort || pageReq.value.sort.length === 0) {
    pageReq.value.sort = '-created_at'
  }
}

const commentMap: Map<string, AdminCommentVO> = new Map<string, AdminCommentVO>()

const get = async () => {
  try {
    generateDefaultSort()
    loading.value = true
    const response: any = await GetComments(pageReq.value)
    const result: IResponse<IPageData<AdminCommentVO>> = response.data
    if (result.code === 0) {
      commentMap.clear()
      data.value = response.data.data?.list || []
      data.value.forEach((commentVO: AdminCommentVO) => {
        commentVO.key = commentVO.id
        commentVO.replies?.forEach((replyVO: AdminCommentVO) => {
          replyVO.fid = commentVO.id
          replyVO.key = commentVO.id + '~' + replyVO.id
          commentMap.set(replyVO.key, replyVO)
        })
        commentMap.set(commentVO.key, commentVO)
      })
      total.value = response.data.data?.totalCount || 0
    }
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
  }
}
get()

const change = (pagination: any, _filters: any, sorter: any) => {
  pageReq.value.pageNo = pagination.current
  pageReq.value.pageSize = pagination.pageSize
  expandedRowKeys.value = []
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

const approveComment = (record: AdminCommentVO) => {
  if (record.type === 'comment') {
    approveCommentById(record.id)
  } else {
    approveReplyById(record.fid || '', record.id)
  }
}

const approveCommentById = async (id: string) => {
  try {
    const response: any = await ApproveCommentById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('审核成功')
    await get()
  } catch (error) {
    console.log(error)
  }
}

const approveReplyById = async (fid: string, id: string) => {
  try {
    const response: any = await ApproveReplyById(fid, id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('审核成功')
    await get()
  } catch (error) {
    console.log(error)
  }
}

const deleteById = (record: AdminCommentVO) => {
  if (record.type === 'comment') {
    deleteCommentById(record.id)
  } else {
    deleteReplyById(record.fid || '', record.id)
  }
}

const deleteCommentById = async (id: string) => {
  try {
    const response: any = await DeleteCommentById(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await get()
  } catch (error) {
    console.log(error)
  }
}

const deleteReplyById = async (fid: string, id: string) => {
  try {
    const response: any = await DeleteReplyById(fid, id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}

const selectedComments = ref<BatchApprovedCommentRequest>({
  comment_ids: [],
  replies: {}
})

type Key = string | number
const selectedRowKeys = ref<Key[]>([])

const onChange = (srks: Key[], selectedRows: AdminCommentVO[]) => {
  selectedRowKeys.value = srks
  selectedComments.value.comment_ids = []
  selectedComments.value.replies = {}
  selectedRows.forEach((row) => {
    if (row.key) {
      if (row.key.includes('~')) {
        // 如果字符串包含 '~'，按第二种数据处理
        // 假设需要分割存储两部分的数据
        const parts = row.key.split('~')
        const commentId = parts[0]
        const replyId = parts[1]
        if (selectedComments.value.replies[commentId]) {
          selectedComments.value.replies[commentId].push(replyId)
        } else {
          selectedComments.value.replies[commentId] = [replyId]
        }
      } else {
        selectedComments.value.comment_ids.push(row.key)
      }
    }
  })
}
const batchApproveComment = async () => {
  try {
    const apiResponse = await batchApproved(selectedComments.value)
    if (apiResponse.data?.code === 0) {
      message.success('批量审核成功')
      await get()
      selectedRowKeys.value = []
      return true
    } else {
      message.error('批量审核失败')
      return false
    }
  } catch (error) {
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 400) {
          message.error('请选中需要审核的评论或回复')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    } else {
      console.log(error)
      message.error('未知错误，批量审核失败')
    }
    return false
  }
}

const approveStatus = ref(-1)

const statusList = [
  { label: '全部', value: -1 },
  { label: '待审核', value: 0 },
  { label: '已审核', value: 1 }
]
const handleChange = (value: number) => {
  pageReq.value.approvalStatus = value === -1 ? undefined : value !== 0
  get()
}

const selection = computed(() => {
  return {
    selectedRowKeys: selectedRowKeys,
    onChange: onChange,
    hideDefaultSelections: true,
    selections: [
      Table.SELECTION_ALL,
      Table.SELECTION_NONE,
      {
        key: 'approval',
        text: '选中审核',
        onSelect: (changableRowKeys: string[]) => {
          let newSelectedRowKeys: string[]
          newSelectedRowKeys = changableRowKeys.filter((_key: string) => {
            return commentMap.get(_key)?.approval_status
          })
          selectedRowKeys.value = newSelectedRowKeys
        }
      },
      {
        key: 'disapproval',
        text: '选中未审核',
        onSelect: (changableRowKeys: string[]) => {
          let newSelectedRowKeys: string[]
          newSelectedRowKeys = changableRowKeys.filter((_key) => {
            return !commentMap.get(_key)?.approval_status
          })
          selectedRowKeys.value = newSelectedRowKeys
        }
      }
    ]
  }
})

const expandedRowKeys = ref<String[]>([])
const expandedRowsChange = (rowKeys: String[]) => {
  expandedRowKeys.value = rowKeys
}

const expandOrHideRows = () => {
  if (expandedRowKeys.value?.length === 0) {
    expandedRowKeys.value = data.value.map((item) => item.key || '')
  } else {
    expandedRowKeys.value = []
  }
}

const batchDeleteComment = async () => {
  const deleteRequest: BatchApprovedCommentRequest = JSON.parse(
    JSON.stringify(selectedComments.value)
  )
  deleteRequest.comment_ids.forEach((commentId: string) => {
    if (deleteRequest.replies[commentId]) {
      delete deleteRequest.replies[commentId]
    }
  })
  try {
    const apiResponse = await batchDelete(deleteRequest)
    if (apiResponse.data?.code === 0) {
      message.success('批量删除成功')
      await get()
      selectedRowKeys.value = []
      return true
    } else {
      message.error('批量删除失败')
      return false
    }
  } catch (error) {
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 400) {
          message.error('请选中需要删除的评论或回复')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    } else {
      console.log(error)
      message.error('未知错误，批量删除失败')
    }
    return false
  }
}

const refresh = () => {
  get()
  selectedRowKeys.value = []
  expandedRowKeys.value = []
}
</script>

<style scoped>
.editable-row-operations a {
  margin-right: 8px;
}
</style>
