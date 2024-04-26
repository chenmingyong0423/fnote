<script setup lang="ts">
import { useWebsiteStore } from '@/stores/website'
import { GetWebSiteMeta, type WebsiteConfigMetaVO } from '@/interfaces/Config'

const websiteStore = useWebsiteStore()
const initWebsiteStore = async () => {
  if (websiteStore.website_name === '' && websiteStore.website_icon === '') {
    try {
      const response: any = await GetWebSiteMeta()
      if (response.data.code === 0) {
        const data: WebsiteConfigMetaVO = response.data.data
        websiteStore.Update(data.website_name || 'fnote', data.website_icon)
      }
    } catch (error) {
      console.log(error)
    }
  }
}

initWebsiteStore()
</script>

<template>
  <RouterView />
</template>

<style scoped></style>
