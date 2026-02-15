import { render, screen } from '@testing-library/react'
import AuthPanel from './AuthPanel'

describe('AuthPanel', () => {
  it('renders auth inputs', () => {
    render(
      <AuthPanel
        baseUrl="https://example.com"
        token={undefined}
        onToken={() => {}}
        clientId=""
        clientSecret=""
        scopes=""
        onClientId={() => {}}
        onClientSecret={() => {}}
        onScopes={() => {}}
      />
    )

    expect(screen.getByLabelText('Client ID')).toBeInTheDocument()
    expect(screen.getByLabelText('Client Secret')).toBeInTheDocument()
    expect(screen.getByLabelText('Scopes (space-separated)')).toBeInTheDocument()
  })
})
