import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    open: true,
    host: true,
    proxy: {
      '/api': {
        target: '',
        changeOrigin: true,
      },
    },
  },
	resolve: {
		alias: {
			"~": path.resolve(__dirname, "./src/"),
		}
	}
})
