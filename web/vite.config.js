import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  // Load env files for local dev, but always let actual process env win (Vercel/Cloud Run)
  const fileEnv = loadEnv(mode, process.cwd(), '');

  const env = {
    PUBLIC_SUPABASE_URL:
      process.env.PUBLIC_SUPABASE_URL ||
      process.env.VITE_PUBLIC_SUPABASE_URL ||
      fileEnv.PUBLIC_SUPABASE_URL ||
      fileEnv.VITE_PUBLIC_SUPABASE_URL ||
      '',
    PUBLIC_SUPABASE_ANON_KEY:
      process.env.PUBLIC_SUPABASE_ANON_KEY ||
      process.env.VITE_PUBLIC_SUPABASE_ANON_KEY ||
      fileEnv.PUBLIC_SUPABASE_ANON_KEY ||
      fileEnv.VITE_PUBLIC_SUPABASE_ANON_KEY ||
      '',
    VITE_CLOUD_RUN_API_URL:
      process.env.VITE_CLOUD_RUN_API_URL ||
      fileEnv.VITE_CLOUD_RUN_API_URL ||
      '',
  };
  
  return {
    plugins: [svelte()],
    server: {
      port: 5173,
      strictPort: false, // Allow Vite to use next available port if 5173 is taken
      proxy: {
        '/api': {
          target: 'http://localhost:8080',
          changeOrigin: true
        }
      }
    },
    // Explicitly define environment variables for client-side access
    define: {
      'import.meta.env.PUBLIC_SUPABASE_URL': JSON.stringify(env.PUBLIC_SUPABASE_URL),
      'import.meta.env.PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(env.PUBLIC_SUPABASE_ANON_KEY),
      'import.meta.env.VITE_PUBLIC_SUPABASE_URL': JSON.stringify(env.PUBLIC_SUPABASE_URL),
      'import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(env.PUBLIC_SUPABASE_ANON_KEY),
      'import.meta.env.VITE_CLOUD_RUN_API_URL': JSON.stringify(env.VITE_CLOUD_RUN_API_URL),
    }
  };
});
