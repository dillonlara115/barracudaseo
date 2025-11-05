# Environment Variables Setup for Vercel

## Required Variables

Set these in Vercel Dashboard → Settings → Environment Variables:

1. **PUBLIC_SUPABASE_URL** - Your Supabase project URL
   - Example: `https://xxxxx.supabase.co`
   - Make sure it's assigned to **All Environments** (Production, Preview, Development)

2. **PUBLIC_SUPABASE_ANON_KEY** - Your Supabase anon key
   - Example: `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - Make sure it's assigned to **All Environments**

## Important Notes

- **Assign to All Environments**: Make sure both variables are assigned to Production, Preview, AND Development environments
- **After Adding Variables**: You MUST redeploy for changes to take effect
- **Case Sensitive**: Variable names are case-sensitive - use exactly `PUBLIC_SUPABASE_URL` and `PUBLIC_SUPABASE_ANON_KEY`

## Troubleshooting

If variables still don't work after redeploying:

1. **Check Build Logs**: Look at the Vercel build logs to see if variables are being read
2. **Clear Cache**: Try redeploying with "Clear Cache" option
3. **Verify Names**: Double-check the variable names match exactly (no typos, correct case)
4. **Check All Environments**: Ensure variables are enabled for the environment you're deploying to

## Verification

After deploying, check the browser console. You should see:
- No "Missing Supabase configuration" errors
- Variables should be accessible via `import.meta.env.PUBLIC_SUPABASE_URL`

