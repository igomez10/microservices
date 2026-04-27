/// <reference types="vite/client" />

type SocialappEnv = {
  readonly VITE_API_BASE_URL?: string
  readonly VITE_DEV_PROXY_TARGET?: string
  readonly VITE_CLIENT_ID?: string
  readonly VITE_CLIENT_SECRET?: string
  readonly VITE_OAUTH_SCOPES?: string
}

interface ImportMetaEnv extends SocialappEnv {}
interface ImportMeta {
  readonly env: ImportMetaEnv
}
