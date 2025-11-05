import { svelte } from '@sveltejs/vite-plugin-svelte';
import { defineConfig, loadEnv } from 'vite';

export default defineConfig(({ mode }) => {
  // Load env file based on `mode` in the current working directory.
  // Set the third parameter to '' to load all env regardless of the `VITE_` prefix.
  const env = loadEnv(mode, process.cwd(), '');
  
  return {
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
    // This ensures they're available at build time (Vercel) and dev time (local)
    // Vite automatically exposes PUBLIC_ and VITE_ prefixed vars, but we define
    // them explicitly here to ensure they work in both environments
    define: {
      'import.meta.env.PUBLIC_SUPABASE_URL': JSON.stringify(env.PUBLIC_SUPABASE_URL || ''),
      'import.meta.env.PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(env.PUBLIC_SUPABASE_ANON_KEY || ''),
      'import.meta.env.VITE_PUBLIC_SUPABASE_URL': JSON.stringify(env.VITE_PUBLIC_SUPABASE_URL || env.PUBLIC_SUPABASE_URL || ''),
      'import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY': JSON.stringify(env.VITE_PUBLIC_SUPABASE_ANON_KEY || env.PUBLIC_SUPABASE_ANON_KEY || ''),
    }
  };
});

