<template>
  <div>
    <PostEditView
      ref="postEditRef"
      :req="postReq"
      :categories="categories"
      :tags="tags"
      @submit="submit"
    ></PostEditView>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import { AddPost, type PostRequest } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'
import PostEditView from '@/views/post/PostEditView.vue'
import originalAxios from 'axios'

const postReq = reactive<PostRequest>({
  id: '',
  author: '',
  title: '',
  summary: '',
  content: '',
  cover_img: '',
  categories: [],
  tempCategories: [],
  tags: [],
  tempTags: [],
  sticky_weight: 0,
  meta_description: '',
  meta_keywords: '',
  is_comment_allowed: true,
  is_displayed: false
})

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])

const submit = async (postReq: PostRequest) => {
  console.log(postReq)
  try {
    const response: any = await AddPost(postReq)
    if (response.code !== 0) {
      message.error(response.message)
      return
    }
    message.success('添加成功')
    postEditRef.value.clearReq()
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
  }
}

const getCategories = async () => {
  try {
    const response: any = await GetSelectedCategories()
    response.data?.list.forEach((item: SelectCategory) => {
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
    response.data?.list.forEach((item: SelectTag) => {
      tags.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}
getTags()
</script>
