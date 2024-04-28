import { defineConfig, presetUno, presetIcons } from 'unocss'
export default defineConfig({
  // 这里可以添加预设、规则等配置
  presets: [
    // 使用预设
    presetUno(),
    presetIcons()
  ],
  shortcuts: [
    { black_border: "border-solid border-#000 border-1" }
  ],
  rules: [
    // 自定义规则
  ]
  // 其他配置项...
})
