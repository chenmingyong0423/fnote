<template>
  <div class="flex gap-x-2">
    <a-button @click="DownloadBackup">导出数据</a-button>
    <a-button @click="RecoveryBackup">导入数据</a-button>
  </div>
</template>

<script lang="ts" setup>
import { message } from 'ant-design-vue'
import { Backup, Recovery } from '@/interfaces/Backup'

const RecoveryBackup = async () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.tar'
  input.onchange = async (e: any) => {
    const file = e.target.files[0]
    const formData = new FormData()
    formData.append('file', file)
    try {
      await Recovery(formData)
      message.success('导入成功')
    } catch (error) {
      console.log(error)
    }
  }
  input.click()
}

const DownloadBackup = async () => {
  try {
    const response: any = await Backup()

    // 首先检查Content-Type来决定是否是文件
    const contentType = response.headers['content-type']
    if (contentType.includes('application/json')) {
      message.error(response.message)
    } else {
      console.log(response.headers)
      // 提取文件名
      const currentDate = new Date().toISOString().slice(0, 10).replace(/-/g, '')
      const filename = `backup_${currentDate}.tar`
      // 创建一个URL指向返回的Blob对象
      console.log(response.data)
      const fileBlob = new Blob([response.data], { type: 'application/tar' })
      const fileURL = window.URL.createObjectURL(fileBlob)
      // 创建一个临时a标签用于下载文件
      const fileLink = document.createElement('a')
      fileLink.href = fileURL
      fileLink.setAttribute('download', filename) // 设定下载文件的名称和格式
      document.body.appendChild(fileLink)
      fileLink.click()

      // 清理操作
      document.body.removeChild(fileLink)
      window.URL.revokeObjectURL(fileURL)
    }
  } catch (error) {
    console.log(error)
  }
}
</script>
