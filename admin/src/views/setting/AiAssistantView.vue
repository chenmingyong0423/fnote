<template>
  <a-card title="AI 助手">
    <template #extra>
      <a-tooltip title="重新加载配置">
        <a-button
          shape="circle"
          :icon="h(ReloadOutlined)"
          :loading="loading"
          :disabled="saving || !!testingProviderId"
          @click="loadConfig"
        />
      </a-tooltip>
    </template>

    <a-alert v-if="loadError" class="config-alert" type="error" show-icon :message="loadError" />

    <a-spin :spinning="loading">
      <a-form
        v-if="loaded"
        ref="formRef"
        :model="form"
        layout="vertical"
        :disabled="busy"
        @finish="save"
      >
        <a-form-item label="启用助手">
          <a-switch
            v-model:checked="form.enabled"
            checked-children="开启"
            un-checked-children="关闭"
          />
        </a-form-item>

        <a-form-item
          label="助手名称"
          name="assistant_name"
          :rules="[{ required: true, whitespace: true, message: '请输入助手名称', max: 80 }]"
        >
          <a-input v-model:value="form.assistant_name" :maxlength="80" />
        </a-form-item>

        <a-form-item label="系统提示词" name="system_prompt">
          <a-textarea
            v-model:value="form.system_prompt"
            :auto-size="{ minRows: 5, maxRows: 14 }"
            :maxlength="32000"
            placeholder="设置助手的回答范围、风格和引用要求"
          />
        </a-form-item>

        <a-form-item label="可用工具">
          <a-checkbox-group v-model:value="form.enabled_tools" class="tool-options">
            <a-checkbox v-for="tool in tools" :key="tool.value" :value="tool.value">
              {{ tool.label }}
            </a-checkbox>
          </a-checkbox-group>
        </a-form-item>

        <div class="providers-heading">
          <h3>模型服务</h3>
          <a-button :disabled="busy || form.providers.length >= 20" @click="addProvider">
            <template #icon><PlusOutlined /></template>
            添加服务
          </a-button>
        </div>

        <a-empty v-if="!form.providers.length" description="尚未添加模型服务" />

        <a-card
          v-for="(provider, index) in form.providers"
          :key="provider.id"
          class="provider-card"
          size="small"
          :title="provider.name || `模型服务 ${index + 1}`"
        >
          <template #extra>
            <a-space>
              <a-button
                html-type="button"
                :loading="testingProviderId === provider.id"
                :disabled="busy"
                @click="testProvider(provider, index)"
                >测试连接</a-button
              >
              <a-radio
                :checked="form.active_provider_id === provider.id"
                @change="form.active_provider_id = provider.id"
              >
                当前使用
              </a-radio>
              <a-popconfirm
                title="删除这个模型服务？"
                description="保存配置后生效。"
                :disabled="busy"
                @confirm="removeProvider(provider.id)"
              >
                <a-button type="text" danger :icon="h(DeleteOutlined)" aria-label="删除模型服务" />
              </a-popconfirm>
            </a-space>
          </template>

          <div class="provider-fields">
            <a-form-item
              label="服务名称"
              :name="['providers', index, 'name']"
              :rules="[{ required: true, whitespace: true, message: '请输入服务名称', max: 80 }]"
            >
              <a-input
                v-model:value="provider.name"
                :maxlength="80"
                placeholder="例如：主要模型服务"
              />
            </a-form-item>
            <a-form-item label="服务类型" :name="['providers', index, 'api_type']">
              <a-select
                v-model:value="provider.api_type"
                :options="protocols"
                @change="changeProtocol(provider)"
              />
            </a-form-item>
            <a-form-item
              label="输出 token 上限"
              :name="['providers', index, 'max_output_tokens']"
              :extra="
                provider.api_type === 'openai_responses'
                  ? '对应 max_output_tokens，包含推理 token；留空使用服务商默认值。'
                  : '对应 max_tokens，包含推理 token；留空使用服务商默认值。'
              "
            >
              <a-input-number
                v-model:value="provider.max_output_tokens"
                :min="16"
                :max="32768"
                :precision="0"
                placeholder="默认"
                style="width: 100%"
              />
            </a-form-item>
            <a-form-item
              v-if="provider.api_type === 'openai_chat_completions'"
              label="DeepSeek 思考模式"
              :name="['providers', index, 'thinking_mode']"
            >
              <a-select
                v-model:value="provider.thinking_mode"
                :options="thinkingModes"
                @change="changeThinkingMode(provider)"
              />
            </a-form-item>
            <a-form-item
              label="推理强度"
              :name="['providers', index, 'reasoning_effort']"
              extra="不同模型支持的选项不同；默认时不发送该参数。"
            >
              <a-select
                v-model:value="provider.reasoning_effort"
                :options="reasoningOptions(provider.api_type)"
                :disabled="busy || provider.thinking_mode === 'disabled'"
              />
            </a-form-item>
            <a-form-item
              label="服务地址"
              :name="['providers', index, 'base_url']"
              :rules="[{ required: true, whitespace: true, message: '请输入服务地址' }]"
            >
              <a-input
                v-model:value="provider.base_url"
                :maxlength="2048"
                placeholder="https://api.example.com/v1"
              />
            </a-form-item>
            <a-form-item
              label="模型名称"
              :name="['providers', index, 'model']"
              :rules="[{ required: true, whitespace: true, message: '请输入模型名称', max: 200 }]"
            >
              <a-input
                v-model:value="provider.model"
                :maxlength="200"
                placeholder="填写服务商提供的模型名称"
              />
            </a-form-item>
          </div>

          <a-form-item
            label="API Key"
            :extra="
              provider.has_api_key
                ? '已配置密钥，留空保留。修改服务地址后，请确认密钥仍适用。'
                : '密钥保存后不会回显；启用此服务前需要填写。'
            "
          >
            <a-input-password
              v-model:value="provider.api_key"
              :maxlength="8192"
              autocomplete="new-password"
              :placeholder="provider.has_api_key ? '输入新密钥以替换，留空保留' : '输入 API Key'"
            />
          </a-form-item>
          <a-alert
            v-if="needsConnectionTest(provider) && !testResults[provider.id]"
            type="info"
            show-icon
            message="请先测试连接；测试通过后才能保存此模型服务。"
          />
          <a-alert
            v-if="testResults[provider.id]"
            :type="testResults[provider.id].success ? 'success' : 'error'"
            show-icon
            :message="`${testResults[provider.id].message} · ${testResults[provider.id].duration_ms} ms`"
          >
            <template v-if="testResults[provider.id].reply" #description>
              <div class="test-reply">{{ testResults[provider.id].reply }}</div>
            </template>
          </a-alert>
        </a-card>

        <a-alert
          v-if="form.enabled && !form.active_provider_id"
          class="config-alert"
          type="warning"
          show-icon
          message="开启助手前，请添加并选择一个模型服务。"
        />

        <a-alert
          v-if="untestedProviders.length"
          class="config-alert"
          type="warning"
          show-icon
          message="新增或修改的模型服务尚未全部通过测试，请先点击对应的「测试连接」。"
        />
        <a-button
          type="primary"
          html-type="submit"
          :loading="saving"
          :disabled="busy || !!untestedProviders.length"
          >保存配置</a-button
        >
      </a-form>
    </a-spin>
  </a-card>
</template>

<script setup lang="ts">
import { computed, h, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { message, type FormInstance } from 'ant-design-vue'
import { DeleteOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import {
  GetAIAssistantConfig,
  UpdateAIAssistantConfig,
  TestAIProvider,
  type AIAPIType,
  type AIProviderRequest,
  type AIProviderTestResult,
  type AIAssistantConfig,
  type AIProvider
} from '@/interfaces/AiAssistant'
import { toErrorMessage } from '@/utils/error'

type ProviderForm = AIProvider & { api_key: string }
type ConfigForm = Omit<AIAssistantConfig, 'providers'> & { providers: ProviderForm[] }

const formRef = ref<FormInstance>()
const loading = ref(false)
const saving = ref(false)
const loaded = ref(false)
const loadError = ref('')
const testingProviderId = ref('')
const testResults = ref<Record<string, AIProviderTestResult>>({})
const savedConnections = ref<Record<string, string>>({})
let testController: AbortController | undefined
const busy = computed(
  () => saving.value || loading.value || !!loadError.value || !!testingProviderId.value
)
const form = ref<ConfigForm>({
  enabled: false,
  assistant_name: '',
  system_prompt: '',
  enabled_tools: [],
  active_provider_id: '',
  providers: []
})

const tools = [
  { value: 'search_posts', label: '搜索公开文章' },
  { value: 'read_post', label: '读取文章正文' },
  { value: 'get_blog_overview', label: '获取博客概况' },
  { value: 'get_author_profile', label: '获取博主公开资料' }
]
const protocols = [
  { value: 'openai_responses', label: 'GPT（OpenAI Responses）' },
  { value: 'openai_chat_completions', label: 'DeepSeek（Chat Completions）' }
]
const thinkingModes = [
  { value: '', label: '默认' },
  { value: 'enabled', label: '开启' },
  { value: 'disabled', label: '关闭' }
]
const reasoningOptions = (apiType: AIAPIType) => [
  { value: '', label: '默认' },
  ...(apiType === 'openai_responses'
    ? ['none', 'minimal', 'low', 'medium', 'high', 'xhigh', 'max']
    : ['none', 'low', 'high', 'max']
  ).map((value) => ({ value, label: value }))
]

const changeProtocol = (provider: ProviderForm) => {
  const defaults = ['https://api.openai.com/v1', 'https://api.deepseek.com']
  if (!provider.base_url || defaults.includes(provider.base_url.replace(/\/+$/, ''))) {
    provider.base_url = provider.api_type === 'openai_responses' ? defaults[0] : defaults[1]
  }
  provider.reasoning_effort = ''
  provider.thinking_mode = ''
}

const changeThinkingMode = (provider: ProviderForm) => {
  if (provider.thinking_mode === 'disabled' || provider.reasoning_effort === 'none') {
    provider.reasoning_effort = ''
  }
}

const providerRequest = (provider: ProviderForm): AIProviderRequest => ({
  id: provider.id,
  name: provider.name,
  api_type: provider.api_type,
  base_url: provider.base_url,
  model: provider.model,
  api_key: provider.api_key.trim() || undefined,
  max_output_tokens: provider.max_output_tokens ?? undefined,
  reasoning_effort: provider.reasoning_effort || undefined,
  thinking_mode: provider.thinking_mode || undefined
})

// Display names do not affect connectivity. Keep test results only while the
// connection parameters (including an entered replacement key) are unchanged.
const connectionSnapshot = (provider: ProviderForm) =>
  JSON.stringify([
    provider.api_type,
    provider.base_url.trim().replace(/\/+$/, ''),
    provider.model.trim(),
    provider.api_key.trim(),
    provider.max_output_tokens ?? null,
    provider.reasoning_effort || '',
    provider.thinking_mode || ''
  ])
const needsConnectionTest = (provider: ProviderForm) =>
  savedConnections.value[provider.id] !== connectionSnapshot(provider)
const untestedProviders = computed(() =>
  form.value.providers.filter(
    (provider) => needsConnectionTest(provider) && !testResults.value[provider.id]?.success
  )
)

const testProvider = async (provider: ProviderForm, index: number) => {
  if (busy.value) return
  try {
    await formRef.value?.validateFields([
      ['providers', index, 'base_url'],
      ['providers', index, 'model']
    ])
  } catch {
    return
  }
  if (busy.value) return
  if (!provider.api_key.trim() && !provider.has_api_key) {
    message.warning('请先填写 API Key')
    return
  }
  testingProviderId.value = provider.id
  delete testResults.value[provider.id]
  testController = new AbortController()
  try {
    const response = await TestAIProvider(providerRequest(provider), testController.signal)
    if (response.data.code !== 0 || !response.data.data)
      throw new Error(response.data.message || '测试失败')
    testResults.value[provider.id] = response.data.data
  } catch (error) {
    testResults.value[provider.id] = {
      success: false,
      message: toErrorMessage(error, '测试失败，请重试'),
      duration_ms: 0
    }
  } finally {
    testingProviderId.value = ''
    testController = undefined
  }
}

const loadConfig = async () => {
  loading.value = true
  loadError.value = ''
  try {
    const response = await GetAIAssistantConfig()
    if (response.data.code !== 0 || !response.data.data) {
      throw new Error(response.data.message || '配置加载失败')
    }
    const config = response.data.data
    form.value = {
      ...config,
      providers: config.providers.map((provider) => ({
        ...provider,
        api_key: '',
        reasoning_effort: provider.reasoning_effort || '',
        thinking_mode: provider.thinking_mode || ''
      }))
    }
    savedConnections.value = Object.fromEntries(
      form.value.providers.map((provider) => [provider.id, connectionSnapshot(provider)])
    )
    testResults.value = {}
    loaded.value = true
    formRef.value?.clearValidate()
  } catch (error) {
    loadError.value = toErrorMessage(error, '配置加载失败，请重试')
  } finally {
    loading.value = false
  }
}

const addProvider = () => {
  const id = crypto.randomUUID()
  form.value.providers.push({
    id,
    name: `模型服务 ${form.value.providers.length + 1}`,
    api_type: 'openai_responses',
    base_url: 'https://api.openai.com/v1',
    model: '',
    has_api_key: false,
    api_key: '',
    reasoning_effort: '',
    thinking_mode: ''
  })
  if (!form.value.active_provider_id) form.value.active_provider_id = id
}

const removeProvider = (id: string) => {
  form.value.providers = form.value.providers.filter((provider) => provider.id !== id)
  if (form.value.active_provider_id === id) form.value.active_provider_id = ''
  formRef.value?.clearValidate()
}

const save = async () => {
  if (busy.value || !loaded.value) return
  if (untestedProviders.value.length) {
    message.warning('请先测试新增或修改的模型服务，全部通过后再保存')
    return
  }
  const active = form.value.providers.find(
    (provider) => provider.id === form.value.active_provider_id
  )
  if (form.value.enabled && (!active || (!active.api_key.trim() && !active.has_api_key))) {
    message.warning('请选择当前使用的模型服务，并配置 API Key')
    return
  }
  saving.value = true
  try {
    const response = await UpdateAIAssistantConfig({
      enabled: form.value.enabled,
      assistant_name: form.value.assistant_name,
      system_prompt: form.value.system_prompt,
      enabled_tools: [...form.value.enabled_tools],
      active_provider_id: form.value.active_provider_id,
      providers: form.value.providers.map(providerRequest)
    })
    if (response.data.code !== 0) throw new Error(response.data.message || '保存失败')
    // Clear plaintext as soon as saving succeeds, even if the refresh fails.
    form.value.providers.forEach((provider) => {
      if (provider.api_key.trim()) provider.has_api_key = true
      provider.api_key = ''
    })
    message.success('AI 助手配置已保存')
    await loadConfig()
  } catch (error) {
    message.error(toErrorMessage(error, '配置保存失败，请重试'))
  } finally {
    saving.value = false
  }
}

watch(
  () =>
    form.value.providers.map((provider) => [provider.id, connectionSnapshot(provider)] as const),
  (current, previous) => {
    const before = new Map(previous)
    const after = new Map(current)
    Object.keys(testResults.value).forEach((id) => {
      if (!after.has(id) || before.get(id) !== after.get(id)) delete testResults.value[id]
    })
  },
  { flush: 'sync' }
)
onBeforeUnmount(() => testController?.abort())
onMounted(loadConfig)
</script>

<style scoped>
.config-alert {
  margin-bottom: 20px;
}
.tool-options {
  display: flex;
  flex-wrap: wrap;
  gap: 12px 20px;
}
.providers-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}
.providers-heading h3 {
  margin: 0;
  font-size: 16px;
}
.provider-card {
  margin-bottom: 20px;
}
.test-reply {
  white-space: pre-wrap;
  overflow-wrap: anywhere;
  max-height: 160px;
  overflow-y: auto;
}
.provider-fields {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0 20px;
}
@media (max-width: 640px) {
  .provider-fields {
    grid-template-columns: minmax(0, 1fr);
  }
  .provider-card :deep(.ant-card-head) {
    flex-wrap: wrap;
  }
  .provider-card :deep(.ant-card-head-wrapper) {
    flex-wrap: wrap;
    gap: 8px;
    padding: 8px 0;
  }
}
</style>
