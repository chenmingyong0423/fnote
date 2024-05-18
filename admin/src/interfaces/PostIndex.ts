import instance from '@/utils/axios'

export const GenerateSitemap = () => {
  return instance({
    url: `/post-index/sitemap`,
    method: 'post'
  })
}
