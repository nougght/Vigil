import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import {VitePWA} from "vite-plugin-pwa";

// https://vite.dev/config/
export default defineConfig({
    server: {
        proxy: {
            '/api': { target: 'http://localhost:8091', changeOrigin: true, ws: true },
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
                    {
                        src: 'logo.png',
                        sizes: '256x256',
                        type: 'image/png'
                    }
                ]
            }
        })
    ],
})
