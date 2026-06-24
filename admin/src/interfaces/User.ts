import instance from '@/utils/axios'
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginVO {
  owner_info: {
    username: string
    picture: string
  }
  token: string
  expiration: number
}

export interface LoginResponse {
  code: number
  data: LoginVO
  message: string
}

export const login = (data: LoginRequest) => {
  return instance({
    url: '/login',
    method: 'post',
    data
  })
}
