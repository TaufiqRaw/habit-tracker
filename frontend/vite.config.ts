import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'
import tailwind from "tailwindcss";
import autoprefixer from "autoprefixer";
import svgr from "vite-plugin-svgr";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react(), svgr()],
  css: {
    postcss: {
      plugins: [tailwind, autoprefixer],
    }
  }
})
