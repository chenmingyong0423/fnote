<template>
  <div>
    <div class="flex h-15 items-center gap-x-2">
      <a-modal
        v-model:open="open"
        title="温馨提示"
        @ok="handleOk"
        :cancelText="'取消'"
        :okText="'确认'"
        :confirm-loading="props.savingDraft"
      >
        <p>检测到没有自定义文章 id，保存草稿之后将会自动生成且后续无法修改，是否继续保存？</p>
      </a-modal>
      <a-input v-model:value="post4Edit.title" addon-before="标题" class="min-w-0 flex-1" />
      <a-input v-model:value="post4Edit.author" addon-before="作者" class="w-56" />
      <a-button
        type="primary"
        @click="visible = true"
        :loading="props.publishing"
        :disabled="props.savingDraft"
        >{{ props.isNewPost ? '发布' : '更新' }}
      </a-button>
      <a-button
        type="primary"
        @click="preSave"
        :loading="props.savingDraft"
        :disabled="props.publishing"
        >保存草稿</a-button
      >
      <span v-if="props.lastSavedAt" class="save-status">
        最近保存 {{ formatSaveTime(props.lastSavedAt) }}
      </span>
      <span v-if="hasUnsavedChanges" class="save-status">有未保存更改</span>
      <span v-if="contentStatsText" class="save-status">{{ contentStatsText }}</span>
      <a-drawer
        v-model:open="visible"
        title="文章元数据"
        placement="right"
        width="720"
        :body-style="{ paddingBottom: '80px' }"
      >
        <template #extra>
          <a-space>
            <a-button @click="visible = false">取消</a-button>
            <a-button type="primary" :loading="props.publishing" @click="submit">提交</a-button>
          </a-space>
        </template>
        <a-form ref="formRef" :model="post4Edit" layout="vertical" name="form_in_drawer">
          <section class="metadata-section">
            <h3 class="metadata-section-title">基础信息</h3>
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
            <a-form-item name="id" label="自定义 id" :rules="postIdRules" :extra="postIdExtra">
              <div class="field-with-action">
                <a-input
                  v-model:value="post4Edit.id"
                  :disabled="!props.isNewPost"
                  placeholder="例如 vue-router-note"
                />
                <a-button type="link" size="small" :disabled="!props.isNewPost" @click="fillPostId">
                  根据标题生成
                </a-button>
              </div>
            </a-form-item>
          </section>

          <section class="metadata-section">
            <h3 class="metadata-section-title">分类与封面</h3>
            <a-form-item
              name="tempCategories"
              label="分类"
              :rules="[{ required: true, message: '请选择分类' }]"
            >
              <div class="taxonomy-tools">
                <a-alert
                  v-if="categoryOptions.length === 0"
                  message="当前还没有分类，可以先快速创建后继续发布。"
                  type="info"
                  show-icon
                />
                <a-button type="link" size="small" @click="openQuickCategory">+ 新建分类</a-button>
              </div>
              <a-select
                v-model:value="post4Edit.tempCategories"
                mode="multiple"
                show-search
                style="width: 100%"
                placeholder="请选择分类"
                :options="categoryOptions"
              ></a-select>
            </a-form-item>
            <a-form-item
              name="tempTags"
              label="标签"
              :rules="[{ required: true, message: '请选择标签' }]"
            >
              <div class="taxonomy-tools">
                <a-alert
                  v-if="tagOptions.length === 0"
                  message="当前还没有标签，可以先快速创建后继续发布。"
                  type="info"
                  show-icon
                />
                <a-button type="link" size="small" @click="openQuickTag">+ 新建标签</a-button>
              </div>
              <a-select
                v-model:value="post4Edit.tempTags"
                mode="multiple"
                show-search
                style="width: 100%"
                placeholder="请选择标签"
                :options="tagOptions"
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
          </section>

          <section class="metadata-section">
            <h3 class="metadata-section-title">发布设置</h3>
            <div class="metadata-settings-grid">
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
            </div>
          </section>

          <section class="metadata-section">
            <h3 class="metadata-section-title">SEO 信息</h3>
            <a-form-item
              name="summary"
              label="文章摘要"
              :rules="[{ required: true, message: '请输入摘要' }]"
              :extra="summaryExtra"
            >
              <div class="field-with-action">
                <a-textarea
                  v-model:value="post4Edit.summary"
                  placeholder="请输入摘要"
                  allow-clear
                />
                <a-button type="link" size="small" @click="fillSummary">从正文生成</a-button>
              </div>
            </a-form-item>
            <a-form-item
              name="meta_description"
              label="seo description"
              :extra="metaDescriptionExtra"
            >
              <div class="field-with-action">
                <a-textarea
                  v-model:value="post4Edit.meta_description"
                  placeholder="请输入描述"
                  allow-clear
                />
                <a-button type="link" size="small" @click="fillMetaDescription">
                  同步摘要
                </a-button>
              </div>
            </a-form-item>
            <a-form-item name="meta_keywords" label="seo keywords">
              <div class="field-with-action">
                <a-input v-model:value="post4Edit.meta_keywords" placeholder="请输入关键字" />
                <a-button type="link" size="small" @click="fillMetaKeywords">
                  从分类标签生成
                </a-button>
              </div>
            </a-form-item>
          </section>
        </a-form>
      </a-drawer>
      <a-modal
        v-model:open="quickCategoryVisible"
        title="新建分类"
        ok-text="创建"
        cancel-text="取消"
        :confirm-loading="quickCategoryLoading"
        @ok="createCategory"
      >
        <TaxonomyCreateForm
          ref="quickCategoryFormRef"
          v-model:name="quickCategoryForm.name"
          v-model:route="quickCategoryForm.route"
          v-model:description="quickCategoryForm.description"
          v-model:enabled="quickCategoryForm.enabled"
          v-model:show-in-nav="quickCategoryForm.show_in_nav"
          type="category"
        />
      </a-modal>
      <a-modal
        v-model:open="quickTagVisible"
        title="新建标签"
        ok-text="创建"
        cancel-text="取消"
        :confirm-loading="quickTagLoading"
        @ok="createTag"
      >
        <TaxonomyCreateForm
          ref="quickTagFormRef"
          v-model:name="quickTagForm.name"
          v-model:route="quickTagForm.route"
          v-model:enabled="quickTagForm.enabled"
          type="tag"
        />
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
    <a-modal
      v-model:visible="textTemplateVisible"
      width="1000px"
      title="文字模板"
      :footer="null"
      :destroyOnClose="true"
      @cancel="closeTextTemplate"
    >
      <TextTemplateListView @insertTemplate="insertTextTemplate" />
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import {
  type PropType,
  computed,
  reactive,
  ref,
  defineEmits,
  onBeforeUnmount,
  onMounted,
  watch
} from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import type { Post4Edit } from '@/interfaces/Post'
import { type FormInstance, message } from 'ant-design-vue'
import {
  AddCategory,
  GetSelectedCategories,
  type CategoryRequest,
  type SelectCategory
} from '@/interfaces/Category'
import { AddTag, GetSelectedTags, type SelectTag, type TagRequest } from '@/interfaces/Tag'
import { FileUpload } from '@/interfaces/File'
import { useUserStore } from '@/stores/user'
import StaticUpload from '@/components/upload/StaticUpload.vue'
import ImageLIstView from '@/views/post/editor/ImageLIstView.vue'
import TextTemplateListView from '@/views/post/editor/TextTemplateListView.vue'
import { toErrorMessage } from '@/utils/error'
import TaxonomyCreateForm from '@/components/form/TaxonomyCreateForm.vue'
import { normalizeSlug, shouldAutoFillSlug, validateSlug } from '@/utils/slug'

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
  },
  publishing: {
    type: Boolean,
    default: false
  },
  savingDraft: {
    type: Boolean,
    default: false
  },
  autoSave: {
    type: Boolean,
    default: false
  },
  lastSavedAt: {
    type: Number,
    default: 0
  },
  baselineKey: {
    type: Number,
    default: 0
  }
})

const imageUrl = ref<string>('')

const formRef = ref<FormInstance>()
const visible = ref(false)
const post4Edit = reactive<Post4Edit>(props.post || ({} as Post4Edit))
const categoryOptions = ref<SelectCategory[]>([])
const tagOptions = ref<SelectTag[]>([])
const autoSaveTimer = ref<ReturnType<typeof setTimeout> | null>(null)
const savedSnapshot = ref('')
const hasUnsavedChanges = ref(false)

watch(
  () => props.categories,
  (categories) => {
    categoryOptions.value = [...categories]
  },
  { immediate: true, deep: true }
)

watch(
  () => props.tags,
  (tags) => {
    tagOptions.value = [...tags]
  },
  { immediate: true, deep: true }
)

const submit = () => {
  if (props.publishing || props.savingDraft) {
    return
  }
  if (formRef.value) {
    formRef.value
      .validateFields()
      .then(async (values) => {
        if (post4Edit.content === '') {
          message.warning('请填写文章内容')
          return
        }
        syncWordCount()
        post4Edit.categories = []
        values.tempCategories.forEach((item: string) => {
          categoryOptions.value.forEach((category) => {
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
          tagOptions.value.forEach((tag) => {
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
  if (props.publishing || props.savingDraft) {
    return
  }
  if (!hasDraftTitle()) {
    message.warning('保存草稿前请先填写标题')
    return
  }
  if (!post4Edit.id || post4Edit.id === '') {
    open.value = true
  } else {
    saveDraft()
  }
}

const saveDraft = (options?: { silent?: boolean }) => {
  if (props.publishing || props.savingDraft) {
    return
  }
  if (hasDraftTitle()) {
    syncWordCount()
    post4Edit.categories = []
    post4Edit.tempCategories?.forEach((item: string) => {
      categoryOptions.value.forEach((category) => {
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
      tagOptions.value.forEach((tag) => {
        if (tag.value === item) {
          post4Edit.tags.push({
            id: tag.id,
            name: tag.value
          })
        }
      })
    })
    // 告诉父组件
    emit('saveDraft', post4Edit, options)
  } else {
    message.warning('保存草稿前请先填写标题')
  }
}

const hasDraftTitle = () => {
  return typeof post4Edit.title === 'string' && post4Edit.title.trim().length > 0
}

const stripMarkdown = (content: string) => {
  return content
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`([^`]+)`/g, '$1')
    .replace(/!\[[^\]]*]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]+)]\([^)]*\)/g, '$1')
    .replace(/<[^>]+>/g, ' ')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^>\s?/gm, '')
    .replace(/^[\s>*+-]*\d+\.\s+/gm, '')
    .replace(/^[\s>*+-]+/gm, '')
    .replace(/[*_~>#|[\]()]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

const countReadableWords = (content: string) => {
  const chineseCount = content.match(/[\u4e00-\u9fa5]/g)?.length || 0
  const latinContent = content.replace(/[\u4e00-\u9fa5]/g, ' ')
  const latinWordCount = latinContent.match(/[a-zA-Z0-9]+(?:[-'][a-zA-Z0-9]+)*/g)?.length || 0
  return chineseCount + latinWordCount
}

const readableContent = computed(() => stripMarkdown(post4Edit.content || ''))

const contentWordCount = computed(() => countReadableWords(readableContent.value))

const readingMinutes = computed(() => {
  if (contentWordCount.value === 0) {
    return 0
  }
  return Math.max(1, Math.ceil(contentWordCount.value / 400))
})

const contentStatsText = computed(() => {
  if (contentWordCount.value === 0) {
    return ''
  }
  return `字数 ${contentWordCount.value} · 约 ${readingMinutes.value} 分钟`
})

const summaryLength = computed(() => post4Edit.summary?.trim().length || 0)

const metaDescriptionLength = computed(() => post4Edit.meta_description?.trim().length || 0)

const summaryExtra = computed(() => {
  if (summaryLength.value === 0) {
    return '建议 80-160 字，便于列表页和搜索摘要展示'
  }
  return `当前 ${summaryLength.value} 字，建议 80-160 字`
})

const metaDescriptionExtra = computed(() => {
  if (metaDescriptionLength.value === 0) {
    return '建议 120-160 字，留空时可同步摘要'
  }
  return `当前 ${metaDescriptionLength.value} 字，建议 120-160 字`
})

const postIdExtra = computed(() => {
  if (!props.isNewPost) {
    return '文章 id 发布后不可修改'
  }
  return '建议使用小写英文、数字和短横线，留空保存草稿时会自动生成'
})

const postIdRules = [
  {
    validator: async (_rule: unknown, value: string) => {
      if (!props.isNewPost || !value) {
        return
      }
      validateSlug(value, '自定义 id')
    }
  }
]

const syncWordCount = () => {
  post4Edit.word_count = contentWordCount.value
}

const truncateText = (content: string, maxLength: number) => {
  if (content.length <= maxLength) {
    return content
  }
  return `${content.slice(0, maxLength).trim()}...`
}

const fillSummary = () => {
  const content = stripMarkdown(post4Edit.content || '')
  if (!content) {
    message.warning('正文为空，无法生成摘要')
    return
  }
  post4Edit.summary = truncateText(content, 160)
}

const fillMetaDescription = () => {
  const content = post4Edit.summary?.trim() || stripMarkdown(post4Edit.content || '')
  if (!content) {
    message.warning('摘要和正文为空，无法生成 SEO 描述')
    return
  }
  post4Edit.meta_description = truncateText(content, 160)
}

const fillMetaKeywords = () => {
  const keywords = [
    ...(post4Edit.tempCategories || []),
    ...(post4Edit.tempTags || []),
    ...(post4Edit.categories || []).map((category) => category.name),
    ...(post4Edit.tags || []).map((tag) => tag.name)
  ]
    .map((keyword) => keyword.trim())
    .filter(Boolean)

  const uniqueKeywords = Array.from(new Set(keywords))
  if (uniqueKeywords.length === 0) {
    message.warning('请先选择分类或标签')
    return
  }
  post4Edit.meta_keywords = uniqueKeywords.join(',')
}

const fillPostId = () => {
  const title = post4Edit.title?.trim() || ''
  if (!title) {
    message.warning('请先填写标题')
    return
  }
  if (!shouldAutoFillSlug(title)) {
    message.warning('中文标题请手动填写自定义 id')
    return
  }
  const id = normalizeSlug(title)
  if (!id) {
    message.warning('当前标题无法生成有效 id')
    return
  }
  post4Edit.id = id
}

const formatSaveTime = (timestamp: number) => {
  return new Date(timestamp).toLocaleTimeString('zh-CN', {
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const clearAutoSaveTimer = () => {
  if (autoSaveTimer.value) {
    clearTimeout(autoSaveTimer.value)
    autoSaveTimer.value = null
  }
}

const scheduleAutoSave = () => {
  clearAutoSaveTimer()
  if (
    !props.autoSave ||
    props.isNewPost ||
    props.publishing ||
    props.savingDraft ||
    !post4Edit.id ||
    !hasDraftTitle()
  ) {
    return
  }

  autoSaveTimer.value = setTimeout(() => {
    saveDraft({ silent: true })
  }, 8000)
}

const getPostSnapshot = () => {
  return JSON.stringify({
    id: post4Edit.id || '',
    title: post4Edit.title || '',
    author: post4Edit.author || '',
    summary: post4Edit.summary || '',
    content: post4Edit.content || '',
    cover_img: post4Edit.cover_img || '',
    meta_description: post4Edit.meta_description || '',
    meta_keywords: post4Edit.meta_keywords || '',
    is_comment_allowed: post4Edit.is_comment_allowed,
    is_displayed: post4Edit.is_displayed,
    sticky_weight: post4Edit.sticky_weight,
    tempCategories: [...(post4Edit.tempCategories || [])].sort(),
    tempTags: [...(post4Edit.tempTags || [])].sort()
  })
}

const markSavedSnapshot = () => {
  savedSnapshot.value = getPostSnapshot()
  hasUnsavedChanges.value = false
}

const updateUnsavedState = () => {
  if (!savedSnapshot.value) {
    markSavedSnapshot()
    return
  }
  hasUnsavedChanges.value = savedSnapshot.value !== getPostSnapshot()
}

watch(
  () => [
    post4Edit.id,
    post4Edit.title,
    post4Edit.author,
    post4Edit.summary,
    post4Edit.content,
    post4Edit.cover_img,
    post4Edit.meta_description,
    post4Edit.meta_keywords,
    post4Edit.is_comment_allowed,
    post4Edit.is_displayed,
    post4Edit.sticky_weight,
    JSON.stringify(post4Edit.tempCategories || []),
    JSON.stringify(post4Edit.tempTags || [])
  ],
  () => {
    updateUnsavedState()
    scheduleAutoSave()
  }
)

watch(
  () => props.lastSavedAt,
  (lastSavedAt) => {
    if (lastSavedAt) {
      markSavedSnapshot()
    }
  }
)

watch(
  () => props.baselineKey,
  () => {
    markSavedSnapshot()
  }
)

const confirmLeave = () => {
  if (!hasUnsavedChanges.value || props.savingDraft || props.publishing) {
    return true
  }
  return window.confirm('当前文章有未保存更改，确认离开吗？')
}

const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  if (!hasUnsavedChanges.value || props.savingDraft || props.publishing) {
    return
  }
  event.preventDefault()
  event.returnValue = ''
}

onMounted(() => {
  markSavedSnapshot()
  window.addEventListener('beforeunload', handleBeforeUnload)
})

onBeforeUnmount(() => {
  clearAutoSaveTimer()
  window.removeEventListener('beforeunload', handleBeforeUnload)
})

onBeforeRouteLeave(() => {
  return confirmLeave()
})

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
  markSavedSnapshot()
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
      message.error(toErrorMessage(error))
    }
  } catch (error) {
    console.log(error)
  }
}

const toolbar = {
  template: {
    title: '素材库',
    icon: 'v-md-icon-toc',
    menus: [
      {
        name: 'image-assets',
        text: '图片素材',
        action(editor: any) {
          visible4Template.value = true
          globalEditor.value = editor
        }
      },
      {
        name: 'text-templates',
        text: '文字模板',
        action(editor: any) {
          textTemplateVisible.value = true
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

const textTemplateVisible = ref(false)

const closeTextTemplate = () => {
  textTemplateVisible.value = false
  globalEditor.value = null
}

const insertTextTemplate = (content: string) => {
  if (!globalEditor.value) return
  globalEditor.value.insert(() => ({
    text: content,
    selected: content
  }))
  closeTextTemplate()
}

const quickCategoryVisible = ref(false)
const quickCategoryLoading = ref(false)
const quickCategoryFormRef = ref<InstanceType<typeof TaxonomyCreateForm>>()
const quickCategoryForm = reactive<CategoryRequest>({
  name: '',
  route: '',
  description: '',
  show_in_nav: true,
  enabled: true
})

const quickTagVisible = ref(false)
const quickTagLoading = ref(false)
const quickTagFormRef = ref<InstanceType<typeof TaxonomyCreateForm>>()
const quickTagForm = reactive<TagRequest>({
  name: '',
  route: '',
  enabled: true
})

const resetQuickCategoryForm = () => {
  quickCategoryForm.name = ''
  quickCategoryForm.route = ''
  quickCategoryForm.description = ''
  quickCategoryForm.show_in_nav = true
  quickCategoryForm.enabled = true
  quickCategoryFormRef.value?.resetRouteTouched()
  quickCategoryFormRef.value?.clearValidate()
}

const resetQuickTagForm = () => {
  quickTagForm.name = ''
  quickTagForm.route = ''
  quickTagForm.enabled = true
  quickTagFormRef.value?.resetRouteTouched()
  quickTagFormRef.value?.clearValidate()
}

const openQuickCategory = () => {
  resetQuickCategoryForm()
  quickCategoryVisible.value = true
}

const openQuickTag = () => {
  resetQuickTagForm()
  quickTagVisible.value = true
}

const refreshCategoryOptions = async () => {
  const response: any = await GetSelectedCategories()
  categoryOptions.value = response.data.data?.list || response.data.data || []
}

const refreshTagOptions = async () => {
  const response: any = await GetSelectedTags()
  tagOptions.value = response.data.data?.list || response.data.data || []
}

const isFormValidateError = (error: unknown) => {
  return typeof error === 'object' && error !== null && 'errorFields' in error
}

const createCategory = async () => {
  if (!quickCategoryFormRef.value) {
    return
  }
  try {
    await quickCategoryFormRef.value.validateFields()
    quickCategoryLoading.value = true
    const response: any = await AddCategory({ ...quickCategoryForm })
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    const createdName = quickCategoryForm.name
    await refreshCategoryOptions()
    const createdOption = categoryOptions.value.find((category) => category.value === createdName)
    if (createdOption) {
      post4Edit.tempCategories = Array.from(
        new Set([...(post4Edit.tempCategories || []), createdOption.value])
      )
    }
    message.success('分类创建成功')
    quickCategoryVisible.value = false
    resetQuickCategoryForm()
  } catch (error) {
    if (isFormValidateError(error)) {
      message.warning('请检查分类信息是否填写正确')
      return
    }
    message.error(toErrorMessage(error, '分类创建失败'))
  } finally {
    quickCategoryLoading.value = false
  }
}

const createTag = async () => {
  if (!quickTagFormRef.value) {
    return
  }
  try {
    await quickTagFormRef.value.validateFields()
    quickTagLoading.value = true
    const response: any = await AddTag({ ...quickTagForm })
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    const createdName = quickTagForm.name
    await refreshTagOptions()
    const createdOption = tagOptions.value.find((tag) => tag.value === createdName)
    if (createdOption) {
      post4Edit.tempTags = Array.from(new Set([...(post4Edit.tempTags || []), createdOption.value]))
    }
    message.success('标签创建成功')
    quickTagVisible.value = false
    resetQuickTagForm()
  } catch (error) {
    if (isFormValidateError(error)) {
      message.warning('请检查标签信息是否填写正确')
      return
    }
    message.error(toErrorMessage(error, '标签创建失败'))
  } finally {
    quickTagLoading.value = false
  }
}
</script>

<style scoped>
.taxonomy-tools {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 8px;
}

.taxonomy-tools :deep(.ant-btn) {
  align-self: flex-start;
  padding-left: 0;
}

.field-with-action {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.field-with-action :deep(.ant-btn) {
  align-self: flex-start;
  padding-left: 0;
}

.metadata-section {
  padding: 4px 0 16px;
  border-bottom: 1px solid var(--app-border);
}

.metadata-section + .metadata-section {
  margin-top: 18px;
}

.metadata-section-title {
  margin: 0 0 16px;
  color: var(--app-text);
  font-size: 15px;
  font-weight: 600;
  line-height: 1.4;
}

.metadata-settings-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  column-gap: 16px;
}

.save-status {
  flex-shrink: 0;
  color: var(--app-text-secondary);
  font-size: 12px;
  white-space: nowrap;
}

@media (max-width: 768px) {
  :deep(.ant-drawer-content-wrapper) {
    width: 100% !important;
  }

  .metadata-settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>
