import { buildBasicAuthHeader, buildUrl, safeJsonParse } from './client'

const decodeBasic = (value: string) => {
  if (typeof atob === 'function') {
    return atob(value)
  }
  const nodeBuffer = (globalThis as { Buffer?: { from: (input: string, encoding: string) => { toString: (encoding: string) => string } } }).Buffer
  return nodeBuffer ? nodeBuffer.from(value, 'base64').toString('utf-8') : value
}

describe('api client helpers', () => {
  it('builds urls with path and query params', () => {
    const url = buildUrl('https://example.com', '/v1/users/{username}', { username: 'john' }, {
      limit: 10,
      offset: 0
    })
    expect(url).toBe('https://example.com/v1/users/john?limit=10&offset=0')
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
})
