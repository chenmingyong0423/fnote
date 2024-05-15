<template>
  <div>
    <a-descriptions title="seo meta" :column="1" bordered>
      <a-descriptions-item label="title - 标题">
        <div>
          <a-input v-if="editable" v-model:value="data.title" style="margin: -5px 0" />
          <template v-else>
            {{ data.title }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="og_title - 社交标题">
        <div>
          <a-input v-if="editable" v-model:value="data.og_title" style="margin: -5px 0" />
          <template v-else>
            {{ data.og_title }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="description - 描述">
        <div>
          <a-input v-if="editable" v-model:value="data.description" style="margin: -5px 0" />
          <template v-else>
            {{ data.description }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="keywords - 关键字">
        <div>
          <a-input v-if="editable" v-model:value="data.keywords" style="margin: -5px 0" />
          <template v-else>
            {{ data.keywords }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="author - 作者">
        <div>
          <a-input v-if="editable" v-model:value="data.author" style="margin: -5px 0" />
          <template v-else>
            {{ data.author }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="robots - 搜索引擎指令">
        <div>
          <a-input v-if="editable" v-model:value="data.robots" style="margin: -5px 0" />
          <template v-else>
            {{ data.robots }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="分享封面">
        <StaticUpload
          v-if="editable"
          :image-url="data.og_image"
          :authorization="userStore.token"
          @update:imageUrl="(value) => (data.og_image = value)"
        />
        <a-image v-else :width="200" :src="serverHost + data.og_image" />
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
  <div></div>
</template>
<script lang="ts" setup>
import { GetSeo, type SeoConfig, UpdateSeo } from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import type { UploadChangeParam, UploadProps } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'

const userStore = useUserStore()

const editable = ref<boolean>(false)

const data = ref<SeoConfig>({
  title: '',
  description: '',
  og_title: '',
  og_image: '',
  keywords: '',
  author: '',
  robots: ''
})
const serverHost = import.meta.env.VITE_API_HOST
const getSeo = async () => {
  try {
    const response = await GetSeo()
    data.value = response.data.data || {}
  } catch (error) {
    console.log(error)
  }
}
getSeo()

const cancel = () => {
  editable.value = false
  getSeo()
}

const save = async () => {
  try {
    const response: any = await UpdateSeo({
      title: data.value.title,
      description: data.value.description,
      og_title: data.value.og_title,
      og_image: data.value.og_image,
      keywords: data.value.keywords,
      author: data.value.author,
      robots: data.value.robots
    })
    if (response.data.code === 0) {
      message.success('保存成功')
      await getSeo()
      editable.value = false
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
