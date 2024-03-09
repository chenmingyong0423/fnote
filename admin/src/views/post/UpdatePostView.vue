<template>
  <div>
    <PostEditView
      ref="postEditRef"
      :req="postReq"
      :categories="categories"
      :tags="tags"
      :is-new-post="false"
      @submit="submit"
    ></PostEditView>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue'
import {
  type Category4Post,
  GetPostDetail,
  type PostRequest,
  type Tag4Post,
  UpdatePost
} from '@/interfaces/Post'
import { message } from 'ant-design-vue'
import { GetSelectedCategories, type SelectCategory } from '@/interfaces/Category'
import { GetSelectedTags, type SelectTag } from '@/interfaces/Tag'
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
  is_displayed: false
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
    const response: any = await GetPostDetail(query.id as string)
    if (response.data.code !== 0) {
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
      postReq.is_displayed = post.is_displayed
      postReq.tags = post.tags
      post.categories.forEach((item: Category4Post) => {
        postReq.tempCategories.push(item.name)
      })
      post.tags.forEach((item: Tag4Post) => {
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

const submit = async (postReq: PostRequest) => {
  try {
    const response: any = await UpdatePost(postReq)
    if (response.data.code !== 0) {
      message.error(response.data.message)
      return
    }
    message.success('更新成功')
    postEditRef.value.clearReq()
    await router.push('/home/post/list')
  } catch (error) {
    console.log(error)
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
