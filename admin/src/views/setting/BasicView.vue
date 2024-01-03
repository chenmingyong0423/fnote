<template>
  <div>
    <a-descriptions title="站点信息" :column="1" bordered>
      <a-descriptions-item label="站点名称">
        <div>
          <a-input v-if="editable" v-model:value="data.name" style="margin: -5px 0" />
          <template v-else>
            {{ data.name }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站点运行时间">
        <div>
          <a-date-picker v-if="editable" v-model:value="liveTime" @change="liveTimeChanged" />
          <template v-else>
            {{ dayjs.unix(data.live_time).format('YYYY-MM-DD') }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站点图标">
        <a-image :width="200" :src="data.icon" />
        <a-upload
          v-if="editable"
          v-model:file-list="fileList"
          name="file"
          action="http://localhost:8080/admin/files/upload"
          @change="handleChange"
          :before-upload="beforeUpload"
          :maxCount="1"
        >
          <a-button>
            <upload-outlined></upload-outlined>
            Click to Upload
          </a-button>
        </a-upload>
      </a-descriptions-item>
    </a-descriptions>
    <div style="margin-top: 10px">
      <a-button v-if="!editable" type="primary" @click="editable = true">编辑</a-button>
      <div v-else>
        <a-button type="primary" @click="cancel" style="margin-right: 5px">取消</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import axios from '@/http/axios'
import type { IBaseResponse, IResponse } from '@/interfaces/Common'
import type { WebsiteConfig } from '@/interfaces/Config'
import { ref } from 'vue'
import dayjs from 'dayjs'
import { type Dayjs } from 'dayjs'
import { message } from 'ant-design-vue'
import type { UploadChangeParam, UploadProps } from 'ant-design-vue'

const editable = ref<boolean>(false)
const liveTime = ref<Dayjs>()

const data = ref<WebsiteConfig>({
  name: '',
  icon: '',
  post_count: 0,
  category_count: 0,
  view_count: 0,
  live_time: 0,
  domain: '',
  records: []
})

const getWebsite = async () => {
  try {
    const response = await axios.get<IResponse<WebsiteConfig>>('/admin/configs/website')
    data.value = response.data.data || {}
    liveTime.value = dayjs(data.value.live_time * 1000)
  } catch (error) {
    console.log(error)
    message.error('获取站点信息失败')
  }
}
getWebsite()

const liveTimeChanged = (date: Dayjs) => {
  liveTime.value = date
  data.value.live_time = Math.floor(date.valueOf() / 1000)
}

const cancel = () => {
  editable.value = false
  getWebsite()
}

const save = async () => {
  try {
    const response = await axios.put<IBaseResponse>('/admin/configs/website', {
      name: data.value.name,
      live_time: data.value.live_time,
      icon: data.value.icon
    })
    if (response.data.code === 200) {
      message.success('保存成功')
      await getWebsite()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('保存失败')
  }
}

// 文件操作
const fileList = ref<UploadProps['fileList']>([])

const handleChange = (info: UploadChangeParam) => {
  if (info.file.status === 'uploading') {
    return
  }
  console.log(info)
  if (info.file.status === 'done') {
    // Get this url from response in real world.
    data.value.icon = info.file.response.data.url
    console.log(data.value.icon)
    message.success('上传成功')
  }
  if (info.file.status === 'error') {
    message.error('upload error')
  }
}

const beforeUpload = (file: UploadProps['fileList'][number]) => {
  const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
  if (!isJpgOrPng) {
    message.error('You can only upload JPG file!')
  }
  const isLt2M = file.size / 1024 / 1024 < 2
  if (!isLt2M) {
    message.error('Image must smaller than 2MB!')
  }
  return isJpgOrPng && isLt2M
}
</script>

<style scoped>
.upload-list-inline :deep(.ant-upload-list-item) {
  float: left;
  width: 200px;
  margin-right: 8px;
}
.upload-list-inline [class*='-upload-list-rtl'] :deep(.ant-upload-list-item) {
  float: right;
}
</style>
