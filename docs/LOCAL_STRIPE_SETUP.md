# Local Stripe Development Setup

## Quick Checklist

### ✅ Environment Variables (Already Set)
Your `.env.local` file has all required Stripe variables:
- `STRIPE_SECRET_KEY` ✅
- `STRIPE_WEBHOOK_SECRET` ✅
- `STRIPE_PRICE_ID_PRO` ✅
- `STRIPE_PRICE_ID_PRO_ANNUAL` ✅
- `STRIPE_PRICE_ID_TEAM_SEAT` ✅
- `STRIPE_SUCCESS_URL` ✅
- `STRIPE_CANCEL_URL` ✅

### ⚠️ Critical: Webhook Forwarding for Local Development

**Stripe webhooks cannot reach `localhost` directly.** You MUST use Stripe CLI to forward webhooks to your local server.

#### Step 1: Install Stripe CLI (if not installed)

```bash
# macOS
brew install stripe/stripe-cli/stripe

# Or download from: https://stripe.com/docs/stripe-cli
```

#### Step 2: Login to Stripe CLI

```bash
stripe login
```

This will open your browser to authenticate with Stripe.

#### Step 3: Forward Webhooks to Local Server

**IMPORTANT:** Run this in a separate terminal window and keep it running:

```bash
stripe listen --forward-to localhost:8080/api/stripe/webhook
```

This will:
- Listen for webhook events from Stripe
- Forward them to your local API server
- Display a webhook signing secret (starts with `whsec_`)

#### Step 4: Update Webhook Secret (if needed)

If Stripe CLI shows a different webhook secret than what's in your `.env.local`, update it:

```bash
# Copy the webhook secret from Stripe CLI output
# It will look like: whsec_xxxxxxxxxxxxx

# Update your .env.local file:
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxx
```

**Note:** The webhook secret from Stripe CLI is different from the one in Stripe Dashboard. Use the CLI secret for local development.

#### Step 5: Restart Your API Server

After updating the webhook secret, restart your API server:

```bash
go run . api --port 8080
```

## Testing the Flow

1. **Start API Server** (Terminal 1):
   ```bash
   go run . api --port 8080
   ```

2. **Start Stripe CLI Webhook Forwarding** (Terminal 2):
   ```bash
   stripe listen --forward-to localhost:8080/api/stripe/webhook
   ```

3. **Start Frontend** (Terminal 3):
   ```bash
   cd web && npm run dev
   ```

4. **Test Payment Flow**:
   - Go to billing page
   - Click "Upgrade to Pro"
   - Use test card: `4242 4242 4242 4242`
   - Complete checkout
   - Check Terminal 2 for webhook events
   - Refresh billing page - should show updated plan

## Verifying Webhooks Are Working

### Check API Server Logs

When a webhook is received, you should see logs like:
```
INFO  HTTP request  method=POST path=/api/stripe/webhook status=200
INFO  Subscription updated  user_id=xxx subscription_id=sub_xxx tier=pro status=active
```

### Check Stripe CLI Output

You should see events like:
```
2025-11-14 02:15:23  --> checkout.session.completed [evt_xxx]
2025-11-14 02:15:24  --> customer.subscription.created [evt_xxx]
2025-11-14 02:15:24  --> customer.subscription.updated [evt_xxx]
```

### Check Database

After payment, verify the subscription was created:

```sql
-- Check subscriptions table
SELECT * FROM subscriptions ORDER BY created_at DESC LIMIT 1;

-- Check profiles table
SELECT id, subscription_tier, subscription_status, stripe_customer_id 
FROM profiles 
WHERE id = 'your-user-id';
```

## Troubleshooting

### Webhook Not Received

1. **Check Stripe CLI is running**: Must be running in separate terminal
2. **Check webhook secret**: Must match the one from Stripe CLI (not Dashboard)
3. **Check API server is running**: Must be on port 8080
4. **Check webhook endpoint**: Must be `/api/stripe/webhook`

### Subscription Not Updating

1. **Check API server logs**: Look for errors in webhook handler
2. **Check database**: Verify subscription record was created
3. **Check profile update**: Verify `subscription_tier` was updated in profiles table
4. **Refresh billing page**: May need to hard refresh (Cmd+Shift+R)

### Billing Page Not Showing Updated Plan

1. **Hard refresh**: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)
2. **Check browser console**: Look for errors loading subscription data
3. **Check API response**: Call `/api/v1/billing/summary` directly and check response
4. **Verify user ID**: Make sure the user_id matches between checkout and webhook

## Stripe FDW Wrapper (Optional)

You mentioned using the Stripe wrapper with Supabase. The FDW wrapper is for **querying** Stripe data from SQL, not for receiving webhooks. Webhooks still need to be handled by your Go backend.

The FDW wrapper setup is documented in:
- `docs/STRIPE_FDW_SETUP.md`
- `supabase/migrations/20250113_setup_stripe_fdw.sql`

This is optional and doesn't affect the webhook flow.

