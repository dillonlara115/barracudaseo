# Where to Add Production Variables

## Quick Answer

**Add production variables to `.env` file** (in the project root).

## File Structure

```
cli-scanner/
├── .env              ← ADD PRODUCTION VARIABLES HERE
├── .env.local       ← Keep local development values here
└── scripts/
    └── update-cloud-run-env.sh  ← Reads from .env
```

## Step-by-Step

### 1. Open `.env` file

```bash
# Edit the .env file in your project root
nano .env
# or
code .env
```

### 2. Add Production Stripe Variables

Add these to your `.env` file:

```bash
# Production Stripe Configuration
STRIPE_SECRET_KEY=sk_live_...              # Your LIVE Stripe secret key
STRIPE_WEBHOOK_SECRET=whsec_...            # Your webhook signing secret
STRIPE_PRICE_ID_PRO=price_...              # Monthly Pro plan price ID
STRIPE_PRICE_ID_PRO_ANNUAL=price_...       # Annual Pro plan price ID
STRIPE_PRICE_ID_TEAM_SEAT=price_...        # Team seat add-on (optional)
STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true

# Production Supabase (if different from local)
PUBLIC_SUPABASE_URL=https://your-prod-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-prod-anon-key

# GCP Configuration
GCP_PROJECT_ID=your-project-id
GCP_REGION=us-central1
```

### 3. Update Cloud Run

After adding variables to `.env`, run:

```bash
./scripts/update-cloud-run-env.sh
```

This script will:
1. Read variables from `.env`
2. Update Cloud Run with those variables
3. **Preserve existing variables** (doesn't replace everything)

## Why `.env` and not `.env.local`?

- **`.env`** = Production values (used by deployment scripts)
- **`.env.local`** = Local development values (overrides `.env` when developing locally)

The scripts (`update-cloud-run-env.sh`, `setup-stripe-secrets.sh`) read from `.env` first, then `.env.local`. For production deployment, use `.env`.

## Example Setup

### `.env` (Production)
```bash
STRIPE_SECRET_KEY=sk_live_51AbCdEf...
STRIPE_PRICE_ID_PRO=price_1SQX6II4GvFkgB3qgsZLKAgN
STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
```

### `.env.local` (Local Development)
```bash
STRIPE_SECRET_KEY=sk_test_51XyZwAb...  # Test key for local dev
STRIPE_PRICE_ID_PRO=price_test123      # Test price ID
STRIPE_SUCCESS_URL=http://localhost:8080/billing?success=true
```

When you run `./scripts/update-cloud-run-env.sh`, it uses values from `.env` (production).

When you develop locally, your code uses `.env.local` values (overrides `.env`).

## Verify It Worked

After running `./scripts/update-cloud-run-env.sh`, check Cloud Run:

```bash
gcloud run services describe barracuda-api \
  --platform managed \
  --region us-central1 \
  --format="value(spec.template.spec.containers[0].env)" | grep STRIPE
```

You should see your Stripe variables listed.

## Summary

| File | Purpose | When Used |
|------|---------|-----------|
| `.env` | **Production values** | Deployment scripts, Cloud Run updates |
| `.env.local` | Local development values | Local development (overrides `.env`) |

**Action**: Add production Stripe variables to `.env` file, then run `./scripts/update-cloud-run-env.sh`

