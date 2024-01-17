<template>
  <div>
    <div class="flex h-15 items-center">
      <a-input v-model:value="postReq.title" addon-before="标题" class="w-59%" />
      <a-input v-model:value="postReq.author" addon-before="作者" class="w-30% ml-1%" />
      <a-button type="primary" @click="visible = true" class="w-9% ml-1%">发布</a-button>
      <a-modal
        v-model:open="visible"
        title="文章元数据"
        ok-text="提交"
        cancel-text="取消"
        @ok="submit"
      >
        <a-form ref="formRef" :model="postReq" name="form_in_modal">
          <a-form-item
            name="title"
            label="标题"
            :rules="[{ required: true, message: '请输入标题' }]"
          >
            {{ postReq.title }}
          </a-form-item>
          <a-form-item
            name="author"
            label="作者"
            :rules="[{ required: true, message: '请输入作者' }]"
          >
            {{ postReq.author }}
          </a-form-item>
          <a-form-item name="id" label="自定义 id">
            <a-input
              v-model:value="postReq.id"
              :disabled="!props.isNewPost"
              placeholder="与文章关联的英文的 id 有助于 seo 优化"
            />
          </a-form-item>
          <a-form-item
            name="tempCategories"
            label="分类"
            :rules="[{ required: true, message: '请选择分类' }]"
          >
            <a-select
              v-model:value="postReq.tempCategories"
              mode="multiple"
              style="width: 100%"
              placeholder="请选择分类"
              :options="categories"
            ></a-select>
          </a-form-item>
          <a-form-item
            name="tempTags"
            label="标签"
            :rules="[{ required: true, message: '请选择标签' }]"
          >
            <a-select
              v-model:value="postReq.tempTags"
              mode="multiple"
              style="width: 100%"
              placeholder="请选择标签"
              :options="tags"
            ></a-select>
          </a-form-item>
          <a-form-item
            name="cover_img"
            label="封面"
            :rules="[{ required: true, message: '请选择封面' }]"
          >
            <a-input v-model:value="postReq.cover_img" placeholder="请输入封面路径" />
            <a-upload
              v-model:file-list="fileList"
              name="file"
              list-type="picture-card"
              class="avatar-uploader m-y-5"
              :show-upload-list="false"
              action="http://localhost:8080/admin/files/upload"
              :before-upload="beforeUpload"
              @change="handleChange"
            >
              <img v-if="imageUrl" :src="imageUrl" alt="avatar" width="250" height="150" />
              <div v-else>
                <loading-outlined v-if="loading"></loading-outlined>
                <plus-outlined v-else></plus-outlined>
                <div class="ant-upload-text">Upload</div>
              </div>
            </a-upload>
          </a-form-item>
          <a-form-item
            name="is_comment_allowed"
            label="开启评论"
            :rules="[{ required: true, message: '请设置评论开关' }]"
          >
            <a-radio-group v-model:value="postReq.is_comment_allowed" name="radioGroup">
              <a-radio :value="false">否</a-radio>
              <a-radio :value="true">是</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="sticky_weight"
            label="置顶状态"
            :rules="[{ required: true, message: '请选择置顶状态' }]"
          >
            <a-radio-group v-model:value="postReq.sticky_weight" name="radioGroup">
              <a-radio :value="0">否</a-radio>
              <a-radio :value="1">是</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="status"
            label="文章状态"
            :rules="[{ required: true, message: '请选择状态' }]"
          >
            <a-radio-group v-model:value="postReq.status" name="radioGroup">
              <a-radio :value="0">草稿</a-radio>
              <a-radio :value="1">已发表</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="summary"
            label="文章摘要"
            :rules="[{ required: true, message: '请输入摘要' }]"
          >
            <a-textarea v-model:value="postReq.summary" placeholder="请输入摘要" allow-clear />
          </a-form-item>
          <a-form-item name="meta_description" label="seo description">
            <a-textarea
              v-model:value="postReq.meta_description"
              placeholder="请输入描述"
              allow-clear
            />
          </a-form-item>
          <a-form-item name="meta_keywords" label="seo keywords">
            <a-input v-model:value="postReq.meta_keywords" placeholder="请输入关键字" />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
    <div>
      <v-md-editor
        v-model="postReq.content"
        height="800px"
        :disabled-menus="[]"
        @upload-image="handleUploadImage"
      ></v-md-editor>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { type PropType, reactive, ref, defineEmits } from 'vue'
import type { PostRequest } from '@/interfaces/Post'
import {
  type FormInstance,
  message,
  type SelectProps,
  type UploadChangeParam,
  type UploadProps
} from 'ant-design-vue'
import axios from '@/http/axios'
import type { IResponse } from '@/interfaces/Common'
import type { SelectCategory } from '@/interfaces/Category'
import type { SelectTag } from '@/interfaces/Tag'
import type { File } from '@/interfaces/File'

const emit = defineEmits(['submit'])

const props = defineProps({
  req: {
    type: Object as PropType<PostRequest>,
    required: true
  },
  categories: {
    type: Array as PropType<SelectCategory[]>,
    default: () => []
  },
  tags: {
    type: Array as PropType<SelectTag[]>,
    default: () => []
  },
  isNewPost: {
    type: Boolean,
    default: true
  }
})

const imageUrl = ref<string>('')

const formRef = ref<FormInstance>()
const visible = ref(false)
const postReq = reactive<PostRequest>(props.req)

const categories = ref<SelectProps['options']>(props.categories)

const tags = ref<SelectProps['options']>(props.tags)

const submit = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        if (postReq.content === '') {
          message.warning('请填写文章内容')
          return
        }
        postReq.categories = []
        postReq.tags = []
        values.tempCategories.forEach((item: string) => {
          categories.value?.forEach((category) => {
            if (category.value === item) {
              postReq.categories.push({
                id: category.id,
                name: category.value
              })
            }
          })
        })
        values.tempTags.forEach((item: string) => {
          tags.value?.forEach((tag) => {
            if (tag.value === item) {
              postReq.tags.push({
                id: tag.id,
                name: tag.value
              })
            }
          })
        })
        // 告诉父组件
        emit('submit', postReq)
      })
      .catch((info) => {
        console.log('Validate Failed:', info)
        message.warning('请检查表单是否填写正确')
      })
  }
}

const clearReq = () => {
  if (formRef.value) {
    formRef.value.resetFields()
    postReq.title = ''
    postReq.author = ''
    postReq.content = ''
    imageUrl.value = ''
    postReq.categories = []
    postReq.tags = []
  }
  visible.value = false
}

defineExpose({
  clearReq
})

// 封面上传
function getBase64(img: Blob, callback: (base64Url: string) => void) {
  const reader = new FileReader()
  reader.addEventListener('load', () => callback(reader.result as string))
  reader.readAsDataURL(img)
}

const fileList = ref([])
const loading = ref<boolean>(false)

const handleChange = (info: UploadChangeParam) => {
  if (info.file.status === 'uploading') {
    loading.value = true
    return
  }
  if (info.file.status === 'done') {
    // Get this url from response in real world.
    getBase64(info.file.originFileObj!, (base64Url: string) => {
      imageUrl.value = base64Url
      loading.value = false
      // Get this url from response in real world.
      postReq.cover_img = info.file.response.data.url
      message.success('上传成功')
    })
  }
  if (info.file.status === 'error') {
    loading.value = false
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

// md 图片上传
const handleUploadImage = async (event, insertImage, files) => {
  console.log(files)
  try {
    const formData = new FormData()
    formData.append('file', files[0])

    const response = await axios.post<IResponse<File>>('/admin/files/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    insertImage({
      url: response?.data?.data?.url,
      desc: response?.data?.data?.file_name
      // width: 'auto',
      // height: 'auto',
    })
  } catch (error) {
    console.log(error)
  }
}
</script>
