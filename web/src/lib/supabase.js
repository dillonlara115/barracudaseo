import { createClient } from '@supabase/supabase-js';

// Get Supabase URL and key from environment variables
// Vite exposes variables prefixed with VITE_ or PUBLIC_
const supabaseUrl = import.meta.env.PUBLIC_SUPABASE_URL || import.meta.env.VITE_PUBLIC_SUPABASE_URL || '';
const supabaseAnonKey = import.meta.env.PUBLIC_SUPABASE_ANON_KEY || import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY || '';

if (!supabaseUrl || !supabaseAnonKey) {
  console.error('Missing Supabase configuration. Please set PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY');
  console.error('Current env values:', {
    PUBLIC_SUPABASE_URL: import.meta.env.PUBLIC_SUPABASE_URL,
    VITE_PUBLIC_SUPABASE_URL: import.meta.env.VITE_PUBLIC_SUPABASE_URL,
    PUBLIC_SUPABASE_ANON_KEY: import.meta.env.PUBLIC_SUPABASE_ANON_KEY ? '***set***' : undefined,
    VITE_PUBLIC_SUPABASE_ANON_KEY: import.meta.env.VITE_PUBLIC_SUPABASE_ANON_KEY ? '***set***' : undefined,
  });
}

// Create Supabase client (will fail gracefully if URL/key are missing)
// Note: Using 'implicit' flow for email magic links because PKCE requires code_verifier
// to be in the same browser session, which doesn't work when opening email links
// in different browsers or after session storage is cleared.
// For production OAuth flows, PKCE is still recommended, but email magic links
// work better with implicit flow.
export const supabase = createClient(
  supabaseUrl || 'https://placeholder.supabase.co',
  supabaseAnonKey || 'placeholder-key',
  {
    auth: {
      persistSession: true,
      autoRefreshToken: true,
      detectSessionInUrl: true,
      // Don't set a default redirectTo - let it use the current URL
      // This allows magic links to work with hash routing
      // Using implicit flow for email magic links (works better than PKCE for email)
      flowType: 'implicit'
    }
  }
);

