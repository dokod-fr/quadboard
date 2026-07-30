import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

export default defineConfig({
  plugins: [svelte()],
  build: {
    outDir: '../internal/http/view', 
    emptyOutDir: false, 
    rollupOptions: {
      output: {
        entryFileNames: 'assets/js/app.min.js', 
        assetFileNames: 'assets/css/style.min.css'
      }
    }
  },
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
    }
  }
});