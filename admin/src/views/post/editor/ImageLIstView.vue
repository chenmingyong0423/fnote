<template>
  <div class="image-assets">
    <aside class="folder-panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">图片分类</div>
          <div class="panel-subtitle">按用途整理常用图片</div>
        </div>
        <a-tooltip title="新增分类">
          <a-button type="text" shape="circle" :icon="h(FolderAddOutlined)" @click="openCreator" />
        </a-tooltip>
      </div>

      <a-spin :spinning="folderLoading">
        <a-empty v-if="!folders.length" :image="simpleImage" description="暂无分类">
          <a-button type="primary" size="small" @click="openCreator">创建分类</a-button>
        </a-empty>
        <div v-else class="folder-list">
          <button
            v-for="item in folders"
            :key="item.id"
            type="button"
            :class="['folder-item', { active: state.selectedMenuItem === item.id }]"
            @click="menuItemChanged(item.id)"
          >
            <span class="folder-label">
              <PictureOutlined />
              <span class="folder-name">{{ item.name }}</span>
            </span>
            <span class="folder-actions" @click.stop>
              <a-tooltip v-if="item.support_edit" title="重命名">
                <EditOutlined @click="preEdit(item)" />
              </a-tooltip>
              <a-popconfirm
                v-if="item.support_delete"
                title="确定删除这个空分类吗？"
                ok-text="删除"
                cancel-text="取消"
                @confirm="deleteAssetFolder(item.id)"
              >
                <a-tooltip title="删除分类"><DeleteOutlined /></a-tooltip>
              </a-popconfirm>
            </span>
          </button>
        </div>
      </a-spin>
    </aside>

    <section class="gallery-panel">
      <header class="gallery-header">
        <div>
          <div class="gallery-title-row">
            <span class="panel-title">{{ selectedFolder?.name || '图片素材' }}</span>
            <a-tag v-if="selectedFolder" color="blue">{{ pagination.totalCount }} 张</a-tag>
          </div>
          <div class="panel-subtitle">选择图片后可插入文章，支持 JPG、PNG，最大 1MB</div>
        </div>
        <a-space v-if="state.selectedMenuItem">
          <SimpleUpload
            @success:imageUrl="uploadAsset"
            :authorization="userStore.token"
            :action="serverHost + '/admin-api/files/upload'"
            label="上传新图片"
            :fileTypes="['image/jpeg', 'image/png']"
            :maxSize="1048576"
          />
          <a-button :icon="h(FileImageOutlined)" @click="visible4ExistImageModal = true">
            从文件库选择
          </a-button>
        </a-space>
      </header>

      <a-spin class="gallery-loading" :spinning="imageLoading">
        <div v-if="images.length" class="image-grid">
          <button
            v-for="image in images"
            :key="image.id"
            type="button"
            :class="['image-card', { selected: isSelected(image.id) }]"
            :aria-label="`选择图片 ${assetName(image)}`"
            @click="selectImage(image.id)"
            @dblclick="selectAndInsert(image.id)"
          >
            <span class="image-stage">
              <img :src="serverHost + image.content" :alt="assetName(image)" />
              <span class="image-overlay">
                <a-tooltip title="查看大图">
                  <span class="preview-button" @click.stop="preview(image.content)">
                    <EyeOutlined />
                  </span>
                </a-tooltip>
              </span>
              <CheckCircleFilled v-if="isSelected(image.id)" class="selected-mark" />
            </span>
            <span class="image-name" :title="assetName(image)">
              {{ assetName(image) }}
            </span>
          </button>
        </div>
        <a-empty
          v-else
          class="gallery-empty"
          :image="simpleImage"
          :description="state.selectedMenuItem ? '该分类还没有图片' : '请先创建或选择分类'"
        />
      </a-spin>

      <footer class="gallery-footer">
        <a-pagination
          v-if="pagination.totalCount > pagination.pageSize"
          v-model:current="pagination.pageNo"
          :page-size="pagination.pageSize"
          :total="pagination.totalCount"
          :show-size-changer="false"
          size="small"
          @change="getImages"
        />
        <span v-else></span>
        <a-space>
          <a-popconfirm
            title="确定删除选中的图片吗？"
            ok-text="删除"
            cancel-text="取消"
            @confirm="deleteAsset"
          >
            <a-button danger :disabled="!state.selectedImgIndex">删除</a-button>
          </a-popconfirm>
          <a-button type="primary" :disabled="!state.selectedImgIndex" @click="insert">
            插入图片
          </a-button>
        </a-space>
      </footer>
    </section>

    <a-modal v-model:open="visible" :title="modalLabel" @ok="handleOk" @cancel="cancel">
      <a-form layout="vertical">
        <a-form-item label="分类名称" required>
          <a-input
            v-model:value="assetFolder.name"
            :maxlength="30"
            show-count
            placeholder="例如：文章配图"
            @press-enter="handleOk"
          />
        </a-form-item>
      </a-form>
    </a-modal>
    <PreviewImg v-model="visible4PreviewImgModal" :image-url="previewImgUrl" />
    <ImageList v-model="visible4ExistImageModal" @insertImg="uploadAsset" />
  </div>
</template>

<script setup lang="ts">
import { computed, defineEmits, h, reactive, ref, watch } from 'vue'
import {
  PictureOutlined,
  FolderAddOutlined,
  EditOutlined,
  DeleteOutlined,
  FileImageOutlined,
  EyeOutlined,
  CheckCircleFilled
} from '@ant-design/icons-vue'
import {
  AddAsset,
  AddAssetFolder,
  type AddAssetFolderRequest,
  type AssetFolderVO,
  type AssetRequest,
  type AssetVO,
  DeleteAsset,
  DeleteAssetFolder,
  EditAssetFolderName,
  GetAssetFolderList,
  GetAssetList
} from '@/interfaces/Asset'
import { Empty, message } from 'ant-design-vue'
import originalAxios from 'axios'
import { useUserStore } from '@/stores/user'
import SimpleUpload from '@/components/upload/SimpleUpload.vue'
import ImageList from '@/components/file/ImageList.vue'
import PreviewImg from '@/components/image/PreviewImg.vue'
import { toErrorMessage } from '@/utils/error'

const state = reactive({
  selectedMenuItem: '',
  selectedImgIndex: ''
})

const folders = ref<AssetFolderVO[]>([])
const folderLoading = ref(false)
const imageLoading = ref(false)
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const selectedFolder = computed(() =>
  folders.value.find((folder) => folder.id === state.selectedMenuItem)
)

const getFolders = async () => {
  folderLoading.value = true
  try {
    const response = await GetAssetFolderList('image', 'post-editor')
    folders.value = response.data.data?.list || []
    if (!folders.value.some((folder) => folder.id === state.selectedMenuItem)) {
      state.selectedMenuItem = folders.value[0]?.id || ''
    }
  } catch (error) {
    console.log(error)
    message.error('获取图片分类失败')
  } finally {
    folderLoading.value = false
  }
}

getFolders()

const images = ref<AssetVO[]>([])
const pagination = reactive({
  pageNo: 1,
  pageSize: 12,
  totalCount: 0
})

const selectImage = (id: string) => {
  state.selectedImgIndex = state.selectedImgIndex == id ? '' : id
}

const isSelected = (id: string) => {
  return state.selectedImgIndex == id
}

const emit = defineEmits(['insertImg'])

const insert = () => {
  const image = images.value.find((item) => item.id == state.selectedImgIndex)
  if (!image) return
  emit('insertImg', markdownImage(assetName(image), image.content))
}

const selectAndInsert = (id: string) => {
  state.selectedImgIndex = id
  insert()
}

const fileName = (path: string) => path.split('/').pop() || '图片'
const assetName = (asset: AssetVO) => asset.title || fileName(asset.content)
const markdownImage = (alt: string, url: string) => {
  const escapedAlt = alt.replace(/\\/g, '\\\\').replace(/\]/g, '\\]')
  return `![${escapedAlt}](${url})`
}

const menuItemChanged = (id: string) => {
  pagination.pageNo = 1
  state.selectedImgIndex = ''
  state.selectedMenuItem = id
}

const visible = ref(false)
const modalLabel = ref('')

const openCreator = () => {
  resetAssetFolder()
  modalLabel.value = '新增分类'
  visible.value = true
}

const assetFolder = reactive<AddAssetFolderRequest>({
  name: '',
  asset_type: 'image',
  type: 'post-editor',
  support_delete: true,
  support_edit: true,
  support_add: true
} as AddAssetFolderRequest)

const handleOk = () => {
  if (modalLabel.value === '新增分类') {
    addAssetFolder()
  } else {
    editAssetFolder()
  }
}

const addAssetFolder = async () => {
  if (!assetFolder.name) {
    message.error('请输入分类名称')
  } else {
    try {
      const response: any = await AddAssetFolder(assetFolder)
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      message.success('添加成功')
      await getFolders()
      assetFolder.name = ''
      visible.value = false
    } catch (error) {
      console.log(error)
      if (originalAxios.isAxiosError(error)) {
        // 这是一个由 axios 抛出的错误
        if (error.response) {
          if (error.response.status === 409) {
            message.error('分类名称重复')
            return
          }
        } else if (error.request) {
          // 请求已发出，但没有收到响应
          console.log('No response received:', error.request)
        } else {
          // 在设置请求时触发了一个错误
          console.log('Error Message:', error.message)
        }
      }
      message.error('添加失败')
    }
  }
}

const preEdit = (assetFolderVO: AssetFolderVO) => {
  modalLabel.value = '编辑分类'
  visible.value = true
  assetFolder.id = assetFolderVO.id
  assetFolder.name = assetFolderVO.name
}

const resetAssetFolder = () => {
  assetFolder.id = ''
  assetFolder.name = ''
  assetFolder.asset_type = 'image'
  assetFolder.type = 'post-editor'
  assetFolder.support_delete = true
  assetFolder.support_edit = true
  assetFolder.support_add = true
}

const editAssetFolder = async () => {
  if (!assetFolder.name) {
    message.error('请输入分类名称')
  } else {
    try {
      const response: any = await EditAssetFolderName(assetFolder.id, assetFolder.name)
      if (response.data.code !== 0) {
        message.error(response.data.message)
        return
      }
      message.success('编辑成功')
      await getFolders()
      resetAssetFolder()
      visible.value = false
    } catch (error) {
      console.log(error)
      if (originalAxios.isAxiosError(error)) {
        // 这是一个由 axios 抛出的错误
        if (error.response) {
          if (error.response.status === 409) {
            message.error('分类名称重复')
            return
          }
        } else if (error.request) {
          // 请求已发出，但没有收到响应
          console.log('No response received:', error.request)
        } else {
          // 在设置请求时触发了一个错误
          console.log('Error Message:', error.message)
        }
      }
      message.error('编辑失败')
    }
  }
}

const cancel = () => {
  resetAssetFolder()
}

const deleteAssetFolder = async (id: string) => {
  try {
    const response: any = await DeleteAssetFolder(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    await getFolders()
  } catch (error) {
    console.log(error)
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 404) {
          message.error('分类不存在')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    }
    message.error('删除失败')
  }
}

const getImages = async () => {
  if (!state.selectedMenuItem) {
    images.value = []
    pagination.totalCount = 0
    return
  }
  state.selectedImgIndex = ''
  imageLoading.value = true
  try {
    const response: any = await GetAssetList(
      state.selectedMenuItem,
      pagination.pageNo,
      pagination.pageSize
    )
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    images.value = response.data.data?.list || []
    pagination.totalCount = response.data.data?.totalCount || 0
  } catch (error) {
    console.log(error)
    message.error('获取图片列表失败')
  } finally {
    imageLoading.value = false
  }
}

watch(
  () => state.selectedMenuItem,
  () => {
    getImages()
  }
)

const serverHost = import.meta.env.VITE_API_HOST
const userStore = useUserStore()

const uploadAsset = async (
  fileId: string,
  fileUrl: string,
  originalFileName = '',
  insertAfterAdd = false
) => {
  try {
    const response: any = await AddAsset(state.selectedMenuItem, {
      title: originalFileName,
      content: fileUrl,
      asset_type: 'image',
      type: 'post-editor',
      metadata: {
        file_id: fileId
      }
    } as AssetRequest)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success(originalFileName ? '已加入素材库' : '上传成功')
    if (insertAfterAdd) {
      emit('insertImg', markdownImage(originalFileName || fileName(fileUrl), fileUrl))
    }
    pagination.pageNo = Math.max(1, Math.ceil((pagination.totalCount + 1) / pagination.pageSize))
    await getImages()
  } catch (error) {
    console.log(error)
    message.error(toErrorMessage(error, '加入素材库失败'))
  }
}

const visible4ExistImageModal = ref(false)

const visible4PreviewImgModal = ref(false)
const previewImgUrl = ref('')
const preview = (imgUrl: string) => {
  previewImgUrl.value = serverHost + imgUrl
  visible4PreviewImgModal.value = true
}

const deleteAsset = async () => {
  try {
    if (state.selectedMenuItem === '' || state.selectedImgIndex === '') {
      message.error('请选择图片')
      return
    }
    const response: any = await DeleteAsset(state.selectedMenuItem, state.selectedImgIndex)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('删除成功')
    state.selectedImgIndex = ''
    if (images.value.length === 1 && pagination.pageNo > 1) {
      pagination.pageNo--
    }
    await getImages()
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}
</script>

<style scoped>
.image-assets {
  display: grid;
  grid-template-columns: 220px minmax(0, 1fr);
  height: 560px;
  overflow: hidden;
  border: 1px solid var(--app-border, #f0f0f0);
  border-radius: 10px;
  background: var(--app-surface, #fff);
}

.folder-panel {
  min-width: 0;
  padding: 18px 14px;
  overflow-y: auto;
  border-right: 1px solid var(--app-border, #f0f0f0);
  background: var(--app-surface-muted, #fafafa);
}

.panel-header,
.gallery-header,
.gallery-footer,
.gallery-title-row,
.folder-label,
.folder-actions {
  display: flex;
  align-items: center;
}

.panel-header,
.gallery-header,
.gallery-footer {
  justify-content: space-between;
}

.panel-header {
  margin-bottom: 14px;
  padding: 0 4px;
}

.panel-title {
  color: var(--app-text, rgba(0, 0, 0, 0.88));
  font-size: 15px;
  font-weight: 600;
}

.panel-subtitle {
  margin-top: 3px;
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.45));
  font-size: 12px;
}

.folder-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.folder-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  min-height: 42px;
  padding: 0 10px;
  border: 1px solid transparent;
  border-radius: 7px;
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.65));
  background: transparent;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}

.folder-item:hover {
  background: var(--app-surface, #fff);
}

.folder-item.active {
  border-color: #91caff;
  color: #1677ff;
  background: #e6f4ff;
}

.folder-label {
  min-width: 0;
  gap: 8px;
}

.folder-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.folder-actions {
  display: none;
  flex-shrink: 0;
  gap: 9px;
}

.folder-item:hover .folder-actions,
.folder-item.active .folder-actions {
  display: flex;
}

.gallery-panel {
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  min-width: 0;
  padding: 18px;
}

.gallery-header {
  gap: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid var(--app-border, #f0f0f0);
}

.gallery-title-row {
  gap: 8px;
}

.gallery-loading {
  min-height: 0;
  overflow-y: auto;
}

.gallery-loading :deep(.ant-spin-container) {
  min-height: 100%;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  padding: 16px 2px;
}

.image-card {
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
    box-shadow 0.2s ease,
    transform 0.2s ease;
}

.image-card:hover {
  border-color: #69b1ff;
  box-shadow: 0 5px 14px rgba(22, 119, 255, 0.12);
  transform: translateY(-1px);
}

.image-card.selected {
  border-color: #1677ff;
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.16);
}

.image-stage {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  aspect-ratio: 1;
  overflow: hidden;
  background-image:
    linear-gradient(45deg, #f0f0f0 25%, transparent 25%),
    linear-gradient(-45deg, #f0f0f0 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, #f0f0f0 75%),
    linear-gradient(-45deg, transparent 75%, #f0f0f0 75%);
  background-position:
    0 0,
    0 6px,
    6px -6px,
    -6px 0;
  background-size: 12px 12px;
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

.image-card:hover .image-overlay {
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
}

.selected-mark {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 1;
  color: #1677ff;
  font-size: 22px;
  filter: drop-shadow(0 1px 2px rgba(255, 255, 255, 0.8));
}

.image-name {
  display: block;
  padding: 8px 9px;
  overflow: hidden;
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.65));
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gallery-empty {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 350px;
}

.gallery-footer {
  min-height: 48px;
  padding-top: 12px;
  border-top: 1px solid var(--app-border, #f0f0f0);
}

@media (max-width: 760px) {
  .image-assets {
    grid-template-columns: 180px minmax(0, 1fr);
  }

  .gallery-header {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
