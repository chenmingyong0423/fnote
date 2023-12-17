export default defineNuxtPlugin((nuxtApp) => {
    // 这里你可以添加你的 localStorage 相关代码
    // 例如，你可以添加到 nuxtApp 的实例上，使其在应用中全局可用
    nuxtApp.provide('localStorage', {
        setItem(key: string, value: string) {
            if (process.client) {
                window.localStorage.setItem(key, value);
            }
        },
        getItem(key: string) {
            return process.client ? window.localStorage.getItem(key) : null;
        },
        // 其他 localStorage 方法...
    });
});
