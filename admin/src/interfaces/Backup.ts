import instance from '@/utils/axios'

export const Backup = () => {
  return instance({
    url: `/backup`,
    method: 'get',
    responseType: 'blob'
  })
}
export const Recovery = (data: any) => {
  return instance({
    url: `/recovery`,
    method: 'post',
    data: data
  })
}
