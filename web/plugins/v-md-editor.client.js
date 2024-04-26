import VMdPreview from "@kangc/v-md-editor/lib/preview";
import githubTheme from "@kangc/v-md-editor/lib/theme/github.js";
// highlightjs 核心代码
import hljs from 'highlight.js/lib/core'
// 按需引入语言包
import c from 'highlight.js/lib/languages/c'
import shell from 'highlight.js/lib/languages/shell'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import python from 'highlight.js/lib/languages/python'
import json from 'highlight.js/lib/languages/json'

// 插件
import createLineNumberPlugin from "@kangc/v-md-editor/lib/plugins/line-number/index";
import createCopyCodePlugin from '@kangc/v-md-editor/lib/plugins/copy-code/index';
import createHighlightLinesPlugin from "@kangc/v-md-editor/lib/plugins/highlight-lines/index";

// 样式
import "@kangc/v-md-editor/lib/style/preview.css";
import "@kangc/v-md-editor/lib/theme/style/github.css";
import "@kangc/v-md-editor/lib/plugins/copy-code/copy-code.css";
import "@kangc/v-md-editor/lib/plugins/highlight-lines/highlight-lines.css";

hljs.registerLanguage('c', c)
hljs.registerLanguage('shell', shell)
hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('python', python)
hljs.registerLanguage('json', json)
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
      .use(createCopyCodePlugin())
      .use(createHighlightLinesPlugin());
    nuxtApp.vueApp.use(VMdPreview);
  }
});
