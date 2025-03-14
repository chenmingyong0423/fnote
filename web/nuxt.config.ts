import {visualizer} from "rollup-plugin-visualizer";

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    vite: {
        plugins: [
            visualizer({
                open: true, // 是否在默认浏览器中自动打开报告
                gzipSize: true, // 显示gzip压缩后的大小
                brotliSize: true, // 显示Brotli压缩后的大小
                filename: "stats.html", // 输出报告的文件名
            }),
        ],
    },
    devtools: {enabled: true},
    modules: ["@unocss/nuxt", "@pinia/nuxt", "dayjs-nuxt", "@nuxtjs/mdc"],
    runtimeConfig: {
        public: {
            domain: process.env.BASE_HOST,
            adminHost: process.env.ADMIN_HOST,
            apiHost: process.env.API_HOST,
        },
    },
    nitro: {
        // 该配置用于服务端请求转发
        routeRules: {
            "/api/**": {
                proxy: process.env.SERVER_HOST + "/**",
            },
            "/static/**": {
                proxy: process.env.SERVER_HOST + "/**",
            },
        },
    },
    unocss: {
        uno: true, // enabled `@unocss/preset-uno`
        icons: true, // enabled `@unocss/preset-icons`
        attributify: false, // enabled `@unocss/preset-attributify`,// core options
        shortcuts: [
            {dark_bg_black: "dark:bg-#03080c"},
            {
                dark_bg_gray:
                    "dark:bg-#207191/10 dark:border-solid dark:border-#1e2227 dark:border-1",
            },
            {
                dark_bg_half_black:
                    "dark:bg-#000/50 dark:border-solid dark:border-#1e2227 dark:border-1",
            },
            {
                dark_bg_full_black:
                    "dark:bg-#000/90 dark:border-solid dark:border-#1e2227 dark:border-1",
            },
            {dark_border: "border-solid border-#1e2227 border-1"},
            {dark_border_2: "border-solid border-#1e2227 border-2"},
            {dark_text_white: "dark:c-#fff"},
            {light_border: "border-1 border-gray-2 border-solid"},
            {gray_border: "border-2 border-#d0d7de border-solid"},
            {light_border_bottom: "border-b-1 border-b-gray-2 border-b-solid"},
            {anchor_border: "border-l-2 border-l-#1e80ff border-l-solid"},
            {c_title_blue: "text-#1890ff"},
            {c_text_black: "text-#111"},
            {c_text_white: "text-#fff"},
            {
                menu_item:
                    "px18 py15 text-14 dark:text-#fff hover:bg-#000/20 active:bg-#000/40 c-#000 cursor-pointer",
            },
            {"flex-center": "flex items-center justify-center"},
            {
                custom_cursor_flow:
                    "hover:drop-shadow-xl hover:translate-y--2 duration-100",
            },
            {custom_icon: "hover:text-blue duration-100"},
            {custom_shadow: "hover:shadow-2xl dark:shadow-white duration-100"},
            {
                custom_bottom_border_gray: "border-b-1 border-b-gray-2 border-b-solid",
            },
            {
                custom_bottom_border_1E80FF:
                    "border-b-2 border-b-#1E80FF border-b-solid",
            },
            {custom_border_gray: "border-1 border-gray-2 border-solid"},
            {custom_border_1E80FF: "border-2 border-#1E80FF border-solid"},
            {custom_bg_gray: "bg-gray-1 b-rounded-2"},
            {
                custom_dark_btn_hover:
                    "dark:hover:bg_transparent dark:hover:border-solid dark:hover:border-solid dark:hover:text-#1e80ff dark:hover:border-#1e80ff",
            },
        ],
        rules: [
            ["footer_shadow", {"box-shadow": " 0 0 10px rgba(0, 0, 0, .5)"}],
            ["bg_transparent", {"background-color": "transparent"}],
            ["border_bottom_blue", {"border-bottom": "2px solid #007fff"}],
            ["custom_shadow_all", {"box-shadow": "0 0 10px rgba(0, 0, 0, 0.5)"}],
        ],
        theme: {
            colors: {
                dtc: "hsla(0,0%,100%,0.7)", //
            },
            hovers: {},
        },
        safelist: [
            "i-fa6-brands-x-twitter",
            "i-fa6-brands-facebook",
            "i-fa6-brands-instagram",
            "i-fa6-brands-youtube",
            "i-fa6-brands-bilibili",
            "i-fa6-brands-qq",
            "i-fa6-brands-github",
            "i-fa6-brands:square-git",
            "i-fa6-brands-weixin",
            "i-fa6-brands-zhihu",
            "i-bi-link-45deg",
        ],
    },
    css: [
        "@/styles/main.css",
        "@/styles/github-markdown.css",
    ],
    plugins: [
        "~/plugins/pinia.js",
        "~/plugins/localStorage.client.ts",
        "~/plugins/qrcode.client.ts",
        "~/plugins/routerGuard.ts",
    ],
    mdc: {
        // highlight: {
        //   theme: {
        //     // Default theme (same as single string)
        //     default: "github-light",
        //     // Theme used if `html.dark`
        //     dark: "github-dark",
        //     // Theme used if `html.sepia`
        //     sepia: "monokai",
        //   },
        //   langs: [
        //     "json",
        //     "js",
        //     "ts",
        //     "html",
        //     "css",
        //     "vue",
        //     "shell",
        //     "mdc",
        //     "md",
        //     "yaml",
        //     "go",
        //     "java",
        //     "protoc",
        //   ],
        // },
        headings: {
            anchorLinks: {
                // Enable/Disable heading anchor links. { h1: true, h2: false }
                h1: true,
            },
        },
    },
});
