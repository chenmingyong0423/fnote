// src/http/axios.ts
import axios, { AxiosError, CanceledError, type InternalAxiosRequestConfig } from 'axios'
import { useUserStore } from '@/stores/user'
import { message } from 'ant-design-vue'
import router from '@/router'
import qs from 'qs'

const instance = axios.create({
  baseURL: import.meta.env.VITE_API_HOST + '/admin-api',
  timeout: 99999,
  paramsSerializer: (params) => {
    // 删除空数组
    Object.keys(params).forEach((key) => {
      if (params[key] === null || params[key] === '' || params[key] === undefined) {
        delete params[key]
      }
    })
    // 使用 qs 来序列化参数
    return qs.stringify(params, { arrayFormat: 'repeat' })
  }
})

const publicPaths = new Set([
  '/login',
  '/configs/check-initialization',
  '/configs/website/meta'
])
const activeRequests = new Set<AbortController>()
let invalidatedToken = ''
let handlingUnauthorized = false

type TrackedRequestConfig = InternalAxiosRequestConfig & {
  sessionController?: AbortController
}

const isPublicRequest = (config: InternalAxiosRequestConfig) =>
  publicPaths.has(String(config.url || '').split('?')[0])

const trackRequest = (config: TrackedRequestConfig) => {
  const controller = new AbortController()
  const signal = config.signal
  if (signal) {
    if (signal.aborted) {
      controller.abort()
    } else if (signal.addEventListener) {
      signal.addEventListener('abort', () => controller.abort(), { once: true })
    }
  }
  config.sessionController = controller
  config.signal = controller.signal
  activeRequests.add(controller)
}

const finishRequest = (config?: TrackedRequestConfig) => {
  if (config?.sessionController) {
    activeRequests.delete(config.sessionController)
  }
}

const pendingRequest = <T>() => new Promise<T>(() => undefined)

const expireSession = (token: string) => {
  if (handlingUnauthorized && invalidatedToken === token) return

  handlingUnauthorized = true
  invalidatedToken = token
  activeRequests.forEach((controller) => controller.abort())
  activeRequests.clear()

  const userStore = useUserStore()
  userStore.clearSession()
  message.warn('登录过期，请重新登录').then((r) => r)
  if (router.currentRoute.value.name !== 'login') {
    router.replace({ name: 'login' }).then((r) => r)
  }
}

// 请求拦截器
instance.interceptors.request.use(
  (config: TrackedRequestConfig) => {
    const userStore = useUserStore()
    const token = userStore.token
    if (token && token !== invalidatedToken) {
      invalidatedToken = ''
      handlingUnauthorized = false
    }
    if (!isPublicRequest(config) && userStore.isSessionExpired) {
      expireSession(token)
      return pendingRequest<InternalAxiosRequestConfig>()
    }
    if (!isPublicRequest(config) && handlingUnauthorized) {
      return pendingRequest<InternalAxiosRequestConfig>()
    }

    trackRequest(config)
    config.headers.set('Authorization', token)
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
    return Promise.reject(error)
  }
)

// 响应拦截器
instance.interceptors.response.use(
  (response) => {
    finishRequest(response.config as TrackedRequestConfig)
    return response
  },
  (error: AxiosError) => {
    finishRequest(error.config as TrackedRequestConfig | undefined)
    if (error instanceof CanceledError || error.code === 'ERR_CANCELED') {
      return pendingRequest()
    }
    // 对响应错误做点什么
    if (!error.response) {
      message.error(error.toString()).then((r) => r)
      return Promise.reject(error)
    }

    const userStore = useUserStore()
    const contentType = String(error.response.headers['content-type'] || '')

    switch (error.response.status) {
      case 401:
        expireSession(userStore.token || invalidatedToken)
        return pendingRequest()
      case 500:
        if (contentType.includes('application/json') && error.response.data) {
          console.log(error)
          const data = error.response.data as { message?: string }
          message.error(data.message || error.message).then((r) => r)
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
