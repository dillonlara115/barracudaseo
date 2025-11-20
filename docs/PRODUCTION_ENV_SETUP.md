# Production Environment Variables Setup

## Overview

You have `.env.local` for local development. For production deployments to Cloud Run, you have two options:

1. **Use `.env` file** - Production values (recommended)
2. **Use `.env.local`** - Can override `.env` values

## File Priority

The scripts load environment variables in this order:
1. `.env` - Production/default values
2. `.env.local` - Local overrides (takes precedence)

This matches how the Go code loads environment variables.

## Option 1: Create Production `.env` File (Recommended)

Create a `.env` file with your production values:

```bash
# Production Supabase Configuration
PUBLIC_SUPABASE_URL=https://your-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-production-anon-key

# Production Stripe Configuration
STRIPE_SECRET_KEY=sk_live_...
STRIPE_WEBHOOK_SECRET=whsec_...
STRIPE_PRICE_ID_PRO=price_...
STRIPE_PRICE_ID_PRO_ANNUAL=price_...
STRIPE_PRICE_ID_TEAM_SEAT=price_...  # Optional
STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true

# GCP Configuration
GCP_PROJECT_ID=your-project-id
GCP_REGION=us-central1
```

**Note**: `.env` is gitignored, so it won't be committed. Keep it secure!

## Option 2: Use `.env.local` for Everything

If you prefer to keep everything in `.env.local`, that works too. The scripts will read from it.

**Note**: Make sure `.env.local` has production values when deploying, or use separate files.

## Recommended Setup

### For Local Development
Keep `.env.local` with your local/test values:
```bash
# .env.local - Local development
PUBLIC_SUPABASE_URL=https://your-dev-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-dev-anon-key
# ... local values
```

### For Production Deployment
Create `.env` with production values:
```bash
# .env - Production
PUBLIC_SUPABASE_URL=https://your-prod-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-prod-anon-key
STRIPE_SECRET_KEY=sk_live_...
# ... production values
```

## Using the Scripts

### Update Cloud Run Environment Variables

```bash
# Scripts will read from .env first, then .env.local
# Make sure your production values are in .env
./scripts/update-cloud-run-env.sh
```

### Setup Stripe Secrets in Secret Manager

```bash
# Reads from .env or .env.local
./scripts/setup-stripe-secrets.sh
```

### Deploy

```bash
# Makefile reads from .env (or .env.local if .env doesn't exist)
make deploy-backend
```

## Quick Start

1. **Copy your production values to `.env`:**
   ```bash
   # You can start from .env.example if it exists
   cp .env.example .env
   
   # Or create it manually with your production values
   ```

2. **Add production Stripe keys:**
   ```bash
   # Edit .env and add:
   STRIPE_SECRET_KEY=sk_live_...
   STRIPE_PRICE_ID_PRO=price_...
   STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
   STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true
   ```

3. **Update Cloud Run:**
   ```bash
   ./scripts/update-cloud-run-env.sh
   ```

## Security Best Practices

1. **Never commit `.env` or `.env.local`** (already in `.gitignore`)
2. **Use Secret Manager for sensitive values** (Stripe keys):
   ```bash
   ./scripts/setup-stripe-secrets.sh
   ```
3. **Keep production values separate** from local development values
4. **Use different Supabase projects** for dev vs production if possible

## Troubleshooting

**Scripts not reading my values?**
- Check that `.env` or `.env.local` exists
- Verify variable names match exactly (case-sensitive)
- Check for syntax errors (no spaces around `=`)

**Wrong values being used?**
- Remember: `.env.local` overrides `.env`
- Check which file has the values you want

**Variables not persisting in Cloud Run?**
- See `docs/CLOUD_RUN_ENV_PERSISTENCE.md`
- Use `./scripts/update-cloud-run-env.sh` which preserves existing vars

