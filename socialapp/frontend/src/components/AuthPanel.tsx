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
  const [response, setResponse] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const canRequestToken = clientId.length > 0 && clientSecret.length > 0

  const handleToken = async () => {
    setError('')
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
    <section>
      <div className="flex-between">
        <h3>Authentication</h3>
        <span className="status-chip">{token ? 'Token set' : 'No token'}</span>
      </div>

      <div className="form-grid mt-16">
        <label className="input-group">
          <span className="input-label">Client ID</span>
          <input className="input input-mono" value={clientId} onChange={(event) => onClientId(event.target.value)} />
        </label>

        <label className="input-group">
          <span className="input-label">Client Secret</span>
          <input
            className="input input-mono"
            type="password"
            value={clientSecret}
            onChange={(event) => onClientSecret(event.target.value)}
          />
        </label>
      </div>

      <label className="input-group mt-16">
        <span className="input-label">Scopes (space-separated)</span>
        <input className="input input-mono" value={scopes} onChange={(event) => onScopes(event.target.value)} />
      </label>

      <p className="small-muted">Token is cached locally in your browser.</p>

      <div className="row mt-16">
        <button className="btn btn-primary" onClick={handleToken} disabled={isLoading || !canRequestToken}>
          {isLoading ? 'Fetching...' : 'Get access token'}
        </button>
        <button className="btn btn-secondary" onClick={() => onToken(undefined)}>
          Clear token
        </button>
      </div>

      {!canRequestToken ? <p className="error-text">Provide client credentials to request a token.</p> : null}
      {error ? <p className="error-text">Error: {error}</p> : null}
      {response ? <pre className="highlight">{response}</pre> : null}
    </section>
  )
}
