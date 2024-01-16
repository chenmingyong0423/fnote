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
import type { PostDetailVO, PostRequest } from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import axios from '@/http/axios'
import type { IListData, IResponse } from '@/interfaces/Common'
import type { SelectCategory } from '@/interfaces/Category'
import type { SelectTag } from '@/interfaces/Tag'
import { useRouter } from 'vue-router'
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

// 获取 query 参数
const router = useRouter()
const query = router.currentRoute.value.params
console.log(query)
postReq.id = query.id as string

const postEditRef = ref()

const categories = ref<SelectCategory[]>([])

const tags = ref<SelectTag[]>([])

const getPostById = async () => {
  try {
    const response = await axios.get<IResponse<PostDetailVO>>(`/admin/posts/${query.id}`)
    if (response.data.code !== 200) {
      message.error(response.data.message)
      return
    }
    const post = response.data.data
    if (post) {
      postReq.id = post.id
      postReq.author = post.author
      postReq.title = post.title
      postReq.summary = post.summary
      postReq.content = post.content
      postReq.cover_img = post.cover_img
      postReq.sticky_weight = post.sticky_weight
      postReq.meta_description = post.meta_description
      postReq.meta_keywords = post.meta_keywords
      postReq.is_comment_allowed = post.is_comment_allowed
      postReq.categories = post.categories
      postReq.tags = post.tags
      post.categories.forEach((item) => {
        postReq.tempCategories.push(item.name)
      })
      post.tags.forEach((item) => {
        postReq.tempTags.push(item.name)
      })
    }
  } catch (error) {
    console.log(error)
  }
}

if (postReq.id) {
  getPostById()
}

const submit = (postReq: PostRequest) => {
  console.log(postReq)
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
