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
import {
  GetNotice,
  type NoticeConfig,
  UpdateNotice,
  UpdateNoticeEnabled
} from '@/interfaces/Config'
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
    const response: any = await GetNotice()
    data.value = response.data.data || {}
  } catch (error) {
    console.log(error)
  }
}
getNotice()

const cancel = () => {
  editable.value = false
  getNotice()
}

const save = async () => {
  try {
    const response: any = await UpdateNotice({
      title: data.value.title,
      content: data.value.content
    })
    if (response.data.code === 0) {
      message.success('保存成功')
      await getNotice()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}

const enabled = async () => {
  try {
    const response: any = await UpdateNoticeEnabled(data.value.enabled)
    if (response.data.code === 0) {
      message.success('更新成功')
      await getNotice()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}
</script>

<style scoped></style>
