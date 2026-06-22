<template>
  <a-modal
    v-model:open="visible"
    title="从文件库选择"
    width="900px"
    :footer="null"
    :destroy-on-close="true"
    @cancel="close"
  >
    <div class="file-picker">
      <div class="picker-toolbar">
        <a-input-search
          v-model:value="keyword"
          class="search-input"
          allow-clear
          placeholder="搜索原始文件名"
          @search="search"
        />
        <a-select v-model:value="fileType" class="type-select" @change="changeFileType">
          <a-select-option value="all">全部格式</a-select-option>
          <a-select-option value="image/jpeg">JPEG</a-select-option>
          <a-select-option value="image/png">PNG</a-select-option>
        </a-select>
      </div>

      <a-alert
        class="picker-tip"
        type="info"
        show-icon
        message="已在素材库中的文件不能重复添加；选择新文件后可直接插入文章。"
      />

      <a-spin :spinning="loading">
        <div v-if="images.length" class="file-grid">
          <button
            v-for="image in images"
            :key="image.file_id"
            type="button"
            :class="[
              'file-card',
              {
                selected: isSelected(image.file_id),
                disabled: isInAssetLibrary(image)
              }
            ]"
            :aria-disabled="isInAssetLibrary(image)"
            @click="selectImage(image)"
          >
            <span class="image-stage">
              <img :src="serverHost + image.url" :alt="displayName(image)" />
              <span class="image-overlay">
                <a-tooltip title="查看大图">
                  <span class="preview-button" @click.stop="preview(image.url)">
                    <EyeOutlined />
                  </span>
                </a-tooltip>
              </span>
              <CheckCircleFilled v-if="isSelected(image.file_id)" class="selected-mark" />
              <a-tag v-if="isInAssetLibrary(image)" class="used-tag" color="default">
                已在素材库
              </a-tag>
            </span>
            <span class="file-info">
              <span class="file-name" :title="displayName(image)">{{ displayName(image) }}</span>
              <span class="file-meta"
                >{{ formatSize(image.file_size) }} · {{ typeLabel(image.file_type) }}</span
              >
            </span>
          </button>
        </div>
        <a-empty v-else class="picker-empty" :image="simpleImage" description="没有找到图片" />
      </a-spin>

      <div class="picker-footer">
        <a-pagination
          v-model:current="pagination.pageNum"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          :show-size-changer="false"
          size="small"
          @change="changePage"
        />
        <a-space>
          <a-button @click="close">取消</a-button>
          <a-button :disabled="!selectedFile" @click="confirm(false)">加入素材库</a-button>
          <a-button type="primary" :disabled="!selectedFile" @click="confirm(true)">
            加入并插入文章
          </a-button>
        </a-space>
      </div>
    </div>

    <PreviewImg v-model="previewVisible" :image-url="previewImageUrl" />
  </a-modal>
</template>

<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue'
import { CheckCircleFilled, EyeOutlined } from '@ant-design/icons-vue'
import { Empty, message } from 'ant-design-vue'
import { type FileVO, GetFileList } from '@/interfaces/File'
import PreviewImg from '@/components/image/PreviewImg.vue'
import { toErrorMessage } from '@/utils/error'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['insertImg', 'update:modelValue'])
const visible = ref(props.modelValue)
const loading = ref(false)
const images = ref<FileVO[]>([])
const selectedFileId = ref('')
const keyword = ref('')
const fileType = ref('all')
const serverHost = import.meta.env.VITE_API_HOST
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const pagination = reactive({ pageNum: 1, pageSize: 12, total: 0 })
const selectedFile = computed(() =>
  images.value.find((image) => image.file_id === selectedFileId.value)
)

watch(
  () => props.modelValue,
  async (value) => {
    visible.value = value
    if (value) {
      selectedFileId.value = ''
      pagination.pageNum = 1
      await loadImages()
    }
  }
)

const loadImages = async () => {
  loading.value = true
  selectedFileId.value = ''
  try {
    const types = fileType.value === 'all' ? ['image/png', 'image/jpeg'] : [fileType.value]
    const response = await GetFileList({
      pageNum: pagination.pageNum,
      pageSize: pagination.pageSize,
      fileType: types,
      keyword: keyword.value.trim()
    })
    images.value = response.data.data?.list || []
    pagination.total = response.data.data?.totalCount || 0
  } catch (error) {
    images.value = []
    pagination.total = 0
    message.error(toErrorMessage(error, '文件库加载失败'))
  } finally {
    loading.value = false
  }
}

const isInAssetLibrary = (image: FileVO) =>
  image.used_in?.some((usage) => usage.type === 'asset') || false

const selectImage = (image: FileVO) => {
  if (isInAssetLibrary(image)) return
  selectedFileId.value = selectedFileId.value === image.file_id ? '' : image.file_id
}

const isSelected = (id: string) => selectedFileId.value === id
const displayName = (image: FileVO) => image.original_file_name || image.file_name
const typeLabel = (type: string) => (type === 'image/jpeg' ? 'JPEG' : 'PNG')

const formatSize = (size: number) => {
  if (size < 1024) return `${size} B`
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`
  return `${(size / 1024 / 1024).toFixed(1)} MB`
}

const search = async () => {
  pagination.pageNum = 1
  await loadImages()
}

const changeFileType = async () => {
  pagination.pageNum = 1
  await loadImages()
}

const changePage = async (page: number) => {
  pagination.pageNum = page
  await loadImages()
}

const confirm = (insertAfterAdd: boolean) => {
  if (!selectedFile.value) return
  const image = selectedFile.value
  emit('insertImg', image.file_id, image.url, displayName(image), insertAfterAdd)
  close()
}

const close = () => {
  selectedFileId.value = ''
  visible.value = false
  emit('update:modelValue', false)
}

const previewVisible = ref(false)
const previewImageUrl = ref('')
const preview = (url: string) => {
  previewImageUrl.value = serverHost + url
  previewVisible.value = true
}

if (visible.value) loadImages()
</script>

<style scoped>
.file-picker {
  min-height: 570px;
}
.picker-toolbar,
.picker-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.search-input {
  width: 320px;
}
.type-select {
  width: 130px;
}
.picker-tip {
  margin: 14px 0;
}
.file-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  min-height: 386px;
}
.file-card {
  min-width: 0;
  padding: 0;
  overflow: hidden;
  border: 1px solid var(--app-border, #d9d9d9);
  border-radius: 8px;
  color: inherit;
  background: var(--app-surface, #fff);
  cursor: pointer;
  text-align: left;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}
.file-card:hover:not(.disabled),
.file-card.selected {
  border-color: #1677ff;
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.14);
}
.file-card.disabled {
  cursor: not-allowed;
  opacity: 0.62;
}
.image-stage {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 120px;
  overflow: hidden;
  background: var(--app-surface-muted, #fafafa);
}
.image-stage img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.image-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  background: rgba(0, 0, 0, 0.42);
  transition: opacity 0.2s ease;
}
.file-card:hover .image-overlay {
  opacity: 1;
}
.preview-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  color: #fff;
  font-size: 17px;
  background: rgba(255, 255, 255, 0.2);
  cursor: pointer;
}
.selected-mark {
  position: absolute;
  top: 7px;
  right: 7px;
  z-index: 1;
  color: #1677ff;
  font-size: 21px;
}
.used-tag {
  position: absolute;
  top: 7px;
  left: 7px;
  margin: 0;
}
.file-info {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 8px 10px;
}
.file-name {
  overflow: hidden;
  font-size: 13px;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.file-meta {
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.45));
  font-size: 11px;
}
.picker-empty {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 386px;
}
.picker-footer {
  margin-top: 16px;
  padding-top: 14px;
  border-top: 1px solid var(--app-border, #f0f0f0);
}
@media (max-width: 760px) {
  .file-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
  .picker-footer {
    align-items: flex-end;
    flex-direction: column;
  }
}
</style>
