<template>
  <a-card title="消息模板">
    <template #extra>
      <a-tooltip title="刷新数据">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="loading"
          @click="loadTemplates"
        />
      </a-tooltip>
    </template>

    <a-alert
      class="template-help"
      type="info"
      show-icon
      message="模板变量使用 {{.VariableName}} 语法，每个模板只能使用列出的变量。"
    />

    <a-table
      row-key="id"
      :columns="columns"
      :data-source="templates"
      :loading="loading"
      :pagination="false"
      :scroll="{ x: 1100 }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'recipient_type'">
          <a-tag :color="record.recipient_type === 0 ? 'blue' : 'green'">
            {{ record.recipient_type === 0 ? '站长' : '用户' }}
          </a-tag>
        </template>

        <template v-else-if="column.key === 'variables'">
          <a-space v-if="record.variables.length" wrap size="small">
            <a-tag v-for="variable in record.variables" :key="variable">
              {{ variablePlaceholder(variable) }}
            </a-tag>
          </a-space>
          <span v-else class="muted">无</span>
        </template>

        <template v-else-if="column.key === 'content'">
          <a-typography-paragraph
            class="content-preview"
            :content="record.content"
            :ellipsis="{ rows: 2, tooltip: record.content }"
          />
        </template>

        <template v-else-if="column.key === 'active'">
          <a-switch
            :checked="record.active"
            :loading="activeUpdatingId === record.id"
            @change="onActiveChange(record, $event)"
          />
        </template>

        <template v-else-if="column.key === 'updated_at'">
          {{ dayjs.unix(record.updated_at).format('YYYY-MM-DD HH:mm:ss') }}
        </template>

        <template v-else-if="column.key === 'operation'">
          <a @click="openEditor(record)">编辑</a>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal
    v-model:open="editorOpen"
    title="编辑消息模板"
    width="760px"
    ok-text="保存"
    cancel-text="取消"
    :confirm-loading="saving"
    @ok="saveTemplate"
    @cancel="closeEditor"
  >
    <a-form ref="formRef" :model="editForm" layout="vertical">
      <a-form-item label="模板标识">
        <a-input :value="editingTemplate?.name" disabled />
      </a-form-item>

      <a-form-item
        label="邮件标题"
        name="title"
        :rules="[{ required: true, whitespace: true, message: '请输入邮件标题' }]"
      >
        <a-input v-model:value="editForm.title" :maxlength="120" show-count />
      </a-form-item>

      <a-form-item
        label="模板内容"
        name="content"
        :rules="[{ required: true, whitespace: true, message: '请输入模板内容' }]"
      >
        <a-textarea
          v-model:value="editForm.content"
          :auto-size="{ minRows: 6, maxRows: 14 }"
          placeholder="请输入模板内容"
        />
      </a-form-item>

      <div class="variable-panel">
        <span class="variable-label">可用变量</span>
        <a-space v-if="editingTemplate?.variables.length" wrap size="small">
          <a-button
            v-for="variable in editingTemplate.variables"
            :key="variable"
            size="small"
            @click="appendVariable(variable)"
          >
            {{ variablePlaceholder(variable) }}
          </a-button>
        </a-space>
        <span v-else class="muted">该模板不需要变量</span>
      </div>
    </a-form>
  </a-modal>
</template>

<script lang="ts" setup>
import { h, reactive, ref } from 'vue'
import type { FormInstance, TableColumnsType } from 'ant-design-vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  GetMessageTemplates,
  type MessageTemplate,
  UpdateMessageTemplate,
  UpdateMessageTemplateActive
} from '@/interfaces/MessageTemplate'
import { toErrorMessage } from '@/utils/error'

document.title = '消息模板 - 后台管理'

const columns: TableColumnsType = [
  { title: '模板标识', dataIndex: 'name', key: 'name', width: 190 },
  { title: '标题', dataIndex: 'title', key: 'title', width: 180 },
  { title: '内容', dataIndex: 'content', key: 'content', width: 320 },
  { title: '收件人', dataIndex: 'recipient_type', key: 'recipient_type', width: 90 },
  { title: '可用变量', dataIndex: 'variables', key: 'variables', width: 220 },
  { title: '启用', dataIndex: 'active', key: 'active', width: 80 },
  { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at', width: 170 },
  { title: '操作', key: 'operation', fixed: 'right', width: 70 }
]

const templates = ref<MessageTemplate[]>([])
const loading = ref(false)
const activeUpdatingId = ref('')
const editorOpen = ref(false)
const saving = ref(false)
const editingTemplate = ref<MessageTemplate>()
const formRef = ref<FormInstance>()
const editForm = reactive({ title: '', content: '' })

const variablePlaceholder = (variable: string) => `{{.${variable}}}`

const loadTemplates = async () => {
  loading.value = true
  try {
    const response = await GetMessageTemplates()
    templates.value = response.data.data?.list || []
  } catch (error) {
    message.error(toErrorMessage(error, '模板列表加载失败'))
  } finally {
    loading.value = false
  }
}

const openEditor = (record: MessageTemplate) => {
  editingTemplate.value = record
  editForm.title = record.title
  editForm.content = record.content
  formRef.value?.clearValidate()
  editorOpen.value = true
}

const closeEditor = () => {
  editorOpen.value = false
  editingTemplate.value = undefined
  editForm.title = ''
  editForm.content = ''
  formRef.value?.clearValidate()
}

const appendVariable = (variable: string) => {
  editForm.content += variablePlaceholder(variable)
}

const saveTemplate = async () => {
  if (!editingTemplate.value) return
  try {
    await formRef.value?.validate()
    saving.value = true
    await UpdateMessageTemplate(editingTemplate.value.id, {
      title: editForm.title.trim(),
      content: editForm.content
    })
    message.success('模板保存成功')
    closeEditor()
    await loadTemplates()
  } catch (error) {
    if ((error as { errorFields?: unknown })?.errorFields) return
    message.error(toErrorMessage(error, '模板保存失败'))
  } finally {
    saving.value = false
  }
}

const changeActive = async (record: MessageTemplate, active: boolean) => {
  activeUpdatingId.value = record.id
  try {
    await UpdateMessageTemplateActive(record.id, active)
    record.active = active
    message.success(active ? '模板已启用' : '模板已停用')
  } catch (error) {
    message.error(toErrorMessage(error, '模板状态更新失败'))
  } finally {
    activeUpdatingId.value = ''
  }
}

const onActiveChange = (record: MessageTemplate, checked: boolean | string | number) => {
  changeActive(record, Boolean(checked))
}

loadTemplates()
</script>

<style scoped>
.template-help {
  margin-bottom: 16px;
}

.content-preview {
  margin-bottom: 0;
  white-space: pre-wrap;
}

.variable-panel {
  padding: 12px;
  border: 1px solid var(--app-border, #f0f0f0);
  border-radius: 6px;
  background: var(--app-surface-muted, #fafafa);
}

.variable-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.muted {
  color: #8c8c8c;
}
</style>
