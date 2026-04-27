import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { apiRequest, buildBasicAuthHeader, buildUrl, safeJsonParse } from './client'

const decodeBasic = (value: string) => {
  if (typeof atob === 'function') {
    return atob(value)
  }
  const nodeBuffer = (globalThis as { Buffer?: { from: (input: string, encoding: string) => { toString: (encoding: string) => string } } }).Buffer
  return nodeBuffer ? nodeBuffer.from(value, 'base64').toString('utf-8') : value
}

describe('api client helpers', () => {
  beforeEach(() => {
    vi.restoreAllMocks()
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  it('builds urls with path and query params', () => {
    const url = buildUrl('https://example.com', '/v1/users/{username}', { username: 'john' }, {
      limit: 10,
      offset: 0
    })
    expect(url).toBe('https://example.com/v1/users/john?limit=10&offset=0')
  })

  it('builds urls from a relative api base path in the browser', () => {
    const url = buildUrl('/api', '/v1/users/{username}', { username: 'john' }, { limit: 20 })
    expect(url).toBe(`${window.location.origin}/api/v1/users/john?limit=20`)
  })

  it('builds basic auth header', () => {
    const header = buildBasicAuthHeader('client', 'secret')
    const encoded = header.replace('Basic ', '')
    expect(decodeBasic(encoded)).toBe('client:secret')
  })

  it('parses json safely when empty', () => {
    expect(safeJsonParse('')).toBeUndefined()
  })

  it('parses json safely when provided', () => {
    expect(safeJsonParse('{"hello":"world"}')).toEqual({ hello: 'world' })
  })

  it('retries network failures with exponential backoff and eventually succeeds', async () => {
    vi.useFakeTimers()
    const fetchMock = vi
      .spyOn(global, 'fetch')
      .mockRejectedValueOnce(new Error('temporary network error'))
      .mockRejectedValueOnce(new Error('temporary network error'))
      .mockResolvedValueOnce(
        new Response(JSON.stringify({ ok: true }), {
          status: 200,
          headers: { 'content-type': 'application/json' }
        })
      )

    const requestPromise = apiRequest<{ ok: boolean }>({
      baseUrl: 'https://example.com',
      path: '/v1/users',
      method: 'GET'
    })

    await vi.advanceTimersByTimeAsync(250)
    await vi.advanceTimersByTimeAsync(500)

    const response = await requestPromise
    expect(fetchMock).toHaveBeenCalledTimes(3)
    expect(response.ok).toBe(true)
    expect(response.data).toEqual({ ok: true })
  })

  it('returns a network error response after exhausting retries', async () => {
    vi.useFakeTimers()
    const fetchMock = vi.spyOn(global, 'fetch').mockRejectedValue(new Error('offline'))

    const requestPromise = apiRequest({
      baseUrl: 'https://example.com',
      path: '/v1/users',
      method: 'GET'
    })

    await vi.advanceTimersByTimeAsync(250)
    await vi.advanceTimersByTimeAsync(500)

    const response = await requestPromise
    expect(fetchMock).toHaveBeenCalledTimes(3)
    expect(response.ok).toBe(false)
    expect(response.status).toBe(0)
    expect(response.statusText).toBe('Network Error')
    expect(response.error).toBe('offline')
  })
})
