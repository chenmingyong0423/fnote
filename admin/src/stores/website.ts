import { defineStore } from 'pinia'

export const useWebsiteStore = defineStore('website', {
  state: () => ({
    website_name: '',
    website_icon: ''
  }),
  actions: {
    Update(websiteName: string, websiteIcon: string) {
      this.website_name = websiteName
      this.website_icon = import.meta.env.VITE_API_HOST + websiteIcon
      // 设置 title 和 meta 的 description 以及 icon
      document.title = websiteName
      document.querySelector('meta[name="description"]')?.setAttribute('content', websiteName)
      document.querySelector('link[rel="icon"]')?.setAttribute('href', this.website_icon)
    }
  }
})
