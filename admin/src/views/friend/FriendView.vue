<template>
  <a-table :columns="columns" :data-source="data" :pagination="pagination" @change="change">
    <template #bodyCell="{ column, text, record }">
      <template v-if="column.dataIndex === 'logo'">
        <a-input
          v-if="editableData[record.id]"
          v-model:value="editableData[record.id][column.dataIndex as keyof FriendReq]"
        >
        </a-input>
        <template v-else>
          <a-image :width="50" :src="serverHost + record.logo" />
        </template>
      </template>
      <template v-if="['name', 'description'].includes(column.dataIndex)">
        <div>
          <a-input
            v-if="editableData[record.id]"
            v-model:value="editableData[record.id][column.dataIndex as keyof FriendReq]"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ text }}
          </template>
        </div>
      </template>
      <template v-if="column.dataIndex === 'status'">
        <a-radio-group
          v-if="editableData[record.id] && record.status != 3 && record.status != 0"
          v-model:value="editableData[record.id][column.dataIndex as keyof FriendReq]"
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
          <a-popconfirm
            v-if="data.length && record.status === 0"
            title="确认接受?"
            @confirm="approved(record.id)"
          >
            <a>接收</a>
          </a-popconfirm>
          <a-modal
            v-model:open="rejectionDialog"
            title="请输入原因"
            @ok="rejected"
          >
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
const serverHost = import.meta.env.VITE_API_HOST
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
    data.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
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
// 删除
const deleteInfo = async (id: string) => {
  try {
    const response: any = await DeleteFriend(id)
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

const updatedId = ref('')

const approved = async (id: string) => {
  try {
    const response: any = await ApproveFriend(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('接受成功')
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
    const response: any = await RejectFriend(updatedId.value,reason.value)
    if (response.data.code !== 0) {
      message.error(response.data.message)
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
    if (response.data.code !== 0) {
      message.error(response.data.message)
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
