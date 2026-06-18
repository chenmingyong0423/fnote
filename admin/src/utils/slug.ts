export const normalizeSlug = (value: string) => {
  return value
    .trim()
    .toLowerCase()
    .replace(/\s+/g, '-')
    .replace(/[^a-z0-9-]/g, '')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '')
}

export const shouldAutoFillSlug = (value: string) => {
  return !/[\u4e00-\u9fa5]/.test(value)
}

export const validateSlug = (value: string, label: string) => {
  if (!value) {
    throw new Error(`请输入${label}`)
  }
  if (!/^[a-z0-9]+(?:-[a-z0-9]+)*$/.test(value)) {
    throw new Error(`${label}仅支持小写英文、数字和短横线，且不能以短横线开头或结尾`)
  }
}

export const createSlugRule = (label: string) => ({
  validator: async (_rule: unknown, value: string) => {
    validateSlug(value, label)
  }
})
