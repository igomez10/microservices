var _a;
import { defineConfig } from '@playwright/test';
var baseURL = (_a = process.env.E2E_BASE_URL) !== null && _a !== void 0 ? _a : 'http://localhost:4173';
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
        command: 'npm run dev -- --host 127.0.0.1 --port 4173',
        url: baseURL,
        reuseExistingServer: !process.env.CI
    }
});
