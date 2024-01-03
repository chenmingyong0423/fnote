<template>
  <div>
    <a-descriptions title="公告配置" :column="1" bordered>
      <a-descriptions-item label="公告标题">
        <div>
          <a-input v-if="editable" v-model:value="data.title" style="margin: -5px 0" />
          <template v-else>
            {{ data.title }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="内容">
        <div>
          <a-input v-if="editable" v-model:value="data.content" style="margin: -5px 0" />
          <template v-else>
            {{ data.content }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="是否显示">
        <div>
          <a-switch v-model:checked="data.enabled" @change="enabled" />
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="发布时间">
        <div>
          {{ dayjs.unix(data.publish_time).format('YYYY-MM-DD HH:mm:ss') }}
        </div>
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
import type { NoticeConfig } from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import dayjs from 'dayjs'

const editable = ref<boolean>(false)
const data = ref<NoticeConfig>({
  title: '',
  content: '',
  enabled: false,
  publish_time: 0
})

const getNotice = async () => {
  try {
    const response = await axios.get<IResponse<NoticeConfig>>('/admin/configs/notice')
    data.value = response.data.data || {}
    console.log(data)
  } catch (error) {
    console.log(error)
    message.error('获取信息失败')
  }
}
getNotice()

const cancel = () => {
  editable.value = false
  getNotice()
}

const save = async () => {
  try {
    const response = await axios.put<IBaseResponse>('/admin/configs/notice', {
      title: data.value.title,
      content: data.value.content
    })
    if (response.data.code === 200) {
      message.success('保存成功')
      await getNotice()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('保存失败')
  }
}

const enabled = async () => {
  try {
    const response = await axios.put<IBaseResponse>('/admin/configs/notice/enabled', {
      enabled: data.value.enabled
    })
    if (response.data.code === 200) {
      message.success('更新成功')
      await getNotice()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
    message.error('更新失败')
  }
}
</script>

<style scoped></style>
