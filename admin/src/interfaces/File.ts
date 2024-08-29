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

export type PageRequest = {
  pageNum: number
  pageSize: number
  fileType: string[]
}

// type FileVO struct {
// 	FileId   string `json:"file_id"`
// 	FileName string `json:"file_name"`
// 	Url      string `json:"url"`
// }

export interface FileVO {
  file_id: string
  file_name: string
  url: string
}

export const GetFileList = (pageRequest: PageRequest) => {
  return instance({
    url: `/files`,
    method: 'get',
    params: pageRequest
  })
}
