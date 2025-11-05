# Vercel Deployment URL

**Production URL:** https://app.barracudaseo.com

## Supabase Configuration

Make sure your Supabase project has these settings:

### Site URL
```
https://app.barracudaseo.com
```

### Redirect URLs
```
https://app.barracudaseo.com
https://app.barracudaseo.com/**
http://localhost:5173
http://localhost:5173/**
```

## Environment Variables in Vercel

The following environment variables should be set in Vercel:

- `PUBLIC_SUPABASE_URL` - Your Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Your Supabase anon key

These are set in: Vercel Dashboard → Project Settings → Environment Variables

