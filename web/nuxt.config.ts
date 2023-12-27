// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    devtools: {enabled: true},
    modules: [
        '@ant-design-vue/nuxt',
        '@unocss/nuxt',
        "@pinia/nuxt",
        'dayjs-nuxt',
    ],
    antd: {},
    unocss: {
        uno: true, // enabled `@unocss/preset-uno`
        icons: true, // enabled `@unocss/preset-icons`
        attributify: false, // enabled `@unocss/preset-attributify`,
        // core options
        shortcuts: [
            {'dark_bg_black': 'dark:bg-#03080c'},
            {'dark_bg_gray': 'dark:bg-#207191/10 dark:border-solid dark:border-#1e2227 dark:border-1'},
            {'dark_bg_full_black': 'dark:bg-#000/90 dark:border-solid dark:border-#1e2227 dark:border-1'},
            {'dark_border': 'border-solid border-#1e2227 border-1'},
            {'dark_text_white': 'dark:c-#fff'},
            {'light_border': 'border-1 border-gray-2 border-solid'},
            {'light_border_bottom': 'border-b-1 border-b-gray-2 border-b-solid'},
            {'c_title_blue': 'text-#1890ff'},
            {'c_text_black': 'text-#111'},
            {'c_text_white': 'text-#fff'},
            {'menu_item': 'px18 py15 text-14 dark:text-#fff hover:bg-#000/20 active:bg-#000/40 c-#000 cursor-pointer'},
            {'flex-center': 'flex items-center justify-center'},
            {'custom_cursor_flow': 'hover:drop-shadow-xl hover:translate-y--2 duration-100'},
            {'custom_icon': 'hover:text-blue duration-100'},
            {'custom_shadow': 'hover:shadow-2xl dark:shadow-white duration-100'},
        ],
        rules: [
            ['footer_shadow', {'box-shadow': ' 0 0 10px rgba(0, 0, 0, .5)'}],
            ['bg_transparent', {'background-color': 'transparent'}],
            ['border_bottom_blue', {'border-bottom': '2px solid #007fff'}],
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
        '~/plugins/v-md-editor.client.js',
    ],
})
