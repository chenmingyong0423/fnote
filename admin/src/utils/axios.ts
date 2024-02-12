// src/http/axios.ts
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import router from '@/router'

const instance = axios.create({
  baseURL: 'http://localhost:8080/admin',
  timeout: 99999
})

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    config.headers.set('Authorization', userStore.token)
    config.headers.set('Content-Type', 'application/json')
    return config
  },
  (error) => {
    message.error(error.toString()).then((r) => r)
    return error
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    // 对响应错误做点什么
    if (!error.response) {
      message.error(error.toString()).then((r) => r)
      return
    }
    const userStore = useUserStore()

    switch (error.response.status) {
      case 401:
        message.warn('登录过期，请重新登录').then((r) => r)
        userStore.token = ''
        localStorage.clear()
        router.push({ path: '/login', replace: true }).then((r) => r)
        break
      case 500:
        message.error(error.toString()).then((r) => r)
        break
      case 404:
        message.error(error.toString()).then((r) => r)
        break
    }
    return Promise.reject(error)
  }
)

export default instance
