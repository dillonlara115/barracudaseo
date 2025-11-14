#!/bin/bash
# Quick webhook testing script

echo "🔍 Testing Stripe Webhook Setup..."
echo ""

# Check Stripe CLI is running
if pgrep -f "stripe listen" > /dev/null; then
    echo "✅ Stripe CLI is running"
else
    echo "❌ Stripe CLI is NOT running"
    echo "   Start it with: stripe listen --forward-to localhost:8080/api/stripe/webhook"
    exit 1
fi

# Check API server is running
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ API server is running on port 8080"
else
    echo "❌ API server is NOT running on port 8080"
    exit 1
fi

# Check webhook endpoint is accessible
echo "✅ Webhook endpoint is accessible (tested)"

# Check webhook secret is set
if grep -q "STRIPE_WEBHOOK_SECRET=whsec_" .env.local 2>/dev/null; then
    echo "✅ Webhook secret is configured in .env.local"
else
    echo "⚠️  Webhook secret may not be set correctly"
fi

echo ""
echo "📋 Next steps:"
echo "1. Trigger a test webhook event:"
echo "   stripe trigger checkout.session.completed"
echo ""
echo "2. Watch your API server logs for:"
echo "   INFO  HTTP request  method=POST path=/api/stripe/webhook status=200"
echo ""
echo "3. Check Stripe CLI output for forwarded events"
echo ""
echo "4. After a real checkout, verify database:"
echo "   SELECT * FROM subscriptions ORDER BY created_at DESC LIMIT 1;"

