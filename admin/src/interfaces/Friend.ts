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
