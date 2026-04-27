import { clearCookie, readCookie, readJwtPayload, writeCookie } from './storage'

function toBase64Url(value: string) {
  return btoa(value).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

describe('storage helpers', () => {
  beforeEach(() => {
    clearCookie('socialapp.token')
    clearCookie('socialapp.username')
  })

  it('writes and reads cookies', () => {
    writeCookie('socialapp.username', 'some2', 120)
    expect(readCookie('socialapp.username')).toBe('some2')
  })

  it('clears cookies', () => {
    writeCookie('socialapp.username', 'some2', 120)
    clearCookie('socialapp.username')
    expect(readCookie('socialapp.username')).toBeUndefined()
  })

  it('reads a jwt payload', () => {
    const header = toBase64Url(JSON.stringify({ alg: 'none', typ: 'JWT' }))
    const payload = toBase64Url(JSON.stringify({ preferred_username: 'some2', sub: '123' }))
    const token = `${header}.${payload}.signature`

    expect(readJwtPayload(token)).toEqual({ preferred_username: 'some2', sub: '123' })
  })

  it('returns undefined for an invalid jwt payload', () => {
    expect(readJwtPayload('not-a-jwt')).toBeUndefined()
  })
})
