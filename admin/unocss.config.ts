import { defineConfig, presetUno, presetIcons } from 'unocss'
import icons from '@iconify-json/bi'
export default defineConfig({
  // 这里可以添加预设、规则等配置
  presets: [
    // 使用预设
    presetUno(),
    presetIcons({
      // 其他选项
      prefix: 'i-',
      extraProperties: {
        display: 'inline-block'
      }
    }),
    icons
  ],
  rules: [
    // 自定义规则
  ]
  // 其他配置项...
})
