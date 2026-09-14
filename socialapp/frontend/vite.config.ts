import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const apiTarget = env.VITE_DEV_PROXY_TARGET || 'http://localhost:8085'

  return {
    plugins: [react()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src')
      }
    },
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: apiTarget,
          // Local Traefik HTTPS uses a self-signed cert in dev.
          // Allow proxying to that endpoint directly to avoid auth-breaking redirects.
          secure: false,
          changeOrigin: true,
          rewrite: (pathValue) => pathValue.replace(/^\/api/, '')
        }
      }
    }
  }
})
