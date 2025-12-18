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
    build: {
      rollupOptions: {
        output: {
          manualChunks: (id) => {
            // Split vendor chunks for better caching and loading performance
            if (id.includes('node_modules')) {
              // Chart.js and related charting libraries (large)
              if (id.includes('chart.js') || id.includes('svelte-chartjs')) {
                return 'vendor-charts';
              }
              // Supabase client (large)
              if (id.includes('@supabase')) {
                return 'vendor-supabase';
              }
              // Markdown parser
              if (id.includes('marked')) {
                return 'vendor-marked';
              }
              // Lucide icons (can be large if many icons are used)
              if (id.includes('lucide-svelte')) {
                return 'vendor-icons';
              }
              // Svelte framework and router
              if (id.includes('svelte') && !id.includes('svelte-chartjs')) {
                return 'vendor-svelte';
              }
              // All other node_modules
              return 'vendor';
            }
          },
          // Increase chunk size warning limit to 600KB (from default 500KB)
          // This gives us some flexibility while still warning about very large chunks
          chunkSizeWarningLimit: 600,
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
