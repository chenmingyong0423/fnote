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
import type { PostRequest } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import axios from '@/http/axios'
import type { IBaseResponse, IListData, IResponse } from '@/interfaces/Common'
import type { SelectCategory } from '@/interfaces/Category'
import type { SelectTag } from '@/interfaces/Tag'
import PostEditView from '@/views/post/PostEditView.vue'

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
  status: 0
})

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])

const submit = async (postReq: PostRequest) => {
  console.log(postReq)
  try {
    const response = await axios.post<IBaseResponse>('/admin/posts', postReq)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    message.success('添加成功')
    postEditRef.value.clearReq()
  } catch (error) {
    console.log(error)
    message.error('添加失败')
  }
}

const getCategories = async () => {
  try {
    const response = await axios.get<IResponse<IListData<SelectCategory>>>(
      '/admin/categories/select'
    )
    response.data.data?.list.forEach((item) => {
      categories.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}
getCategories()

const getTags = async () => {
  try {
    const response = await axios.get<IResponse<IListData<SelectTag>>>('/admin/tags/select')
    response.data.data?.list.forEach((item) => {
      tags.value?.push(item)
    })
  } catch (error) {
    console.log(error)
  }
}
getTags()
</script>
