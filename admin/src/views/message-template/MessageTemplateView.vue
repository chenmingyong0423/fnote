<template>
  <a-card title="消息模板">
    <template #extra>
      <a-space>
        <a-button type="primary" @click="openCreator">新增模板</a-button>
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="loadTemplates"
          />
        </a-tooltip>
      </a-space>
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
        <template v-if="column.key === 'type'">
          <a-tooltip :title="record.type">
            <span>{{ typeLabel(record.type) }}</span>
          </a-tooltip>
        </template>

        <template v-else-if="column.key === 'recipient_type'">
          <a-tag :color="record.recipient_type === 0 ? 'blue' : 'green'">
            {{ record.recipient_type === 0 ? '站长' : '用户' }}
          </a-tag>
        </template>

        <template v-else-if="column.key === 'is_default'">
          <a-tag :color="record.is_default ? 'gold' : 'default'">
            {{ record.is_default ? '默认' : '普通' }}
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
            :disabled="record.is_default"
            :loading="activeUpdatingId === record.id"
            @change="onActiveChange(record, $event)"
          />
        </template>

        <template v-else-if="column.key === 'updated_at'">
          {{ dayjs.unix(record.updated_at).format('YYYY-MM-DD HH:mm:ss') }}
        </template>

        <template v-else-if="column.key === 'operation'">
          <a-space>
            <a @click="openEditor(record)">编辑</a>
            <a-popconfirm
              v-if="!record.is_default"
              title="确定将该模板设为默认吗？"
              @confirm="setDefault(record)"
            >
              <a>设为默认</a>
            </a-popconfirm>
            <a-popconfirm
              v-if="!record.is_default"
              title="确定删除该模板吗？"
              @confirm="deleteTemplate(record)"
            >
              <a class="danger-link">删除</a>
            </a-popconfirm>
          </a-space>
        </template>
      </template>
    </a-table>
  </a-card>

  <a-modal
    v-model:open="editorOpen"
    :title="editingTemplate ? '编辑消息模板' : '新增消息模板'"
    width="760px"
    ok-text="保存"
    cancel-text="取消"
    :confirm-loading="saving"
    @ok="saveTemplate"
    @cancel="closeEditor"
  >
    <a-form ref="formRef" :model="editForm" layout="vertical">
      <a-form-item
        label="模板类型"
        name="type"
        :rules="[{ required: true, message: '请选择模板类型' }]"
      >
        <a-select
          v-model:value="editForm.type"
          :options="typeOptions"
          :disabled="Boolean(editingTemplate)"
        />
      </a-form-item>

      <a-form-item
        label="模板名称"
        name="name"
        :rules="[{ required: true, whitespace: true, message: '请输入模板名称' }]"
      >
        <a-input v-model:value="editForm.name" :maxlength="60" show-count />
      </a-form-item>

      <a-form-item v-if="!editingTemplate" label="初始状态">
        <a-switch
          v-model:checked="editForm.active"
          checked-children="启用"
          un-checked-children="停用"
        />
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
        <a-space v-if="currentVariables.length" wrap size="small">
          <a-button
            v-for="variable in currentVariables"
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
import { computed, h, reactive, ref } from 'vue'
import type { FormInstance, TableColumnsType } from 'ant-design-vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import {
  GetMessageTemplates,
  type MessageTemplate,
  CreateMessageTemplate,
  DeleteMessageTemplate,
  SetDefaultMessageTemplate,
  UpdateMessageTemplate,
  UpdateMessageTemplateActive
} from '@/interfaces/MessageTemplate'
import { toErrorMessage } from '@/utils/error'

document.title = '消息模板 - 后台管理'

const columns: TableColumnsType = [
  { title: '模板类型', dataIndex: 'type', key: 'type', width: 190 },
  { title: '模板名称', dataIndex: 'name', key: 'name', width: 130 },
  { title: '标题', dataIndex: 'title', key: 'title', width: 180 },
  { title: '内容', dataIndex: 'content', key: 'content', width: 320 },
  { title: '收件人', dataIndex: 'recipient_type', key: 'recipient_type', width: 90 },
  { title: '可用变量', dataIndex: 'variables', key: 'variables', width: 220 },
  { title: '启用', dataIndex: 'active', key: 'active', width: 80 },
  { title: '默认模板', dataIndex: 'is_default', key: 'is_default', width: 100 },
  { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at', width: 170 },
  { title: '操作', key: 'operation', fixed: 'right', width: 210 }
]

const templates = ref<MessageTemplate[]>([])
const loading = ref(false)
const activeUpdatingId = ref('')
const editorOpen = ref(false)
const saving = ref(false)
const editingTemplate = ref<MessageTemplate>()
const formRef = ref<FormInstance>()
const editForm = reactive({ type: '', name: '', title: '', content: '', active: true })

const typeLabels: Record<string, string> = {
  comment: '新评论通知',
  'user-comment-approval': '评论审核通过',
  'user-comment-disapproval': '评论审核驳回',
  'user-comment-reply': '评论收到回复',
  friend: '新友链申请',
  'friend-approval': '友链申请通过',
  'friend-rejection': '友链申请驳回'
}

const typeLabel = (type: string) => typeLabels[type] || type

const typeOptions = computed(() => {
  const types = new Set(templates.value.map((item) => item.type))
  return Array.from(types).map((type) => ({ label: typeLabel(type), value: type }))
})

const currentVariables = computed(
  () => templates.value.find((item) => item.type === editForm.type)?.variables || []
)

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
  editForm.type = record.type
  editForm.name = record.name
  editForm.title = record.title
  editForm.content = record.content
  editForm.active = record.active
  formRef.value?.clearValidate()
  editorOpen.value = true
}

const openCreator = () => {
  editingTemplate.value = undefined
  editForm.type = typeOptions.value[0]?.value || ''
  editForm.name = ''
  editForm.title = ''
  editForm.content = ''
  editForm.active = true
  formRef.value?.clearValidate()
  editorOpen.value = true
}

const closeEditor = () => {
  editorOpen.value = false
  editingTemplate.value = undefined
  editForm.type = ''
  editForm.name = ''
  editForm.title = ''
  editForm.content = ''
  formRef.value?.clearValidate()
}

const appendVariable = (variable: string) => {
  editForm.content += variablePlaceholder(variable)
}

const saveTemplate = async () => {
  try {
    await formRef.value?.validate()
    saving.value = true
    if (editingTemplate.value) {
      await UpdateMessageTemplate(editingTemplate.value.id, {
        name: editForm.name.trim(),
        title: editForm.title.trim(),
        content: editForm.content
      })
    } else {
      await CreateMessageTemplate({
        type: editForm.type,
        name: editForm.name.trim(),
        title: editForm.title.trim(),
        content: editForm.content,
        active: editForm.active
      })
    }
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

const setDefault = async (record: MessageTemplate) => {
  try {
    await SetDefaultMessageTemplate(record.id)
    message.success('默认模板设置成功')
    await loadTemplates()
  } catch (error) {
    message.error(toErrorMessage(error, '默认模板设置失败'))
  }
}

const deleteTemplate = async (record: MessageTemplate) => {
  try {
    await DeleteMessageTemplate(record.id)
    message.success('模板删除成功')
    await loadTemplates()
  } catch (error) {
    message.error(toErrorMessage(error, '模板删除失败'))
  }
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

.danger-link {
  color: #ff4d4f;
}
</style>
