import { useState } from 'react'
import { apiRequest, formatJson } from '@/api/client'

export type AuthPanelProps = {
  baseUrl: string
  token?: string
  onToken: (token: string | undefined) => void
  clientId: string
  clientSecret: string
  scopes: string
  onClientId: (value: string) => void
  onClientSecret: (value: string) => void
  onScopes: (value: string) => void
}

export default function AuthPanel({
  baseUrl,
  token,
  onToken,
  clientId,
  clientSecret,
  scopes,
  onClientId,
  onClientSecret,
  onScopes
}: AuthPanelProps) {
  const [status, setStatus] = useState('')
  const [response, setResponse] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const canRequestToken = clientId.length > 0 && clientSecret.length > 0

  const handleToken = async () => {
    setError('')
    setStatus('')
    setResponse('')
    setIsLoading(true)
    try {
      const res = await apiRequest<{ access_token?: string }>({
        baseUrl,
        path: '/v1/oauth/token',
        method: 'POST',
        basicAuth: {
          clientId,
          clientSecret
        },
        queryParams: scopes ? { scope: scopes } : undefined
      })

      setStatus(`${res.status} ${res.statusText}`)
      if (res.data) {
        setResponse(formatJson(res.data))
      }
      if (res.data?.access_token) {
        onToken(res.data.access_token)
      }
      if (res.error) {
        setError(res.error)
      }
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div>
      <h2>Authentication</h2>
      <div className="grid">
        <label className="field">
          Client ID
          <input value={clientId} onChange={(event) => onClientId(event.target.value)} />
        </label>
        <label className="field">
          Client Secret
          <input
            type="password"
            value={clientSecret}
            onChange={(event) => onClientSecret(event.target.value)}
          />
        </label>
        <label className="field">
          Scopes (space-separated)
          <input value={scopes} onChange={(event) => onScopes(event.target.value)} />
        </label>
      </div>
      <div className="muted">Token is cached locally in your browser.</div>
      <div className="row">
        <button className="button" onClick={handleToken} disabled={isLoading || !canRequestToken}>
          {isLoading ? 'Fetching…' : 'Get access token'}
        </button>
        <button className="button secondary" onClick={() => onToken(undefined)}>
          Clear token
        </button>
        <span className="pill">{token ? 'Token set' : 'No token'}</span>
      </div>
      {status ? <div className="muted">Status: {status}</div> : null}
      {!canRequestToken ? <div className="error">Provide client credentials to request a token.</div> : null}
      {error ? <div className="error">Error: {error}</div> : null}
      {response ? <pre className="highlight">{response}</pre> : null}
    </div>
  )
}
