# Debugging Checkout 500 Error

If you're seeing a 500 error when clicking the "Upgrade" button, follow these steps to identify the root cause.

## Step 1: Check Application Logs

The HTTP request log you showed only indicates a 500 status. To see the actual error message, check the application logs:

```bash
# View recent logs with error details
gcloud run services logs read barracuda-api \
  --platform managed \
  --region us-central1 \
  --limit 100 \
  --format json | jq '.[] | select(.severity == "ERROR") | {timestamp, textPayload, jsonPayload}'

# Or view all recent logs
gcloud run services logs read barracuda-api \
  --platform managed \
  --region us-central1 \
  --limit 100
```

Look for log entries with:
- `"Failed to create checkout session"`
- `"Stripe not configured"`
- `"Checkout URLs not configured"`
- `"Failed to get user profile"`
- `"Failed to get user email"`
- `"Failed to create Stripe customer"`

## Step 2: Common Error Causes

### 1. Missing Stripe Secret Key
**Error in logs:** `"Stripe secret key not configured"`  
**Fix:** Set `STRIPE_SECRET_KEY` environment variable in Cloud Run

### 2. Missing Success/Cancel URLs
**Error in logs:** `"Stripe success/cancel URLs not configured"`  
**Fix:** Set both `STRIPE_SUCCESS_URL` and `STRIPE_CANCEL_URL` environment variables

### 3. Database Connection Issue
**Error in logs:** `"Failed to get user profile"` or `"Failed to create user profile"`  
**Fix:** Check Supabase connection and service role key

### 4. Auth API Issue
**Error in logs:** `"Failed to get user email"`  
**Fix:** Verify `PUBLIC_SUPABASE_URL` and `SUPABASE_SERVICE_ROLE_KEY` are correct

### 5. Stripe API Error
**Error in logs:** `"Failed to create checkout session"` with Stripe error details  
**Fix:** Check Stripe error code and message in logs

## Step 3: Verify Environment Variables

Check which Stripe-related environment variables are set:

```bash
gcloud run services describe barracuda-api \
  --platform managed \
  --region us-central1 \
  --format="value(spec.template.spec.containers[0].env)" | grep STRIPE
```

Required variables:
- `STRIPE_SECRET_KEY` - Your Stripe secret key (starts with `sk_live_...`)
- `STRIPE_PRICE_ID_PRO` - Monthly Pro plan price ID
- `STRIPE_PRICE_ID_PRO_ANNUAL` - Annual Pro plan price ID  
- `STRIPE_SUCCESS_URL` - Redirect URL after successful checkout
- `STRIPE_CANCEL_URL` - Redirect URL if checkout is canceled

## Step 4: Test the Fix

After updating environment variables:

1. Wait 10-30 seconds for Cloud Run to update
2. Try the upgrade button again
3. Check logs again if it still fails

## Step 5: Enhanced Error Messages

The code has been updated to include more detailed error logging. When you check logs, you should now see:
- User ID
- Price ID being used
- Which specific configuration is missing
- Stripe error codes and messages (if applicable)

## Quick Fix Script

If you have your Stripe keys ready, you can update Cloud Run directly:

```bash
# Set your region
export GCP_REGION=us-central1

# Update with your actual values
gcloud run services update barracuda-api \
  --platform managed \
  --region $GCP_REGION \
  --update-env-vars="STRIPE_SECRET_KEY=sk_live_...,STRIPE_PRICE_ID_PRO=price_...,STRIPE_PRICE_ID_PRO_ANNUAL=price_...,STRIPE_SUCCESS_URL=https://app.barracudaseo.com/billing?success=true,STRIPE_CANCEL_URL=https://app.barracudaseo.com/billing?canceled=true"
```

## Still Having Issues?

If the error persists after checking logs and verifying configuration:

1. **Check the exact error message** from application logs (not just HTTP logs)
2. **Verify Stripe keys are live keys** (`sk_live_...`) not test keys (`sk_test_...`)
3. **Check Stripe Dashboard** to ensure price IDs are correct
4. **Verify URLs** are accessible and use HTTPS in production
5. **Check Supabase connection** - ensure service role key is valid

