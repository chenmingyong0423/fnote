<template>
  <a-card title="文件管理">
    <template #extra>
      <a-button :icon="h(ReloadOutlined)" :loading="loading" @click="loadFiles">刷新</a-button>
    </template>

    <div class="file-toolbar">
      <SimpleUpload
        :action="serverHost + '/admin-api/files/upload'"
        :authorization="userStore.token"
        :file-types="['image/jpeg', 'image/png']"
        :max-size="5 * 1024 * 1024"
        label="上传图片"
        @success:imageUrl="handleUploadSuccess"
      />
      <a-select v-model:value="fileType" class="w-45" @change="changeFileType">
        <a-select-option value="all">全部文件</a-select-option>
        <a-select-option value="image">图片</a-select-option>
        <a-select-option value="image/png">PNG</a-select-option>
        <a-select-option value="image/jpeg">JPEG</a-select-option>
      </a-select>
    </div>

    <a-table
      :columns="columns"
      :data-source="files"
      :loading="loading"
      :pagination="pagination"
      row-key="file_id"
      bordered
      @change="changePage"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'preview'">
          <a-image
            v-if="record.file_type.startsWith('image/')"
            :src="serverHost + record.url"
            :width="64"
            :height="64"
            class="file-preview"
          />
          <a-tag v-else>{{ record.file_type || '未知类型' }}</a-tag>
        </template>

        <template v-else-if="column.key === 'name'">
          <div>{{ record.original_file_name || record.file_name }}</div>
          <a-typography-text type="secondary" copyable>{{ record.file_name }}</a-typography-text>
        </template>

        <template v-else-if="column.key === 'size'">
          {{ formatFileSize(record.file_size) }}
        </template>

        <template v-else-if="column.key === 'usage'">
          <span v-if="record.used_in.length === 0">未使用</span>
          <div v-else class="file-usages">
            <a-tag v-for="usage in record.used_in" :key="`${usage.type}:${usage.id}`" color="blue">
              {{ usage.type }}：{{ usage.id }}
            </a-tag>
          </div>
        </template>

        <template v-else-if="column.key === 'created_at'">
          {{ formatTime(record.created_at) }}
        </template>

        <template v-else-if="column.key === 'operation'">
          <a :href="serverHost + record.url" target="_blank" rel="noopener noreferrer">打开</a>
          <a-divider type="vertical" />
          <a-tooltip :title="record.used_in.length > 0 ? '文件正在被使用，不能删除' : ''">
            <a-popconfirm
              title="确认删除文件？物理文件也会一并删除。"
              :disabled="record.used_in.length > 0"
              @confirm="deleteFile(record)"
            >
              <a-button
                type="link"
                danger
                size="small"
                :disabled="record.used_in.length > 0"
                :loading="deletingFileId === record.file_id"
              >
                删除
              </a-button>
            </a-popconfirm>
          </a-tooltip>
        </template>
      </template>
    </a-table>
  </a-card>
</template>

<script lang="ts" setup>
import { computed, h, ref } from 'vue'
import dayjs from 'dayjs'
import { message, type TableColumnType, type TableProps } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import SimpleUpload from '@/components/upload/SimpleUpload.vue'
import { DeleteFile, type FileVO, GetFileList } from '@/interfaces/File'
import { useUserStore } from '@/stores/user'
import { toErrorMessage } from '@/utils/error'

document.title = '文件管理 - 后台管理'

const serverHost = import.meta.env.VITE_API_HOST
const userStore = useUserStore()
const files = ref<FileVO[]>([])
const loading = ref(false)
const deletingFileId = ref('')
const fileType = ref('all')
const pageNum = ref(1)
const pageSize = ref(10)
const total = ref(0)

const columns: TableColumnType<FileVO>[] = [
  { title: '预览', key: 'preview', width: 100 },
  { title: '文件名', key: 'name', width: 260 },
  { title: '类型', dataIndex: 'file_type', key: 'file_type', width: 150 },
  { title: '大小', key: 'size', width: 110 },
  { title: '使用情况', key: 'usage', width: 300 },
  { title: '上传时间', key: 'created_at', width: 180 },
  { title: '操作', key: 'operation', fixed: 'right', width: 150 }
]

const pagination = computed(() => ({
  current: pageNum.value,
  pageSize: pageSize.value,
  total: total.value,
  showSizeChanger: true,
  showTotal: (value: number) => `共 ${value} 个文件`
}))

const selectedFileTypes = () => {
  if (fileType.value === 'all') return []
  if (fileType.value === 'image') return ['image/png', 'image/jpeg']
  return [fileType.value]
}

const loadFiles = async () => {
  loading.value = true
  try {
    const response: any = await GetFileList({
      pageNum: pageNum.value,
      pageSize: pageSize.value,
      fileType: selectedFileTypes()
    })
    files.value = response.data.data?.list || []
    total.value = response.data.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '文件列表加载失败'))
  } finally {
    loading.value = false
  }
}

const changePage: TableProps<FileVO>['onChange'] = (page) => {
  pageNum.value = Number(page.current || 1)
  pageSize.value = Number(page.pageSize || 10)
  loadFiles()
}

const changeFileType = () => {
  pageNum.value = 1
  loadFiles()
}

const handleUploadSuccess = () => {
  message.success('上传成功')
  pageNum.value = 1
  loadFiles()
}

const deleteFile = async (file: FileVO) => {
  deletingFileId.value = file.file_id
  try {
    await DeleteFile(file.file_id)
    message.success('删除成功')
    if (files.value.length === 1 && pageNum.value > 1) pageNum.value--
    await loadFiles()
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    deletingFileId.value = ''
  }
}

const formatFileSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

const formatTime = (timestamp: number) => dayjs.unix(timestamp).format('YYYY-MM-DD HH:mm:ss')

loadFiles()
</script>

<style scoped>
.file-toolbar {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 16px;
}

.file-preview {
  object-fit: contain;
}

.file-usages {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}
</style>
