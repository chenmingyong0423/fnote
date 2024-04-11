// src/http/axios.ts
import axios from 'axios'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import router from '@/router'

const instance = axios.create({
  baseURL: import.meta.env.API_BASE_URL || 'http://localhost:8080/admin',
  timeout: 99999
})

// 请求拦截器
instance.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    config.headers.set('Authorization', userStore.token)
    // 判断body里是否有 file 参数，有则设置请求头为 multipart/form-data
    if (config.data instanceof FormData) {
      config.headers.set('Content-Type', 'multipart/form-data')
    } else {
      config.headers.set('Content-Type', 'application/json')
    }
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
    return response
  },
  (error) => {
    // 对响应错误做点什么
    if (!error.response) {
      message.error(error.toString()).then((r) => r)
      return
    }

    const userStore = useUserStore()
    const contentType = error.response.headers['content-type']

    switch (error.response.status) {
      case 401:
        message.warn('登录过期，请重新登录').then((r) => r)
        userStore.token = ''
        localStorage.clear()
        router.push({ path: '/login', replace: true }).then((r) => r)
        break
      case 500:
        if (contentType && contentType.includes('application/json') && error.response.data) {
          console.log(error)
          message.error(error.response.data.message).then((r) => r)
          return
        } else {
          message.error(error.toString()).then((r) => r)
        }
        break
      case 404:
        if (contentType && contentType.includes('application/json')) {
          console.log(error)
          message.error(error.response.data.message).then((r) => r)
          return
        } else {
          message.error(error.toString()).then((r) => r)
        }
        break
    }
    return Promise.reject(error)
  }
)

export default instance
