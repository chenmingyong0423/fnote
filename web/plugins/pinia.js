import piniaPluginPersistedState from 'pinia-plugin-persistedstate'
import { defineStore, storeToRefs } from 'pinia'

export default defineNuxtPlugin((nuxtApp) => {
    nuxtApp.$pinia.use(piniaPluginPersistedState)
    nuxtApp.provide('defineStore', defineStore)
    nuxtApp.provide('storeToRefs', storeToRefs)
})