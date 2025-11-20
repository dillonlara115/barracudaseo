# Deployment Environment Variables Guide

This guide explains which environment variables need to be set in **Cloud Run** (backend) vs **Vercel** (frontend).

## Quick Answer

**Cloud Run (Backend API):** ✅ YES - Add email-related variables  
**Vercel (Frontend):** ❌ NO - Frontend doesn't send emails

## Cloud Run Environment Variables

The backend API handles email sending, so these variables must be set in **Google Cloud Run**:

### Required for Email Functionality

```bash
# Email provider (optional - defaults to 'supabase')
EMAIL_PROVIDER=resend

# Resend API key (if using Resend API directly)
RESEND_API_KEY=re_your-api-key-here

# Sender email address (must match verified domain in Resend)
EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com

# Base URL for invite links
APP_URL=https://app.barracudaseo.com
```

### How to Add to Cloud Run

**Option 1: Using the update script (Recommended)**

Add to your `.env` file:
```bash
EMAIL_PROVIDER=resend
RESEND_API_KEY=re_your-api-key-here
EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com
APP_URL=https://app.barracudaseo.com
```

Then run:
```bash
./scripts/update-cloud-run-env.sh --production
```

**Option 2: Using gcloud CLI**

```bash
gcloud run services update barracuda-api \
  --region us-central1 \
  --update-env-vars="EMAIL_PROVIDER=resend,RESEND_API_KEY=re_...,EMAIL_FROM_ADDRESS=noreply@mail.barracudaseo.com,APP_URL=https://app.barracudaseo.com"
```

**Option 3: Using Google Cloud Console**

1. Go to Cloud Run → Select your service
2. Click "Edit & Deploy New Revision"
3. Go to "Variables & Secrets" tab
4. Add environment variables:
   - `EMAIL_PROVIDER`
   - `RESEND_API_KEY` (consider using Secret Manager for this)
   - `EMAIL_FROM_ADDRESS`
   - `APP_URL`

### Using Secret Manager for Sensitive Values

For `RESEND_API_KEY`, consider using Secret Manager:

```bash
# Create secret
echo -n "re_your-api-key-here" | gcloud secrets create resend-api-key \
  --data-file=- \
  --replication-policy="automatic"

# Update Cloud Run to use secret
gcloud run services update barracuda-api \
  --region us-central1 \
  --update-secrets="RESEND_API_KEY=resend-api-key:latest"
```

## Vercel Environment Variables

**NO email-related variables needed** - The frontend doesn't send emails directly.

Vercel only needs:
- `PUBLIC_SUPABASE_URL`
- `PUBLIC_SUPABASE_ANON_KEY`
- `VITE_CLOUD_RUN_API_URL` (optional, if frontend calls backend API)

## Summary Table

| Variable | Cloud Run | Vercel | Notes |
|----------|-----------|--------|-------|
| `EMAIL_PROVIDER` | ✅ Yes | ❌ No | Backend only |
| `RESEND_API_KEY` | ✅ Yes | ❌ No | Backend only (use Secret Manager) |
| `EMAIL_FROM_ADDRESS` | ✅ Yes | ❌ No | Backend only |
| `APP_URL` | ✅ Yes | ❌ No | Backend only (for invite links) |
| `PUBLIC_SUPABASE_URL` | ✅ Yes | ✅ Yes | Both need it |
| `PUBLIC_SUPABASE_ANON_KEY` | ✅ Yes | ✅ Yes | Both need it |

## Alternative: Using Supabase SMTP

If you configure Resend SMTP in Supabase Dashboard instead of using Resend API directly:

- ✅ **No environment variables needed** in Cloud Run
- ✅ Configure SMTP in Supabase Dashboard:
  - Host: `smtp.resend.com`
  - Port: `465`
  - User: `resend`
  - Password: Your Resend API key
  - Sender: `noreply@mail.barracudaseo.com`

See [EMAIL_CONFIGURATION.md](./EMAIL_CONFIGURATION.md) for details.

