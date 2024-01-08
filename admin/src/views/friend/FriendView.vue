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
          v-if="editableData[record.id]"
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
            >{{ statusConvert(record.status) }}</a-tag
          >
        </template>
      </template>
      <template v-if="column.dataIndex === 'create_time'">
        {{ dayjs.unix(text).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
      <template v-else-if="column.dataIndex === 'operation'">
        <div class="editable-row-operations">
          <a-popconfirm
            v-if="data.length && record.status === 0"
            title="确认接受？"
            @confirm="approved(record.id)"
          >
            <a>接受</a>
          </a-popconfirm>

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
import type { Friend, FriendReq } from '@/interfaces/Friend'
import axios from '@/http/axios'
import type { IBaseResponse, IPageData, IResponse, PageRequest } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'

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
    const response = await axios.get<IResponse<IPageData<Friend>>>('/admin/friends', {
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
// 删除
const deleteInfo = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.delete<IBaseResponse>(`/admin/friends/${id}`)
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

const approved = async (id: string) => {
  try {
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/friends/${id}/approval`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('接受成功')
    await get()
  } catch (error) {
    console.log(error)
    message.error('接受失败')
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
    // 提交 body 参数 values
    const response = await axios.put<IBaseResponse>(`/admin/friends/${id}`, editableDatum)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await get()
  } catch (error) {
    console.log(error)
    message.error('更新失败')
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
