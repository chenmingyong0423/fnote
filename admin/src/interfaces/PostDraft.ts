import instance from '@/utils/axios'

export interface Category4Post {
  id: string
  name: string
}

export interface Tag4Post {
  id: string
  name: string
}

export interface PostDraftRequest {
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
  created_at: number
}

export interface PostDraftDetail {
  id: string
  post_id: string
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
  word_count: number
  meta_keywords: string
  is_comment_allowed: boolean
  created_at: number
}

export const SavePostDraft = (data: PostDraftRequest) => {
  return instance({
    url: '/post-draft',
    method: 'post',
    data: data
  })
}

export const GetPostDraftDetail = (id: string) => {
  return instance({
    url: `/post-draft/${id}`,
    method: 'get'
  })
}