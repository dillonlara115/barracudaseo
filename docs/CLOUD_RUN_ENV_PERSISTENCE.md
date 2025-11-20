# Persisting Environment Variables in Cloud Run

The issue: When you deploy with `--set-env-vars`, it **replaces** all environment variables. If you don't include all variables in the deployment command, they get removed.

## Best Practice: Use Secret Manager for Sensitive Values

**Recommended approach**: Store sensitive values (like Stripe keys) in Secret Manager, which persists across deployments.

### Step 1: Store Stripe Secrets in Secret Manager

```bash
# Store Stripe secret key
echo -n "sk_live_..." | gcloud secrets create stripe-secret-key \
  --data-file=- \
  --replication-policy="automatic"

# Store Stripe webhook secret
echo -n "whsec_..." | gcloud secrets create stripe-webhook-secret \
  --data-file=- \
  --replication-policy="automatic"
```

### Step 2: Update Your Deployment to Use Secrets

When deploying, reference secrets instead of plain environment variables:

```bash
gcloud run deploy barracuda-api \
  --image $GCP_REGION-docker.pkg.dev/$GCP_PROJECT_ID/barracuda/barracuda-api:latest \
  --platform managed \
  --region $GCP_REGION \
  --allow-unauthenticated \
  --set-env-vars="PUBLIC_SUPABASE_URL=...,PUBLIC_SUPABASE_ANON_KEY=...,STRIPE_PRICE_ID_PRO=...,STRIPE_PRICE_ID_PRO_ANNUAL=...,STRIPE_SUCCESS_URL=...,STRIPE_CANCEL_URL=..." \
  --set-secrets="SUPABASE_SERVICE_ROLE_KEY=supabase-service-role-key:latest,STRIPE_SECRET_KEY=stripe-secret-key:latest,STRIPE_WEBHOOK_SECRET=stripe-webhook-secret:latest" \
  --memory=512Mi \
  --cpu=1 \
  --timeout=300 \
  --max-instances=10 \
  --port=8080
```

**Benefits:**
- ✅ Secrets persist across deployments
- ✅ More secure (not visible in deployment commands)
- ✅ Can update secrets without redeploying
- ✅ Version controlled (can reference specific secret versions)

## Alternative: Use `--update-env-vars` Instead of `--set-env-vars`

If you prefer to keep using environment variables (not secrets), use `--update-env-vars` which **merges** instead of replacing:

```bash
# This merges with existing variables instead of replacing them
gcloud run services update barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --update-env-vars="STRIPE_SECRET_KEY=sk_live_...,STRIPE_SUCCESS_URL=..."
```

**Note**: `--update-env-vars` only works with `gcloud run services update`, not `gcloud run deploy`.

## Solution: Update Deployment Scripts

The Makefile and deployment scripts should:
1. Use `--update-env-vars` when updating existing services
2. Include ALL environment variables when doing a fresh deployment
3. Prefer Secret Manager for sensitive values

## Recommended Workflow

### Initial Setup (One-Time)

1. **Store secrets in Secret Manager:**
```bash
echo -n "sk_live_..." | gcloud secrets create stripe-secret-key --data-file=-
echo -n "whsec_..." | gcloud secrets create stripe-webhook-secret --data-file=-
```

2. **Deploy with all variables:**
```bash
make deploy-backend  # Uses .env file and includes all variables
```

### Regular Deployments

**Option A: Use the update script (recommended for env vars)**
```bash
# Add/update variables in .env file
./scripts/update-cloud-run-env.sh
```

**Option B: Always include all variables in deployment**
```bash
# Make sure .env has all variables, then:
make deploy-backend
```

**Option C: Use Secret Manager (recommended for sensitive values)**
```bash
# Secrets persist automatically, just deploy:
make deploy-backend
```

## Updating Secrets Without Redeploying

If you need to update a secret value:

```bash
# Update the secret
echo -n "new_secret_value" | gcloud secrets versions add stripe-secret-key --data-file=-

# Cloud Run will automatically use the latest version
# Or specify a version explicitly:
gcloud run services update barracuda-api \
  --update-secrets="STRIPE_SECRET_KEY=stripe-secret-key:2" \
  --region $GCP_REGION
```

## Verifying Environment Variables

Check what's currently set:

```bash
gcloud run services describe barracuda-api \
  --platform managed \
  --region us-central1 \
  --format="value(spec.template.spec.containers[0].env)"
```

## Troubleshooting

### Variables Disappear After Deployment?

**Cause**: Using `--set-env-vars` without including all variables.

**Fix**: 
- Use `--update-env-vars` instead, OR
- Always include ALL variables in deployment command, OR
- Use Secret Manager for values that shouldn't change

### Can't Update Variables?

**Cause**: Using `--set-env-vars` with `gcloud run deploy`.

**Fix**: Use `gcloud run services update` with `--update-env-vars`:

```bash
gcloud run services update barracuda-api \
  --update-env-vars="NEW_VAR=value" \
  --region us-central1
```

### Secrets Not Working?

**Check:**
1. Secret exists: `gcloud secrets list`
2. Service account has access: `gcloud run services describe barracuda-api --format="value(spec.template.spec.serviceAccountName)"`
3. Grant access if needed: `gcloud secrets add-iam-policy-binding stripe-secret-key --member="serviceAccount:PROJECT_NUMBER-compute@developer.gserviceaccount.com" --role="roles/secretmanager.secretAccessor"`

