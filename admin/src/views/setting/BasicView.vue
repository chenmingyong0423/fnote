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
      <a-descriptions-item label="站点 logo">
        <div>
          <StaticUpload
            v-if="editable"
            :image-url="data.website_icon"
            :authorization="userStore.token"
            @update:imageUrl="(value) => (data.website_icon = value)"
          />
          <a-image v-else :width="200" :src="apiHost + data.website_icon" />
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长昵称">
        <div>
          <a-input v-if="editable" v-model:value="data.website_owner" style="margin: -5px 0" />
          <template v-else>
            {{ data.website_owner }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长简介">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.website_owner_profile"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ data.website_owner_profile }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长邮箱">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.website_owner_email"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ data.website_owner_email }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站长照片">
        <div>
          <StaticUpload
            v-if="editable"
            :image-url="data.website_owner_avatar"
            :authorization="userStore.token"
            @update:imageUrl="(value) => (data.website_owner_avatar = value)"
          />
          <a-image v-else :width="200" :src="apiHost + data.website_owner_avatar" />
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="站点运行时间">
        <div>
          <a-date-picker v-if="editable" v-model:value="liveTime" @change="liveTimeChanged" />
          <template v-else>
            {{ dayjs.unix(data.website_runtime).format('YYYY-MM-DD') }}
          </template>
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
    <div class="text-4 font-bold my-2">备案信息</div>
    <div class="flex flex-col">
      <div class="flex">
        <a-input v-model:value="record"></a-input>
        <a-button @click="pushRecord">添加</a-button>
      </div>
      <div
        class="flex p-3 border-b-1 border-b-solid border-b-gray-2"
        v-for="(item, index) in data.website_records"
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
import {
  AddRecord,
  DeleteRecord,
  GetWebSite,
  UpdateWebSite,
  type WebsiteConfig
} from '@/interfaces/Config'
import { ref } from 'vue'
import dayjs from 'dayjs'
import { type Dayjs } from 'dayjs'
import { message } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'

const userStore = useUserStore()

const editable = ref<boolean>(false)
const liveTime = ref<Dayjs>()
const apiHost = import.meta.env.VITE_API_HOST;

const data = ref<WebsiteConfig>({
  website_name: '',
  website_icon: '',
  website_owner: '',
  website_owner_profile: '',
  website_owner_avatar: '',
  website_owner_email: '',
  website_runtime: 0,
  website_records: []
})

const getWebsite = async () => {
  try {
    const response: any = await GetWebSite()
    if (response.data.code === 0) {
      data.value = response.data.data || data.value
      liveTime.value = dayjs(data.value.website_runtime * 1000)
    }
  } catch (error) {
    console.log(error)
  }
}
getWebsite()

const liveTimeChanged = (date: Dayjs) => {
  liveTime.value = date
  data.value.website_runtime = Math.floor(date.valueOf() / 1000)
}

const cancel = () => {
  editable.value = false
  getWebsite()
}

const save = async () => {
  try {
    const response: any = await UpdateWebSite({
      website_name: data.value.website_name,
      website_icon: data.value.website_icon,
      website_owner: data.value.website_owner,
      website_owner_profile: data.value.website_owner_profile,
      website_owner_avatar: data.value.website_owner_avatar,
      website_owner_email: data.value.website_owner_email,
      website_runtime: data.value.website_runtime
    })
    if (response.data.code === 0) {
      message.success('保存成功')
      await getWebsite()
      editable.value = false
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}

const record = ref<string>('')

const pushRecord = async () => {
  if (record.value === '') {
    message.warning('请输入备案信息')
    return
  }
  try {
    const response: any = await AddRecord(record.value)
    if (response.data.code === 0) {
      message.success('添加成功')
      await getWebsite()
      record.value = ''
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
}

const pullRecord = async (r: string) => {
  try {
    const response: any = await DeleteRecord(r)
    if (response.data.code === 0) {
      message.success('删除成功')
      await getWebsite()
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    console.log(error)
  }
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
