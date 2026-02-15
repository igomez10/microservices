export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE'

export type ApiResponse<T = unknown> = {
  ok: boolean
  status: number
  statusText: string
  data?: T
  error?: string
  rawText?: string
  headers: Record<string, string>
}

export type ApiRequestOptions = {
  baseUrl: string
  path: string
  method: HttpMethod
  pathParams?: Record<string, string | number>
  queryParams?: Record<string, string | number | undefined>
  body?: unknown
  token?: string
  basicAuth?: { clientId: string; clientSecret: string }
  redirect?: RequestRedirect
}

export function buildBasicAuthHeader(clientId: string, clientSecret: string) {
  const raw = `${clientId}:${clientSecret}`
  const nodeBuffer = (globalThis as { Buffer?: { from: (input: string, encoding: string) => { toString: (encoding: string) => string } } }).Buffer
  const encoded =
    typeof btoa === 'function'
      ? btoa(raw)
      : nodeBuffer
        ? nodeBuffer.from(raw, 'utf-8').toString('base64')
        : raw
  return `Basic ${encoded}`
}

export function buildUrl(
  baseUrl: string,
  path: string,
  pathParams?: Record<string, string | number>,
  queryParams?: Record<string, string | number | undefined>
) {
  const base = baseUrl.endsWith('/') ? baseUrl : `${baseUrl}/`
  let resolvedPath = path.startsWith('/') ? path.slice(1) : path
  if (pathParams) {
    for (const [key, value] of Object.entries(pathParams)) {
      resolvedPath = resolvedPath.replace(`{${key}}`, encodeURIComponent(String(value)))
    }
  }

  const url = new URL(resolvedPath, base)
  if (queryParams) {
    for (const [key, value] of Object.entries(queryParams)) {
      if (value !== undefined && value !== '') {
        url.searchParams.set(key, String(value))
      }
    }
  }

  return url.toString()
}

export async function apiRequest<T>(options: ApiRequestOptions): Promise<ApiResponse<T>> {
  const url = buildUrl(options.baseUrl, options.path, options.pathParams, options.queryParams)
  const headers = new Headers({
    Accept: 'application/json'
  })

  if (options.body !== undefined) {
    headers.set('Content-Type', 'application/json')
  }

  if (options.token) {
    headers.set('Authorization', `Bearer ${options.token}`)
  } else if (options.basicAuth) {
    headers.set('Authorization', buildBasicAuthHeader(options.basicAuth.clientId, options.basicAuth.clientSecret))
  }

  const response = await fetch(url, {
    method: options.method,
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
    redirect: options.redirect
  })

  const contentType = response.headers.get('content-type') ?? ''
  let data: T | undefined
  let rawText: string | undefined
  let error: string | undefined

  if (contentType.includes('application/json')) {
    try {
      data = (await response.json()) as T
    } catch (err) {
      error = `Failed to parse JSON: ${(err as Error).message}`
    }
  } else {
    rawText = await response.text()
  }

  if (!response.ok && !error) {
    error = rawText || (data ? JSON.stringify(data) : response.statusText)
  }

  return {
    ok: response.ok,
    status: response.status,
    statusText: response.statusText,
    data,
    rawText,
    error,
    headers: Object.fromEntries(response.headers.entries())
  }
}

export function safeJsonParse(value: string) {
  if (!value.trim()) {
    return undefined
  }
  return JSON.parse(value)
}

export function formatJson(value: unknown) {
  return JSON.stringify(value, null, 2)
}
