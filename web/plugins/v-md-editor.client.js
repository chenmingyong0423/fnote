import VMdPreview from '@kangc/v-md-editor/lib/preview';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import hljs from 'highlight.js';

// 插件
import createLineNumberPlugin from '@kangc/v-md-editor/lib/plugins/line-number/index';
import createCopyCodePreview from '@kangc/v-md-editor/lib/plugins/copy-code/preview';
import createHighlightLinesPlugin from '@kangc/v-md-editor/lib/plugins/highlight-lines/index';

// 样式
import '@kangc/v-md-editor/lib/style/preview.css';
import '@kangc/v-md-editor/lib/theme/style/github.css';
import '@kangc/v-md-editor/lib/plugins/copy-code/copy-code.css';
import '@kangc/v-md-editor/lib/plugins/highlight-lines/highlight-lines.css';

export default defineNuxtPlugin((nuxtApp) => {
    if (process.client) {
        VMdPreview.use(githubTheme, {
            Hljs: hljs,
            extend(md) {
                // md为 markdown-it 实例，可以在此处进行修改配置,并使用 plugin 进行语法扩展
                // md.set(option).use(plugin);
                // 加 class 标签
            },
        })
            .use(createLineNumberPlugin())
            .use(createCopyCodePreview())
            .use(createHighlightLinesPlugin());
        nuxtApp.vueApp.use(VMdPreview);
    }
});
