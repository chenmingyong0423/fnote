<template>
  <div>
    <a-descriptions title="站点信息" :column="1" bordered>
      <a-descriptions-item label="站点名称">
        <div>
          <a-input v-if="editable" v-model:value="data.website_name" style="margin: -5px 0" />
          <template v-else>
            {{ data.website_name }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长昵称">
        <div>
          <a-input v-if="editable" v-model:value="data.owner_name" style="margin: -5px 0" />
          <template v-else>
            {{ data.owner_name }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长简介">
        <div>
          <a-input v-if="editable" v-model:value="data.owner_profile" style="margin: -5px 0" />
          <template v-else>
            {{ data.owner_profile }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长照片">
        <div>
          <a-image :width="200" :src="data.owner_picture" />
          <a-upload
            v-if="editable"
            v-model:file-list="fileList4picture"
            name="file"
            action="http://localhost:8080/admin/files/upload"
            @change="handleChange4picture"
            :before-upload="beforeUpload4picture"
            :maxCount="1"
          >
            <a-button>
              <upload-outlined></upload-outlined>
              Click to Upload
            </a-button>
          </a-upload>
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
    <div class="text-4 font-bold my-2">备案信息</div>
    <div class="flex flex-col">
      <div class="flex">
        <a-input v-model:value="record"></a-input> <a-button @click="pushRecord">添加</a-button>
      </div>
      <div
        class="flex p-3 border-b-1 border-b-solid border-b-gray-2"
        v-for="(item, index) in data.records"
        :key="index"
      >
        <div v-html="item"></div>
        <a-popconfirm class="ml-auto" title="确定取消？" @confirm="pullRecord(item)">
          <a-button type="primary" danger>删除</a-button>
        </a-popconfirm>
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
  website_name: '',
  icon: '',
  live_time: 0,
  records: [],
  owner_name: '',
  owner_profile: '',
  owner_picture: ''
})

const getWebsite = async () => {
  try {
    const response = await axios.get<IResponse<WebsiteConfig>>('/admin/configs/website')
    data.value = response.data.data || data.value
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
      website_name: data.value.website_name,
      live_time: data.value.live_time,
      icon: data.value.icon,
      owner_name: data.value.owner_name,
      owner_profile: data.value.owner_profile,
      owner_picture: data.value.owner_name
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

const record = ref<string>('')

const pushRecord = async () => {
  if (record.value === '') {
    message.warning('请输入备案信息')
    return
  }
  try {
    const response = await axios.post<IBaseResponse>('/admin/configs/website/records', {
      record: record.value
    })
    if (response.data.code === 200) {
      message.success('添加成功')
      await getWebsite()
      record.value = ''
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('添加失败')
  }
}

const pullRecord = async (item: string) => {
  try {
    const response = await axios.delete<IBaseResponse>(
      `/admin/configs/website/records?record=${item}`
    )
    if (response.data.code === 200) {
      message.success('删除成功')
      await getWebsite()
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('删除失败')
  }
}

// picture 文件操作

// 文件操作
const fileList4picture = ref<UploadProps['fileList']>([])

const handleChange4picture = (info: UploadChangeParam) => {
  if (info.file.status === 'uploading') {
    return
  }
  console.log(info)
  if (info.file.status === 'done') {
    // Get this url from response in real world.
    data.value.owner_picture = info.file.response.data.url
    message.success('上传成功')
  }
  if (info.file.status === 'error') {
    message.error('upload error')
  }
}

const beforeUpload4picture = (file: UploadProps['fileList'][number]) => {
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
