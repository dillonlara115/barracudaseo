# Environment Files Guide

## File Structure

Your project supports two environment files:

1. **`.env`** - Production/default values
2. **`.env.local`** - Local development overrides

## Loading Order

Environment variables are loaded in this order:
1. `.env` (if exists)
2. `.env.local` (if exists) - **overrides** `.env` values

This matches how the Go code (`cmd/api.go`) loads environment variables.

## Usage

### For Local Development

Keep your local/test values in `.env.local`:
```bash
# .env.local
PUBLIC_SUPABASE_URL=https://your-dev-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-dev-anon-key
# ... local values
```

### For Production Deployment

Create `.env` with production values:
```bash
# .env
PUBLIC_SUPABASE_URL=https://your-prod-project.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-prod-anon-key
STRIPE_SECRET_KEY=sk_live_...
# ... production values
```

## Scripts Behavior

All scripts (`update-cloud-run-env.sh`, `setup-stripe-secrets.sh`) and the Makefile (`make deploy-backend`) will:

1. Load `.env` first (if it exists)
2. Then load `.env.local` (if it exists), overriding any matching variables

## Examples

### Scenario 1: Local Development Only

You only have `.env.local`:
```bash
# .env.local
PUBLIC_SUPABASE_URL=https://dev-project.supabase.co
```

Scripts will read from `.env.local` ✅

### Scenario 2: Production Deployment

You have both files:
```bash
# .env (production)
PUBLIC_SUPABASE_URL=https://prod-project.supabase.co
STRIPE_SECRET_KEY=sk_live_...

# .env.local (local overrides)
PUBLIC_SUPABASE_URL=https://dev-project.supabase.co
```

When running scripts locally, `.env.local` values will be used (overrides `.env`).

**For production deployment**, you can:
- Temporarily rename `.env.local` to `.env.local.backup`
- Run deployment (uses `.env` production values)
- Restore `.env.local`

Or better: Use the scripts which read from `.env` for production values.

### Scenario 3: Production Values Only

You only have `.env`:
```bash
# .env
PUBLIC_SUPABASE_URL=https://prod-project.supabase.co
STRIPE_SECRET_KEY=sk_live_...
```

Scripts will read from `.env` ✅

## Best Practice

**Recommended setup:**

1. **Create `.env`** with production values (for Cloud Run deployment)
2. **Create `.env.local`** with local development values
3. **Use `.env.local`** for local development (it overrides `.env`)
4. **Use `.env`** for production deployments

This way:
- Production values are in `.env` (used by deployment scripts)
- Local overrides are in `.env.local` (used when developing locally)
- Both files are gitignored (secure)

## Quick Commands

```bash
# Check what files exist
ls -la | grep "\.env"

# Create .env from example (if exists)
cp .env.example .env

# Update Cloud Run (reads from .env, then .env.local)
./scripts/update-cloud-run-env.sh

# Deploy (reads from .env, then .env.local)
make deploy-backend
```

## Troubleshooting

**Scripts using wrong values?**
- Check which file has the values you want
- Remember: `.env.local` overrides `.env`
- For production, ensure values are in `.env`

**Variables not loading?**
- Check file exists: `ls -la .env .env.local`
- Check syntax (no spaces around `=`)
- Check variable names match exactly

**Want to use only production values?**
- Temporarily rename `.env.local`: `mv .env.local .env.local.backup`
- Run your script
- Restore: `mv .env.local.backup .env.local`

