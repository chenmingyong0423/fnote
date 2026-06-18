<template>
  <a-card title="robots.txt">
    <template #extra>
      <a-space>
        <a-tag :color="exists ? 'success' : 'default'">
          {{ exists ? '已生成' : '未生成' }}
        </a-tag>
        <a-tooltip title="刷新内容">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="loadRobots"
          />
        </a-tooltip>
        <a-tooltip title="复制内容">
          <a-button
            shape="circle"
            :icon="h(CopyOutlined)"
            :disabled="!content"
            @click="copyRobots"
          />
        </a-tooltip>
      </a-space>
    </template>

    <a-alert
      class="mb-4"
      type="info"
      show-icon
      message="公开访问地址：/robots.txt"
      description="可直接编辑抓取规则；一键生成会覆盖当前内容，并自动关联 /sitemap.xml。"
    />

    <a-textarea
      v-model:value="content"
      class="robots-editor"
      :rows="18"
      :maxlength="102400"
      show-count
      placeholder="输入 robots.txt 内容"
    />

    <div class="robots-actions">
      <a-popconfirm
        title="生成默认内容会覆盖当前编辑内容，确定继续吗？"
        ok-text="确定生成"
        cancel-text="取消"
        @confirm="generateRobots"
      >
        <a-button :loading="generating">生成默认内容</a-button>
      </a-popconfirm>
      <a-button type="primary" :loading="saving" @click="saveRobots">保存</a-button>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { h, onMounted, ref } from 'vue'
import { CopyOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import {
  GenerateRobotsTxt,
  GetRobotsTxt,
  SaveRobotsTxt,
  type RobotsTxtData
} from '@/interfaces/PostIndex'
import type { IBaseResponse } from '@/interfaces/Common'

const content = ref('')
const exists = ref(false)
const loading = ref(false)
const saving = ref(false)
const generating = ref(false)

const loadRobots = async () => {
  try {
    loading.value = true
    const response = await GetRobotsTxt()
    const result = response.data as IBaseResponse & { data: RobotsTxtData }
    content.value = result.data?.content || ''
    exists.value = result.data?.exists || false
  } catch (error) {
    console.error(error)
    message.error('robots.txt 加载失败')
  } finally {
    loading.value = false
  }
}

const generateRobots = async () => {
  try {
    generating.value = true
    const response = await GenerateRobotsTxt()
    const result = response.data as IBaseResponse & { data: RobotsTxtData }
    content.value = result.data.content
    exists.value = true
    message.success('robots.txt 生成成功')
  } catch (error) {
    console.error(error)
    message.error('robots.txt 生成失败')
  } finally {
    generating.value = false
  }
}

const saveRobots = async () => {
  try {
    saving.value = true
    await SaveRobotsTxt(content.value)
    exists.value = true
    message.success('robots.txt 保存成功')
  } catch (error) {
    console.error(error)
    message.error('robots.txt 保存失败')
  } finally {
    saving.value = false
  }
}

const copyRobots = async () => {
  await navigator.clipboard.writeText(content.value)
  message.success('复制成功')
}

onMounted(loadRobots)
</script>

<style scoped>
.robots-editor {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  line-height: 1.65;
  resize: vertical;
}

.robots-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}

@media (max-width: 576px) {
  .robots-actions {
    flex-direction: column-reverse;
  }

  .robots-actions :deep(.ant-btn) {
    width: 100%;
  }
}
</style>
