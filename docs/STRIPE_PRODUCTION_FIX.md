# Quick Fix: Stripe "Not Configured" Error in Production

If you're seeing the error `"Failed to create checkout session: Error: Stripe not configured"` in production, it means the Stripe environment variables are not set in your Cloud Run service.

## Quick Fix (Choose One Method)

### Method 1: Using the Update Script (Recommended)

1. **Add Stripe variables to your `.env` file** (or export them):

```bash
# Required for Stripe checkout
STRIPE_SECRET_KEY=sk_live_...  # Your production Stripe secret key
STRIPE_WEBHOOK_SECRET=whsec_...  # Your webhook signing secret
STRIPE_PRICE_ID_PRO=price_...  # Monthly Pro plan price ID
STRIPE_PRICE_ID_PRO_ANNUAL=price_...  # Annual Pro plan price ID
STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true
STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true

# Optional
STRIPE_PRICE_ID_TEAM_SEAT=price_...  # Team seat add-on (if using)
```

2. **Run the update script**:

```bash
./scripts/update-cloud-run-env.sh
```

This will update your Cloud Run service with all environment variables from your `.env` file.

### Method 2: Using gcloud CLI Directly

```bash
# Set your region (if not already set)
export GCP_REGION=us-central1  # or your region

# Update Cloud Run with Stripe variables
gcloud run services update barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --update-env-vars="STRIPE_SECRET_KEY=sk_live_...,STRIPE_WEBHOOK_SECRET=whsec_...,STRIPE_PRICE_ID_PRO=price_...,STRIPE_PRICE_ID_PRO_ANNUAL=price_...,STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true,STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true"
```

### Method 3: Using Google Cloud Console

1. Go to [Google Cloud Console](https://console.cloud.google.com)
2. Navigate to **Cloud Run** → **barracuda-api**
3. Click **Edit & Deploy New Revision**
4. Go to **Variables & Secrets** tab
5. Add the following environment variables:
   - `STRIPE_SECRET_KEY` = `sk_live_...`
   - `STRIPE_WEBHOOK_SECRET` = `whsec_...`
   - `STRIPE_PRICE_ID_PRO` = `price_...`
   - `STRIPE_PRICE_ID_PRO_ANNUAL` = `price_...`
   - `STRIPE_SUCCESS_URL` = `https://app.barracudaseo.com/billing?success=true`
   - `STRIPE_CANCEL_URL` = `https://app.barracudaseo.com/billing?canceled=true`
6. Click **Deploy**

## Getting Your Stripe Keys

### 1. Stripe Secret Key

1. Go to [Stripe Dashboard](https://dashboard.stripe.com)
2. Make sure you're in **Live mode** (toggle in top right)
3. Go to **Developers** → **API keys**
4. Copy the **Secret key** (starts with `sk_live_...`)

⚠️ **Important**: Use `sk_live_...` for production, not `sk_test_...`

### 2. Webhook Secret

1. Go to **Developers** → **Webhooks**
2. Find your webhook endpoint (or create one)
3. Click on the webhook
4. Copy the **Signing secret** (starts with `whsec_...`)

### 3. Price IDs

1. Go to **Products** in Stripe Dashboard
2. Click on your product (e.g., "Barracuda Pro Monthly")
3. Copy the **Price ID** (starts with `price_...`)

## Verify the Fix

After updating the environment variables:

1. **Wait a few seconds** for Cloud Run to update
2. **Test the upgrade flow**:
   - Go to your app's billing page
   - Click "Upgrade to Pro"
   - You should be redirected to Stripe checkout (not see an error)

3. **Check logs** (if still having issues):
```bash
gcloud run services logs read barracuda-api \
  --platform managed \
  --region us-central1 \
  --limit 50
```

Look for:
- ✅ `"Stripe initialized"` - Good!
- ❌ `"Stripe integration disabled"` - Variables still not set

## Required Variables Summary

| Variable | Description | Example |
|----------|-------------|---------|
| `STRIPE_SECRET_KEY` | **Required** - Stripe API secret key | `sk_live_51...` |
| `STRIPE_WEBHOOK_SECRET` | **Required** - Webhook signing secret | `whsec_...` |
| `STRIPE_PRICE_ID_PRO` | **Required** - Monthly Pro plan price ID | `price_1SQX6I...` |
| `STRIPE_PRICE_ID_PRO_ANNUAL` | **Required** - Annual Pro plan price ID | `price_1SQX6I...` |
| `STRIPE_SUCCESS_URL` | **Required** - Redirect URL after successful checkout | `https://app.barracudaseo.com/billing?success=true` |
| `STRIPE_CANCEL_URL` | **Required** - Redirect URL if checkout canceled | `https://app.barracudaseo.com/billing?canceled=true` |
| `STRIPE_PRICE_ID_TEAM_SEAT` | Optional - Team seat add-on price ID | `price_1SQX9L...` |

## Troubleshooting

### Still Getting "Stripe not configured" Error?

1. **Verify variables are set**:
```bash
gcloud run services describe barracuda-api \
  --platform managed \
  --region us-central1 \
  --format="value(spec.template.spec.containers[0].env)"
```

2. **Check for typos** in variable names (they're case-sensitive)

3. **Make sure you're using live keys** (`sk_live_...`) not test keys (`sk_test_...`)

4. **Redeploy the service** to ensure variables are loaded:
```bash
gcloud run services update barracuda-api \
  --platform managed \
  --region us-central1
```

### Variables Not Persisting?

If variables keep getting reset, check:
- You're updating the correct Cloud Run service
- You're using the correct region
- You have permissions to update the service

## Next Steps

After fixing the configuration:
1. ✅ Test the upgrade flow end-to-end
2. ✅ Verify webhooks are working (check Stripe Dashboard → Webhooks)
3. ✅ Test subscription cancellation flow
4. ✅ Monitor logs for any Stripe-related errors

For more details, see [STRIPE_SETUP.md](./STRIPE_SETUP.md).

