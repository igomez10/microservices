export function readStored<T>(key: string, fallback: T): T {
  if (typeof window === 'undefined') {
    return fallback
  }
  const raw = window.localStorage.getItem(key)
  if (!raw) {
    return fallback
  }
  try {
    return JSON.parse(raw) as T
  } catch {
    return fallback
  }
}

export function writeStored<T>(key: string, value: T) {
  if (typeof window === 'undefined') {
    return
  }
  window.localStorage.setItem(key, JSON.stringify(value))
}

export function clearStored(key: string) {
  if (typeof window === 'undefined') {
    return
  }
  window.localStorage.removeItem(key)
}

export function readCookie(name: string) {
  if (typeof document === 'undefined') {
    return undefined
  }

  const cookie = document.cookie
    .split('; ')
    .find((entry) => entry.startsWith(`${encodeURIComponent(name)}=`))

  if (!cookie) {
    return undefined
  }

  const [, value = ''] = cookie.split('=')
  return decodeURIComponent(value)
}

export function writeCookie(name: string, value: string, maxAgeSeconds = 3600) {
  if (typeof document === 'undefined') {
    return
  }

  document.cookie = `${encodeURIComponent(name)}=${encodeURIComponent(value)}; Path=/; Max-Age=${maxAgeSeconds}; SameSite=Lax`
}

export function clearCookie(name: string) {
  if (typeof document === 'undefined') {
    return
  }

  document.cookie = `${encodeURIComponent(name)}=; Path=/; Max-Age=0; SameSite=Lax`
}

export function readJwtPayload(token: string) {
  const [, payload] = token.split('.')
  if (!payload) {
    return undefined
  }

  try {
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
    const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, '=')
    const decoded =
      typeof atob === 'function'
        ? atob(padded)
        : Buffer.from(padded, 'base64').toString('utf-8')
    return JSON.parse(decoded) as Record<string, unknown>
  } catch {
    return undefined
  }
}
