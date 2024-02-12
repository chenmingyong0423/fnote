import instance from '@/utils/axios'
import type { PageRequest } from '@/interfaces/Common'

export interface Friend {
  id: string
  name: string
  url: string
  logo: string
  description: string
  status: number
  create_time: number
}

export interface FriendReq {
  name: string
  logo: string
  description: string
  show: boolean
}

export const GetFriends = (pageReq: PageRequest) => {
  return instance({
    url: '/friends',
    method: 'get',
    params: pageReq
  })
}

export const DeleteFriend = (id: string) => {
  return instance({
    url: `/friends/${id}`,
    method: 'delete'
  })
}

// blogUrl 后续要删除
export const ApproveFriend = (id: string, blogUrl: string) => {
  return instance({
    url: `/friends/${id}/approval`,
    method: 'put',
    data: {
      host: blogUrl
    }
  })
}

export const RejectFriend = (id: string, blogUrl: string, reason: string) => {
  return instance({
    url: `/friends/${id}/rejection`,
    method: 'put',
    data: {
      host: blogUrl,
      reason: reason
    }
  })
}

export const UpdateFriend = (id: string, friend: FriendReq) => {
  return instance({
    url: `/friends/${id}`,
    method: 'put',
    data: friend
  })
}
