import { defineStore } from 'pinia'
import type { WebsiteConfigMetaVO } from '@/interfaces/Config'

export const useWebsiteStore = defineStore('website', {
  state: () => ({
    website_name: '',
    website_icon: ''
  }),
  actions: {
    Update(data : WebsiteConfigMetaVO) {
      this.website_name = data.website_name
      this.website_icon = import.meta.env.VITE_API_HOST + data.website_icon
      // 设置 title 和 meta 的 description 以及 icon
      document.title = data.website_name + '- 后台管理'
      document
        .querySelector('meta[name="description"]')
        ?.setAttribute('content', data.website_name + '- 后台管理')
      document.querySelector('link[rel="icon"]')?.setAttribute('href', this.website_icon)
    }
  }
})
