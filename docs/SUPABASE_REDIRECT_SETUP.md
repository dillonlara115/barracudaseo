# Supabase Redirect URL Configuration

## Problem
Email confirmation links are pointing to `http://localhost:3000` instead of your Vercel deployment URL.

## Solution

You need to configure Supabase to use your production URL. Here's how:

### Step 1: Update Site URL in Supabase Dashboard

1. Go to your Supabase project dashboard: https://supabase.com/dashboard
2. Navigate to **Settings** → **Authentication**
3. Scroll down to **URL Configuration**
4. Update the **Site URL** field:
   - Change from: `http://localhost:3000`
   - Change to: `https://app.barracudaseo.com`

### Step 2: Add Redirect URLs

In the same **URL Configuration** section, add your redirect URLs:

1. **Redirect URLs** - Add these URLs (one per line):
   ```
   https://app.barracudaseo.com
   https://app.barracudaseo.com/**
   http://localhost:5173
   http://localhost:5173/**
   ```

   Note: Include both production (Vercel) and local development URLs.

### Step 3: Verify Email Templates (Optional)

1. Go to **Authentication** → **Email Templates**
2. Check the **Confirm signup** template
3. The redirect URL should automatically use your Site URL setting

### Step 4: Test

1. Try signing up again from your Vercel deployment
2. Check the confirmation email
3. The link should now point to your Vercel URL instead of localhost

## Important Notes

- **Site URL**: This is the default redirect URL for all auth flows
- **Redirect URLs**: These are allowed redirect URLs (wildcards are supported)
- **Wildcards**: Use `**` to allow all subpaths (e.g., `https://app.vercel.app/**`)

## Multiple Environments

If you have multiple environments (production, preview, staging), add all of them:

```
https://app.barracudaseo.com
https://app.barracudaseo.com/**
https://barracuda-web-*.vercel.app
https://barracuda-web-*.vercel.app/**
```

The `*` wildcard matches any branch name for preview deployments.

