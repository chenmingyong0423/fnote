<template>
  <div class="text-templates">
    <aside class="folder-panel">
      <div class="panel-header">
        <div>
          <div class="panel-title">模板分类</div>
          <div class="panel-subtitle">按用途查找文字模板</div>
        </div>
      </div>

      <a-spin :spinning="folderLoading">
        <a-empty v-if="!folders.length" :image="simpleImage" description="暂无模板分类" />
        <div v-else class="folder-list">
          <button
            v-for="folder in folders"
            :key="folder.id"
            type="button"
            :class="['folder-item', { active: folder.id === selectedFolderId }]"
            @click="selectFolder(folder.id)"
          >
            <FileTextOutlined />
            <span class="folder-name">{{ folder.name }}</span>
          </button>
        </div>
      </a-spin>
    </aside>

    <section class="template-panel">
      <header class="template-header">
        <div>
          <div class="template-title-row">
            <span class="panel-title">{{ selectedFolder?.name || '文字模板' }}</span>
            <a-tag v-if="selectedFolder" color="blue">{{ pagination.total }} 条</a-tag>
          </div>
          <div class="panel-subtitle">选择后插入到编辑器光标位置，双击可直接插入</div>
        </div>
      </header>

      <a-spin class="template-loading" :spinning="assetLoading">
        <div v-if="assets.length" class="template-list">
          <button
            v-for="asset in assets"
            :key="asset.id"
            type="button"
            :class="['template-card', { selected: asset.id === selectedAssetId }]"
            @click="selectedAssetId = asset.id"
            @dblclick="selectAndInsert(asset)"
          >
            <span class="template-card-main">
              <span class="template-name">{{ asset.title }}</span>
              <span class="template-description">
                {{ asset.description || contentSummary(asset.content) }}
              </span>
            </span>
            <a-tooltip title="Markdown 预览">
              <span class="preview-button" @click.stop="openPreview(asset)">
                <EyeOutlined />
              </span>
            </a-tooltip>
            <CheckCircleFilled v-if="asset.id === selectedAssetId" class="selected-mark" />
          </button>
        </div>
        <a-empty
          v-else
          class="template-empty"
          :image="simpleImage"
          :description="selectedFolderId ? '该分类暂无文字模板' : '暂无可用的模板分类'"
        />
      </a-spin>

      <footer class="template-footer">
        <a-pagination
          v-if="pagination.total > pagination.pageSize"
          v-model:current="pagination.current"
          :page-size="pagination.pageSize"
          :total="pagination.total"
          :show-size-changer="false"
          size="small"
          @change="loadAssets"
        />
        <span v-else></span>
        <a-button type="primary" :disabled="!selectedAsset" @click="insertSelected">
          插入模板
        </a-button>
      </footer>
    </section>

    <a-modal
      v-model:open="previewVisible"
      :title="previewAsset?.title || '模板预览'"
      width="820px"
      :footer="null"
      :destroy-on-close="true"
    >
      <v-md-editor :model-value="previewContent" mode="preview" height="520px" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { CheckCircleFilled, EyeOutlined, FileTextOutlined } from '@ant-design/icons-vue'
import { Empty, message } from 'ant-design-vue'
import {
  GetAssetFolderList,
  GetAssetList,
  type AssetFolderVO,
  type AssetVO
} from '@/interfaces/Asset'
import { toErrorMessage } from '@/utils/error'

const emit = defineEmits<{
  insertTemplate: [content: string]
}>()

const folders = ref<AssetFolderVO[]>([])
const assets = ref<AssetVO[]>([])
const selectedFolderId = ref('')
const selectedAssetId = ref('')
const folderLoading = ref(false)
const assetLoading = ref(false)
const previewVisible = ref(false)
const previewAsset = ref<AssetVO>()
const simpleImage = Empty.PRESENTED_IMAGE_SIMPLE
const apiHost = import.meta.env.VITE_API_HOST?.replace(/\/$/, '') || ''
const pagination = reactive({ current: 1, pageSize: 8, total: 0 })

const selectedFolder = computed(() =>
  folders.value.find((folder) => folder.id === selectedFolderId.value)
)
const selectedAsset = computed(() =>
  assets.value.find((asset) => asset.id === selectedAssetId.value)
)
const previewContent = computed(() => {
  const content = previewAsset.value?.content || ''
  if (!apiHost) return content
  return content
    .replace(/(!\[[^\]]*\]\(\s*)(\/static\/)/g, `$1${apiHost}$2`)
    .replace(/(<img\b[^>]*\bsrc\s*=\s*["'])(\/static\/)/gi, `$1${apiHost}$2`)
})

const loadFolders = async () => {
  folderLoading.value = true
  try {
    const response = await GetAssetFolderList('text', 'post-editor')
    folders.value = response.data.data?.list || []
    selectedFolderId.value = folders.value[0]?.id || ''
  } catch (error) {
    message.error(toErrorMessage(error, '文字模板分类加载失败'))
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
  selectedAssetId.value = ''
  try {
    const response = await GetAssetList(
      selectedFolderId.value,
      pagination.current,
      pagination.pageSize
    )
    assets.value = response.data.data?.list || []
    pagination.total = response.data.data?.totalCount || 0
  } catch (error) {
    message.error(toErrorMessage(error, '文字模板加载失败'))
  } finally {
    assetLoading.value = false
  }
}

const selectFolder = (id: string) => {
  if (id === selectedFolderId.value) return
  pagination.current = 1
  selectedFolderId.value = id
}

const contentSummary = (content: string) =>
  content
    .replace(/[#>*_`\[\]()~-]/g, '')
    .replace(/\s+/g, ' ')
    .trim() || '暂无内容摘要'

const openPreview = (asset: AssetVO) => {
  previewAsset.value = asset
  previewVisible.value = true
}

const insertSelected = () => {
  if (selectedAsset.value) emit('insertTemplate', selectedAsset.value.content)
}

const selectAndInsert = (asset: AssetVO) => {
  selectedAssetId.value = asset.id
  emit('insertTemplate', asset.content)
}

watch(selectedFolderId, loadAssets)
loadFolders()
</script>

<style scoped>
.text-templates {
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

.folder-list,
.template-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.folder-item {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  min-height: 42px;
  padding: 0 10px;
  border: 1px solid transparent;
  border-radius: 7px;
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.65));
  background: transparent;
  cursor: pointer;
  text-align: left;
}

.folder-item:hover {
  background: var(--app-surface, #fff);
}

.folder-item.active {
  border-color: #91caff;
  color: #1677ff;
  background: #e6f4ff;
}

.folder-name,
.template-name,
.template-description {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.template-panel {
  display: grid;
  grid-template-rows: auto minmax(0, 1fr) auto;
  min-width: 0;
  padding: 18px;
}

.template-header {
  padding-bottom: 16px;
  border-bottom: 1px solid var(--app-border, #f0f0f0);
}

.template-title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.template-loading {
  min-height: 0;
  overflow-y: auto;
}

.template-list {
  padding: 16px 2px;
}

.template-card {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  min-height: 66px;
  padding: 11px 44px 11px 14px;
  border: 1px solid var(--app-border, #d9d9d9);
  border-radius: 8px;
  color: inherit;
  background: var(--app-surface, #fff);
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
}

.template-card:hover {
  border-color: #69b1ff;
  box-shadow: 0 4px 12px rgba(22, 119, 255, 0.1);
}

.template-card.selected {
  border-color: #1677ff;
  background: #f0f7ff;
  box-shadow: 0 0 0 2px rgba(22, 119, 255, 0.12);
}

.template-card-main {
  display: flex;
  flex: 1;
  flex-direction: column;
  min-width: 0;
  gap: 5px;
}

.template-name {
  color: var(--app-text, rgba(0, 0, 0, 0.88));
  font-weight: 600;
}

.template-description {
  color: var(--app-text-secondary, rgba(0, 0, 0, 0.45));
  font-size: 12px;
}

.preview-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border-radius: 6px;
  color: #8c8c8c;
}

.preview-button:hover {
  color: #1677ff;
  background: #e6f4ff;
}

.selected-mark {
  position: absolute;
  top: 8px;
  right: 8px;
  color: #1677ff;
  font-size: 18px;
}

.template-empty {
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 350px;
}

.template-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 48px;
  padding-top: 12px;
  border-top: 1px solid var(--app-border, #f0f0f0);
}

@media (max-width: 760px) {
  .text-templates {
    grid-template-columns: 180px minmax(0, 1fr);
  }
}
</style>
