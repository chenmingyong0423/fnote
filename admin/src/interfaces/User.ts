import instance from '@/utils/axios'
export interface LoginRequest {
  username: string
  password: string
}

export interface LoginVO {
  token: string
  expiration: number
}

export const login = (data: LoginRequest) => {
  return instance({
    url: '/login',
    method: 'post',
    data
  })
}
