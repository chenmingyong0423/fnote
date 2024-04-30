import instance from '@/utils/axios'
import type { PageRequest } from '@/interfaces/Common'

export interface AdminCommentVO {
  id: string
  post_info: PostInfo
  content: string
  user_info: UserInfo4Comment
  replies: AdminCommentVO[]
  approval_status: boolean
  created_at: number
  updated_at: number
  reply_to_id: string
  type: string
  key?: string
  fid?: string
}

export interface PostInfo {
  post_id: string
  post_title: string
  post_url: string
}

export const GetComments = (req: PageRequest) => {
  return instance({
    url: `/comments`,
    method: 'get',
    params: req
  })
}

export interface UserInfo4Comment {
  name: string
  email: string
  ip: string
  website: string
  picture: string
}

export const ApproveCommentById = (id: string) => {
  return instance({
    url: `/comments/${id}/approval`,
    method: 'put'
  })
}

export const ApproveReplyById = (fid: string, id: string) => {
  return instance({
    url: `/comments/${fid}/replies/${id}/approval`,
    method: 'put'
  })
}

export const DeleteCommentById = (id: string) => {
  return instance({
    url: `/comments/${id}`,
    method: 'delete'
  })
}

export const DeleteReplyById = (fid: string, id: string) => {
  return instance({
    url: `/comments/${fid}/replies/${id}`,
    method: 'delete'
  })
}

export interface BatchApprovedCommentRequest {
  comment_ids: string[]
  replies: {
    [key: string]: string[]
  }
}

export const batchApproved = (req: BatchApprovedCommentRequest) => {
  return instance({
    url: `/comments/batch-approval`,
    method: 'put',
    data: req
  })
}
