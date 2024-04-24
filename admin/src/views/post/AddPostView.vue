<template>
  <div>
    <PostEditView
      ref="postEditRef"
      :categories="categories"
      :tags="tags"
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

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])

const saveDraft = async (post4Edit: Post4Edit) => {
  const postDraftReq = {} as PostDraftRequest
  Object.assign(postDraftReq, post4Edit)
  delete (postDraftReq as any).tempCategories
  delete (postDraftReq as any).tempTags
  try {
    const res: any = await SavePostDraft(postDraftReq)
    if (res.data.code === 0) {
      console.log(res.data)
      message.success('保存成功')
      await router.push(`/home/post/draft/${res.data.data.id}`)
    } else {
      message.error(res.data.message)
    }
  } catch (error) {
    message.error(error)
  }
}

const submit = async (post4Edit: Post4Edit) => {
  console.log()
  const postReq = {} as PostRequest
  Object.assign(postReq, post4Edit)
  delete (postReq as any).tempCategories
  delete (postReq as any).tempTags
  try {
    const response: any = await AddPost(postReq)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('添加成功')
    postEditRef.value.clearReq()
  } catch (error) {
    console.log(error)
    if (originalAxios.isAxiosError(error)) {
      // 这是一个由 axios 抛出的错误
      if (error.response) {
        if (error.response.data.status === 409) {
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
