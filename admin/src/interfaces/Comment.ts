import instance from '@/utils/axios'
import type { PageRequest } from '@/interfaces/Common'

export interface Comment {
  id: string
  post_info: {
    post_id: string
    post_title: string
    post_url: string
  }
  content: string
  user_info: {
    username: string
    email: string
    website?: string
    ip: string
  }
  fid?: string
  type: number
}

export const GetComments = (req: PageRequest) => {
  return instance({
    url: `/comments`,
    method: 'get',
    params: req
  })
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

export const DisapproveCommentById = (id: string, reason: string) => {
  return instance({
    url: `/comments/${id}/disapproval`,
    method: 'put',
    data: {
      reason: reason
    }
  })
}

export const DisapproveReplyById = (fid: string, id: string, reason: string) => {
  return instance({
    url: `/comments/${fid}/replies/${id}/disapproval`,
    method: 'put',
    data: {
      reason: reason
    }
  })
}

export const UpdateCommentStatusById = (id: string, status: number) => {
  return instance({
    url: `/comments/${id}/status`,
    method: 'put',
    data: {
      status: status
    }
  })
}

export const UpdateReplyStatusById = (fid: string, id: string, status: number) => {
  return instance({
    url: `/comments/${fid}/replies/${id}/status`,
    method: 'put',
    data: {
      status: status
    }
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
