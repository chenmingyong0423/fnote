<template>
  <a-card title="站点信息">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getWebsite"
          />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
      <a-descriptions :column="1" bordered>
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
            <a-image
              v-else
              :width="200"
              :src="data.website_icon === '' ? '' : apiHost + data.website_icon"
            />
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
        <a-descriptions-item label="站长照片">
          <div>
            <StaticUpload
              v-if="editable"
              :image-url="data.website_owner_avatar"
              :authorization="userStore.token"
              @update:imageUrl="(value) => (data.website_owner_avatar = value)"
            />
            <a-image
              v-else
              :width="200"
              :src="data.website_owner_avatar === '' ? '' : apiHost + data.website_owner_avatar"
            />
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
    </a-spin>
    <div style="margin-top: 10px">
      <a-button v-if="!editable" @click="editable = true">编辑</a-button>
      <div v-else>
        <a-button type="primary" @click="cancel" style="margin-right: 5px">取消</a-button>
        <a-button type="primary" @click="save">保存</a-button>
      </div>
    </div>
    <section class="record-section">
      <div class="record-heading">
        <h3>备案信息</h3>
        <a-badge :count="data.website_records.length" show-zero color="#1677ff" />
      </div>

      <div class="record-composer">
        <a-textarea
          v-model:value="record"
          placeholder="输入备案信息或 HTML 链接"
          :auto-size="{ minRows: 2, maxRows: 5 }"
          @keydown.ctrl.enter.prevent="pushRecord"
        />
        <a-button
          type="primary"
          :loading="addingRecord"
          :disabled="!normalizedRecord || recordDuplicated"
          @click="pushRecord"
        >
          <template #icon><PlusOutlined /></template>
          添加
        </a-button>
      </div>

      <a-alert
        v-if="recordDuplicated"
        class="record-alert"
        type="warning"
        show-icon
        message="该备案信息已存在"
      />

      <div v-if="normalizedRecord && !recordDuplicated" class="record-preview">
        <span>预览</span>
        <div v-html="normalizedRecord"></div>
      </div>

      <a-empty
        v-if="!data.website_records.length"
        class="record-empty"
        :image="Empty.PRESENTED_IMAGE_SIMPLE"
        description="暂无备案信息"
      />
      <div v-else class="record-list">
        <div
          v-for="(item, index) in data.website_records"
          :key="`${index}-${item}`"
          class="record-item"
        >
          <div class="record-content" v-html="item"></div>
          <a-popconfirm title="确定删除这条备案信息？" @confirm="pullRecord(item)">
            <a-tooltip title="删除">
              <a-button
                type="text"
                danger
                shape="circle"
                :icon="h(DeleteOutlined)"
                :loading="deletingRecord === item"
              />
            </a-tooltip>
          </a-popconfirm>
        </div>
      </div>
    </section>
  </a-card>
</template>
<script lang="ts" setup>
import {
  AddRecord,
  DeleteRecord,
  GetWebSite,
  UpdateWebSite,
  type WebsiteConfig
} from '@/interfaces/Config'
import { computed, h, ref } from 'vue'
import dayjs from 'dayjs'
import { type Dayjs } from 'dayjs'
import { Empty, message } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import { DeleteOutlined, PlusOutlined, ReloadOutlined } from '@ant-design/icons-vue'
import { toErrorMessage } from '@/utils/error'

document.title = '博客设置 - 后台管理'

const userStore = useUserStore()

const editable = ref<boolean>(false)
const liveTime = ref<Dayjs>()
const apiHost = import.meta.env.VITE_API_HOST

const data = ref<WebsiteConfig>({
  website_name: '',
  website_icon: '',
  website_owner: '',
  website_owner_profile: '',
  website_owner_avatar: '',
  website_runtime: 0,
  website_records: []
})

const loading = ref<boolean>(false)

const getWebsite = async () => {
  try {
    loading.value = true
    const response: any = await GetWebSite()
    if (response.data.code === 0) {
      data.value = response.data.data
        ? { ...response.data.data, website_records: response.data.data.website_records || [] }
        : data.value
      liveTime.value = dayjs(data.value.website_runtime * 1000)
    }
  } catch (error) {
    console.log(error)
  } finally {
    loading.value = false
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
const addingRecord = ref(false)
const deletingRecord = ref('')
const normalizedRecord = computed(() => record.value.trim())
const recordDuplicated = computed(() => data.value.website_records.includes(normalizedRecord.value))

const pushRecord = async () => {
  if (!normalizedRecord.value) {
    message.warning('请输入备案信息')
    return
  }
  if (recordDuplicated.value) {
    message.warning('该备案信息已存在')
    return
  }
  try {
    addingRecord.value = true
    const response: any = await AddRecord(normalizedRecord.value)
    if (response.data.code === 0) {
      message.success('添加成功')
      await getWebsite()
      record.value = ''
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    message.error(toErrorMessage(error, '添加失败'))
  } finally {
    addingRecord.value = false
  }
}

const pullRecord = async (r: string) => {
  try {
    deletingRecord.value = r
    const response: any = await DeleteRecord(r)
    if (response.data.code === 0) {
      message.success('删除成功')
      await getWebsite()
    } else {
      message.error(response.data.message)
    }
  } catch (error) {
    message.error(toErrorMessage(error, '删除失败'))
  } finally {
    deletingRecord.value = ''
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

.record-section {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.record-heading,
.record-composer,
.record-item,
.record-preview {
  display: flex;
  align-items: center;
}

.record-heading {
  gap: 10px;
  margin-bottom: 14px;
}

.record-heading h3 {
  margin: 0;
  font-size: 16px;
  letter-spacing: 0;
}

.record-composer {
  align-items: flex-start;
  gap: 10px;
}

.record-composer .ant-btn {
  flex: 0 0 auto;
}

.record-alert,
.record-preview {
  margin-top: 10px;
}

.record-preview {
  align-items: flex-start;
  gap: 12px;
  padding: 10px 12px;
  background: #fafafa;
  border-radius: 6px;
  overflow-wrap: anywhere;
}

.record-preview > span {
  flex: 0 0 auto;
  color: rgba(0, 0, 0, 0.45);
  font-size: 12px;
}

.record-list {
  margin-top: 14px;
  border-top: 1px solid #f0f0f0;
}

.record-item {
  min-height: 52px;
  gap: 12px;
  padding: 8px 4px;
  border-bottom: 1px solid #f0f0f0;
}

.record-content {
  min-width: 0;
  flex: 1;
  overflow-wrap: anywhere;
}

.record-empty {
  margin: 24px 0 8px;
}

@media (max-width: 576px) {
  .record-composer {
    flex-direction: column;
  }

  .record-composer .ant-btn {
    width: 100%;
  }
}
</style>
