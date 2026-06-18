import instance from '@/utils/axios'

export const GenerateSitemap = () => {
  return instance({
    url: `/post-index/sitemap`,
    method: 'post'
  })
}

export interface SitemapData {
  content: string
  exists: boolean
}

export const GetSitemap = () => {
  return instance.get('/post-index/sitemap')
}

export interface RobotsTxtData {
  content: string
  exists: boolean
}

export const GetRobotsTxt = () => {
  return instance.get('/post-index/robots')
}

export const SaveRobotsTxt = (content: string) => {
  return instance.put('/post-index/robots', { content })
}

export const GenerateRobotsTxt = () => {
  return instance.post('/post-index/robots/generate')
}
