// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: [
    '@unocss/nuxt',
    "@pinia/nuxt",
    '@element-plus/nuxt'
  ],

  unocss: {
    uno: true, // enabled `@unocss/preset-uno`
    icons: true, // enabled `@unocss/preset-icons`
    attributify: true, // enabled `@unocss/preset-attributify`,
    // core options
    shortcuts: [
      { 'dark_bg_black': 'dark:bg-#000' },
      { 'dark_bg_gray': 'dark:bg-#202020' },
      { 'dark_text_white': 'dark:c-#fff' },
    ],
    rules: [
      ['footer_shadow', { 'box-shadow': ' 0 0 10px rgba(0, 0, 0, .5)' }]
    ],
    safelist: [],
  },
  pinia: {
    autoImports: [
      // 自动引入 `defineStore(), storeToRefs()`
      "defineStore",
      "storeToRefs"
    ],
  },
  elementPlus: {
    /** Options */
  },
  css: [
    '@/styles/main.css'
  ]
})
