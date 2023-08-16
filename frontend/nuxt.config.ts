// https://nuxt.com/docs/api/configuration/nuxt-config
export default {
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
    shortcuts: [],
    rules: [],
    safelist: [],
  },
  pinia: {
    autoImports: [
      // 自动引入 `defineStore(), storeToRefs()`
      "defineStore",
      "storeToRefs"
    ],
  },
  elementPlus: { /** Options */ },
  css: [
    '@/styles/main.css'
  ]
}
