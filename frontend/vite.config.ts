/// <reference types="vitest/config" />
import { createReadStream } from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import tailwindcss from '@tailwindcss/vite'
import { tanstackRouter } from '@tanstack/router-plugin/vite'
import react from '@vitejs/plugin-react'
import { defineConfig, type Plugin } from 'vite'

// Serves MSW's service worker straight from the package, so it never ends up in a build, and answers
// any /api request the worker didn't intercept with an error instead of the app's index.html.
function mockMode(): Plugin {
  const workerPath = fileURLToPath(import.meta.resolve('msw/mockServiceWorker.js'))
  return {
    name: 'mock-mode',
    apply: 'serve',
    configureServer(server) {
      server.middlewares.use('/mockServiceWorker.js', (_request, response) => {
        response.setHeader('Content-Type', 'text/javascript')
        createReadStream(workerPath).pipe(response)
      })
      server.middlewares.use('/api', (_request, response) => {
        response.statusCode = 503
        response.setHeader('Content-Type', 'application/json')
        response.end(JSON.stringify({ error: 'The mock API is not active. Reload the page.' }))
      })
    },
  }
}

export default defineConfig(({ mode }) => ({
  // The router plugin must come before react(); it generates src/routeTree.gen.ts from src/routes.
  plugins: [
    tanstackRouter({ target: 'react', autoCodeSplitting: true }),
    react(),
    tailwindcss(),
    mode === 'mock' && mockMode(),
  ],
  resolve: {
    alias: {
      '@': path.resolve(import.meta.dirname, './src'),
    },
  },
  server: {
    port: 8081,
    strictPort: true,
    // In mock mode the service worker answers /api, so nothing may reach a real backend.
    proxy:
      mode === 'mock'
        ? undefined
        : {
            '/api': {
              target: 'http://127.0.0.1:8080',
              changeOrigin: true,
            },
          },
  },
  test: {
    environment: 'jsdom',
    env: { TZ: 'UTC' },
    unstubGlobals: true,
    setupFiles: ['./src/test/setup.ts'],
  },
}))
