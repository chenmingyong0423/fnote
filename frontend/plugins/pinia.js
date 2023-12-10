import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import { defineStore, storeToRefs } from 'pinia'

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.$pinia.use(piniaPluginPersistedstate)
    nuxtApp.provide('defineStore', defineStore)
    nuxtApp.provide('storeToRefs', storeToRefs)
})