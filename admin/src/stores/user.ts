import { ref } from 'vue'
import { defineStore } from 'pinia'
import { login, type LoginRequest } from '@/interfaces/User'
import { message } from 'ant-design-vue'

export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: {
      username: '',
      picture: ''
    },
    token: localStorage.getItem('token') || '',
    isLoggedIn: !!localStorage.getItem('token'),
    initialization: false
  }),
  actions: {
    async loginIn(req: LoginRequest): Promise<boolean> {
      try {
        const res: any = await login(req)
        if (res.data.code === 0) {
          this.token = res.data.data.token || ''
          this.isLoggedIn = true
          localStorage.setItem('token', this.token)
          message.success('登录成功').then((r) => r)
          return true
        }
        if (res.data.code === 100001) {
          message.error('用户名或密码错误').then((r) => r)
          return false
        }
        return false
      } catch (error) {
        console.error(error)
        return false
      }
    }
  }
})
