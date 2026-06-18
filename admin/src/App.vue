<script setup lang="ts">
import { useWebsiteStore } from '@/stores/website'
import { useThemeStore } from '@/stores/theme'
import { GetWebSiteMeta, type WebsiteConfigMetaVO } from '@/interfaces/Config'
import { computed } from 'vue'
import { theme as antTheme } from 'ant-design-vue'
import zhCN from 'ant-design-vue/es/locale/zh_CN'
import { BulbOutlined, HighlightOutlined } from '@ant-design/icons-vue'

const websiteStore = useWebsiteStore()
const themeStore = useThemeStore()
const isDark = computed(() => themeStore.mode === 'dark')
const themeConfig = computed(() => ({
  algorithm: isDark.value ? antTheme.darkAlgorithm : antTheme.defaultAlgorithm
}))
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
  <a-config-provider :locale="zhCN" :theme="themeConfig">
    <RouterView />
    <a-tooltip :title="isDark ? '切换到亮色模式' : '切换到暗黑模式'">
      <a-button
        class="theme-toggle"
        shape="circle"
        :aria-label="isDark ? '切换到亮色模式' : '切换到暗黑模式'"
        @click="themeStore.toggleTheme"
      >
        <template #icon>
          <BulbOutlined v-if="isDark" />
          <HighlightOutlined v-else />
        </template>
      </a-button>
    </a-tooltip>
  </a-config-provider>
</template>

<style scoped>
.theme-toggle {
  position: fixed;
  top: 16px;
  right: 20px;
  z-index: 1100;
  box-shadow: var(--app-shadow);
}
</style>
