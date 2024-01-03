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

export interface PostRequest {
  id: string
  author: string
  title: string
  summary: string
  content: string
  cover_img: string
  categories: Category4Post[]
  tempCategories: string[]
  tags: Tag4Post[]
  tempTags: string[]
  status: number;
  sticky_weight: number
  meta_description: string
  meta_keywords: string
  is_comment_allowed: boolean
}
