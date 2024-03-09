import instance from '@/utils/axios'
import type { PostRequest } from '@/interfaces/Post'

export interface File {
  file_id: string
  file_name: string
  url: string
}

export const UpdatePost = (post: PostRequest) => {
  return instance({
    url: `/files/`,
    method: 'put',
    data: post
  })
}
