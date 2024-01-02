export interface Tag {
  id: string
  name: string
  route: string
  description: string
  disabled: boolean
  create_time: number
  update_time: number
}

export interface TagRequest {
  name: string
  route: string
  description: string
  disabled: boolean
}