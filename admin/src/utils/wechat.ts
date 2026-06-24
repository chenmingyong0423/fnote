import VMdEditor from '@kangc/v-md-editor'

const apiHost = String(import.meta.env.VITE_API_HOST || '').replace(/\/$/, '')

const rootStyle = [
  'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei",Arial,sans-serif',
  'font-size:17px',
  'line-height:1.9',
  'color:#2b2f36',
  'word-break:break-word'
].join(';')

const accentColor = '#2f8cff'
const accentSoftColor = '#eaf4ff'
const accentLighterColor = '#f6fbff'

const elementStyles: Record<string, string> = {
  h1: `margin:34px 0 26px;padding:0 0 10px;border-bottom:2px solid ${accentColor};text-align:center;font-size:26px;line-height:1.35;font-weight:700;color:${accentColor}`,
  h2: `margin:30px 0 18px;padding:4px 0 4px 14px;border-left:5px solid ${accentColor};font-size:22px;line-height:1.4;font-weight:700;color:${accentColor};background:linear-gradient(90deg, ${accentSoftColor}, rgba(255,255,255,0))`,
  h3: `margin:26px 0 14px;font-size:19px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h4: `margin:22px 0 12px;font-size:18px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h5: `margin:20px 0 10px;font-size:17px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h6: 'margin:18px 0 8px;font-size:16px;line-height:1.45;font-weight:700;color:#4b89a8',
  p: 'margin:0 0 18px;font-size:17px;line-height:1.9;color:#2b2f36',
  a: `color:${accentColor};text-decoration:none`,
  strong: `font-weight:700;color:${accentColor}`,
  em: 'font-style:italic',
  blockquote:
    `margin:24px 0;padding:18px 20px;border:1px solid #d8eaff;border-radius:10px;background:linear-gradient(135deg, ${accentLighterColor} 0%, #ffffff 58%, #ffffff 100%);box-shadow:0 10px 28px rgba(47,140,255,0.12);color:#2b2f36`,
  ul: 'margin:0 0 18px;padding-left:1.4em',
  ol: 'margin:0 0 18px;padding-left:1.4em',
  li: 'margin:6px 0;line-height:1.9',
  img: 'display:block;max-width:100%;height:auto;margin:16px auto;border-radius:4px',
  pre: `margin:0;padding:34px 16px 28px;overflow-x:auto;border-radius:8px;border:1px solid #d8eaff;background:${accentLighterColor};color:#1f6f9d;font-size:14px;line-height:1.7;white-space:pre-wrap;tab-size:2`,
  code: `padding:2px 6px;border-radius:4px;background:${accentSoftColor};color:#1677ff;font-family:Menlo,Consolas,Monaco,monospace;font-size:0.92em`,
  table: 'width:100%;margin:18px 0;border-collapse:collapse;font-size:14px;line-height:1.6',
  th: `padding:8px 10px;border:1px solid #cae4f4;background:${accentLighterColor};font-weight:700;text-align:left;color:${accentColor}`,
  td: 'padding:8px 10px;border:1px solid #cae4f4;text-align:left',
  hr: `height:1px;margin:26px 0;border:0;background:${accentColor}`
}

const blockquoteParagraphStyle = 'margin:0;font-size:17px;line-height:1.9;color:#2b2f36'
const preCodeStyle =
  'display:block;padding:0;background:transparent;color:inherit;font-family:Menlo,Consolas,Monaco,monospace;font-size:14px;line-height:1.7;white-space:pre-wrap;tab-size:2'
const codeBlockWrapperStyle = 'position:relative;margin:18px 0'
const codeLanguageBadgeStyle = `position:absolute;right:12px;top:8px;z-index:1;padding:1px 8px;border-radius:999px;background:#ffffff;color:${accentColor};font-size:12px;font-family:Menlo,Consolas,Monaco,monospace;line-height:1.7;border:1px solid #d8eaff`
const codeSignatureBadgeStyle =
  'position:absolute;right:12px;bottom:7px;z-index:1;color:#8a96a3;font-size:12px;line-height:1.7'

const highlightStyles: Array<{ classes: string[]; style: string }> = [
  {
    classes: ['hljs-keyword', 'hljs-selector-tag', 'hljs-doctag', 'hljs-meta', 'hljs-template-tag'],
    style: 'color:#d73a49;font-weight:600'
  },
  {
    classes: ['hljs-string', 'hljs-regexp', 'hljs-symbol', 'hljs-bullet', 'hljs-addition'],
    style: 'color:#22863a'
  },
  {
    classes: ['hljs-title', 'hljs-section', 'hljs-selector-id', 'hljs-function'],
    style: 'color:#6f42c1;font-weight:600'
  },
  {
    classes: ['hljs-attr', 'hljs-attribute', 'hljs-variable', 'hljs-template-variable'],
    style: 'color:#005cc5'
  },
  {
    classes: ['hljs-number', 'hljs-literal', 'hljs-type', 'hljs-built_in', 'hljs-builtin-name'],
    style: 'color:#e36209'
  },
  {
    classes: ['hljs-comment', 'hljs-quote', 'hljs-deletion'],
    style: 'color:#6a737d;font-style:italic'
  }
]

type CopyWechatOptions = {
  author?: string
}

const getSignatureText = (author?: string) => {
  const normalizedAuthor = author?.trim()

  if (!normalizedAuthor) {
    return '程序员陈明勇'
  }

  return normalizedAuthor.startsWith('程序员') ? normalizedAuthor : `程序员${normalizedAuthor}`
}

const languageNames: Record<string, string> = {
  bash: 'Shell',
  c: 'C',
  cpp: 'C++',
  css: 'CSS',
  go: 'Go',
  golang: 'Go',
  html: 'HTML',
  java: 'Java',
  javascript: 'JavaScript',
  js: 'JavaScript',
  json: 'JSON',
  markdown: 'Markdown',
  md: 'Markdown',
  python: 'Python',
  py: 'Python',
  shell: 'Shell',
  sh: 'Shell',
  sql: 'SQL',
  ts: 'TypeScript',
  typescript: 'TypeScript',
  xml: 'XML',
  yaml: 'YAML',
  yml: 'YAML'
}

const getCodeBlockLanguage = (pre: Element) => {
  const code = pre.querySelector('code')
  const classNames = [...pre.classList, ...(code ? [...code.classList] : [])]
  const languageClass = classNames.find((className) => /^(language|lang)-/.test(className))
  const rawLanguage = languageClass?.replace(/^(language|lang)-/, '').trim().toLowerCase()

  if (!rawLanguage) {
    return 'Code'
  }

  return languageNames[rawLanguage] || rawLanguage.toUpperCase()
}

const formatLanguageName = (language?: string) => {
  const rawLanguage = language?.trim().toLowerCase()

  if (!rawLanguage) {
    return ''
  }

  return languageNames[rawLanguage] || rawLanguage.toUpperCase()
}

const extractFenceLanguages = (markdown: string) => {
  return [...markdown.matchAll(/^```([^\s`]*)[^\n]*$/gm)].map((match) =>
    formatLanguageName(match[1])
  )
}

const escapeHtml = (content: string) =>
  content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

const renderMarkdown = (markdown: string) => {
  try {
    const parser = (VMdEditor as any).vMdParser
    const html = parser?.parse?.(markdown)

    if (typeof html === 'string') {
      return html
    }
  } catch (error) {
    console.warn('Markdown parser is unavailable, falling back to plain text rendering.', error)
  }

  return markdown
    .split(/\n{2,}/)
    .map((paragraph) => `<p>${escapeHtml(paragraph).replace(/\n/g, '<br>')}</p>`)
    .join('')
}

const stripOpeningQrCodeBlock = (markdown: string) => {
  const normalizedMarkdown = markdown.replace(/\r\n/g, '\n')

  return normalizedMarkdown.replace(
    /^\s*>[ \t]*扫码关注公众号，手机阅读更方便[ \t]*\n(?:>[ \t]*\n)?[ \t]*>[ \t]*!\[[^\]]*]\([^)]*wx-gzh-qrcode[^)]*\)[ \t]*(?:\n+|$)/,
    ''
  )
}

const toAbsoluteUrl = (url: string) => {
  if (!url || /^(https?:|data:|blob:|mailto:|#)/i.test(url)) {
    return url
  }

  if (url.startsWith('//')) {
    return `${window.location.protocol}${url}`
  }

  const baseUrl = apiHost || window.location.origin
  return new URL(url, `${baseUrl}/`).href
}

const applyStyle = (element: Element, style: string) => {
  const currentStyle = element.getAttribute('style')
  element.setAttribute('style', currentStyle ? `${currentStyle};${style}` : style)
}

const applyHighlightStyle = (element: Element) => {
  if (!element.closest('pre')) {
    return
  }

  const highlightStyle = highlightStyles.find((item) =>
    item.classes.some((className) => element.classList.contains(className))
  )

  if (highlightStyle) {
    applyStyle(element, highlightStyle.style)
  }
}

const decorateCodeBlock = (pre: Element, signatureText: string, language?: string) => {
  if (pre.parentElement?.getAttribute('data-wechat-code-block') === 'true') {
    return
  }

  const wrapper = document.createElement('section')
  wrapper.setAttribute('data-wechat-code-block', 'true')
  wrapper.setAttribute('style', codeBlockWrapperStyle)

  const languageBadge = document.createElement('section')
  languageBadge.textContent = language || getCodeBlockLanguage(pre)
  languageBadge.setAttribute('style', codeLanguageBadgeStyle)

  const signatureBadge = document.createElement('section')
  signatureBadge.textContent = signatureText
  signatureBadge.setAttribute('style', codeSignatureBadgeStyle)

  pre.parentNode?.insertBefore(wrapper, pre)
  wrapper.append(languageBadge, pre, signatureBadge)
}

const normalizeElement = (element: Element, signatureText: string, codeLanguages: string[]) => {
  const tagName = element.tagName.toLowerCase()
  const isCodeInPre = tagName === 'code' && Boolean(element.closest('pre'))
  const style = isCodeInPre ? preCodeStyle : elementStyles[tagName]

  if (style) {
    applyStyle(element, style)
  }

  applyHighlightStyle(element)

  if (tagName === 'blockquote') {
    const marker = document.createElement('section')
    marker.innerHTML = `<span style="display:inline-block;padding:1px 8px;border-radius:999px;background:${accentSoftColor};color:${accentColor};font-size:12px;font-weight:700;letter-spacing:0.5px;line-height:1.8">${escapeHtml(signatureText)}</span><span style="display:inline-block;width:44px;height:2px;margin-left:8px;vertical-align:middle;border-radius:999px;background:${accentColor}"></span>`
    marker.setAttribute('style', 'margin:0 0 12px;line-height:1')
    element.prepend(marker)
  }

  if (tagName === 'p' && element.closest('blockquote')) {
    element.setAttribute('style', blockquoteParagraphStyle)
  }

  if (tagName === 'img') {
    const src = element.getAttribute('src') || ''
    element.setAttribute('src', toAbsoluteUrl(src))
  }

  if (tagName === 'a') {
    const href = element.getAttribute('href') || ''
    element.setAttribute('href', toAbsoluteUrl(href))
  }

  if (tagName === 'pre') {
    element.querySelectorAll('code').forEach((code) => {
      code.setAttribute('style', preCodeStyle)
    })
    decorateCodeBlock(element, signatureText, codeLanguages.shift())
  }

  element.removeAttribute('class')
}

const buildWechatHtml = (markdown: string, options: CopyWechatOptions = {}) => {
  const template = document.createElement('template')
  template.innerHTML = renderMarkdown(markdown)
  const signatureText = getSignatureText(options.author)
  const codeLanguages = extractFenceLanguages(markdown)
  template.content
    .querySelectorAll('*')
    .forEach((element) => normalizeElement(element, signatureText, codeLanguages))

  const root = document.createElement('section')
  root.setAttribute('style', rootStyle)
  root.append(...Array.from(template.content.childNodes))

  return root.outerHTML
}

const getPlainText = (html: string) => {
  const container = document.createElement('div')
  container.innerHTML = html
  return container.innerText.trim()
}

const copyHtmlWithSelection = async (html: string, plainText: string) => {
  const container = document.createElement('div')
  container.contentEditable = 'true'
  container.style.position = 'fixed'
  container.style.left = '-9999px'
  container.style.top = '0'
  container.innerHTML = html
  document.body.appendChild(container)

  const selection = window.getSelection()
  const range = document.createRange()
  range.selectNodeContents(container)
  selection?.removeAllRanges()
  selection?.addRange(range)

  const copied = document.execCommand('copy')
  selection?.removeAllRanges()
  document.body.removeChild(container)

  if (!copied) {
    await navigator.clipboard.writeText(plainText)
  }
}

const writeRichClipboard = async (html: string, plainText: string) => {
  if (navigator.clipboard?.write && typeof ClipboardItem !== 'undefined') {
    await navigator.clipboard.write([
      new ClipboardItem({
        'text/html': new Blob([html], { type: 'text/html' }),
        'text/plain': new Blob([plainText], { type: 'text/plain' })
      })
    ])
    return
  }

  await copyHtmlWithSelection(html, plainText)
}

export const copyMarkdownAsWechat = async (markdown: string, options: CopyWechatOptions = {}) => {
  const html = buildWechatHtml(stripOpeningQrCodeBlock(markdown), options)
  await writeRichClipboard(html, getPlainText(html))
}
