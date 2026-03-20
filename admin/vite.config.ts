import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import { visualizer } from 'rollup-plugin-visualizer'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig(({ mode }) => {
  const isDocker = mode === 'docker'

  return {
    plugins: [
      vue(),
      UnoCSS(),

      // 🚨 关键：Docker 时禁用 visualizer
      !isDocker &&
        visualizer({
          open: true,
          gzipSize: true,
          brotliSize: true
        }),

      Components({
        resolvers: [
          AntDesignVueResolver({
            importStyle: false
          })
        ]
      })
    ].filter(Boolean),

    base: `/${process.env.VITE_BUILD_DIR || ''}`,

    build: {
      outDir: `dist/${process.env.VITE_BUILD_DIR || ''}`
    },

    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    }
  }
})
