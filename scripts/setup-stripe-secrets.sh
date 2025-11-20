#!/bin/bash
# Setup Stripe secrets in Google Cloud Secret Manager
# This is the recommended approach for sensitive values that should persist

set -e

# Load from .env.local first (local overrides), then .env (production defaults)
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi
if [ -f .env.local ]; then
    export $(cat .env.local | grep -v '^#' | xargs)
fi

# Get project from gcloud config or environment
GCP_PROJECT_ID=${GCP_PROJECT_ID:-$(gcloud config get-value project 2>/dev/null)}

if [ -z "$GCP_PROJECT_ID" ]; then
    echo "Error: GCP_PROJECT_ID not set"
    exit 1
fi

echo "Setting up Stripe secrets in Secret Manager"
echo "Project: $GCP_PROJECT_ID"
echo ""

# Check if secrets already exist
SECRET_EXISTS=$(gcloud secrets list --filter="name:stripe-secret-key" --format="value(name)" 2>/dev/null || echo "")

if [ -n "$SECRET_EXISTS" ]; then
    echo "Secret 'stripe-secret-key' already exists."
    read -p "Update it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -z "$STRIPE_SECRET_KEY" ]; then
            echo "Error: STRIPE_SECRET_KEY not set in .env file"
            exit 1
        fi
        echo -n "$STRIPE_SECRET_KEY" | gcloud secrets versions add stripe-secret-key --data-file=-
        echo "✓ Updated stripe-secret-key"
    fi
else
    if [ -z "$STRIPE_SECRET_KEY" ]; then
        echo "Error: STRIPE_SECRET_KEY not set in .env file"
        exit 1
    fi
    echo -n "$STRIPE_SECRET_KEY" | gcloud secrets create stripe-secret-key \
        --data-file=- \
        --replication-policy="automatic"
    echo "✓ Created stripe-secret-key"
fi

# Webhook secret
WEBHOOK_EXISTS=$(gcloud secrets list --filter="name:stripe-webhook-secret" --format="value(name)" 2>/dev/null || echo "")

if [ -n "$WEBHOOK_EXISTS" ]; then
    echo "Secret 'stripe-webhook-secret' already exists."
    read -p "Update it? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        if [ -z "$STRIPE_WEBHOOK_SECRET" ]; then
            echo "Warning: STRIPE_WEBHOOK_SECRET not set, skipping..."
        else
            echo -n "$STRIPE_WEBHOOK_SECRET" | gcloud secrets versions add stripe-webhook-secret --data-file=-
            echo "✓ Updated stripe-webhook-secret"
        fi
    fi
else
    if [ -n "$STRIPE_WEBHOOK_SECRET" ]; then
        echo -n "$STRIPE_WEBHOOK_SECRET" | gcloud secrets create stripe-webhook-secret \
            --data-file=- \
            --replication-policy="automatic"
        echo "✓ Created stripe-webhook-secret"
    else
        echo "Warning: STRIPE_WEBHOOK_SECRET not set, skipping..."
    fi
fi

echo ""
echo "✓ Stripe secrets setup complete!"
echo ""
echo "Next steps:"
echo "1. Update your deployment to use secrets:"
echo "   --set-secrets=\"STRIPE_SECRET_KEY=stripe-secret-key:latest,STRIPE_WEBHOOK_SECRET=stripe-webhook-secret:latest\""
echo ""
echo "2. Or update existing service:"
echo "   gcloud run services update barracuda-api \\"
echo "     --update-secrets=\"STRIPE_SECRET_KEY=stripe-secret-key:latest,STRIPE_WEBHOOK_SECRET=stripe-webhook-secret:latest\" \\"
echo "     --region us-central1"

