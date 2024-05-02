export interface IPageData<T> {
  pageNo: number
  pageSize: number
  totalPage: number
  totalCount: number
  list: T[]
}

export interface IListData<T> {
  list: T[]
}

export interface IBaseResponse {
  code: number
  message: string
}

export interface IResponse<T> {
  code: number
  data?: T
  message: string
}

export type PageRequest = {
  pageNo: number
  pageSize: number
  approvalStatus?: boolean
  sort?: string
}
