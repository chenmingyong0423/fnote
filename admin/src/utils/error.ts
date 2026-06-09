import axios from 'axios'

export function toErrorMessage(error: unknown, fallback = '操作失败，请稍后再试') {
  if (axios.isAxiosError(error)) {
    const data = error.response?.data as { message?: unknown } | undefined
    if (typeof data?.message === 'string' && data.message.trim()) {
      return data.message
    }
    if (error.message) {
      return error.message
    }
  }

  if (error instanceof Error && error.message) {
    return error.message
  }

  if (typeof error === 'string' && error.trim()) {
    return error
  }

  return fallback
}
