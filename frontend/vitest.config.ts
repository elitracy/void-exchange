import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// Separate from vite.config.ts because the Tailwind Vite plugin has no
// role in a jsdom test run and only slows test startup down.
export default defineConfig({
    plugins: [react()],
    test: {
        environment: 'jsdom',
        setupFiles: ['./src/test/setup.ts'],
        globals: true,
    },
})
