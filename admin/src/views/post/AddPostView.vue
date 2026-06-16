<template>
  <div>
    <PostEditView
      ref="postEditRef"
      :categories="categories"
      :tags="tags"
      :publishing="publishing"
      :saving-draft="savingDraft"
      @publish="submit"
      @saveDraft="saveDraft"
    ></PostEditView>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { AddPost, type Post4Edit, type PostRequest } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'
import PostEditView from '@/views/post/PostEditView.vue'
import originalAxios from 'axios'
import { type PostDraftRequest, SavePostDraft } from '@/interfaces/PostDraft'
import router from '@/router'
import { toErrorMessage } from '@/utils/error'

document.title = '文章发布 - 后台管理'

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])
const publishing = ref(false)
const savingDraft = ref(false)

const saveDraft = async (post4Edit: Post4Edit) => {
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
      message.success('保存成功')
      await router.push(`/home/post/draft/${res.data.data.id}`)
    } else {
      message.error(res.data.message)
    }
  } catch (error) {
    message.error(toErrorMessage(error))
  } finally {
    savingDraft.value = false
  }
}

const submit = async (post4Edit: Post4Edit) => {
  if (publishing.value || savingDraft.value) {
    return
  }
  console.log()
  const postReq = {} as PostRequest
  Object.assign(postReq, post4Edit)
  delete (postReq as any).tempCategories
  delete (postReq as any).tempTags
  try {
    publishing.value = true
    const response: any = await AddPost(postReq)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('添加成功')
    postEditRef.value.clearReq()
    await router.push('/home/post/list')
  } catch (error) {
    console.log(error)
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.status === 409) {
          message.error('id 重复')
          return
        }
      } else if (error.request) {
        // 请求已发出，但没有收到响应
        console.log('No response received:', error.request)
      } else {
        // 在设置请求时触发了一个错误
        console.log('Error Message:', error.message)
      }
    }
    message.error('添加失败')
  } finally {
    publishing.value = false
  }
}

const getCategories = async () => {
  try {
    const response: any = await GetSelectedCategories()
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
    const response: any = await GetSelectedTags()
    response.data.data?.list.forEach((item: SelectTag) => {
      tags.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}
getTags()
</script>
