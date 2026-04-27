# Socialapp Frontend

A social-network style UI for Socialapp (not a Swagger console). It exposes all OpenAPI features through feeds, profiles, admin tooling, and URL utilities.

## Requirements
- Node.js 18+
- Backend running at `http://localhost:8086` (default)

## Setup
```sh
cd frontend
npm install
npm run generate:api
```

You can also generate types from the repo root:
```sh
make generate-frontend-types
```

## Dev Server
```sh
VITE_DEV_PROXY_TARGET=http://localhost:8085 npm run dev
```

The frontend defaults to `/api` in the browser and uses the Vite dev proxy.
Set `VITE_DEV_PROXY_TARGET` to choose which backend host `/api/*` forwards to.

## Tests
### Unit
```sh
npm run test:unit
```

### Integration (live backend)
```sh
INTEGRATION_API_BASE_URL=http://localhost:8086 \
INTEGRATION_CLIENT_ID=... \
INTEGRATION_CLIENT_SECRET=... \
npm run test:integration
```

### UI (Playwright)
```sh
cd frontend
npx playwright install
E2E_API_BASE_URL=http://localhost:8086 \
E2E_CLIENT_ID=... \
E2E_CLIENT_SECRET=... \
npm run test:e2e
```

For local UI tests with the proxy, you can omit `E2E_API_BASE_URL` (defaults to `http://localhost:4173/api`).

## Environment Variables
- `VITE_API_BASE_URL` (optional browser-visible override, default: `/api`)
- `VITE_DEV_PROXY_TARGET` (dev-server proxy target, for example `http://localhost:8085`)
- `VITE_CLIENT_ID`
- `VITE_CLIENT_SECRET`
- `VITE_OAUTH_SCOPES`

### Integration Tests
- `INTEGRATION_API_BASE_URL`
- `INTEGRATION_CLIENT_ID`
- `INTEGRATION_CLIENT_SECRET`
- `INTEGRATION_SCOPES`

### UI Tests
- `E2E_BASE_URL` (default: `http://localhost:4173`)
- `E2E_API_BASE_URL`
- `E2E_CLIENT_ID`
- `E2E_CLIENT_SECRET`
