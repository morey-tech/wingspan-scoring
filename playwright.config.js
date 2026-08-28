// @ts-check
const { defineConfig, devices } = require('@playwright/test');
const os = require('os');
const path = require('path');

// Every run gets a throwaway SQLite file. The game-end endpoint writes a row on each
// calculation, so tests must never be pointed at a real database.
const testDbPath = path.join(os.tmpdir(), `wingspan-e2e-${process.pid}.db`);

module.exports = defineConfig({
    testDir: './tests/e2e',
    fullyParallel: false, // The app keeps a single game in localStorage.
    workers: 1,
    forbidOnly: !!process.env.CI,
    retries: process.env.CI ? 1 : 0,
    reporter: process.env.CI ? [['github'], ['html', { open: 'never' }]] : [['list']],

    use: {
        baseURL: 'http://127.0.0.1:8080',
        trace: 'on-first-retry',
    },

    projects: [
        {
            name: 'chromium',
            use: { ...devices['Desktop Chrome'] },
        },
    ],

    webServer: {
        command: 'go run .',
        url: 'http://127.0.0.1:8080/api/version',
        reuseExistingServer: false,
        timeout: 120 * 1000,
        env: { DB_PATH: testDbPath },
    },
});
