import instance from '@/utils/axios'

export interface IPost {
  id: string
  cover_img: string
  title: string
  summary: string
  categories: Category4Post[]
  tags: Tag4Post[]
  create_time: number
  update_time: number
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
  categories?: string[]
  tags?: string[]
}

export interface Post4Edit extends PostRequest{
  tempCategories: string[]
  tempTags: string[]
  created_at: number
}

export interface PostRequest{
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

export const GetPostDetail = (id: string) => {
  return instance({
    url: `/posts/${id}`,
    method: 'get'
  })
}

export const UpdatePost = (post: PostRequest) => {
  return instance({
    url: `/posts`,
    method: 'put',
    data: post
  })
}
