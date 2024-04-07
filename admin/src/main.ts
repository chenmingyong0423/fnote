import './assets/main.css'

import { createApp } from 'vue'
import { createPinia } from 'pinia'

import App from './App.vue'
import router from './router'
import Antd, { message } from 'ant-design-vue'

import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import '@kangc/v-md-editor/lib/theme/style/github.css'
import 'uno.css'
import 'virtual:uno.css'

// highlightjs
import hljs from 'highlight.js'
import { GetSeo, isInit } from '@/interfaces/Config'
import { useUserStore } from '@/stores/user'

VMdEditor.use(githubTheme, {
  Hljs: hljs
})

const app = createApp(App)

app.use(VMdEditor)
app.use(Antd)
app.use(createPinia())
app.use(router)

const userStore = useUserStore();


try {
  const res: any =  await isInit()
  if (res.data.code === 0) {
    userStore.isInit = res.data.data.initStatus
  }
} catch (error) {
  console.log(error)
}

app.mount('#app')
