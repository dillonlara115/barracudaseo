# Testing Stripe Webhooks Locally

## Quick Verification Checklist

### 1. Check Stripe CLI is Running

```bash
# Check if Stripe CLI is running
ps aux | grep "stripe listen" | grep -v grep

# Should show something like:
# dillonlara  27773  stripe listen --forward-to localhost:8080/api/stripe/webhook
```

If not running, start it:
```bash
stripe listen --forward-to localhost:8080/api/stripe/webhook
```

### 2. Verify Webhook Secret

When Stripe CLI starts, it displays a webhook signing secret. Make sure this matches your `.env.local`:

```bash
# Check your .env.local
grep STRIPE_WEBHOOK_SECRET .env.local

# Should match the secret shown by Stripe CLI (starts with whsec_)
```

**Important:** The webhook secret from Stripe CLI is different from the one in Stripe Dashboard. Use the CLI secret for local development.

### 3. Test Webhook Endpoint Directly

```bash
# Test the webhook endpoint is accessible
curl -X POST http://localhost:8080/api/stripe/webhook \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}'

# Should return an error about missing signature (expected)
# This confirms the endpoint is reachable
```

### 4. Trigger a Test Webhook Event

Use Stripe CLI to trigger a test event:

```bash
# Trigger a checkout.session.completed event
stripe trigger checkout.session.completed

# Trigger a customer.subscription.created event
stripe trigger customer.subscription.created

# Trigger a customer.subscription.updated event
stripe trigger customer.subscription.updated
```

### 5. Check API Server Logs

Watch your API server logs for webhook events. You should see:

```
INFO  HTTP request  method=POST path=/api/stripe/webhook status=200
INFO  Subscription updated  user_id=xxx subscription_id=sub_xxx tier=pro status=active
```

### 6. Check Stripe CLI Output

The Stripe CLI terminal should show events being forwarded:

```
2025-11-14 19:15:23  --> checkout.session.completed [evt_xxx]
2025-11-14 19:15:24  --> customer.subscription.created [evt_xxx]
2025-11-14 19:15:24  --> customer.subscription.updated [evt_xxx]
```

## Detailed Testing Steps

### Step 1: Verify Stripe CLI Connection

```bash
# Start Stripe CLI (if not already running)
stripe listen --forward-to localhost:8080/api/stripe/webhook

# You should see:
# > Ready! Your webhook signing secret is whsec_xxxxxxxxxxxxx
# > Forwarding events to http://localhost:8080/api/stripe/webhook
```

### Step 2: Update Webhook Secret

Copy the webhook secret from Stripe CLI and update your `.env.local`:

```bash
# Edit .env.local and update:
STRIPE_WEBHOOK_SECRET=whsec_xxxxxxxxxxxxx  # Use the secret from Stripe CLI

# Restart your API server after updating
```

### Step 3: Test with a Real Checkout Flow

1. **Start all services:**
   ```bash
   # Terminal 1: API Server
   go run . api --port 8080
   
   # Terminal 2: Stripe CLI
   stripe listen --forward-to localhost:8080/api/stripe/webhook
   
   # Terminal 3: Frontend
   cd web && npm run dev
   ```

2. **Complete a test checkout:**
   - Go to billing page
   - Click "Upgrade to Pro"
   - Use test card: `4242 4242 4242 4242`
   - Complete checkout

3. **Watch for webhook events:**
   - Check Terminal 2 (Stripe CLI) for forwarded events
   - Check Terminal 1 (API Server) for processing logs

### Step 4: Verify Database Updates

After a successful webhook, check the database:

```sql
-- Check subscriptions table
SELECT * FROM subscriptions 
ORDER BY created_at DESC 
LIMIT 5;

-- Check profiles table
SELECT id, subscription_tier, subscription_status, stripe_customer_id, stripe_subscription_id
FROM profiles
WHERE stripe_customer_id IS NOT NULL
ORDER BY updated_at DESC
LIMIT 5;
```

### Step 5: Test Webhook Error Handling

Test that webhook signature validation works:

```bash
# Test with invalid signature (should fail)
curl -X POST http://localhost:8080/api/stripe/webhook \
  -H "Content-Type: application/json" \
  -H "Stripe-Signature: invalid_signature" \
  -d '{"type": "test.event", "data": {}}'

# Should return 400 Bad Request
```

## Common Issues and Solutions

### Issue: Webhook Secret Mismatch

**Symptoms:**
- API server logs show "Webhook signature verification failed"
- Webhook returns 400 Bad Request

**Solution:**
1. Get the webhook secret from Stripe CLI output
2. Update `STRIPE_WEBHOOK_SECRET` in `.env.local`
3. Restart API server

### Issue: Webhook Not Received

**Symptoms:**
- Stripe CLI shows events forwarded
- But API server doesn't log anything

**Solution:**
1. Check API server is running on port 8080
2. Check webhook endpoint is `/api/stripe/webhook`
3. Check CORS settings (webhook endpoint should allow all origins)

### Issue: Subscription Not Created

**Symptoms:**
- Webhook received successfully (200 OK)
- But subscription not in database

**Solution:**
1. Check API server logs for errors
2. Verify `getUserIDByStripeCustomerID` is finding the user
3. Check database connection
4. Verify Stripe customer ID matches between checkout and webhook

### Issue: Profile Not Updated

**Symptoms:**
- Subscription created in `subscriptions` table
- But `profiles.subscription_tier` still shows "free"

**Solution:**
1. Check database trigger `sync_subscription_to_profile_trigger` exists
2. Verify trigger function `sync_subscription_to_profile()` is working
3. Check for errors in database logs

## Automated Testing Script

Create a test script to verify webhook setup:

```bash
#!/bin/bash
# test-webhook.sh

echo "Testing Stripe Webhook Setup..."

# Check Stripe CLI is running
if ! pgrep -f "stripe listen" > /dev/null; then
    echo "❌ Stripe CLI is not running"
    echo "   Start it with: stripe listen --forward-to localhost:8080/api/stripe/webhook"
    exit 1
fi
echo "✅ Stripe CLI is running"

# Check API server is running
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ API server is not running on port 8080"
    exit 1
fi
echo "✅ API server is running"

# Check webhook endpoint is accessible
if ! curl -s -X POST http://localhost:8080/api/stripe/webhook \
    -H "Content-Type: application/json" \
    -d '{"test": "data"}' | grep -q "signature"; then
    echo "⚠️  Webhook endpoint may not be configured correctly"
else
    echo "✅ Webhook endpoint is accessible"
fi

# Trigger a test event
echo ""
echo "Triggering test webhook event..."
stripe trigger checkout.session.completed

echo ""
echo "✅ Test complete. Check your API server logs for webhook processing."
```

## Monitoring Webhooks in Real-Time

### Watch API Server Logs

```bash
# If using structured logging, filter for webhook events
# Or just watch all logs:
tail -f /path/to/api-server.log | grep -i webhook
```

### Watch Stripe CLI Output

The Stripe CLI terminal shows all forwarded events in real-time. Keep it visible while testing.

### Database Monitoring

```sql
-- Watch for new subscriptions
SELECT * FROM subscriptions 
WHERE created_at > NOW() - INTERVAL '1 minute'
ORDER BY created_at DESC;

-- Watch for profile updates
SELECT id, subscription_tier, updated_at 
FROM profiles 
WHERE updated_at > NOW() - INTERVAL '1 minute'
ORDER BY updated_at DESC;
```

## Success Criteria

Your webhook setup is working correctly if:

1. ✅ Stripe CLI is forwarding events
2. ✅ API server receives and processes webhooks (200 OK)
3. ✅ Subscriptions are created in database
4. ✅ Profiles are updated with subscription info
5. ✅ Billing page shows updated subscription after payment

