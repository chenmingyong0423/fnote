import type { PageRequest } from '@/interfaces/Post'
import instance from '@/utils/axios'

export interface ICategory {
  id: string
  name: string
  route: string
  description: string
  enabled: boolean
  show_in_nav: boolean
  post_count: number
  create_time: number
  update_time: number
}

export interface CategoryRequest {
  name: string
  route: string
  description: string
  enabled: boolean
  show_in_nav: boolean
}

export interface SelectCategory {
  id: string
  value: string
  label: string
}

export const GetCategories = (pageReq: PageRequest) => {
  return instance({
    url: '/categories',
    method: 'get',
    params: pageReq
  })
}

export const AddCategory = (category: CategoryRequest) => {
  return instance({
    url: '/categories',
    method: 'post',
    data: category
  })
}

export const ChangeCategoryEnabled = (id: string, enabled: boolean) => {
  return instance({
    url: `/categories/${id}/enabled`,
    method: 'put',
    data: {
      enabled: enabled
    }
  })
}

export const ChangeCategoryShowInNav = (id: string, showInNav: boolean) => {
  return instance({
    url: `/categories/${id}/navigation`,
    method: 'put',
    data: {
      show_in_nav: showInNav
    }
  })
}

export const DeleteCategory = (id: string) => {
  return instance({
    url: `/categories/${id}`,
    method: 'delete'
  })
}

export interface UpdateCategoryRequest {
  description: string
}

export const UpdateCategory = (id: string, category: UpdateCategoryRequest) => {
  return instance({
    url: `/categories/${id}`,
    method: 'put',
    data: category
  })
}

export interface SelectCategory {
  id: string
  value: string
  label: string
}

export const GetSelectedCategories = () => {
  return instance({
    url: '/categories/select',
    method: 'get'
  })
}
