import VMdPreview from '@kangc/v-md-editor/lib/preview';
import '@kangc/v-md-editor/lib/style/preview.css';
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js';
import '@kangc/v-md-editor/lib/theme/style/github.css';
// 代码行号
import createLineNumbertPlugin from '@kangc/v-md-editor/lib/plugins/line-number/index';
// highlightjs
import hljs from 'highlight.js';
// 代码复制
import createCopyCodePlugin from '@kangc/v-md-editor/lib/plugins/copy-code/index';
import '@kangc/v-md-editor/lib/plugins/copy-code/copy-code.css';
// 代码行高亮
import createHighlightLinesPlugin from '@kangc/v-md-editor/lib/plugins/highlight-lines/index';
import '@kangc/v-md-editor/lib/plugins/highlight-lines/highlight-lines.css';
export default defineNuxtPlugin((nuxtApp) => {
    if (process.client) {
        VMdPreview.use(githubTheme, { Hljs: hljs })
                  .use(createLineNumbertPlugin())
                  .use(createCopyCodePlugin())
                  .use(createHighlightLinesPlugin());
    
        nuxtApp.vueApp.use(VMdPreview);
      }
})




