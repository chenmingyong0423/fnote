<template>
  <div>
    <div class="flex h-15 items-center">
      <a-modal
        v-model:open="open"
        title="温馨提示"
        @ok="handleOk"
        :cancelText="'取消'"
        :okText="'确认'"
      >
        <p>检测到没有自定义文章 id，保存草稿之后将会自动生成且后续无法修改，是否继续保存？</p>
      </a-modal>
      <a-input v-model:value="post4Edit.title" addon-before="标题" class="w-59%" />
      <a-input v-model:value="post4Edit.author" addon-before="作者" class="w-30% ml-1%" />
      <a-button type="primary" @click="visible = true" class="w-9% ml-1%"
        >{{ props.isNewPost ? '发布' : '更新' }}
      </a-button>
      <a-button type="primary" @click="preSave" class="w-9% ml-1%">保存草稿</a-button>
      <a-modal
        v-model:open="visible"
        title="文章元数据"
        ok-text="提交"
        cancel-text="取消"
        @ok="submit"
      >
        <a-form ref="formRef" :model="post4Edit" name="form_in_modal">
          <a-form-item
            name="title"
            label="标题"
            :rules="[{ required: true, message: '请输入标题' }]"
          >
            {{ post4Edit.title }}
          </a-form-item>
          <a-form-item
            name="author"
            label="作者"
            :rules="[{ required: true, message: '请输入作者' }]"
          >
            {{ post4Edit.author }}
          </a-form-item>
          <a-form-item name="id" label="自定义 id">
            <a-input
              v-model:value="post4Edit.id"
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
              v-model:value="post4Edit.tempCategories"
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
              v-model:value="post4Edit.tempTags"
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
              :image-url="post4Edit.cover_img"
              @update:imageUrl="(value) => (post4Edit.cover_img = value)"
              :authorization="userStore.token"
            />
          </a-form-item>
          <a-form-item
            name="is_comment_allowed"
            label="开启评论"
            :rules="[{ required: true, message: '请设置评论开关' }]"
          >
            <a-radio-group v-model:value="post4Edit.is_comment_allowed" name="radioGroup">
              <a-radio :value="false">否</a-radio>
              <a-radio :value="true">是</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="sticky_weight"
            label="置顶状态"
            :rules="[{ required: true, message: '请选择置顶状态' }]"
          >
            <a-radio-group v-model:value="post4Edit.sticky_weight" name="radioGroup">
              <a-radio :value="0">否</a-radio>
              <a-radio :value="1">是</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="is_displayed"
            label="文章状态"
            :rules="[{ required: true, message: '请选择状态' }]"
          >
            <a-radio-group v-model:value="post4Edit.is_displayed" name="radioGroup">
              <a-radio :value="false">隐藏</a-radio>
              <a-radio :value="true">显示</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item
            name="summary"
            label="文章摘要"
            :rules="[{ required: true, message: '请输入摘要' }]"
          >
            <a-textarea v-model:value="post4Edit.summary" placeholder="请输入摘要" allow-clear />
          </a-form-item>
          <a-form-item name="meta_description" label="seo description">
            <a-textarea
              v-model:value="post4Edit.meta_description"
              placeholder="请输入描述"
              allow-clear
            />
          </a-form-item>
          <a-form-item name="meta_keywords" label="seo keywords">
            <a-input v-model:value="post4Edit.meta_keywords" placeholder="请输入关键字" />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
    <div>
      <v-md-editor
        v-model="post4Edit.content"
        height="800px"
        :disabled-menus="[]"
        @upload-image="handleUploadImage"
        @save="preSave"
        left-toolbar="undo redo clear | h bold italic strikethrough quote | ul ol table hr | link image code | save | template"
        :toolbar="toolbar"
      />
    </div>
    <a-modal
      v-model:visible="visible4Template"
      width="1000px"
      title="图片素材"
      @ok="handleOk4Template"
      :footer="null"
      :destroyOnClose="true"
    >
      <ImageLIstView @insertImg="insertImg" />
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { type PropType, reactive, ref, defineEmits } from 'vue'
import type { Post4Edit } from '@/interfaces/Post'
import { type FormInstance, message } from 'ant-design-vue'
import type { SelectCategory } from '@/interfaces/Category'
import type { SelectTag } from '@/interfaces/Tag'
import { FileUpload } from '@/interfaces/File'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import ImageLIstView from '@/views/post/editor/ImageLIstView.vue'

const emit = defineEmits(['publish', 'saveDraft'])
const userStore = useUserStore()

const props = defineProps({
  post: {
    type: Object as PropType<Post4Edit>,
    default: () => {
      return {
        is_displayed: true,
        sticky_weight: 0,
        is_comment_allowed: true
      }
    }
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
const post4Edit = reactive<Post4Edit>(props.post || ({} as Post4Edit))
const submit = () => {
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        if (post4Edit.content === '') {
          message.warning('请填写文章内容')
          return
        }
        post4Edit.categories = []
        values.tempCategories.forEach((item: string) => {
          props.categories.forEach((category) => {
            if (category.value === item) {
              post4Edit.categories.push({
                id: category.id,
                name: category.value
              })
            }
          })
        })
        post4Edit.tags = []
        values.tempTags.forEach((item: string) => {
          props.tags.forEach((tag) => {
            if (tag.value === item) {
              post4Edit.tags.push({
                id: tag.id,
                name: tag.value
              })
            }
          })
        })
        // 告诉父组件
        emit('publish', post4Edit)
      })
      .catch((info) => {
        console.log('Validate Failed:', info)
        message.warning('请检查表单是否填写正确')
      })
  }
}

const open = ref<boolean>(false)

const handleOk = () => {
  saveDraft()
  open.value = false
}

const preSave = () => {
  if (!post4Edit.id || post4Edit.id === '') {
    open.value = true
  } else {
    saveDraft()
  }
}

const saveDraft = () => {
  if (!!post4Edit.title && !!post4Edit.author && !!post4Edit.content) {
    post4Edit.categories = []
    post4Edit.tempCategories?.forEach((item: string) => {
      props.categories.forEach((category) => {
        if (category.value === item) {
          post4Edit.categories.push({
            id: category.id,
            name: category.value
          })
        }
      })
    })
    post4Edit.tags = []
    post4Edit.tempTags?.forEach((item: string) => {
      props.tags.forEach((tag) => {
        if (tag.value === item) {
          post4Edit.tags.push({
            id: tag.id,
            name: tag.value
          })
        }
      })
    })
    // 告诉父组件
    emit('saveDraft', post4Edit)
  } else {
    message.warning('保存草稿时，标题和作者以及内容必填。')
  }
}

const clearReq = () => {
  if (formRef.value) {
    formRef.value.resetFields()
    post4Edit.title = ''
    post4Edit.author = ''
    post4Edit.content = ''
    imageUrl.value = ''
    post4Edit.categories = []
    post4Edit.tags = []
    post4Edit.tempCategories = []
    post4Edit.tempTags = []
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
      const res: any = await FileUpload(formData)
      if (res.data.code !== 0) {
        message.error(res.data.message)
        return
      }
      insertImage({
        url: res.data.data.url,
        desc: '请在此添加图片描述'
      })
    } catch (error) {
      message.error(error)
    }
  } catch (error) {
    console.log(error)
  }
}

const toolbar = {
  template: {
    title: '模板',
    icon: 'v-md-icon-tip',
    menus: [
      {
        name: 'personal-images',
        text: '图片素材',
        action(editor: any) {
          visible4Template.value = true
          globalEditor.value = editor
        }
      }
    ]
  }
}

const visible4Template = ref<boolean>(false)

const handleOk4Template = (e: MouseEvent) => {
  console.log(e)
  visible4Template.value = false
  globalEditor.value = null
}

const globalEditor = ref<any>(null)

const insertImg = (content: string) => {
  // @ts-ignore
  globalEditor.value.insert(function () {
    return {
      text: content,
      selected: content
    }
  })
  visible4Template.value = false
  globalEditor.value = null
}
</script>
