<template>
  <div>
    <a-descriptions title="站点设置" :column="1" bordered>
      <a-descriptions-item label="站点名称">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.name"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ data.name }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站点运行时间">
        <div>
          <a-date-picker
            v-if="editable"
            v-model:value="liveTime" @change="liveTimeChanged" />
          <template v-else>
            {{ dayjs.unix(data.live_time).format('YYYY-MM-DD') }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站点图标">
        <a-image
          :width="200"
          :src="data.icon"
        />
      </a-descriptions-item>
    </a-descriptions>
    <div style="margin-top: 10px;">
      <a-button v-if="!editable" type="primary" @click="editable=true">编辑</a-button>
      <div v-else>
        <a-button type="primary" @click="cancel" style="margin-right: 5px;">取消</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import axios from '@/http/axios'
import type { IResponse } from '@/interfaces/Common'
import type { WebsiteConfig } from '@/interfaces/Config'
import { ref } from 'vue'
import dayjs from 'dayjs'
import { type Dayjs } from 'dayjs'
import { message } from 'ant-design-vue'

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
    const response = await axios.put<IResponse<WebsiteConfig>>('/admin/configs/website', {
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
</script>