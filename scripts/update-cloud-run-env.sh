#!/bin/bash
# Update Cloud Run environment variables

set -e

# Load from .env if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Check required variables
if [ -z "$PUBLIC_SUPABASE_URL" ] || [ -z "$PUBLIC_SUPABASE_ANON_KEY" ]; then
    echo "Error: PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY must be set"
    echo "Either:"
    echo "  1. Set them in your .env file, or"
    echo "  2. Export them: export PUBLIC_SUPABASE_URL=... export PUBLIC_SUPABASE_ANON_KEY=..."
    exit 1
fi

# Get project and region from gcloud config or environment
GCP_PROJECT_ID=${GCP_PROJECT_ID:-$(gcloud config get-value project 2>/dev/null)}
GCP_REGION=${GCP_REGION:-us-central1}

echo "Updating Cloud Run service: barracuda-api"
echo "Project: $GCP_PROJECT_ID"
echo "Region: $GCP_REGION"
echo ""

# Build environment variables string
ENV_VARS="PUBLIC_SUPABASE_URL=$PUBLIC_SUPABASE_URL,PUBLIC_SUPABASE_ANON_KEY=$PUBLIC_SUPABASE_ANON_KEY"

# Add Stripe variables if they're set
if [ -n "$STRIPE_SECRET_KEY" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY"
fi
if [ -n "$STRIPE_WEBHOOK_SECRET" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_WEBHOOK_SECRET=$STRIPE_WEBHOOK_SECRET"
fi
if [ -n "$STRIPE_PRICE_ID_PRO" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_PRO=$STRIPE_PRICE_ID_PRO"
fi
if [ -n "$STRIPE_PRICE_ID_PRO_ANNUAL" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_PRO_ANNUAL=$STRIPE_PRICE_ID_PRO_ANNUAL"
fi
if [ -n "$STRIPE_PRICE_ID_TEAM_SEAT" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_TEAM_SEAT=$STRIPE_PRICE_ID_TEAM_SEAT"
fi
if [ -n "$STRIPE_SUCCESS_URL" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_SUCCESS_URL=$STRIPE_SUCCESS_URL"
fi
if [ -n "$STRIPE_CANCEL_URL" ]; then
    ENV_VARS="$ENV_VARS,STRIPE_CANCEL_URL=$STRIPE_CANCEL_URL"
fi

gcloud run services update barracuda-api \
    --platform managed \
    --region $GCP_REGION \
    --update-env-vars="$ENV_VARS" \
    --quiet

echo ""
echo "✓ Environment variables updated!"
echo ""
echo "Service URL:"
gcloud run services describe barracuda-api \
    --platform managed \
    --region $GCP_REGION \
    --format="value(status.url)"

