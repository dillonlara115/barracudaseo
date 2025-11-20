# Quick Guide: Persisting Environment Variables in Cloud Run

## The Problem

When you deploy with `--set-env-vars`, it **replaces** all environment variables. If you don't include all variables in the deployment command, they get removed.

## Solution 1: Use Secret Manager (Recommended for Sensitive Values)

**Best for**: Stripe keys, API keys, and other sensitive values

### Setup Secrets (One-Time)

```bash
# Add Stripe keys to your .env file, then run:
./scripts/setup-stripe-secrets.sh
```

This creates secrets in Secret Manager that persist across deployments.

### Update Deployment to Use Secrets

Update your `Makefile` or deployment command to reference secrets:

```bash
--set-secrets="SUPABASE_SERVICE_ROLE_KEY=supabase-service-role-key:latest,STRIPE_SECRET_KEY=stripe-secret-key:latest,STRIPE_WEBHOOK_SECRET=stripe-webhook-secret:latest"
```

**Benefits:**
- ✅ Persists across deployments automatically
- ✅ More secure
- ✅ Can update without redeploying

## Solution 2: Use Update Script (For Non-Sensitive Values)

**Best for**: URLs, price IDs, and other non-sensitive config

### Update Variables Without Redeploying

```bash
# Add variables to .env file, then:
./scripts/update-cloud-run-env.sh
```

This uses `--update-env-vars` which **merges** with existing variables (doesn't replace them).

**Note**: This only updates variables you specify - existing ones remain unchanged.

## Solution 3: Always Include All Variables in Deployment

**Best for**: When you want everything in one place

Make sure your `.env` file has ALL variables, then deploy:

```bash
make deploy-backend
```

The Makefile reads from `.env` and includes all variables in the deployment.

## Recommended Workflow

### Initial Setup

1. **Store sensitive values in Secret Manager:**
   ```bash
   ./scripts/setup-stripe-secrets.sh
   ```

2. **Add non-sensitive values to .env:**
   ```bash
   STRIPE_PRICE_ID_PRO=price_...
   STRIPE_PRICE_ID_PRO_ANNUAL=price_...
   STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
   STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true
   ```

3. **Deploy with secrets:**
   ```bash
   make deploy-backend
   ```

### Regular Updates

**To update Stripe keys:**
```bash
# Update secret (no redeploy needed)
echo -n "new_key" | gcloud secrets versions add stripe-secret-key --data-file=-
```

**To update other variables:**
```bash
# Edit .env file, then:
./scripts/update-cloud-run-env.sh
```

**To deploy new code:**
```bash
make deploy-backend  # Secrets persist automatically
```

## Quick Reference

| Task | Command |
|------|---------|
| Setup Stripe secrets | `./scripts/setup-stripe-secrets.sh` |
| Update env vars | `./scripts/update-cloud-run-env.sh` |
| Deploy with all vars | `make deploy-backend` |
| Check current vars | `gcloud run services describe barracuda-api --format="value(spec.template.spec.containers[0].env)"` |
| Update a secret | `echo -n "value" \| gcloud secrets versions add SECRET_NAME --data-file=-` |

## Troubleshooting

**Variables disappear after deployment?**
- Use `./scripts/update-cloud-run-env.sh` instead of deploying
- Or ensure `.env` has ALL variables before deploying

**Secrets not working?**
- Check secret exists: `gcloud secrets list`
- Verify service account has access
- See `docs/CLOUD_RUN_ENV_PERSISTENCE.md` for details

