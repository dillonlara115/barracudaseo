#!/bin/bash
# Update Cloud Run environment variables
# This script uses --update-env-vars which MERGES with existing variables
# (doesn't replace them, so existing vars persist)
#
# Usage:
#   ./scripts/update-cloud-run-env.sh           # Uses .env only (production-safe default)
#   ./scripts/update-cloud-run-env.sh --local   # Uses .env, then .env.local (for testing)

set -e

# Check for --local flag (default is production mode - skip .env.local)
LOAD_LOCAL=false
if [ "$1" == "--local" ]; then
    LOAD_LOCAL=true
    echo "Local mode: Will load .env.local overrides"
fi

# Load from .env (production defaults)
if [ -f .env ]; then
    # Filter out comments (lines starting with #) and empty lines
    # Also remove inline comments (everything after # on a line)
    export $(cat .env | grep -v '^#' | grep -v '^[[:space:]]*$' | sed 's/#.*$//' | xargs)
fi

# Load from .env.local (local overrides) only if --local flag is set
# By default, skip .env.local to prevent accidental local config deployment
if [ "$LOAD_LOCAL" = true ] && [ -f .env.local ]; then
    echo "Loading local overrides from .env.local..."
    # Filter out comments and empty lines, remove inline comments
    export $(cat .env.local | grep -v '^#' | grep -v '^[[:space:]]*$' | sed 's/#.*$//' | xargs)
elif [ "$LOAD_LOCAL" = false ]; then
    echo "Production mode: Skipping .env.local (use --local flag to include it)"
fi

# Get project and region from gcloud config or environment
GCP_PROJECT_ID=${GCP_PROJECT_ID:-$(gcloud config get-value project 2>/dev/null)}
GCP_REGION=${GCP_REGION:-us-central1}

SERVICE_NAME=${SERVICE_NAME:-barracuda-api}

echo "Updating Cloud Run service: $SERVICE_NAME"
echo "Project: $GCP_PROJECT_ID"
echo "Region: $GCP_REGION"
echo ""
echo "Note: This merges with existing variables (doesn't replace them)"
echo ""

# Build environment variables string (only include variables that are set)
ENV_VARS=""

# Required variables
if [ -n "$PUBLIC_SUPABASE_URL" ]; then
    ENV_VARS="PUBLIC_SUPABASE_URL=$PUBLIC_SUPABASE_URL"
fi
if [ -n "$PUBLIC_SUPABASE_ANON_KEY" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,PUBLIC_SUPABASE_ANON_KEY=$PUBLIC_SUPABASE_ANON_KEY"
    else
        ENV_VARS="PUBLIC_SUPABASE_ANON_KEY=$PUBLIC_SUPABASE_ANON_KEY"
    fi
fi

# Stripe variables
if [ -n "$STRIPE_SECRET_KEY" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY"
    else
        ENV_VARS="STRIPE_SECRET_KEY=$STRIPE_SECRET_KEY"
    fi
fi
if [ -n "$STRIPE_WEBHOOK_SECRET" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_WEBHOOK_SECRET=$STRIPE_WEBHOOK_SECRET"
    else
        ENV_VARS="STRIPE_WEBHOOK_SECRET=$STRIPE_WEBHOOK_SECRET"
    fi
fi
if [ -n "$STRIPE_PRICE_ID_PRO" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_PRO=$STRIPE_PRICE_ID_PRO"
    else
        ENV_VARS="STRIPE_PRICE_ID_PRO=$STRIPE_PRICE_ID_PRO"
    fi
fi
if [ -n "$STRIPE_PRICE_ID_PRO_ANNUAL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_PRO_ANNUAL=$STRIPE_PRICE_ID_PRO_ANNUAL"
    else
        ENV_VARS="STRIPE_PRICE_ID_PRO_ANNUAL=$STRIPE_PRICE_ID_PRO_ANNUAL"
    fi
fi
if [ -n "$STRIPE_PRICE_ID_TEAM_SEAT" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_PRICE_ID_TEAM_SEAT=$STRIPE_PRICE_ID_TEAM_SEAT"
    else
        ENV_VARS="STRIPE_PRICE_ID_TEAM_SEAT=$STRIPE_PRICE_ID_TEAM_SEAT"
    fi
fi
if [ -n "$STRIPE_SUCCESS_URL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_SUCCESS_URL=$STRIPE_SUCCESS_URL"
    else
        ENV_VARS="STRIPE_SUCCESS_URL=$STRIPE_SUCCESS_URL"
    fi
fi
if [ -n "$STRIPE_CANCEL_URL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,STRIPE_CANCEL_URL=$STRIPE_CANCEL_URL"
    else
        ENV_VARS="STRIPE_CANCEL_URL=$STRIPE_CANCEL_URL"
    fi
fi

# Email variables (for Resend/Elastic Email)
if [ -n "$EMAIL_PROVIDER" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,EMAIL_PROVIDER=$EMAIL_PROVIDER"
    else
        ENV_VARS="EMAIL_PROVIDER=$EMAIL_PROVIDER"
    fi
fi
if [ -n "$RESEND_API_KEY" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,RESEND_API_KEY=$RESEND_API_KEY"
    else
        ENV_VARS="RESEND_API_KEY=$RESEND_API_KEY"
    fi
fi
if [ -n "$EMAIL_FROM_ADDRESS" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,EMAIL_FROM_ADDRESS=$EMAIL_FROM_ADDRESS"
    else
        ENV_VARS="EMAIL_FROM_ADDRESS=$EMAIL_FROM_ADDRESS"
    fi
fi
if [ -n "$APP_URL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,APP_URL=$APP_URL"
    else
        ENV_VARS="APP_URL=$APP_URL"
    fi
fi
if [ -n "$BETA_INVITE_CODE" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,BETA_INVITE_CODE=$BETA_INVITE_CODE"
    else
        ENV_VARS="BETA_INVITE_CODE=$BETA_INVITE_CODE"
    fi
fi
# GSC OAuth variables
if [ -n "$GSC_CLIENT_ID" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,GSC_CLIENT_ID=$GSC_CLIENT_ID"
    else
        ENV_VARS="GSC_CLIENT_ID=$GSC_CLIENT_ID"
    fi
fi
if [ -n "$GSC_CLIENT_SECRET" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,GSC_CLIENT_SECRET=$GSC_CLIENT_SECRET"
    else
        ENV_VARS="GSC_CLIENT_SECRET=$GSC_CLIENT_SECRET"
    fi
fi
if [ -n "$GSC_REDIRECT_URL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,GSC_REDIRECT_URL=$GSC_REDIRECT_URL"
    else
        ENV_VARS="GSC_REDIRECT_URL=$GSC_REDIRECT_URL"
    fi
fi

# DataForSEO variables
if [ -n "$DATAFORSEO_LOGIN" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,DATAFORSEO_LOGIN=$DATAFORSEO_LOGIN"
    else
        ENV_VARS="DATAFORSEO_LOGIN=$DATAFORSEO_LOGIN"
    fi
fi
if [ -n "$DATAFORSEO_PASSWORD" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,DATAFORSEO_PASSWORD=$DATAFORSEO_PASSWORD"
    else
        ENV_VARS="DATAFORSEO_PASSWORD=$DATAFORSEO_PASSWORD"
    fi
fi
if [ -n "$DATAFORSEO_BASE_URL" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,DATAFORSEO_BASE_URL=$DATAFORSEO_BASE_URL"
    else
        ENV_VARS="DATAFORSEO_BASE_URL=$DATAFORSEO_BASE_URL"
    fi
fi

# OpenAI (app-wide key for AI features)
if [ -n "$OPENAI_API_KEY" ]; then
    if [ -n "$ENV_VARS" ]; then
        ENV_VARS="$ENV_VARS,OPENAI_API_KEY=$OPENAI_API_KEY"
    else
        ENV_VARS="OPENAI_API_KEY=$OPENAI_API_KEY"
    fi
fi

if [ -z "$ENV_VARS" ]; then
    echo "Error: No environment variables to update."
    echo "Set variables in your .env file or export them."
    exit 1
fi

echo "Updating variables:"
echo "$ENV_VARS" | tr ',' '\n' | sed 's/^/  - /'
echo ""

gcloud run services update $SERVICE_NAME \
    --platform managed \
    --region $GCP_REGION \
    --update-env-vars="$ENV_VARS" \
    --quiet

echo ""
echo "✓ Environment variables updated!"
echo ""
echo "Service URL:"
gcloud run services describe $SERVICE_NAME \
    --platform managed \
    --region $GCP_REGION \
    --format="value(status.url)"
