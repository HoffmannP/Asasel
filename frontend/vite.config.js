import { sveltekit } from '@sveltejs/kit/vite'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [sveltekit()],
  build: {
    rollupOptions: {
      output: {
        manualChunks (id) {
          if (
            id.includes('svelte/src/internal/server/context.js') ||
            id.includes('svelte/src/index-server.js') ||
            id.includes('@sveltejs/kit/src/runtime/app/state/server.js')
          ) {
            return 'svelte-server'
          }

          if (id.includes('node_modules')) {
            return 'vendor'
          }
        }
      }
    }
  }
})
