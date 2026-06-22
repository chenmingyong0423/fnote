<template>
  <a-card title="文字素材">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="folderLoading"
          @click="loadFolders"
        />
      </a-tooltip>
    </template>

    <a-alert
      class="asset-help"
      type="info"
      show-icon
      message="按分类管理可重复使用的文字模板，点击“复制”即可将正文放入剪贴板。"
    />

    <div class="asset-layout">
      <aside class="folder-panel">
        <div class="panel-header">
          <span class="panel-title">素材分类</span>
          <a-button type="link" size="small" @click="openFolderCreator">新增</a-button>
        </div>

        <a-spin :spinning="folderLoading">
          <a-empty v-if="!folders.length" :image="simpleImage" description="暂无分类">
            <a-button type="primary" size="small" @click="openFolderCreator">创建分类</a-button>
          </a-empty>
          <div v-else class="folder-list">
            <button
              v-for="folder in folders"
              :key="folder.id"
              type="button"
              :class="['folder-item', { active: folder.id === selectedFolderId }]"
              @click="selectFolder(folder.id)"
            >
              <span class="folder-name">{{ folder.name }}</span>
              <span class="folder-actions" @click.stop>
                <EditOutlined v-if="folder.support_edit" @click="openFolderEditor(folder)" />
                <a-popconfirm
                  v-if="folder.support_delete"
                  title="确定删除这个空分类吗？"
                  @confirm="removeFolder(folder)"
                >
                  <DeleteOutlined />
                </a-popconfirm>
              </span>
            </button>
          </div>
        </a-spin>
      </aside>

      <section class="template-panel">
        <div class="panel-header">
          <div>
            <span class="panel-title">{{ selectedFolder?.name || '文字模板' }}</span>
            <span v-if="selectedFolder" class="template-count">共 {{ pagination.total }} 条</span>
          </div>
          <a-button type="primary" :disabled="!selectedFolderId" @click="openAssetCreator">
            新增文字模板
          </a-button>
        </div>

        <a-table
          row-key="id"
          :columns="columns"
          :data-source="assets"
          :loading="assetLoading"
          :pagination="pagination"
          :scroll="{ x: 900 }"
          @change="onTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'content'">
              <a-typography-paragraph
                class="content-preview"
                :content="record.content"
                :ellipsis="{ rows: 2, tooltip: record.content }"
              />
            </template>
            <template v-else-if="column.key === 'description'">
              <span class="description-text">{{ record.description || '—' }}</span>
            </template>
            <template v-else-if="column.key === 'operation'">
              <a-space>
                <a @click="openAssetViewer(record)">查看</a>
                <a @click="copyContent(record)">复制</a>
                <a @click="openAssetEditor(record)">编辑</a>
                <a-popconfirm title="确定删除这个文字模板吗？" @confirm="removeAsset(record)">
                  <a class="danger-link">删除</a>
                </a-popconfirm>
              </a-space>
            </template>
          </template>
          <template #emptyText>
            <a-empty
              :description="selectedFolderId ? '该分类暂无文字模板' : '请先选择或创建分类'"
            />
          </template>
        </a-table>
      </section>
    </div>
  </a-card>

  <a-modal
    v-model:open="folderEditorOpen"
    :title="editingFolder ? '编辑分类' : '新增分类'"
    ok-text="保存"
    cancel-text="取消"
    :confirm-loading="folderSaving"
    @ok="saveFolder"
  >
    <a-form ref="folderFormRef" :model="folderForm" layout="vertical">
      <a-form-item
        label="分类名称"
        name="name"
        :rules="[{ required: true, whitespace: true, message: '请输入分类名称' }]"
      >
        <a-input
          v-model:value="folderForm.name"
          :maxlength="30"
          show-count
          @press-enter="saveFolder"
        />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal
    v-model:open="assetEditorOpen"
    :title="editingAsset ? '编辑文字模板' : '新增文字模板'"
    width="720px"
    ok-text="保存"
    cancel-text="取消"
    :confirm-loading="assetSaving"
    @ok="saveAsset"
  >
    <a-form ref="assetFormRef" :model="assetForm" layout="vertical">
      <a-form-item
        label="模板标题"
        name="title"
        :rules="[{ required: true, whitespace: true, message: '请输入模板标题' }]"
      >
        <a-input v-model:value="assetForm.title" :maxlength="80" show-count />
      </a-form-item>
      <a-form-item
        label="模板正文"
        name="content"
        :rules="[{ required: true, whitespace: true, message: '请输入模板正文' }]"
      >
        <a-textarea
          v-model:value="assetForm.content"
          :auto-size="{ minRows: 8, maxRows: 18 }"
          show-count
        />
      </a-form-item>
      <a-form-item label="备注" name="description">
        <a-textarea
          v-model:value="assetForm.description"
          :maxlength="200"
          :auto-size="{ minRows: 2, maxRows: 4 }"
          show-count
        />
      </a-form-item>
    </a-form>
  </a-modal>

  <a-modal
    v-model:open="assetViewerOpen"
    title="查看文字模板"
    width="900px"
    :footer="null"
    :destroy-on-close="true"
  >
    <a-descriptions class="viewer-meta" :column="1" size="small" bordered>
      <a-descriptions-item label="模板标题">{{ viewingAsset?.title }}</a-descriptions-item>
      <a-descriptions-item label="所属分类">{{ selectedFolder?.name }}</a-descriptions-item>
      <a-descriptions-item label="备注">{{ viewingAsset?.description || '—' }}</a-descriptions-item>
    </a-descriptions>
    <v-md-editor :model-value="previewContent" mode="preview" height="560px" />
  </a-modal>
</template>

<script lang="ts" setup>
import { computed, h, reactive, ref } from 'vue'
import type { FormInstance, TableColumnsType, TablePaginationConfig } from 'ant-design-vue'
import { Empty, message } from 'ant-design-vue'
import { DeleteOutlined, EditOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import {
  AddAsset,
  AddAssetFolder,
  DeleteAsset,
  DeleteAssetFolder,
  EditAssetFolderName,
  GetAssetFolderList,
  GetAssetList,
  UpdateAsset,
  type AssetFolderVO,
  type AssetRequest,
  type AssetVO
} from '@/interfaces/Asset'
import { toErrorMessage } from '@/utils/error'

document.title = '文字素材 - 后台管理'

const ASSET_TYPE = 'text'
const USE_TYPE = 'post-editor'
const apiHost = import.meta.env.VITE_API_HOST?.replace(/\/$/, '') || ''
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const columns: TableColumnsType = [
  { title: '标题', dataIndex: 'title', key: 'title', width: 180 },
  { title: '正文', dataIndex: 'content', key: 'content', width: 400 },
  { title: '备注', dataIndex: 'description', key: 'description', width: 180 },
  { title: '操作', key: 'operation', fixed: 'right', width: 190 }
]

const folders = ref<AssetFolderVO[]>([])
const assets = ref<AssetVO[]>([])
const selectedFolderId = ref('')
const folderLoading = ref(false)
const assetLoading = ref(false)
const folderEditorOpen = ref(false)
const assetEditorOpen = ref(false)
const assetViewerOpen = ref(false)
const folderSaving = ref(false)
const assetSaving = ref(false)
const editingFolder = ref<AssetFolderVO>()
const editingAsset = ref<AssetVO>()
const viewingAsset = ref<AssetVO>()
const folderFormRef = ref<FormInstance>()
const assetFormRef = ref<FormInstance>()
const folderForm = reactive({ name: '' })
const assetForm = reactive({ title: '', content: '', description: '' })
const pagination = reactive<TablePaginationConfig>({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true
})
const selectedFolder = computed(() =>
  folders.value.find((folder) => folder.id === selectedFolderId.value)
)
const previewContent = computed(() => {
  const content = viewingAsset.value?.content || ''
  if (!apiHost) return content

  return content
    .replace(/(!\[[^\]]*\]\(\s*)(\/static\/)/g, `$1${apiHost}$2`)
    .replace(/(<img\b[^>]*\bsrc\s*=\s*["'])(\/static\/)/gi, `$1${apiHost}$2`)
})

const loadFolders = async () => {
  folderLoading.value = true
  try {
    const response = await GetAssetFolderList(ASSET_TYPE, USE_TYPE)
    folders.value = response.data.data?.list || []
    if (!folders.value.some((folder) => folder.id === selectedFolderId.value)) {
      selectedFolderId.value = folders.value[0]?.id || ''
      pagination.current = 1
    }
    await loadAssets()
  } catch (error) {
    folders.value = []
    selectedFolderId.value = ''
    assets.value = []
    pagination.total = 0
    message.error(toErrorMessage(error, '文字素材分类加载失败'))
  } finally {
    folderLoading.value = false
  }
}

const loadAssets = async () => {
  if (!selectedFolderId.value) {
    assets.value = []
    pagination.total = 0
    return
  }
  assetLoading.value = true
  try {
    const response = await GetAssetList(
      selectedFolderId.value,
      pagination.current || 1,
      pagination.pageSize || 10
    )
    assets.value = response.data.data?.list || []
    pagination.total = response.data.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '文字模板加载失败'))
  } finally {
    assetLoading.value = false
  }
}

const selectFolder = async (id: string) => {
  if (id === selectedFolderId.value) return
  selectedFolderId.value = id
  pagination.current = 1
  await loadAssets()
}

const openFolderCreator = () => {
  editingFolder.value = undefined
  folderForm.name = ''
  folderFormRef.value?.clearValidate()
  folderEditorOpen.value = true
}

const openFolderEditor = (folder: AssetFolderVO) => {
  editingFolder.value = folder
  folderForm.name = folder.name
  folderFormRef.value?.clearValidate()
  folderEditorOpen.value = true
}

const saveFolder = async () => {
  try {
    await folderFormRef.value?.validate()
    folderSaving.value = true
    if (editingFolder.value) {
      await EditAssetFolderName(editingFolder.value.id, folderForm.name.trim())
    } else {
      await AddAssetFolder({
        id: '',
        name: folderForm.name.trim(),
        asset_type: ASSET_TYPE,
        type: USE_TYPE,
        support_delete: true,
        support_edit: true,
        support_add: true
      })
    }
    folderEditorOpen.value = false
    message.success('分类保存成功')
    await loadFolders()
  } catch (error) {
    if ((error as { errorFields?: unknown })?.errorFields) return
    message.error(toErrorMessage(error, '分类保存失败'))
  } finally {
    folderSaving.value = false
  }
}

const removeFolder = async (folder: AssetFolderVO) => {
  try {
    await DeleteAssetFolder(folder.id)
    message.success('分类删除成功')
    await loadFolders()
  } catch (error) {
    message.error(toErrorMessage(error, '分类删除失败，请先清空分类中的模板'))
  }
}

const openAssetCreator = () => {
  editingAsset.value = undefined
  Object.assign(assetForm, { title: '', content: '', description: '' })
  assetFormRef.value?.clearValidate()
  assetEditorOpen.value = true
}

const openAssetEditor = (asset: AssetVO) => {
  editingAsset.value = asset
  Object.assign(assetForm, {
    title: asset.title,
    content: asset.content,
    description: asset.description
  })
  assetFormRef.value?.clearValidate()
  assetEditorOpen.value = true
}

const openAssetViewer = (asset: AssetVO) => {
  viewingAsset.value = asset
  assetViewerOpen.value = true
}

const assetRequest = (): AssetRequest => ({
  id: editingAsset.value?.id || '',
  title: assetForm.title.trim(),
  content: assetForm.content,
  description: assetForm.description.trim(),
  asset_type: ASSET_TYPE,
  type: USE_TYPE,
  metadata: {}
})

const saveAsset = async () => {
  if (!selectedFolderId.value) return
  try {
    await assetFormRef.value?.validate()
    assetSaving.value = true
    if (editingAsset.value) {
      await UpdateAsset(selectedFolderId.value, editingAsset.value.id, assetRequest())
    } else {
      await AddAsset(selectedFolderId.value, assetRequest())
    }
    assetEditorOpen.value = false
    message.success('文字模板保存成功')
    await loadAssets()
  } catch (error) {
    if ((error as { errorFields?: unknown })?.errorFields) return
    message.error(toErrorMessage(error, '文字模板保存失败'))
  } finally {
    assetSaving.value = false
  }
}

const removeAsset = async (asset: AssetVO) => {
  try {
    await DeleteAsset(selectedFolderId.value, asset.id)
    if (assets.value.length === 1 && (pagination.current || 1) > 1)
      pagination.current = (pagination.current || 1) - 1
    message.success('文字模板删除成功')
    await loadAssets()
  } catch (error) {
    message.error(toErrorMessage(error, '文字模板删除失败'))
  }
}

const copyContent = async (asset: AssetVO) => {
  try {
    await navigator.clipboard.writeText(asset.content)
    message.success('模板正文已复制')
  } catch {
    message.error('复制失败，请手动复制')
  }
}

const onTableChange = async (page: TablePaginationConfig) => {
  pagination.current = page.current || 1
  pagination.pageSize = page.pageSize || 10
  await loadAssets()
}

loadFolders()
</script>

<style scoped>
.asset-help {
  margin-bottom: 16px;
}
.asset-layout {
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr);
  min-height: 520px;
  border: 1px solid var(--app-border, #f0f0f0);
  border-radius: 8px;
  overflow: hidden;
}
.folder-panel {
  padding: 16px;
  border-right: 1px solid var(--app-border, #f0f0f0);
  background: var(--app-surface-muted, #fafafa);
}
.template-panel {
  min-width: 0;
  padding: 16px;
}
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 32px;
  margin-bottom: 16px;
}
.panel-title {
  font-size: 16px;
  font-weight: 600;
}
.template-count {
  margin-left: 10px;
  color: #8c8c8c;
  font-size: 13px;
}
.folder-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.folder-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 9px 10px;
  border: 0;
  border-radius: 6px;
  color: inherit;
  background: transparent;
  cursor: pointer;
  text-align: left;
}
.folder-item:hover,
.folder-item.active {
  color: #1677ff;
  background: #e6f4ff;
}
.folder-name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.folder-actions {
  display: none;
  flex-shrink: 0;
  gap: 10px;
}
.folder-item:hover .folder-actions,
.folder-item.active .folder-actions {
  display: flex;
}
.content-preview {
  margin-bottom: 0;
  white-space: pre-wrap;
}
.description-text {
  color: #595959;
}
.danger-link {
  color: #ff4d4f;
}
.viewer-meta {
  margin-bottom: 16px;
}
@media (max-width: 900px) {
  .asset-layout {
    grid-template-columns: 1fr;
  }
  .folder-panel {
    border-right: 0;
    border-bottom: 1px solid var(--app-border, #f0f0f0);
  }
}
</style>
