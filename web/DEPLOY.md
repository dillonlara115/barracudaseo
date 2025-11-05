# Deploy to Vercel

## Quick Deploy

### Option 1: Via Vercel Dashboard (Easiest)

1. **Go to Vercel**: https://vercel.com/new
2. **Import your Git repository** (GitHub/GitLab/Bitbucket)
3. **Configure project**:
   - Root Directory: `web`
   - Framework Preset: Vite
   - Build Command: `npm run build` (auto-detected)
   - Output Directory: `dist` (auto-detected)
4. **Add Environment Variables** (Settings → Environment Variables):
   - `PUBLIC_SUPABASE_URL` = Your Supabase URL
   - `PUBLIC_SUPABASE_ANON_KEY` = Your Supabase anon key
   - `VITE_CLOUD_RUN_API_URL` = `https://barracuda-api-7paxg34svq-uc.a.run.app`
5. **Deploy**

### Option 2: Via Vercel CLI

```bash
# Install Vercel CLI (if not installed)
npm install -g vercel

# Navigate to web directory
cd web

# Deploy
vercel

# Follow prompts:
# - Link to existing project or create new
# - Confirm settings
# - Deploy

# For production deployment:
vercel --prod
```

## Environment Variables

Set these in Vercel Dashboard → Project Settings → Environment Variables:

- `PUBLIC_SUPABASE_URL` - Your Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Your Supabase anon key  
- `VITE_CLOUD_RUN_API_URL` - Your Cloud Run API URL (optional)

**Note:** Vite exposes variables prefixed with `VITE_` or `PUBLIC_` to the client-side code.

## Build Configuration

The `vercel.json` file is already configured with:
- SPA routing (all routes → index.html)
- Asset caching headers
- Build commands

## After Deployment

1. Visit your Vercel deployment URL
2. The frontend will load (may show errors for API calls until Supabase is integrated)
3. Next step: Add Supabase client integration

## Troubleshooting

- **Build fails**: Check Node.js version (Vercel uses 18.x by default)
- **404 on routes**: Verify `vercel.json` rewrites are configured
- **Environment variables not working**: Make sure they're prefixed with `VITE_` or `PUBLIC_`

