import instance from '@/utils/axios'

export interface IPost {
  id: string
  cover_img: string
  title: string
  summary: string
  categories: Category4Post[]
  tags: Tag4Post[]
  created_at: number
  updated_at: number
}

export interface PostDetailVO {
  id: string
  author: string
  title: string
  summary: string
  content: string
  cover_img: string
  categories: Category4Post[]
  tags: Tag4Post[]
  is_displayed: boolean
  sticky_weight: number
  meta_description: string
  meta_keywords: string
  is_comment_allowed: boolean
}

export interface Category4Post {
  id: string
  name: string
}

export interface Tag4Post {
  id: string
  name: string
}

export type PageRequest = {
  pageNo: number
  pageSize: number
  sortField?: string
  sortOrder?: string
  keyword?: string
  category_filter?: string[]
  tag_filter?: string[]
}

export interface Post4Edit extends PostRequest {
  tempCategories?: string[]
  tempTags?: string[]
  created_at?: number
}

export interface PostRequest {
  id: string
  author: string
  title: string
  summary: string
  content: string
  cover_img: string
  categories: Category4Post[]
  tags: Tag4Post[]
  is_displayed: boolean
  sticky_weight: number
  meta_description: string
  meta_keywords: string
  is_comment_allowed: boolean
}

export const AddPost = (post: PostRequest) => {
  return instance({
    url: '/posts',
    method: 'post',
    data: post
  })
}

export const GetPost = (pageReq: PageRequest) => {
  return instance({
    url: '/posts',
    method: 'get',
    params: pageReq
  })
}

export const DeletePost = (id: string) => {
  return instance({
    url: `/posts/${id}`,
    method: 'delete'
  })
}

export const ChangePostDisplayStatus = (id: string, isDisplayed: boolean) => {
  return instance({
    url: `/posts/${id}/display`,
    method: 'put',
    data: {
      is_displayed: isDisplayed
    }
  })
}

export const ChangeCommentAllowedStatus = (id: string, isCommentAllowed: boolean) => {
  return instance({
    url: `/posts/${id}/comment-allowed`,
    method: 'put',
    data: {
      is_comment_allowed: isCommentAllowed
    }
  })
}

export const PublishPost = (post: PostRequest) => {
  return instance({
    url: `/post-draft/${post.id}/publish`,
    method: 'post',
    data: post
  })
}

export const DeletePostDraftById = (id: string) => {
  return instance({
    url: `/post-draft/${id}`,
    method: 'delete'
  })
}

export const GetPostDraft = (pageReq: PageRequest) => {
  return instance({
    url: '/post-draft',
    method: 'get',
    params: pageReq
  })
}
