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
