# Local Magic Link Setup Guide

This guide will help you set up magic link authentication for local development with Supabase.

## Prerequisites

1. **Docker** must be running (Supabase uses Docker)
2. **Supabase CLI** installed (`supabase --version` to check)
3. **Node.js** installed for the web app

## Step 1: Start Supabase Locally

```bash
cd /home/dillon/Sites/cli-scanner
supabase start
```

**Expected Output:**
```
Started supabase local development setup.

         API URL: http://127.0.0.1:54321
     GraphQL URL: http://127.0.0.1:54321/graphql/v1
          DB URL: postgresql://postgres:postgres@127.0.0.1:54322/postgres
      Studio URL: http://127.0.0.1:54323
    Inbucket URL: http://127.0.0.1:54324
      JWT secret: super-secret-jwt-token-with-at-least-32-characters-long
        anon key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
service_role key: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Important:** Copy the `API URL` and `anon key` - you'll need these for environment variables.

## Step 2: Configure Environment Variables

Create a `.env.local` file in the `web/` directory (or update your existing `.env` file):

```bash
cd web
cat > .env.local << EOF
PUBLIC_SUPABASE_URL=http://127.0.0.1:54321
PUBLIC_SUPABASE_ANON_KEY=<paste-anon-key-from-supabase-start>
EOF
```

**Replace `<paste-anon-key-from-supabase-start>`** with the actual anon key from Step 1.

### Alternative: Check Existing Environment Variables

If you already have a `.env` file, make sure it points to local Supabase:

```bash
# Check current values
cd web
cat .env.local 2>/dev/null || cat .env 2>/dev/null || echo "No .env file found"
```

**Required values for local development:**
- `PUBLIC_SUPABASE_URL=http://127.0.0.1:54321`
- `PUBLIC_SUPABASE_ANON_KEY=<your-local-anon-key>`

## Step 3: Verify Supabase is Running

### Check Supabase Status

```bash
cd /home/dillon/Sites/cli-scanner
supabase status
```

### Check Inbucket (Email Testing)

Open in your browser: **http://localhost:54324**

You should see the Inbucket interface. This is where all magic link emails will appear.

### Check Supabase Studio

Open in your browser: **http://localhost:54323**

This is the Supabase admin interface where you can:
- View database tables
- Check authentication logs
- View email templates

## Step 4: Start Your Web App

```bash
cd web
npm run dev
```

The app should start on **http://localhost:5173** (or the port shown in terminal).

## Step 5: Test Magic Link Flow

### Test Sign Up

1. **Open app**: http://localhost:5173
2. **Enter email**: `test@example.com`
3. **Enter name** (if signup form)
4. **Click**: "Send magic link"
5. **Check console**: Should see `✅ Magic link sent successfully`
6. **Open Inbucket**: http://localhost:54324
7. **Find email**: Look for email to `test@example.com`
8. **Click magic link** in the email
9. **Verify**: Should redirect and log you in

### Test Sign In

1. **Sign out** (if logged in)
2. **Enter email**: `test@example.com`
3. **Click**: "Send magic link"
4. **Check Inbucket**: Should see new email
5. **Click magic link**: Should log you in

## Troubleshooting

### Emails Not Appearing in Inbucket

**Problem**: Magic link sent but no email in Inbucket

**Solutions:**

1. **Verify Supabase is running:**
   ```bash
   supabase status
   ```
   If not running, start it:
   ```bash
   supabase start
   ```

2. **Check environment variables:**
   ```bash
   cd web
   # Make sure these point to local Supabase
   echo $PUBLIC_SUPABASE_URL  # Should be http://127.0.0.1:54321
   ```

3. **Restart web app** after changing environment variables:
   ```bash
   # Stop the dev server (Ctrl+C)
   # Then restart
   npm run dev
   ```

4. **Check browser console** for errors:
   - Open DevTools (F12)
   - Look for Supabase connection errors
   - Check Network tab for failed requests

5. **Verify Inbucket is accessible:**
   - Open: http://localhost:54324
   - If it doesn't load, Supabase might not be fully started
   - Wait a minute and try `supabase status` again

### "Missing Supabase configuration" Error

**Problem**: Console shows "Missing Supabase configuration"

**Solution:**

1. **Create `.env.local` file** in `web/` directory:
   ```bash
   cd web
   cat > .env.local << EOF
   PUBLIC_SUPABASE_URL=http://127.0.0.1:54321
   PUBLIC_SUPABASE_ANON_KEY=<your-anon-key>
   EOF
   ```

2. **Get anon key** from Supabase:
   ```bash
   cd /home/dillon/Sites/cli-scanner
   supabase status | grep "anon key"
   ```

3. **Restart dev server** after creating `.env.local`

### Magic Link Doesn't Work After Clicking

**Problem**: Click magic link but doesn't log in

**Solutions:**

1. **Check redirect URLs** in `supabase/config.toml`:
   ```toml
   [auth]
   site_url = "http://127.0.0.1:3000"
   additional_redirect_urls = [
     "http://127.0.0.1:3000/#/",
     "http://localhost:5173/#/",
     "http://localhost:5173/auth/confirm"
   ]
   ```

2. **Restart Supabase** after changing config:
   ```bash
   supabase stop
   supabase start
   ```

3. **Check browser console** for errors when clicking link

4. **Verify the redirect URL** in the magic link email matches your app URL

### Docker Permission Errors

**Problem**: `permission denied while trying to connect to the Docker daemon socket`

**Solution:**

```bash
# Add your user to docker group (requires logout/login)
sudo usermod -aG docker $USER

# Or run with sudo (not recommended for regular use)
sudo supabase start
```

## Quick Verification Checklist

- [ ] Supabase is running (`supabase status` shows all services)
- [ ] Inbucket accessible at http://localhost:54324
- [ ] `.env.local` file exists in `web/` directory
- [ ] Environment variables point to `http://127.0.0.1:54321`
- [ ] Web app running (`npm run dev`)
- [ ] Can request magic link (no console errors)
- [ ] Email appears in Inbucket
- [ ] Clicking magic link logs you in

## Getting Help

If you're still having issues:

1. **Check Supabase logs:**
   ```bash
   supabase logs
   ```

2. **Check browser console** for detailed error messages

3. **Verify config** matches examples in:
   - `supabase/config.toml`
   - `docs/MAGIC_LINK_AUTH.md`
   - `docs/PKCE_MAGIC_LINK_SETUP.md`

4. **Restart everything:**
   ```bash
   # Stop Supabase
   supabase stop
   
   # Start Supabase
   supabase start
   
   # Restart web app
   cd web
   npm run dev
   ```

## Next Steps

Once local magic links are working:

1. **Test all auth flows** (signup, signin, signout)
2. **Test with different emails**
3. **Check email templates** in Supabase Studio
4. **Review production setup** in `docs/PRODUCTION_MAGIC_LINK_SETUP.md`
