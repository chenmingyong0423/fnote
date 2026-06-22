import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import { visualizer } from 'rollup-plugin-visualizer'
import Components from 'unplugin-vue-components/vite'
import { AntDesignVueResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig(({ mode }) => {
  const isDocker = mode === 'docker'
  const env = loadEnv(mode, process.cwd(), '')
  const apiHost = env.VITE_API_HOST

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

    server: {
      proxy: apiHost
        ? {
            '/static': {
              target: apiHost,
              changeOrigin: true
            }
          }
        : undefined
    },

    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url))
      }
    }
  }
})
