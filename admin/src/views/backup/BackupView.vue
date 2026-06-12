<template>
  <a-card title="数据备份">
    <div class="flex flex-col gap-y-4">
      <div class="flex gap-x-2">
        <a-button
          type="primary"
          :loading="exporting"
          :disabled="importing"
          @click="downloadBackup"
        >
          {{ exporting ? '导出中...' : '导出数据' }}
        </a-button>
        <a-button
          :loading="importing"
          :disabled="exporting"
          @click="selectBackupFile"
        >
          {{ importing ? '导入中...' : '导入数据' }}
        </a-button>
      </div>

      <a-alert
        v-if="runningTip"
        type="info"
        show-icon
        banner
        :message="runningTip"
      />
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue'
import { message } from 'ant-design-vue'
import { Backup, Recovery } from '@/interfaces/Backup'
import { toErrorMessage } from '@/utils/error'

document.title = '备份 - 后台管理'

const exporting = ref(false)
const importing = ref(false)

const runningTip = computed(() => {
  if (exporting.value) return '正在导出数据，请不要关闭页面或重复点击。'
  if (importing.value) return '正在导入数据，请不要关闭页面或重复点击。'
  return ''
})

function getBackupFilename(contentDisposition?: string) {
  if (!contentDisposition) return `fnote-backup-${Date.now()}.zip`

  const utf8Filename = contentDisposition.match(/filename\*=UTF-8''([^;]+)/i)
  if (utf8Filename?.[1]) {
    return decodeURIComponent(utf8Filename[1])
  }

  const filename = contentDisposition.match(/filename="?([^"]+)"?/i)
  return filename?.[1] || `fnote-backup-${Date.now()}.zip`
}

function downloadBlob(data: BlobPart, filename: string) {
  const fileBlob = new Blob([data], { type: 'application/zip' })
  const fileURL = window.URL.createObjectURL(fileBlob)
  const fileLink = document.createElement('a')
  fileLink.href = fileURL
  fileLink.setAttribute('download', filename)
  document.body.appendChild(fileLink)
  fileLink.click()
  document.body.removeChild(fileLink)
  window.URL.revokeObjectURL(fileURL)
}

const selectBackupFile = () => {
  if (importing.value || exporting.value) return

  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.zip'
  input.onchange = async (event: Event) => {
    const target = event.target as HTMLInputElement
    const file = target.files?.[0]
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)

    importing.value = true
    const hide = message.loading('正在导入数据，请稍候...', 0)
    try {
      await Recovery(formData)
      message.success('导入成功')
    } catch (error) {
      message.error(toErrorMessage(error, '导入失败，请稍后再试'))
    } finally {
      hide()
      importing.value = false
    }
  }
  input.click()
}

const downloadBackup = async () => {
  if (exporting.value || importing.value) return

  exporting.value = true
  const hide = message.loading('正在导出数据，请稍候...', 0)
  try {
    const response = await Backup()
    const contentType = String(response.headers['content-type'] || '')

    if (contentType.includes('application/json')) {
      message.error(response.data?.message || '导出失败，请稍后再试')
      return
    }

    const filename = getBackupFilename(String(response.headers['content-disposition'] || ''))
    downloadBlob(response.data, filename)
    message.success('导出成功')
  } catch (error) {
    message.error(toErrorMessage(error, '导出失败，请稍后再试'))
  } finally {
    hide()
    exporting.value = false
  }
}
</script>
