export function resolvePublicUrl(value: string): string
export function resolvePublicUrl(value?: string | null): string | undefined
export function resolvePublicUrl(value?: string | null) {
  if (!value) return undefined

  try {
    return new URL(value).toString()
  } catch {
    const base = (
      process.env.NEXT_PUBLIC_SERVER_HOST ||
      process.env.BASE_HOST ||
      ''
    ).replace(/\/$/, '')

    if (!base) return value

    return `${base}${value.startsWith('/') ? '' : '/'}${value}`
  }
}
