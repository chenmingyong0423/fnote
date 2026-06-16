<template>
  <div>
    <PostEditView
      ref="postEditRef"
      :post="post4Edit"
      :categories="categories"
      :tags="tags"
      :is-new-post="false"
      :publishing="publishing"
      :saving-draft="savingDraft"
      :auto-save="autoSaveEnabled"
      :last-saved-at="lastSavedAt"
      @publish="submit"
      @saveDraft="saveDraft"
    ></PostEditView>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { type Category4Post, type Post4Edit, type Tag4Post, PublishPost } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'
import { useRoute, useRouter } from 'vue-router'
import PostEditView from '@/views/post/PostEditView.vue'
import {
  GetPostDraftDetail,
  type PostDraftDetail,
  type PostDraftRequest,
  SavePostDraft
} from '@/interfaces/PostDraft'
import axios from 'axios'
import { toErrorMessage } from '@/utils/error'

document.title = '编辑草稿 - 后台管理'

const route = useRoute()
const router = useRouter()
const post4Edit = reactive<Post4Edit>({
  id: '',
  author: '',
  title: '',
  summary: '',
  content: '',
  cover_img: '',
  categories: [],
  tags: [],
  is_displayed: true,
  sticky_weight: 0,
  meta_description: '',
  meta_keywords: '',
  is_comment_allowed: true,
  tempCategories: [],
  tempTags: [],
  created_at: 0
})

const createdAt = ref<Number>(0)
const id = route.params.id as string

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])
const publishing = ref(false)
const savingDraft = ref(false)
const autoSaveEnabled = ref(false)
const lastSavedAt = ref(0)

type SaveDraftOptions = {
  silent?: boolean
}

const saveDraft = async (post4Edit: Post4Edit, options: SaveDraftOptions = {}) => {
  if (savingDraft.value || publishing.value) {
    return
  }
  const postDraftReq = {} as PostDraftRequest
  Object.assign(postDraftReq, post4Edit)
  delete (postDraftReq as any).tempCategories
  delete (postDraftReq as any).tempTags
  try {
    savingDraft.value = true
    const res: any = await SavePostDraft(postDraftReq)
    if (res.data.code === 0) {
      console.log(res.data)
      lastSavedAt.value = Date.now()
      if (!options.silent) {
        message.success('保存成功')
        await getPostDraftById(res.data.data.id)
      }
    } else {
      message.error(res.data.message)
    }
  } catch (error) {
    message.error(toErrorMessage(error))
  } finally {
    savingDraft.value = false
  }
}

const getPostDraftById = async (id: string) => {
  try {
    autoSaveEnabled.value = false
    const response: any = await GetPostDraftDetail(id)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    const postDraft: PostDraftDetail = response.data.data
    if (postDraft) {
      post4Edit.id = postDraft.id
      post4Edit.author = postDraft.author
      post4Edit.title = postDraft.title
      post4Edit.summary = postDraft.summary
      post4Edit.content = postDraft.content
      post4Edit.cover_img = postDraft.cover_img
      post4Edit.sticky_weight = postDraft.sticky_weight
      post4Edit.meta_description = postDraft.meta_description
      post4Edit.meta_keywords = postDraft.meta_keywords
      post4Edit.is_comment_allowed = postDraft.is_comment_allowed
      post4Edit.categories = postDraft.categories
      post4Edit.is_displayed = postDraft.is_displayed
      createdAt.value = postDraft.created_at
      post4Edit.tags = postDraft.tags
      post4Edit.created_at = postDraft.created_at
      post4Edit.tempCategories = []
      postDraft.categories.forEach((item: Category4Post) => {
        post4Edit.tempCategories?.push(item.name)
      })
      post4Edit.tempTags = []
      postDraft.tags.forEach((item: Tag4Post) => {
        post4Edit.tempTags?.push(item.name)
      })
      autoSaveEnabled.value = true
    }
  } catch (error) {
    if (axios.isAxiosError(error)) {
      if (!error.response) {
        message.error(toErrorMessage(error))
        return
      }
      switch (error.response.status) {
        case 404:
          message.error('查询草稿失败，id 不存在')
          await router.push('/home/post/list')
          break
      }
    } else {
      message.error(toErrorMessage(error))
    }
  }
}

getPostDraftById(id)

const submit = async (post4Edit: Post4Edit) => {
  if (publishing.value || savingDraft.value) {
    return
  }
  const postReq = {} as Post4Edit
  Object.assign(postReq, post4Edit)
  delete postReq.tempCategories
  delete postReq.tempTags
  delete postReq.created_at
  try {
    publishing.value = true
    const response: any = await PublishPost(postReq)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('发布成功')
    postEditRef.value.clearReq()
    await router.push('/home/post/list')
  } catch (error) {
    message.error(toErrorMessage(error))
  } finally {
    publishing.value = false
  }
}

const getCategories = async () => {
  try {
    const response = await GetSelectedCategories()
    response.data.data?.list.forEach((item: SelectCategory) => {
      categories.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}

getCategories()

const getTags = async () => {
  try {
    const response = await GetSelectedTags()
    response.data.data?.list.forEach((item: SelectTag) => {
      tags.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}

getTags()
</script>
