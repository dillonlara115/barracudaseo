# Vercel Deployment Guide

This guide walks you through deploying the Barracuda frontend to Vercel.

## Prerequisites

1. **Vercel account** (sign up at https://vercel.com)
2. **Vercel CLI** installed (optional, but helpful)
3. **Cloud Run API URL** (from your deployment)
4. **Supabase credentials** (URL and anon key)

## Step 1: Install Vercel CLI (Optional)

```bash
npm install -g vercel
```

Or use npx without installing:
```bash
npx vercel
```

## Step 2: Set Up Frontend Supabase Integration

Before deploying, we need to:
1. Install Supabase client library
2. Update the frontend to use Supabase instead of local API
3. Add authentication UI

## Step 3: Configure Environment Variables

Create a `.env.local` file in the `web/` directory (or use Vercel dashboard):

```bash
cd web
cat > .env.local << EOF
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key
VITE_CLOUD_RUN_API_URL=https://barracuda-api-7paxg34svq-uc.a.run.app
EOF
```

**Note:** Vite requires the `VITE_` prefix for client-side variables, OR you can use `PUBLIC_` prefix which Vite also supports.

## Step 4: Deploy to Vercel

### Option A: Using Vercel CLI

```bash
cd web
vercel
```

Follow the prompts:
- Link to existing project or create new
- Confirm project settings
- Deploy

### Option B: Using Vercel Dashboard

1. Go to https://vercel.com/new
2. Import your Git repository
3. Set root directory to `web/`
4. Configure build settings:
   - Framework Preset: Vite
   - Build Command: `npm run build`
   - Output Directory: `dist`
5. Add environment variables:
   - `PUBLIC_SUPABASE_URL`
   - `PUBLIC_SUPABASE_ANON_KEY`
   - `VITE_CLOUD_RUN_API_URL` (optional, if you need to call Cloud Run API)

## Step 5: Configure Vercel Settings

Create `vercel.json` in the `web/` directory:

```json
{
  "buildCommand": "npm run build",
  "outputDirectory": "dist",
  "devCommand": "npm run dev",
  "installCommand": "npm install",
  "framework": "vite",
  "rewrites": [
    {
      "source": "/(.*)",
      "destination": "/index.html"
    }
  ]
}
```

## Environment Variables in Vercel

After deploying, add these in Vercel Dashboard → Project Settings → Environment Variables:

- `PUBLIC_SUPABASE_URL` = Your Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` = Your Supabase anon key
- `VITE_CLOUD_RUN_API_URL` = `https://barracuda-api-7paxg34svq-uc.a.run.app` (optional)

## Testing the Deployment

1. Visit your Vercel deployment URL
2. Test authentication (login/signup)
3. Test data fetching from Supabase
4. Verify API calls to Cloud Run (if needed)

## Next Steps

After Vercel deployment:
1. Update CORS settings in Cloud Run to allow your Vercel domain
2. Test end-to-end flow: Sign up → Create project → Upload crawl → View results
3. Set up custom domain (optional)

## Troubleshooting

### CORS Errors
If you see CORS errors, update Cloud Run CORS settings or configure them in the API server.

### Environment Variables Not Working
- Make sure variables start with `VITE_` or `PUBLIC_` for Vite
- Redeploy after adding new variables
- Check Vercel build logs for variable issues

### Build Failures
- Check `web/package.json` has all dependencies
- Verify Node.js version (Vercel uses Node 18+ by default)
- Check build logs in Vercel dashboard

