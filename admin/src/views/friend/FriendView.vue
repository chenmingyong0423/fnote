<template>
  <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="change">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'logo'">
        <a-input
          v-if="editableData[record.id]"
          v-model:value="editableData[record.id][column.dataIndex]"
        >
        </a-input>
        <template v-else>
          <a-image :width="50" :src="record.logo" />
        </template>
      </template>
      <template v-if="['name', 'description'].includes(column.dataIndex)">
        <div>
          <a-input
            v-if="editableData[record.id]"
            v-model:value="editableData[record.id][column.dataIndex]"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ text }}
          </template>
        </div>
      </template>
      <template v-if="column.dataIndex === 'status'">
        <a-radio-group
          v-if="
            editableData[record.id] &&
            editableData[record.id][column.dataIndex] != 3 &&
            editableData[record.id][column.dataIndex] != 0
          "
          v-model:value="editableData[record.id][column.dataIndex]"
        >
          <a-radio :value="2">隐藏</a-radio>
          <a-radio :value="1">展示</a-radio>
        </a-radio-group>
        <template v-else>
          <a-tag
            :color="
              record.status === 0 ? 'processing' : record.status === 1 ? 'success' : 'warning'
            "
            >{{ statusConvert(record.status) }}
          </a-tag>
        </template>
      </template>
      <template v-if="column.dataIndex === 'create_time'">
        {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <div class="editable-row-operations">
          <a-modal
            v-model:open="approvalDialog"
            title="请输入博客地址，默认值为当前域名。"
            @ok="approved"
          >
            <a-input v-model:value="blogUrl" placeholder="请输入博客地址，默认值为当前域名。" />
          </a-modal>

          <span v-if="data.length && record.status === 0" @click="openApprovalDialog(record.id)">
            <a>接受</a>
          </span>

          <a-modal
            v-model:open="rejectionDialog"
            title="请输入博客地址，默认值为当前域名。"
            @ok="rejected"
          >
            <a-input v-model:value="blogUrl" placeholder="请输入博客地址，默认值为当前域名。" />
            <a-input v-model:value="reason" placeholder="请输入审核不通过的原因。" />
          </a-modal>
          <span v-if="data.length && record.status === 0" @click="openRejectionDialog(record.id)">
            <a>拒绝</a>
          </span>

          <span v-if="editableData[record.id]">
            <a-typography-link @click="save(record.id)">保存</a-typography-link>
            <a-popconfirm title="确定取消？" @confirm="cancel(record.id)">
              <a>取消</a>
            </a-popconfirm>
          </span>
          <span v-else>
            <a @click="edit(record.id)">编辑</a>
          </span>

          <a-popconfirm v-if="data.length" title="确认删除？" @confirm="deleteInfo(record.id)">
            <a>删除</a>
          </a-popconfirm>
        </div>
      </template>
    </template>
  </a-table>
</template>

<script lang="ts" setup>
import { computed, reactive, ref, type UnwrapRef } from 'vue'
import {
  ApproveFriend,
  DeleteFriend,
  type Friend,
  type FriendReq,
  GetFriends,
  RejectFriend
} from '@/interfaces/Friend'
import type { PageRequest } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'
import { UpdateFriend } from '@/interfaces/Friend'

const columns = [
  {
    title: 'logo',
    dataIndex: 'logo',
    key: 'logo'
  },
  {
    title: '站点链接',
    key: 'url',
    dataIndex: 'url'
  },
  {
    title: '站点名称',
    dataIndex: 'name',
    key: 'name'
  },
  {
    title: '站点描述',
    key: 'description',
    dataIndex: 'description'
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

const data = ref<Friend[]>([])
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
    const response: any = await GetFriends(pageReq.value)
    data.value = response.data?.list || []
    total.value = response.data?.totalCount || 0
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
// 删除
const deleteInfo = async (id: string) => {
  try {
    const response: any = await DeleteFriend(id)
    if (response.code !== 0) {
      message.error(response.message)
      return
    }
    message.success('删除成功')
    await get()
  } catch (error) {
    console.log(error)
  }
}

const blogUrl = ref(window.location.host)
const approvalDialog = ref(false)

const updatedId = ref('')

const openApprovalDialog = (id: string) => {
  updatedId.value = id
  approvalDialog.value = true
}

const approved = async () => {
  try {
    const response: any = await ApproveFriend(updatedId.value, blogUrl.value)
    if (response.code !== 0) {
      message.error(response.message)
      return
    }
    message.success('接受成功')
    approvalDialog.value = false
    updatedId.value = ''
    await get()
  } catch (error) {
    console.log(error)
  }
}

const openRejectionDialog = (id: string) => {
  updatedId.value = id
  rejectionDialog.value = true
}

const rejectionDialog = ref(false)
const reason = ref('')
const rejected = async () => {
  try {
    const response: any = await RejectFriend(updatedId.value, blogUrl.value, reason.value)
    if (response.code !== 0) {
      message.error(response.message)
      return
    }
    message.success('拒绝成功')
    await get()
    rejectionDialog.value = false
    updatedId.value = ''
  } catch (error) {
    console.log(error)
  }
}

// 编辑
const editableData: UnwrapRef<Record<string, FriendReq>> = reactive({})
const edit = (id: string) => {
  editableData[id] = cloneDeep(data.value.filter((item) => id === item.id)[0])
}

const save = async (id: string) => {
  const editableDatum = editableData[id]
  try {
    const response: any = await UpdateFriend(id, editableDatum)
    if (response.code !== 0) {
      message.error(response.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await get()
  } catch (error) {
    console.log(error)
  }
}
const cancel = (key: string) => {
  delete editableData[key]
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
</script>

<style scoped>
.collection-create-form_last-form-item {
  margin-bottom: 0;
}

.editable-row-operations a {
  margin-right: 8px;
}
</style>
