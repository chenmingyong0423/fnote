export interface ICategory {
  id: string
  name: string
  route: string
  description: string
  enabled: boolean
  show_in_nav: boolean
  post_count: number
  create_time: number
  update_time: number
}

export interface CategoryRequest {
  name: string
  route: string
  description: string
  enabled: boolean
  show_in_nav: boolean
}

export interface SelectCategory {
  id: string
  value: string
  label: string
}
