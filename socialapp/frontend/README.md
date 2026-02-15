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
VITE_API_BASE_URL=http://localhost:8086 npm run dev
```

If you omit `VITE_API_BASE_URL`, the app defaults to `http://localhost:5173/api` and uses the Vite dev proxy
to forward `/api/*` requests to `http://localhost:8086`.

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
- `VITE_API_BASE_URL` (default: `http://localhost:8086`)
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
