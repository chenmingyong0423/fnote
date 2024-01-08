<template>
  <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="change">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'post'">
        <a :href="record.post_info.post_url" target="_blank">{{ record.post_info.post_url }}</a>
      </template>
      <template v-if="column.dataIndex === 'type'">
        <a-tag color="success">{{ record.type == 0 ? '评论' : '回复' }}</a-tag>
      </template>
      <template v-if="column.dataIndex === 'replied_content'">
        {{ record.replied_content }}
      </template>
      <template v-if="column.dataIndex === 'status'">
        <a-tag
          :color="record.status === 0 ? 'processing' : record.status === 1 ? 'success' : 'warning'"
          >{{ statusConvert(record.status) }}</a-tag
        >
      </template>
      <template v-if="column.dataIndex === 'create_time'">
        {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <div class="editable-row-operations">
          <a-popconfirm
            v-if="data.length && record.status === 0"
            title="确认通过？"
            @confirm="approveComment(record)"
          >
            <a>通过</a>
          </a-popconfirm>
          <a-popconfirm
            v-if="data.length && record.status === 0"
            title="确认驳回？"
            @confirm="openDisapproveDialog(record)"
          >
            <a>驳回</a>
          </a-popconfirm>
          <a-popconfirm
            v-if="data.length && record.status === 2"
            title="确认显示？"
            @confirm="updateStatus(record, 1)"
          >
            <a>显示</a>
          </a-popconfirm>
          <a-popconfirm
            v-if="data.length && record.status === 1"
            title="确认隐藏？"
            @confirm="updateStatus(record, 2)"
          >
            <a>隐藏</a>
          </a-popconfirm>
          <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deleteById(record)">
            <a>删除</a>
          </a-popconfirm>
        </div>
      </template>
    </template>
  </a-table>
  <a-modal v-model:open="disapproveDialog" title="驳回原因" @ok="disapproveComment">
    <a-input v-model:value="reason" placeholder="请输入审核不通过的原因。" />
  </a-modal>
</template>
<script setup lang="ts">
import dayjs from 'dayjs'
import { computed, ref } from 'vue'
import type { IBaseResponse, IPageData, IResponse, PageRequest } from '@/interfaces/Common'
import axios from '@/http/axios'
import type { Comment } from '@/interfaces/Comment'
import { message } from 'ant-design-vue'

const domain: string = window.location.host

const columns = [
  {
    title: '文章',
    dataIndex: 'post',
    key: 'post'
  },
  {
    title: '类型',
    key: 'type',
    dataIndex: 'type'
  },
  {
    title: '内容',
    dataIndex: 'content',
    key: 'content'
  },
  {
    title: '状态',
    key: 'status',
    dataIndex: 'status'
  },
  {
    title: '提交时间',
    key: 'create_time',
    dataIndex: 'create_time'
  },
  {
    title: 'operation',
    dataIndex: 'operation'
  }
]

const data = ref<Comment[]>([])
const pageReq = ref<PageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'create_time',
  sortOrder: 'desc'
} as PageRequest)

const total = ref(0)

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const get = async () => {
  try {
    const response = await axios.get<IResponse<IPageData<Comment>>>('/admin/comments', {
      params: pageReq.value
    })
    data.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
  }
}
get()
const change = (pg, filters, sorter, { currentDataSource }) => {
  pageReq.value.pageNo = pg.current
  pageReq.value.pageSize = pg.pageSize
  get()
}
const statusConvert = (status: number) => {
  switch (status) {
    case 0:
      return '未审核'
    case 1:
      return '显示'
    case 2:
      return '隐藏'
    default:
      return '审核不通过'
  }
}

const approveComment = (record: Comment) => {
  if (record.type === 0) {
    approveCommentById(record.id)
  } else {
    approveReplyById(record.fid || '', record.id)
  }
}

const approveCommentById = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/comments/${id}/approval`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('审核成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('审核失败')
  }
}

const approveReplyById = async (fid: string, id: string) => {
  try {
    const response = await axios.put<IBaseResponse>(`/admin/comments/${fid}/replies/${id}/approval`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('审核成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('审核失败')
  }
}

const disapproveDialog = ref(false)
const comment = ref<Comment>()
const reason = ref('')

const openDisapproveDialog = (record: Comment) => {
  disapproveDialog.value = true
  comment.value = record
}
const disapproveComment = () => {
  if (comment.value?.type === 0) {
    disapproveCommentById(comment.value.id!)
  } else {
    disapproveReplyById(comment.value?.fid || '', comment.value?.id!)
  }
}

const disapproveCommentById = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/comments/${id}/disapproval`, {
      reason: reason.value
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('驳回成功')
    await get()
    reason.value = ''
    disapproveDialog.value = false
  } catch (error) {
    console.log(error)
    message.error('驳回失败')
  }
}

const disapproveReplyById = async (fid: string, id: string) => {
  try {
    const response = await axios.put<IBaseResponse>(
      `/admin/comments/${fid}/replies/${id}/disapproval`,
      {
        reason: reason.value
      }
    )
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('驳回成功')
    await get()
    reason.value = ''
    disapproveDialog.value = false
  } catch (error) {
    console.log(error)
    message.error('驳回失败')
  }
}

const updateStatus = (record: Comment, status: number) => {
  if (record.type === 0) {
    updateStatusById(record.id, status)
  } else {
    updateReplyStatusById(record.fid, record.id, status)
  }
}

const updateStatusById = async (id: string, status: number) => {
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/comments/${id}/status`, {
      status: status
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('更新失败')
  }
}

const updateReplyStatusById = async (fid: string, id: string, status: number) => {
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/comments/${fid}/replies/${id}/status`, {
      status: status
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('更新失败')
  }
}

const deleteById = (record: Comment) => {
  if (record.type === 0) {
    deleteCommentById(record.id)
  } else {
    deleteReplyById(record.fid, record.id)
  }
}

const deleteCommentById = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.delete<IBaseResponse>(`/admin/comments/${id}`)
    if (response.data.code !== 200) {
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

const deleteReplyById = async (fid: string, id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.delete<IBaseResponse>(`/admin/comments/${fid}/replies/${id}`)
    if (response.data.code !== 200) {
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
</script>

<style scoped>
.editable-row-operations a {
  margin-right: 8px;
}
</style>
