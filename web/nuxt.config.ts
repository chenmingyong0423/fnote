// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    devtools: {enabled: true},
    modules: [
        '@ant-design-vue/nuxt',
        '@unocss/nuxt'
    ],
    antd: {
        components: ['Flex']
    },
    css: []
})
