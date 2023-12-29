import QrcodeVue from 'qrcode.vue'

export default defineNuxtPlugin((nuxtApp) => {
    if (process.client) {
        nuxtApp.vueApp.component('QrcodeVue', QrcodeVue)
    }
});