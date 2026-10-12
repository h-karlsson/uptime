import { defineConfig } from '@hey-api/openapi-ts'

export default defineConfig({
  input: '../api-docs/openapi.yaml',
  output: 'src/client',
  plugins: [
    { name: '@hey-api/client-fetch', baseUrl: false, throwOnError: true },
    '@hey-api/sdk',
    '@hey-api/typescript',
  ],
})
