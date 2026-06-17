<template>
  <a-card title="友链列表">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button shape="circle" :icon="h(ReloadOutlined)" :loading="loading" @click="refresh" />
      </a-tooltip>
    </template>

    <div class="friend-toolbar">
      <a-input-search
        v-model:value="pageReq.keyword"
        allow-clear
        class="friend-search"
        placeholder="搜索站点名称"
        @search="searchFriends"
        @pressEnter="searchFriends"
      />
      <a-segmented
        v-model:value="statusFilter"
        :options="statusFilterOptions"
        @change="changeStatusFilter"
      />
    </div>

    <a-spin :spinning="loading">
      <a-table
        :columns="columns"
        :data-source="data"
        :pagination="pagination"
        :scroll="{ x: 1180 }"
        row-key="id"
        bordered
        @change="change"
      >
        <template #bodyCell="{ column, record, text }">
          <template v-if="column.key === 'logo'">
            <a-input
              v-if="editableData[record.id]"
              v-model:value="editableData[record.id].logo"
              placeholder="Logo 地址"
            />
            <a-image v-else :width="56" :height="56" :src="record.logo" class="friend-logo" />
          </template>

          <template v-else-if="column.key === 'site'">
            <div class="friend-site">
              <a-input
                v-if="editableData[record.id]"
                v-model:value="editableData[record.id].name"
                placeholder="站点名称"
              />
              <template v-else>
                <a :href="record.url" target="_blank" rel="noopener noreferrer">
                  {{ record.name }}
                </a>
                <span class="friend-url">{{ record.url }}</span>
              </template>
              <a-input
                v-if="editableData[record.id]"
                v-model:value="editableData[record.id].url"
                placeholder="站点链接"
              />
            </div>
          </template>

          <template v-else-if="column.key === 'description'">
            <a-textarea
              v-if="editableData[record.id]"
              v-model:value="editableData[record.id].description"
              :auto-size="{ minRows: 2, maxRows: 3 }"
              placeholder="站点描述"
            />
            <a-tooltip v-else :title="record.description">
              <div class="friend-description">{{ record.description }}</div>
            </a-tooltip>
          </template>

          <template v-else-if="column.key === 'status'">
            <a-radio-group
              v-if="editableData[record.id] && canEditStatus(record)"
              v-model:value="editableData[record.id].status"
              button-style="solid"
            >
              <a-radio-button :value="1">展示</a-radio-button>
              <a-radio-button :value="2">隐藏</a-radio-button>
            </a-radio-group>
            <a-tag v-else :color="statusMeta(record.status).color">
              {{ statusMeta(record.status).label }}
            </a-tag>
          </template>

          <template v-else-if="column.key === 'created_at'">
            {{ formatTime(Number(text)) }}
          </template>

          <template v-else-if="column.key === 'operation'">
            <div class="friend-actions">
              <template v-if="record.status === 0">
                <a-popconfirm title="确认接受？" @confirm="approveFriend(record.id)">
                  <a-button type="link" size="small" :loading="isApproving(record.id)">
                    接受
                  </a-button>
                </a-popconfirm>
                <a-button type="link" size="small" danger @click="openRejectionDialog(record.id)">
                  拒绝
                </a-button>
              </template>

              <template v-if="editableData[record.id]">
                <a-button
                  type="link"
                  size="small"
                  :loading="isSaving(record.id)"
                  @click="save(record.id)"
                >
                  保存
                </a-button>
                <a-popconfirm title="确定取消？" @confirm="cancel(record.id)">
                  <a-button type="link" size="small">取消</a-button>
                </a-popconfirm>
              </template>
              <a-button v-else type="link" size="small" @click="edit(record)">编辑</a-button>

              <a-popconfirm title="确认删除？" @confirm="deleteInfo(record.id)">
                <a-button type="link" size="small" danger :loading="isDeleting(record.id)">
                  删除
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-spin>

    <a-modal
      v-model:open="rejectionDialog"
      title="拒绝友链申请"
      ok-text="确认拒绝"
      cancel-text="取消"
      :confirm-loading="rejecting"
      @ok="rejectFriend"
      @cancel="closeRejectionDialog"
    >
      <a-form layout="vertical">
        <a-form-item label="拒绝原因" required>
          <a-textarea
            v-model:value="reason"
            :auto-size="{ minRows: 3, maxRows: 5 }"
            placeholder="请输入审核不通过的原因"
            :maxlength="120"
            show-count
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </a-card>
</template>

<script lang="ts" setup>
import { computed, h, reactive, ref, type UnwrapRef } from 'vue'
import {
  ApproveFriend,
  DeleteFriend,
  type Friend,
  type FriendReq,
  GetFriends,
  RejectFriend,
  UpdateFriend
} from '@/interfaces/Friend'
import type { PageRequest } from '@/interfaces/Common'
import { message } from 'ant-design-vue'
import type { TableColumnType, TableProps } from 'ant-design-vue'
import { cloneDeep } from 'lodash-es'
import dayjs from 'dayjs'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { toErrorMessage } from '@/utils/error'

document.title = '友链列表 - 后台管理'

type StatusFilter = -1 | 0 | 1 | 2 | 3
interface FriendPageRequest extends PageRequest {
  keyword?: string
  sortField?: string
  sortOrder?: string
  status?: number
}

const columns: TableColumnType<Friend>[] = [
  {
    title: 'Logo',
    dataIndex: 'logo',
    key: 'logo',
    width: 110
  },
  {
    title: '站点',
    key: 'site',
    dataIndex: 'url',
    width: 260
  },
  {
    title: '站点描述',
    key: 'description',
    dataIndex: 'description',
    width: 280
  },
  {
    title: '状态',
    key: 'status',
    dataIndex: 'status',
    width: 160
  },
  {
    title: '提交时间',
    key: 'created_at',
    dataIndex: 'created_at',
    sorter: (f1: Friend, f2: Friend) => f1.created_at - f2.created_at,
    defaultSortOrder: 'descend',
    sortDirections: ['descend', 'ascend'],
    width: 180
  },
  {
    title: '操作',
    key: 'operation',
    dataIndex: 'operation',
    fixed: 'right',
    width: 220
  }
]

const data = ref<Friend[]>([])
const pageReq = ref<FriendPageRequest>({
  pageNo: 1,
  pageSize: 5,
  sortField: 'created_at',
  sortOrder: 'DESC',
  keyword: ''
} as PageRequest)
const total = ref(0)
const loading = ref(false)
const approvingIds = ref<Set<string>>(new Set())
const deletingIds = ref<Set<string>>(new Set())
const savingIds = ref<Set<string>>(new Set())
const rejecting = ref(false)
const updatedId = ref('')
const rejectionDialog = ref(false)
const reason = ref('')
const statusFilter = ref<StatusFilter>(-1)
const editableData: UnwrapRef<Record<string, FriendReq>> = reactive({})

const statusFilterOptions = [
  { label: '全部', value: -1 },
  { label: '待审核', value: 0 },
  { label: '展示中', value: 1 },
  { label: '隐藏', value: 2 },
  { label: '已拒绝', value: 3 }
]

const pagination = computed(() => ({
  total: total.value,
  current: pageReq.value.pageNo,
  pageSize: pageReq.value.pageSize
}))

const statusMeta = (status: number) => {
  switch (status) {
    case 0:
      return { label: '待审核', color: 'processing' }
    case 1:
      return { label: '展示中', color: 'success' }
    case 2:
      return { label: '隐藏', color: 'warning' }
    case 3:
      return { label: '已拒绝', color: 'error' }
    default:
      return { label: '未知', color: 'default' }
  }
}

const canEditStatus = (record: Friend) => record.status === 1 || record.status === 2

const formatTime = (timestamp?: number) => {
  if (!timestamp) {
    return '-'
  }
  return dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')
}

const setPending = (source: typeof approvingIds, id: string, pending: boolean) => {
  const next = new Set(source.value)
  if (pending) {
    next.add(id)
  } else {
    next.delete(id)
  }
  source.value = next
}

const isApproving = (id: string) => approvingIds.value.has(id)

const isDeleting = (id: string) => deletingIds.value.has(id)

const isSaving = (id: string) => savingIds.value.has(id)

const clearEditing = () => {
  Object.keys(editableData).forEach((id) => {
    delete editableData[id]
  })
}

const get = async () => {
  try {
    loading.value = true
    const response: any = await GetFriends(pageReq.value)
    if (response.data.code !== 0) {
      message.error(response.data.message || '友链列表加载失败')
      return
    }
    data.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '友链列表加载失败'))
  } finally {
    loading.value = false
  }
}

get()

const change: TableProps<Friend>['onChange'] = (pg, _filters, sorter: any) => {
  pageReq.value.pageNo = Number(pg.current || 1)
  pageReq.value.pageSize = Number(pg.pageSize || 5)
  if (sorter.field && sorter.order) {
    pageReq.value.sortField = sorter.field
    pageReq.value.sortOrder = sorter.order === 'ascend' ? 'ASC' : 'DESC'
  } else {
    pageReq.value.sortField = 'created_at'
    pageReq.value.sortOrder = 'DESC'
  }
  clearEditing()
  get()
}

const searchFriends = () => {
  pageReq.value.pageNo = 1
  clearEditing()
  get()
}

const changeStatusFilter = (value: StatusFilter) => {
  pageReq.value.pageNo = 1
  pageReq.value.status = value === -1 ? undefined : value
  clearEditing()
  get()
}

const refresh = () => {
  clearEditing()
  get()
}

const deleteInfo = async (id: string) => {
  if (isDeleting(id)) {
    return
  }

  try {
    setPending(deletingIds, id, true)
    const response: any = await DeleteFriend(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    delete editableData[id]
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    setPending(deletingIds, id, false)
  }
}

const approveFriend = async (id: string) => {
  if (isApproving(id)) {
    return
  }

  try {
    setPending(approvingIds, id, true)
    const response: any = await ApproveFriend(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('接受成功')
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '接受失败'))
  } finally {
    setPending(approvingIds, id, false)
  }
}

const openRejectionDialog = (id: string) => {
  updatedId.value = id
  reason.value = ''
  rejectionDialog.value = true
}

const closeRejectionDialog = () => {
  rejectionDialog.value = false
  updatedId.value = ''
  reason.value = ''
}

const rejectFriend = async () => {
  if (!reason.value.trim()) {
    message.warning('请输入拒绝原因')
    return
  }

  try {
    rejecting.value = true
    const response: any = await RejectFriend(updatedId.value, reason.value.trim())
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('拒绝成功')
    closeRejectionDialog()
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '拒绝失败'))
  } finally {
    rejecting.value = false
  }
}

const edit = (record: Friend) => {
  editableData[record.id] = cloneDeep({
    name: record.name,
    url: record.url,
    logo: record.logo,
    description: record.description,
    status: record.status
  })
}

const save = async (id: string) => {
  if (isSaving(id)) {
    return
  }

  const editableDatum = editableData[id]
  if (
    !editableDatum.name.trim() ||
    !editableDatum.url.trim() ||
    !editableDatum.logo.trim() ||
    !editableDatum.description.trim()
  ) {
    message.warning('请完整填写站点名称、站点链接、Logo 和描述')
    return
  }
  try {
    new URL(editableDatum.url.trim())
  } catch {
    message.warning('请输入正确的站点链接')
    return
  }

  try {
    setPending(savingIds, id, true)
    const response: any = await UpdateFriend(id, {
      ...editableDatum,
      name: editableDatum.name.trim(),
      url: editableDatum.url.trim(),
      logo: editableDatum.logo.trim(),
      description: editableDatum.description.trim()
    })
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    delete editableData[id]
    await get()
  } catch (error) {
    message.error(toErrorMessage(error, '更新失败'))
  } finally {
    setPending(savingIds, id, false)
  }
}

const cancel = (id: string) => {
  delete editableData[id]
}
</script>

<style scoped>
.friend-toolbar {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.friend-search {
  width: 240px;
}

.friend-logo {
  object-fit: cover;
  border-radius: 6px;
}

.friend-site {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.friend-url {
  overflow: hidden;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.friend-description {
  display: -webkit-box;
  overflow: hidden;
  line-height: 1.5;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.friend-actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
}

.friend-actions :deep(.ant-btn) {
  padding: 0;
}
</style>
