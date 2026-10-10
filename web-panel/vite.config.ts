import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import {VitePWA} from "vite-plugin-pwa";

// https://vite.dev/config/
export default defineConfig({
    server: {
        proxy: {
            '/api': {target: 'http://localhost:8091', changeOrigin: true, ws: true},
        },
    },
    plugins: [react(),
        VitePWA({
            registerType: 'autoUpdate', // Автоматически обновляет service worker
            manifest: {
                name: 'Vigil PWA',
                short_name: 'Vigil',
                description: 'Приложение Vigil',
                theme_color: '#ffffff',
                icons: [
                    {src: 'pwa-64x64.png', sizes: '64x64', type: 'image/png'},
                    {src: 'pwa-192x192.png', sizes: '192x192', type: 'image/png'},
                    {src: 'pwa-512x512.png', sizes: '512x512', type: 'image/png'},
                    {src: 'maskable-icon-512x512.png', sizes: '512x512', type: 'image/png', purpose: 'maskable'},
                ],
            }
        })
    ],
})
