export interface Friend {
  id: string
  name: string
  url: string
  logo: string
  description: string
  show: boolean
  accepted: boolean
  create_time: number
}

export interface FriendReq {
  name: string
  logo: string
  description: string
  show: boolean
}
