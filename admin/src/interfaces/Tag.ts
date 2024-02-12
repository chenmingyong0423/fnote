import type { PageRequest } from '@/interfaces/Common'
import instance from '@/utils/axios'

export interface Tag {
  id: string
  name: string
  route: string
  post_count: number
  enabled: boolean
  create_time: number
  update_time: number
}

export interface TagRequest {
  name: string
  route: string
  enabled: boolean
}

export interface SelectTag {
  id: string
  value: string
  label: string
}

export const GetTag = (pageReq: PageRequest) => {
  return instance({
    url: '/tags',
    method: 'get',
    params: pageReq
  })
}

export const AddTag = (tag: TagRequest) => {
  return instance({
    url: '/tags',
    method: 'post',
    data: tag
  })
}

export const ChangeTagEnabled = (id: string, enabled: boolean) => {
  return instance({
    url: `/tags/${id}/enabled`,
    method: 'put',
    data: {
      enabled: enabled
    }
  })
}

export const DeleteTag = (id: string) => {
  return instance({
    url: `/tags/${id}`,
    method: 'delete'
  })
}

export const GetSelectedTags = () => {
  return instance({
    url: '/tags/select',
    method: 'get'
  })
}
