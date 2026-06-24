import { defineStore } from 'pinia'
import { login, type LoginRequest, type LoginResponse, type LoginVO } from '@/interfaces/User'
import { message } from 'ant-design-vue'
import type { AxiosResponse } from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    ownerInfo: {
      username: '',
      picture: ''
    },
    token: localStorage.getItem('token') || '',
    tokenExpiration: Number(localStorage.getItem('token-expiration')) || 0,
    isLoggedIn: !!localStorage.getItem('token'),
    initialization: false
  }),
  getters: {
    isSessionExpired: (state) =>
      Boolean(state.token) &&
      state.tokenExpiration > 0 &&
      state.tokenExpiration <= Math.floor(Date.now() / 1000)
  },
  actions: {
    clearSession() {
      this.token = ''
      this.tokenExpiration = 0
      this.isLoggedIn = false
      this.ownerInfo = { username: '', picture: '' }
      localStorage.removeItem('token')
      localStorage.removeItem('token-expiration')
    },
    async loginIn(req: LoginRequest): Promise<boolean> {
      try {
        const res: AxiosResponse<LoginResponse> = await login(req)

        const body = res.data

        if (body.code === 40101) {
          message.error('用户名或密码错误').then((r) => r)
          return false
        }

        if (body.code === 0) {
          this.token = body.data.token || ''
          this.tokenExpiration = Number(body.data.expiration) || 0
          this.isLoggedIn = true
          this.ownerInfo = body.data.owner_info

          localStorage.setItem('token', this.token)
          localStorage.setItem('token-expiration', String(this.tokenExpiration))

          message.success('登录成功').then((r) => r)
          return true
        }

        return false
      } catch (error) {
        console.error(error)
        return false
      }
    }
  }
})
