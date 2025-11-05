import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  // Explicitly define environment variables for client-side access
  // This ensures they're available at build time in Vercel
  define: {
    'import.meta.env.PUBLIC_SUPABASE_URL': JSON.stringify(process.env.PUBLIC_SUPABASE_URL || ''),
    'import.meta.env.PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(process.env.PUBLIC_SUPABASE_ANON_KEY || ''),
    'import.meta.env.VITE_PUBLIC_SUPABASE_URL': JSON.stringify(process.env.VITE_PUBLIC_SUPABASE_URL || process.env.PUBLIC_SUPABASE_URL || ''),
    'import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(process.env.VITE_PUBLIC_SUPABASE_ANON_KEY || process.env.PUBLIC_SUPABASE_ANON_KEY || ''),
  }
});

