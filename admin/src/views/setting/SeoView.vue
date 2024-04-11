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
      <a-descriptions-item label="baidu_site_verification - 百度站点验证">
        <div>
          <a-input
            v-if="editable"
            v-model:value="data.baidu_site_verification"
            style="margin: -5px 0"
          />
          <template v-else>
            {{ data.baidu_site_verification }}
          </template>
        </div>
      </a-descriptions-item>
      <a-descriptions-item label="分享封面">
        <a-image :width="200" :src="data.og_image" />
        <a-upload
          v-if="editable"
          v-model:file-list="fileList"
          name="file"
          action="http://localhost:8080/admin/files/upload"
          @change="handleChange"
          :before-upload="beforeUpload"
          :headers="{ Authorization: userStore.token }"
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
import { GetSeo, type SeoConfig, UpdateSeo } from '@/interfaces/Config'
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import type { UploadChangeParam, UploadProps } from 'ant-design-vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

const editable = ref<boolean>(false)

const data = ref<SeoConfig>({
  title: '',
  description: '',
  og_title: '',
  og_image: '',
  baidu_site_verification: '',
  keywords: '',
  author: '',
  robots: ''
})

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
      baidu_site_verification: data.value.baidu_site_verification,
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

// 文件操作
const fileList = ref<UploadProps['fileList']>([])

const handleChange = (info: UploadChangeParam) => {
  if (info.file.status === 'uploading') {
    return
  }
  console.log(info)
  if (info.file.status === 'done') {
    // Get this url from response in real world.
    data.value.og_image = info.file.response.data.data.url
    message.success('上传成功')
  }
  if (info.file.status === 'error') {
    message.error('upload error')
  }
}

// const beforeUpload = (file: UploadProps['fileList'][number]) => { = -!无力吐槽， 官网的写法，ts 无法保证类型安全
const beforeUpload = (file: any) => {
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
