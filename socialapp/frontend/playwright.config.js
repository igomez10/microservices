var _a, _b;
import { defineConfig } from '@playwright/test';
var baseURL = (_a = process.env.E2E_BASE_URL) !== null && _a !== void 0 ? _a : 'http://localhost:4173';
var apiProxyTarget = (_b = process.env.E2E_API_BASE_URL) !== null && _b !== void 0 ? _b : 'https://socialapp.gomezignacio.com';
export default defineConfig({
    testDir: './tests/e2e',
    timeout: 60000,
    expect: {
        timeout: 10000
    },
    use: {
        baseURL: baseURL,
        trace: 'on-first-retry',
        ignoreHTTPSErrors: true
    },
    webServer: {
        command: "VITE_DEV_PROXY_TARGET=".concat(apiProxyTarget, " npm run dev -- --host 127.0.0.1 --port 4173"),
        url: baseURL,
        reuseExistingServer: !process.env.CI
    }
});
