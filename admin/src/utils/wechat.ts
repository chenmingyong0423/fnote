import VMdEditor from '@kangc/v-md-editor'
import githubTheme from '@kangc/v-md-editor/lib/theme/github.js'
import hljs from 'highlight.js/lib/core'
import c from 'highlight.js/lib/languages/c'
import css from 'highlight.js/lib/languages/css'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import python from 'highlight.js/lib/languages/python'
import shell from 'highlight.js/lib/languages/shell'
import typescript from 'highlight.js/lib/languages/typescript'
import xml from 'highlight.js/lib/languages/xml'

hljs.registerLanguage('c', c)
hljs.registerLanguage('css', css)
hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('json', json)
hljs.registerLanguage('python', python)
hljs.registerLanguage('shell', shell)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('xml', xml)

VMdEditor.use(githubTheme, {
  Hljs: hljs
})

const apiHost = String(import.meta.env.VITE_API_HOST || '').replace(/\/$/, '')

const leftAlignedTextStyle =
  'text-align:left;text-align-last:left;letter-spacing:0;word-spacing:0'
const rootStyle = [
  'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI","PingFang SC","Hiragino Sans GB","Microsoft YaHei",Arial,sans-serif',
  'font-size:17px',
  'line-height:1.9',
  'word-break:normal',
  'overflow-wrap:break-word'
].join(';')

const accentColor = '#2f8cff'
const accentSoftColor = 'rgba(47,140,255,0.10)'
const accentLighterColor = 'rgba(47,140,255,0.06)'
const accentBorderColor = 'rgba(47,140,255,0.28)'
const subtleBorderColor = 'rgba(47,140,255,0.18)'
const mutedTextColor = '#7a8794'

const elementStyles: Record<string, string> = {
  h1: `margin:34px 0 26px;padding:0 0 10px;border-bottom:2px solid ${accentColor};text-align:center;text-align-last:center;font-size:26px;line-height:1.35;font-weight:700;color:${accentColor}`,
  h2: `margin:30px 0 18px;padding:4px 0 4px 14px;border-left:5px solid ${accentColor};font-size:22px;line-height:1.4;font-weight:700;color:${accentColor};background:linear-gradient(90deg, ${accentSoftColor}, transparent)`,
  h3: `margin:26px 0 14px;font-size:19px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h4: `margin:22px 0 12px;font-size:18px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h5: `margin:20px 0 10px;font-size:17px;line-height:1.45;font-weight:700;color:${accentColor}`,
  h6: 'margin:18px 0 8px;font-size:16px;line-height:1.45;font-weight:700;color:#4b89a8',
  p: `margin:0 0 18px;font-size:17px;line-height:1.9;${leftAlignedTextStyle}`,
  a: `color:${accentColor};text-decoration:none`,
  strong: `font-weight:700;color:${accentColor}`,
  em: 'font-style:italic',
  blockquote:
    `margin:24px 0;padding:18px 20px;border:1px solid ${accentBorderColor};border-radius:10px;background:${accentLighterColor};box-shadow:0 2px 6px rgba(47,140,255,0.16),0 8px 18px rgba(47,140,255,0.14);${leftAlignedTextStyle}`,
  ul: 'margin:0 0 18px;padding-left:1.4em',
  ol: 'margin:0 0 18px;padding-left:1.4em',
  li: 'margin:6px 0;line-height:1.9',
  img: 'display:block;max-width:100%;height:auto;margin:16px auto;border-radius:4px',
  pre: `margin:0;padding:18px 16px 4px;overflow-x:auto;border:0;background:transparent;color:inherit;font-size:14px;line-height:1.7;white-space:pre-wrap;tab-size:2`,
  code: `display:inline;padding:2px 6px;border-radius:4px;background:${accentSoftColor};color:#1677ff;font-family:Menlo,Consolas,Monaco,monospace;font-size:0.92em;line-height:1.6;white-space:normal;vertical-align:baseline`,
  table: 'width:100%;margin:18px 0;border-collapse:collapse;font-size:14px;line-height:1.6',
  th: `padding:8px 10px;border:1px solid ${subtleBorderColor};background:${accentLighterColor};font-weight:700;text-align:left;color:${accentColor}`,
  td: `padding:8px 10px;border:1px solid ${subtleBorderColor};text-align:left`,
  hr: `height:1px;margin:26px 0;border:0;background:${accentColor}`
}

const blockquoteParagraphStyle =
  `margin:0;font-size:17px;line-height:1.9;${leftAlignedTextStyle}`
const listItemParagraphStyle = `margin:0;font-size:17px;line-height:1.9;${leftAlignedTextStyle}`
const preCodeStyle =
  'display:block;padding:0;background:transparent;color:inherit;font-family:Menlo,Consolas,Monaco,monospace;font-size:14px;line-height:1.7;white-space:pre-wrap;tab-size:2'
const codeLineStyle =
  'min-height:1.7em;font-family:Menlo,Consolas,Monaco,monospace;font-size:14px;line-height:1.7;white-space:nowrap'
const codeBlockWrapperStyle =
  `margin:18px 0;overflow:hidden;border:1px solid ${accentBorderColor};border-radius:8px;background:${accentLighterColor}`
const codeSignatureBadgeStyle =
  `display:block;padding:2px 12px 10px;text-align:right;color:${mutedTextColor};font-size:12px;line-height:1.2`

const highlightStyles: Array<{ classes: string[]; style: string }> = [
  {
    classes: ['hljs-keyword', 'hljs-selector-tag', 'hljs-doctag', 'hljs-meta', 'hljs-template-tag'],
    style: 'color:#e25563;font-weight:700'
  },
  {
    classes: ['hljs-string', 'hljs-regexp', 'hljs-symbol', 'hljs-bullet', 'hljs-addition'],
    style: 'color:#2fb344'
  },
  {
    classes: ['hljs-title', 'hljs-section', 'hljs-selector-id', 'hljs-function'],
    style: 'color:#8b5cf6;font-weight:700'
  },
  {
    classes: ['hljs-attr', 'hljs-attribute', 'hljs-variable', 'hljs-template-variable'],
    style: 'color:#2f8cff;font-weight:600'
  },
  {
    classes: ['hljs-number', 'hljs-literal', 'hljs-type', 'hljs-built_in', 'hljs-builtin-name'],
    style: 'color:#f59e0b;font-weight:600'
  },
  {
    classes: ['hljs-comment', 'hljs-quote', 'hljs-deletion'],
    style: 'color:#8b949e;font-style:italic'
  },
  {
    classes: ['hljs-params', 'hljs-property', 'hljs-name', 'hljs-tag'],
    style: 'color:#38bdf8'
  },
  {
    classes: ['hljs-operator', 'hljs-punctuation'],
    style: 'color:#94a3b8'
  }
]

const getHighlightStyle = (element: Element) => {
  return highlightStyles.find((item) =>
    item.classes.some((className) => element.classList.contains(className))
  )?.style
}

type CopyWechatOptions = {
  author?: string
}

type CodeFence = {
  language?: string
  code: string
}

type PreparedMarkdown = {
  markdown: string
  fences: CodeFence[]
}

const getSignatureText = (author?: string) => {
  const normalizedAuthor = author?.trim()

  if (!normalizedAuthor) {
    return '程序员陈明勇'
  }

  return normalizedAuthor.startsWith('程序员') ? normalizedAuthor : `程序员${normalizedAuthor}`
}

const languageAliases: Record<string, string> = {
  bash: 'shell',
  golang: 'go',
  html: 'xml',
  js: 'javascript',
  jsx: 'javascript',
  py: 'python',
  sh: 'shell',
  ts: 'typescript',
  tsx: 'typescript'
}

const getFenceLanguage = (info: string) => {
  return info.trim().split(/\s+/)[0]?.replace(/[{}]/g, '').toLowerCase() || ''
}

const normalizeHighlightLanguage = (language?: string) => {
  const rawLanguage = language?.trim().toLowerCase()

  if (!rawLanguage) {
    return ''
  }

  return languageAliases[rawLanguage] || rawLanguage
}

const prepareCodeFences = (markdown: string): PreparedMarkdown => {
  const fences: CodeFence[] = []
  const lines = markdown.replace(/\r\n?/g, '\n').split('\n')
  const outputLines: string[] = []
  let index = 0

  while (index < lines.length) {
    const openingMatch = lines[index].match(/^( {0,3})(`{3,}|~{3,})(.*)$/)

    if (!openingMatch) {
      outputLines.push(lines[index])
      index += 1
      continue
    }

    const fenceMarker = openingMatch[2]
    const language = getFenceLanguage(openingMatch[3] || '')
    const fenceChar = fenceMarker[0]
    const minFenceLength = fenceMarker.length
    const closingPattern = new RegExp(`^ {0,3}\\${fenceChar}{${minFenceLength},}[ \\t]*$`)
    const codeLines: string[] = []
    index += 1

    while (index < lines.length && !closingPattern.test(lines[index])) {
      codeLines.push(lines[index])
      index += 1
    }

    const fenceId = fences.length
    fences.push({
      language,
      code: codeLines.join('\n')
    })
    outputLines.push(`<section data-wechat-code-id="${fenceId}"></section>`)

    if (index < lines.length) {
      index += 1
    }
  }

  return {
    markdown: outputLines.join('\n'),
    fences
  }
}

const escapeHtml = (content: string) =>
  content
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#39;')

const escapeCodeText = (content: string) =>
  escapeHtml(content)
    .replace(/\t/g, '&nbsp;&nbsp;')
    .replace(/ {2}/g, '&nbsp;&nbsp;')
    .replace(/ /g, '&nbsp;')

const highlightCode = (code: string, language?: string) => {
  const highlightLanguage = normalizeHighlightLanguage(language)

  if (highlightLanguage && hljs.getLanguage(highlightLanguage)) {
    try {
      return hljs.highlight(code, { language: highlightLanguage, ignoreIllegals: true }).value
    } catch (error) {
      console.warn('Code highlight failed, falling back to auto detection.', error)
    }
  }

  try {
    return hljs.highlightAuto(code).value
  } catch (error) {
    console.warn('Code highlight failed, falling back to escaped plain text.', error)
  }

  return escapeHtml(code)
}

const preserveCodeWhitespace = (html: string) => {
  const template = document.createElement('template')
  template.innerHTML = html

  const preserveNode = (node: Node) => {
    if (node.nodeType === Node.TEXT_NODE) {
      const text = node.textContent || ''
      const fragment = document.createDocumentFragment()
      const parts = text.split('\n')

      parts.forEach((part, index) => {
        if (index > 0) {
          fragment.appendChild(document.createElement('br'))
        }
        if (part) {
          const span = document.createElement('span')
          span.innerHTML = escapeCodeText(part)
          fragment.appendChild(span)
        }
      })
      node.parentNode?.replaceChild(fragment, node)
      return
    }

    Array.from(node.childNodes).forEach(preserveNode)
  }

  Array.from(template.content.childNodes).forEach(preserveNode)
  return template.innerHTML
}

const codeHtmlToLines = (html: string) => {
  return preserveCodeWhitespace(html)
    .split(/<br\s*\/?>/i)
    .map((line) => line || '&nbsp;')
}

const createCodeBlockContent = (fence?: CodeFence) => {
  if (!fence) {
    return null
  }

  const content = document.createElement('section')
  content.setAttribute('style', preCodeStyle)

  codeHtmlToLines(highlightCode(fence.code, fence.language)).forEach((lineHtml) => {
    const line = document.createElement('section')
    line.setAttribute('style', codeLineStyle)
    line.innerHTML = lineHtml
    line.querySelectorAll('*').forEach((token) => {
      const highlightStyle = getHighlightStyle(token)
      if (highlightStyle) {
        token.setAttribute('style', highlightStyle)
      }
      token.removeAttribute('class')
    })
    content.appendChild(line)
  })

  return content
}

const createCodeBlockElement = (signatureText: string, fence: CodeFence) => {
  const wrapper = document.createElement('section')
  wrapper.setAttribute('data-wechat-code-block', 'true')
  wrapper.setAttribute('style', codeBlockWrapperStyle)

  const pre = document.createElement('pre')
  pre.setAttribute('style', elementStyles.pre)
  const codeContent = createCodeBlockContent(fence)
  if (codeContent) {
    pre.appendChild(codeContent)
  }

  const signatureBadge = document.createElement('section')
  signatureBadge.textContent = signatureText
  signatureBadge.setAttribute('style', codeSignatureBadgeStyle)

  wrapper.append(pre, signatureBadge)
  return wrapper
}

const renderMarkdown = (markdown: string) => {
  try {
    const parser = (VMdEditor as any).vMdParser
    const markdownParser = parser?.themeConfig?.markdownParser
    const originalBreaks = markdownParser?.options?.breaks
    const canConfigureBreaks = markdownParser?.set && typeof originalBreaks === 'boolean'

    if (canConfigureBreaks) {
      markdownParser.set({ breaks: false })
    }

    let html: unknown
    try {
      html = parser?.parse?.(markdown)
    } finally {
      if (canConfigureBreaks) {
        markdownParser.set({ breaks: originalBreaks })
      }
    }

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

const getMarkdownLinkText = (element: Element) => {
  const text = element.textContent?.trim()
  return text || element.getAttribute('href') || ''
}

const replaceTextLinkWithMarkdownSource = (element: Element) => {
  if (element.querySelector('img')) {
    return false
  }

  const href = element.getAttribute('href') || ''
  const text = getMarkdownLinkText(element)
  const linkSource = document.createElement('span')
  linkSource.innerHTML = `[${escapeHtml(text)}](${escapeHtml(href)})`
  linkSource.setAttribute(
    'style',
    'display:inline;white-space:normal;overflow-wrap:anywhere;word-break:break-word;vertical-align:baseline'
  )

  element.replaceWith(linkSource)
  return true
}

const applyStyle = (element: Element, style: string) => {
  const currentStyle = element.getAttribute('style')
  element.setAttribute('style', currentStyle ? `${currentStyle};${style}` : style)
}

const applyHighlightStyle = (element: Element) => {
  if (!element.closest('pre')) {
    return
  }

  const highlightStyle = getHighlightStyle(element)

  if (highlightStyle) {
    applyStyle(element, highlightStyle)
  }
}

const startsWithInlineContinuation = (content?: string | null) => {
  const text = content?.trimStart() || ''
  return /^[：:、，,；;）)\/\\]/.test(text)
}

const getNextMeaningfulSibling = (node: Node) => {
  let sibling = node.nextSibling

  while (sibling && sibling.nodeType === Node.TEXT_NODE && !sibling.textContent?.trim()) {
    sibling = sibling.nextSibling
  }

  return sibling
}

const isListItemBlockChild = (node: Node) => {
  if (!(node instanceof HTMLElement)) {
    return false
  }

  return [
    'address',
    'article',
    'aside',
    'blockquote',
    'div',
    'dl',
    'figure',
    'h1',
    'h2',
    'h3',
    'h4',
    'h5',
    'h6',
    'hr',
    'ol',
    'p',
    'pre',
    'section',
    'table',
    'ul'
  ].includes(node.tagName.toLowerCase())
}

const wrapInlineListItemContent = (element: Element) => {
  const nodes = Array.from(element.childNodes)
  let paragraph: HTMLParagraphElement | null = null

  const closeParagraph = () => {
    paragraph = null
  }

  nodes.forEach((node) => {
    if (!node.parentNode) {
      return
    }

    if (node.nodeType === Node.TEXT_NODE && !node.textContent?.trim()) {
      if (paragraph) {
        paragraph.appendChild(node)
      }
      return
    }

    if (isListItemBlockChild(node)) {
      closeParagraph()
      return
    }

    if (!paragraph) {
      paragraph = document.createElement('p')
      paragraph.setAttribute('style', listItemParagraphStyle)
      element.insertBefore(paragraph, node)
    }

    paragraph.appendChild(node)
  })
}

const mergeListItemColonDescriptions = (element: Element) => {
  wrapInlineListItemContent(element)

  const paragraphs = Array.from(element.children).filter(
    (child) => child.tagName.toLowerCase() === 'p'
  )

  paragraphs.forEach((paragraph) => {
    const next = getNextMeaningfulSibling(paragraph)

    if (!(next instanceof HTMLElement) || next.tagName.toLowerCase() !== 'p') {
      return
    }

    if (!startsWithInlineContinuation(next.textContent)) {
      return
    }

    while (next.firstChild) {
      paragraph.appendChild(next.firstChild)
    }
    next.remove()
  })

  element.querySelectorAll('br').forEach((lineBreak) => {
    const next = getNextMeaningfulSibling(lineBreak)

    if (next?.nodeType === Node.TEXT_NODE && startsWithInlineContinuation(next.textContent)) {
      lineBreak.remove()
    }
  })

}

const decorateCodeBlock = (pre: Element, signatureText: string, fence?: CodeFence) => {
  if (pre.parentElement?.getAttribute('data-wechat-code-block') === 'true') {
    return
  }

  const codeContent = createCodeBlockContent(fence || { code: pre.textContent || '' })
  if (codeContent) {
    pre.replaceChildren(codeContent)
  }

  const wrapper = document.createElement('section')
  wrapper.setAttribute('data-wechat-code-block', 'true')
  wrapper.setAttribute('style', codeBlockWrapperStyle)

  const signatureBadge = document.createElement('section')
  signatureBadge.textContent = signatureText
  signatureBadge.setAttribute('style', codeSignatureBadgeStyle)

  pre.parentNode?.insertBefore(wrapper, pre)
  wrapper.append(pre, signatureBadge)
}

const normalizeElement = (element: Element, signatureText: string, codeFences: CodeFence[]) => {
  const tagName = element.tagName.toLowerCase()
  const codeFenceId = element.getAttribute('data-wechat-code-id')

  if (codeFenceId !== null) {
    const fence = codeFences[Number(codeFenceId)]
    if (fence) {
      element.replaceWith(createCodeBlockElement(signatureText, fence))
    } else {
      element.remove()
    }
    return
  }

  if (tagName === 'li') {
    mergeListItemColonDescriptions(element)
  }

  const isCodeInPre = tagName === 'code' && Boolean(element.closest('pre'))

  if (tagName === 'code' && !isCodeInPre) {
    const inlineCode = document.createElement('span')
    inlineCode.innerHTML = element.innerHTML
    inlineCode.setAttribute('style', elementStyles.code)
    element.replaceWith(inlineCode)
    return
  }

  const style = isCodeInPre ? preCodeStyle : elementStyles[tagName]

  if (style) {
    applyStyle(element, style)
  }

  applyHighlightStyle(element)

  if (tagName === 'blockquote') {
    const marker = document.createElement('section')
    marker.innerHTML = `<span style="display:inline-block;padding:1px 8px;border-radius:999px;background:${accentSoftColor};color:${accentColor};font-size:12px;font-weight:700;letter-spacing:0.5px;line-height:1.8">${escapeHtml(signatureText)}</span>`
    marker.setAttribute('style', 'margin:0 0 12px;line-height:1')
    element.prepend(marker)
  }

  if (tagName === 'p' && element.closest('blockquote')) {
    element.setAttribute('style', blockquoteParagraphStyle)
  }

  if (tagName === 'p' && element.closest('li')) {
    element.setAttribute('style', listItemParagraphStyle)
  }

  if (tagName === 'img') {
    const src = element.getAttribute('src') || ''
    element.setAttribute('src', toAbsoluteUrl(src))
  }

  if (tagName === 'a') {
    if (replaceTextLinkWithMarkdownSource(element)) {
      return
    }

    const href = element.getAttribute('href') || ''
    element.setAttribute('href', toAbsoluteUrl(href))
  }

  if (tagName === 'pre') {
    element.querySelectorAll('code').forEach((code) => {
      code.setAttribute('style', preCodeStyle)
    })
    decorateCodeBlock(element, signatureText, codeFences.shift())
  }

  element.removeAttribute('class')
}

const buildWechatHtml = (markdown: string, options: CopyWechatOptions = {}) => {
  const preparedMarkdown = prepareCodeFences(markdown)
  const template = document.createElement('template')
  template.innerHTML = renderMarkdown(preparedMarkdown.markdown)
  const signatureText = getSignatureText(options.author)
  const codeFences = preparedMarkdown.fences
  template.content
    .querySelectorAll('*')
    .forEach((element) => normalizeElement(element, signatureText, codeFences))

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
