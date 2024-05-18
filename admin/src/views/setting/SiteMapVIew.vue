<template>
  <a-card title="sitemap">
    <template #extra>
      <div class="flex gap-x-3">
        <a-tooltip title="刷新数据">
          <a-button
            shape="circle"
            :icon="h(ReloadOutlined)"
            :loading="loading"
            @click="getSitemap"
          />
        </a-tooltip>
        <a-tooltip title="复制 sitemap 的内容">
          <a-button shape="circle" :icon="h(CopyrightOutlined)" @click="copySitemap" />
        </a-tooltip>
      </div>
    </template>
    <a-spin :spinning="loading">
      <a-collapse v-model:activeKey="activeKey">
        <a-collapse-panel key="1" header="sitemap">
          <pre class="max-h-500px overflow-auto">{{ sitemap }}</pre>
        </a-collapse-panel>
      </a-collapse>
    </a-spin>
    <a-button @click="generateNewSitemap" class="mt-5">刷新 sitemap </a-button>
  </a-card>
</template>

<script lang="ts" setup>
import { h, ref } from 'vue'
import { ReloadOutlined, CopyrightOutlined } from '@ant-design/icons-vue'
import { GenerateSitemap } from '@/interfaces/PostIndex'
import { message } from 'ant-design-vue'
import type { IBaseResponse } from '@/interfaces/Common'

const loading = ref(false)
const sitemap = ref('')
const activeKey = ref(['1'])
const serverHost = import.meta.env.VITE_API_HOST

const getSitemap = async () => {
  try {
    const sitemapUrl = serverHost + '/static/sitemap.xml'
    const response = await fetch(`${sitemapUrl}?t=${new Date().getTime()}`)
    sitemap.value = await response.text()
    loading.value = true
  } catch (e) {
    console.log(e)
  } finally {
    loading.value = false
  }
}
getSitemap()

const generateNewSitemap = async () => {
  try {
    loading.value = true
    const res: any = await GenerateSitemap()
    const apiRes: IBaseResponse = res.data
    if (apiRes.code === 0) {
      console.log(res.data)
      message.success('生成成功')
      await getSitemap()
    } else {
      message.error(res.data.message)
    }
  } catch (e) {
    console.log(e)
  } finally {
    loading.value = false
  }
}

const copySitemap = async () => {
  await navigator.clipboard.writeText(sitemap.value)
  message.success('复制成功！', 2000)
}
</script>
