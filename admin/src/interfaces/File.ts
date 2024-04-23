import instance from '@/utils/axios'
import type { PostRequest } from '@/interfaces/Post'

export interface File {
  file_id: string
  file_name: string
  url: string
}

export const FileUpload = (formData: FormData) => {
  return instance({
    url: `/files/upload`,
    method: 'post',
    data: formData
  })
}
