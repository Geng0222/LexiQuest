// vite.config.js
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig({
  plugins: [vue()],
  base: '/LexiQuest/', 
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0',  
    port: 5173,        
    strictPort: true,  
    open: false,       
    proxy: {
      '/api': {
        target: 'http://localhost:28080',
        changeOrigin: true,
        secure: false,
      }
    }
  },
  build: {
    outDir: 'dist',  // **Vue 打包後的文件會放到 Go 伺服器的 `dist/` 內**
    emptyOutDir: true,   // **確保每次打包都清除舊的 `dist/`**
  }
});
