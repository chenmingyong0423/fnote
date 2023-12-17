// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    devtools: {enabled: true},
    modules: [
        '@ant-design-vue/nuxt',
        '@unocss/nuxt',
        "@pinia/nuxt",
    ],
    antd: {},
    unocss: {
        uno: true, // enabled `@unocss/preset-uno`
        icons: true, // enabled `@unocss/preset-icons`
        attributify: false, // enabled `@unocss/preset-attributify`,
        // core options
        shortcuts: [
            {'dark_bg_black': 'dark:bg-#000'},
            {'dark_bg_gray': 'dark:bg-#202020'},
            {'dark_text_white': 'dark:c-#fff'},
            {'c_title_blue': 'text-#1890ff'},
            {'c_text_black': 'text-#111'},
            {'c_text_white': 'text-#fff'},
            {'menu_item': 'px18 py15 text-14 dark:text-#fff hover:bg-#000/20 active:bg-#000/40 c-#000 cursor-pointer'},
        ],
        rules: [
            ['footer_shadow', {'box-shadow': ' 0 0 10px rgba(0, 0, 0, .5)'}],
            ['bg_transparent', {'background-color': 'transparent'}],
            ['border_color_white', {'border-color': '#fff'}],
            ['border_color_gray', {'border-color': '#86909c'}],
        ],
        theme: {
            colors: {
                'dtc': 'hsla(0,0%,100%,0.7)', //
            },
            hovers:{

            }
        },
        safelist: [],
    },
    css: [
        '@/styles/main.css'
    ],
    plugins: [
        '~/plugins/pinia.js',
        '~/plugins/localStorage.client.ts',
    ],
})
