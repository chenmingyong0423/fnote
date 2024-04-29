<template>
  <a-card title="评论列表">
    <div>
      <div>
        状态：
        <a-select
          ref="select"
          v-model:value="pageReq.status"
          style="width: 120px"
          :options="[{ value: 1, label: '1' }]"
          @focus="focus"
          @change="handleChange"
        ></a-select>
      </div>
    </div>
    <a-table
      :columns="columns"
      :data-source="data"
      :pagination="pagination"
      @change="change"
      :row-selection="rowSelection"
      childrenColumnName="replies"
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
  GetComments
} from '@/interfaces/Comment'
import { message } from 'ant-design-vue'

document.title = '评论列表 - 后台管理'

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
    dataIndex: 'created_at'
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
  pageSize: 5,
  sortField: 'create_time',
  sortOrder: 'desc',
  status: 1
} as PageRequest)

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const get = async () => {
  try {
    const response: any = await GetComments(pageReq.value)
    const result: IResponse<IPageData<AdminCommentVO>> = response.data
    if (result.code === 0) {
      data.value = response.data.data?.list || []
      data.value.forEach((commentVO: AdminCommentVO) => {
        commentVO.key = commentVO.id
        commentVO.replies?.forEach((replyVO: AdminCommentVO) => {
          replyVO.fid = commentVO.id
          replyVO.key = replyVO.id
        })
      })
      total.value = response.data.data?.totalCount || 0
      console.log(data.value)
    }
  } catch (error) {
    console.log(error)
  }
}
get()
const change = (pg: any) => {
  pageReq.value.pageNo = pg.current
  pageReq.value.pageSize = pg.pageSize
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

const rowSelection = ref({
  checkStrictly: false,
  onChange: (selectedRowKeys: (string | number)[], selectedRows: Comment[]) => {
    console.log(`selectedRowKeys: ${selectedRowKeys}`, 'selectedRows: ', selectedRows)
  },
  onSelect: (record: Comment, selected: boolean, selectedRows: Comment[]) => {
    console.log(record, selected, selectedRows)
  },
  onSelectAll: (selected: boolean, selectedRows: Comment[], changeRows: Comment[]) => {
    console.log(selected, selectedRows, changeRows)
  }
})

const focus = () => {
  console.log('focus')
}

const handleChange = (value: string) => {
  console.log(`selected ${value}`)
}
</script>

<style scoped>
.editable-row-operations a {
  margin-right: 8px;
}
</style>
