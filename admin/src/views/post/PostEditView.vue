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
              :options="props.categories"
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
              :options="props.tags"
            ></a-select>
          </a-form-item>
          <a-form-item
            name="cover_img"
            label="封面"
            :rules="[{ required: true, message: '请选择封面' }]"
          >
            <StaticUpload
              :image-url="postReq.cover_img"
              @update:imageUrl="(value) => (postReq.cover_img = value)"
              :authorization="userStore.token"
            />
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
            name="is_displayed"
            label="文章状态"
            :rules="[{ required: true, message: '请选择状态' }]"
          >
            <a-radio-group v-model:value="postReq.is_displayed" name="radioGroup">
              <a-radio :value="false">隐藏</a-radio>
              <a-radio :value="true">显示</a-radio>
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
import axios from '@/utils/axios'
import type { IResponse } from '@/interfaces/Common'
import type { SelectCategory } from '@/interfaces/Category'
import type { SelectTag } from '@/interfaces/Tag'
import { type File, FileUpload } from '@/interfaces/File'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'

const emit = defineEmits(['submit'])
const userStore = useUserStore()

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
        values.tempCategories.forEach((item: string) => {
          props.categories.forEach((category) => {
            if (category.value === item) {
              postReq.categories.push({
                id: category.id,
                name: category.value
              })
            }
          })
        })
        postReq.tags = []
        values.tempTags.forEach((item: string) => {
          props.tags.forEach((tag) => {
            if (tag.value === item) {
              postReq.tags.push({
                id: tag.id,
                name: tag.value
              })
            }
          })
        })
        console.log(postReq)
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
    postReq.tempCategories = []
    postReq.tempTags = []
  }
  visible.value = false
}

defineExpose({
  clearReq
})

// md 图片上传
const handleUploadImage = async (event: any, insertImage: any, files: any) => {
  try {

    const formData = new FormData()
    formData.append('file', files[0])
    try {
      const res : any = await FileUpload(formData)
      if (res.data.code !== 0) {
        message.error(res.data.message)
        return
      }
      insertImage({
        url: import.meta.env.VITE_API_HOST + res.data.data.url,
        desc: '请在此添加图片描述',
      })
    } catch (error) {
      message.error(error)
    }
  } catch (error) {
    console.log(error)
  }
}
</script>
