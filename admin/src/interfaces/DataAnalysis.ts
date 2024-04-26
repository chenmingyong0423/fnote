import instance from '@/utils/axios'


export interface TodayTrafficStatsVO {
  view_count: number
  user_view_count: number
  comment_count: number
  like_count: number
}

export interface TrafficStatsVO {
  view_count: number
  comment_count: number
  like_count: number
}

export const GetTodayTrafficStats = () => {
  return instance({
    url: '/data-analysis/traffic/today',
    method: 'get'
  })
}

export const GetTrafficStats = () => {
  return instance({
    url: '/data-analysis/traffic',
    method: 'get'
  })
}